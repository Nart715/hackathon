#!/bin/bash

echo "CLEAN RECEIVE SERVICE OLD BIN"
rm -rf ./receiveservicebin

echo "BUILD DEPENDENCIES"
go mod tidy && go mod vendor

echo "BUILD RECEIVE SERVICE"
go build -o receiveservicebin main.go

echo "START RECEIVE SERVICE"
./receiveservicebin
