package http

import (
	"component-master/config"
	"component-master/middleware"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HttpServer struct {
	AppName string
	Conf    *config.ServerInfo
	app     *fiber.App
}

func (r *HttpServer) Start() {
	if r == nil || r.app == nil {
		slog.Error("http server is nil")
		return
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverShutdown := make(chan struct{})

	go func() {
		err := r.app.Listen(fmt.Sprintf("%s:%d", r.Conf.Host, r.Conf.Port))
		if err != nil {
			slog.Error("listen error", "err", err)
		}
		close(serverShutdown)
	}()

	<-shutdown
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.app.ShutdownWithContext(ctx); err != nil {
		slog.Info(fmt.Sprintf("Server forced to shutdown: %v\n", err))
	}

	<-serverShutdown
	slog.Info("Server gracefully stopped")

}

func (r *HttpServer) InitHttpServer() {
	app := fiber.New(r.ConfigFiber(r.Conf))
	app.Use(middleware.CorsFilter())
	r.app = app
}

func (r *HttpServer) App() *fiber.App {
	return r.app
}

func (r *HttpServer) ConfigFiber(conf *config.ServerInfo) fiber.Config {
	return fiber.Config{
		AppName:           "Fiber App",
		EnablePrintRoutes: true,
		ReadTimeout:       time.Duration(conf.ConnectTimeOut) * time.Millisecond,
		WriteTimeout:      time.Duration(conf.ConnectTimeOut) * time.Millisecond,
	}
}
