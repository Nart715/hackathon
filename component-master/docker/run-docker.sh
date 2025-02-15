#!/bin/bash

status=$1

echo $status

docker_compose=docker-compose.yaml

if [ "$status" = "up" ]; then
    docker-compose -f $docker_compose --project-name=init-database up -d --remove-orphans
fi

if [ "$status" = "down" ]; then
    docker-compose -f $docker_compose --project-name=init-database down
fi

