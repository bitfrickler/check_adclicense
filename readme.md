# check_adclicense

## Description

check_adclicense is a Nagios check plugin to monitor the remaining days on your Citrix ADC license.
It uses the [Nitro API](https://docs.citrix.com/en-us/citrix-adc/current-release/nitro-api.html) to access license information.

For details on how to format a Nagios-specific range, please refer to the [Nagios documentation](https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT)

This project uses the very cool [nagiosplugin](https://pkg.go.dev/github.com/olorin/nagiosplugin) for Go.

## Usage

```cli
  -critical string
        The range for critical status. For specification please refer to the Nagios docs.
  -hostname string
        Hostname of the Citrix ADC server
  -password string
        Password to access the Nitro API
  -secure
        Use HTTPS to access the Nitro API
  -testvalue float
        Pass a value to override (used for testing) (default -1)
  -username string
        Username to access the Nitro API
  -warning string
        The range for warning status. For specification please refer to the Nagios docs.
```

## Precompiled executables

The bin folder contains executables for AMD64 on Windows, Linux and Darwin (Mac OS).
The subfolder upx contains a lzma-compressed version using UPX. Please be aware that AV scanners might flag these as false positives. When in doubt, just use the non-UPX version or compress the file yourself:

```cli
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/linux-amd64/check_adclicense main.go
upx --lzma -o bin/linux-amd64/upx/check_adclicense bin/linux-amd64/check_adclicense
```

## Examples

Use HTTP instead of HTTPS to avoid problems with untrusted (i.e. self-signed) certificates.
Specify a warning threshold for anything from 11 to 60 days. The threshold for a critical status is any value of 10 or less (inclusive of endpoints).

```cli
check_adclicense --hostname=adc_1 --username=nagios --password=nagios123 --secure=0 --warning=@11:60 --critical=@~:10
```

Use the same example as before but do not actually access the ADC (parameters hostname, username, password and secure may just as well be omitted). The test value will be used instead.

```cli
check_adclicense --hostname=adc_1 --username=nagios --password=nagios123 --secure=0 --warning=@11:60 --critical=@~:10 --testvalue=9
```
