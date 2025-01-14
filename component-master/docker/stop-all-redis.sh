#!/bin/bash

echo "REDIS RUNNING"
lsof -i -P -n | grep redis | awk '{print $2}'

kill -9 $(lsof -i -P -n | grep redis | awk '{print $2}')

lsof -i -P -n | grep redis | awk '{print $2}'
