#!/bin/bash

echo "BUILD DEPENDENCIES"
go mod tidy && go mod vendor

echo "BUILD SEND SERVICE"
go build -o sent-service main.go

echo "START SEND SERVICE"
./sent-service
