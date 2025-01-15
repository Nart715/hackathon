#!/bin/bash

echo "BUILD DEPENDENCIES"
go mod tidy && go mod vendor

echo "BUILD RECEIVE SERVICE"
go build -o receive-service main.go

echo "START RECEIVE SERVICE"
./receive-service
