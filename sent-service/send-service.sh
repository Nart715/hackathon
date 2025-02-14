#!/bin/bash

echo "CLEAN SEND SERVICE OLD BIN"
rm -rf ./sentservicebin

echo "BUILD DEPENDENCIES"
go mod tidy && go mod vendor

echo "BUILD SEND SERVICE"
go build -o sentservicebin main.go

echo "START SEND SERVICE"
./sentservicebin
