package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/olorin/nagiosplugin"
	"github.com/thedevsaddam/gojsonq"
)

func main() {

	var daysRemaining float64
	var prefix string = "https://"

	hostname := flag.String("hostname", "", "Hostname of the Citrix ADC server")
	username := flag.String("username", "", "Username to access the Nitro API")
	password := flag.String("password", "", "Password to access the Nitro API")
	warning := flag.String("warning", "", "The range for warning status. For specification please refer to the Nagios docs.")
	critical := flag.String("critical", "", "The range for critical status. For specification please refer to the Nagios docs.")
	secure := flag.Bool("secure", false, "Use HTTPS to access the Nitro API")
	testvalue := flag.Float64("testvalue", -1.0, "Pass a value to override (used for testing)")

	flag.Parse()

	if !*secure {
		prefix = "http://"
	}

	if *testvalue >= 0 {

		daysRemaining = *testvalue

	} else {

		url := fmt.Sprintf("%s%s/nitro/v1/config/nslicense", prefix, *hostname)
		json := downloadJson(url, *username, *password)
		result := gojsonq.New().FromString(json).Find("nslicense.daystoexpiration")
		daysRemaining, _ = strconv.ParseFloat(result.(string), 64)

	}

	warningRange, err := nagiosplugin.ParseRange(*warning)

	if err != nil {
		fmt.Printf("Error parsing warning range: %s\r\n", err)
		os.Exit(3)
	}

	criticalRange, _ := nagiosplugin.ParseRange(*critical)

	if err != nil {
		fmt.Printf("Error parsing critical range: %s\r\n", err)
		os.Exit(3)
	}

	check := nagiosplugin.NewCheck()
	defer check.Finish()

	//check.AddPerfDatum("days_remaining", "", daysRemaining, 0.0, math.Inf(1), -60, -10)

	check.AddResult(nagiosplugin.OK, fmt.Sprintf("Your license will expire in %v day(s)", daysRemaining))

	if criticalRange.Check(daysRemaining) {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Your license will expire in less than %v days", criticalRange.End))
	}

	if warningRange.Check(daysRemaining) {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Your license will expire in less than %v days", warningRange.End))
	}
}

func downloadJson(url, username, password string) string {

	var client http.Client

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Could not download JSON from %s: %s", url, err)
		os.Exit(3)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		return bodyString
	}

	return ""
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
