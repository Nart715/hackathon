package client

import (
	"component-master/config"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	doOne sync.Once
	conn  *grpc.ClientConn

	grpcReadTimeout  time.Duration
	grpcWriteTimeout time.Duration
)

const (
	// second unit
	WriteTimeOutDefault = 120
	ReadTimeOutDefault  = 120
	ContextTimeout      = 120 * time.Second
)

type ClientConfig struct {
	Host            string
	Port            int
	EnableTLS       bool
	Timeout         time.Duration
	KeepAliveParams keepalive.ClientParameters
}

func ContextwithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	go cancelContext(ContextTimeout, cancel)
	return ctx
}

func cancelContext(timeout time.Duration, cancel context.CancelFunc) {
	time.Sleep(timeout)
	cancel()
	slog.Info("context canceled")
}

func InitConnection(conf config.GrpcConfigClient) (*grpc.ClientConn, error) {
	cfg := mapClientConfig(conf)
	opts := []grpc.DialOption{
		grpc.WithTimeout(120 * time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientInterceptor()),
	}

	grpcAddress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	slog.Info("grpc client address", "address", grpcAddress)
	// conn, err := grpc.NewClient(grpcAddress, opts...)
	conn, err := grpc.NewClient(grpcAddress, opts...)
	if err != nil {
		slog.Error("grpc client connection error", "err", err)
		return nil, err
	}

	return conn, err
}

func mapClientConfig(conf config.GrpcConfigClient) ClientConfig {
	if conf.ReadTimeOut < 1000 || conf.WriteTimeOut < 1000 {
		conf.ReadTimeOut = ReadTimeOutDefault
		conf.WriteTimeOut = WriteTimeOutDefault
	}
	return ClientConfig{
		Host:      conf.Host,
		Port:      conf.Port,
		EnableTLS: false,
		Timeout:   time.Duration(conf.ReadTimeOut) * time.Second,
		KeepAliveParams: keepalive.ClientParameters{
			Time:                time.Duration(conf.WriteTimeOut) * time.Second,
			Timeout:             time.Duration(conf.WriteTimeOut) * time.Second,
			PermitWithoutStream: true,
		}}
}

func clientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		slog.Info("Sending request",
			"method", method,
			"request", req)

		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			slog.Error("Request failed",
				"method", method,
				"error", err)
		} else {
			slog.Info("Request completed",
				"method", method,
				"response", reply)
		}

		return err
	}
}
