# back-end

Require: golang 1.23.0

### Install grpc

1. brew install protobuf | apt install -y protobuf-compiler
2. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
3. go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

1. go to component-master, run: go mod tidy && go mod vendor
2. go to receive-service, run: `sh receive.sh start`
3. go to sent-service, run: `sh sent-service.sh start`

#### Generator proto file:
- go to component-master
```
make grpc

```

### Build redis cluster

go to component-master/docker

```
sh install-redis-cluster.sh start
```

### Start kafka service before run receive service
