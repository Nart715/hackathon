package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

func StringToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func ContextwithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	go cancelContext(ContextTimeout, cancel)
	return ctx
}

func cancelContext(timeout time.Duration, cancel context.CancelFunc) {
	time.Sleep(timeout)
	cancel()
}

func String(i int64) string {
	return strconv.Itoa(int(i))
}

func StringPattern(pattern string, arg ...any) string {
	return fmt.Sprintf(pattern, arg...)
}
