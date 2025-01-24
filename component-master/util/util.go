package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const (
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ContextTimeout = 120 * time.Millisecond
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found. Using exists environment variables")
	}
}

func SliceToJson(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	resp, err := json.Marshal(slice)
	if err != nil {
		return ""
	}
	return string(resp)
}

func StructToJson(structAny any) string {
	resp, err := json.Marshal(structAny)
	if err != nil {
		return ""
	}
	return string(resp)
}

func UUID() string {
	id := uuid.New()
	return id.String()
}

func UUIDFunc() func() string {
	return func() string {
		return UUID()
	}
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

func String(pattern string, args ...any) string {
	return fmt.Sprintf(pattern, args...)
}
