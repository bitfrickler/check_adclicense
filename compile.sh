#!/bin/bash

rm bin/windows-amd64/upx/check_adclicense
rm bin/linux-amd64/upx/check_adclicense
rm bin/darwin-amd64/upx/check_adclicense

echo Compiling for linux-amd64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/linux-amd64/check_adclicense main.go

echo Compiling for windows-amd64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/windows-amd64/check_adclicense main.go

echo Compiling for darwin-amd64
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/darwin-amd64/check_adclicense main.go

echo Packing executables with upx
echo 
upx --lzma -o bin/windows-amd64/upx/check_adclicense bin/windows-amd64/check_adclicense
upx --lzma -o bin/darwin-amd64/upx/check_adclicense bin/darwin-amd64/check_adclicense
upx --lzma -o bin/linux-amd64/upx/check_adclicense bin/linux-amd64/check_adclicense

echo 
echo Done!!