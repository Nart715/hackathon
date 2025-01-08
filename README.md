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

### Grpc generate

1. go to component-master
```
make grpc

```


### Docker commands
1. Docker stop all containers: `docker stop $(docker ps -a -q)`
2. Docker remove all containers: `docker rm $(docker ps -a -q)`
