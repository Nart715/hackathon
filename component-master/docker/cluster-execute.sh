#!/bin/bash


docker exec -it init-database-redis-node-1-1 redis-cli -p 6379 cluster info

docker exec -it $(docker ps -q -f name=redis-node-1) redis-cli -p 6379 cluster nodes

# src/redis-cli --cluster create 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 127.0.0.1:7006 --cluster-replicas 1
