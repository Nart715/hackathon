# back-end

Keep it simple

#### Tech stack
- Golang(1.23.0)
- Docker
- Viper
- Grpc
- Redis
- Fiber V2
- Slog (Build-in)
- Nats
- MySQL

### Run test

### Install grpc

1. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
2.

### Grpc generate

1. brew install protobuf | apt install -y protobuf-compiler
2. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
3. go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
make grpc

```

### Build redis cluster

go to component-master/docker

```
sh install-redis-cluster.sh start
```

port 7001 7002 7003 7004 7005 7006


### Docker commands
1. Docker stop all containers: `docker stop $(docker ps -a -q)`
2. Docker remove all containers: `docker rm $(docker ps -a -q)`


### fix cluster: CLUSTERDOWN Hash slot not served

redis-cli --cluster fix localhost:7001
