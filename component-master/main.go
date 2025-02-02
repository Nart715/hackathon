package main

import (
	"component-master/config"
	"component-master/infra/redis"
	"component-master/proto/account"
	"component-master/shared/idgen"
	"component-master/util"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var (
	redisClient *redis.RedisClient

	createAccountUrl = "http://localhost:8081/api/v1/player/created-account"
	depositUrl       = "http://localhost:8081/api/v1/player/balance-change"
	bettingUrl       = "http://localhost:8081/api/v1/player/balance-change"
	requestCreateUrl = &account.CreateAccountRequest{}
	depositRequest   = &account.BalanceChangeRequest{}
)

func NewHttpClient() http.Client {
	return http.Client{
		Timeout: util.ContextTimeout,
	}
}

func NewHttpHeader() http.Header {
	return http.Header{
		"Content-Type": []string{"application/json"},
	}
}

func init() {
	util.LoadEnv()
}

func main() {
	execute(create1000kAccounts)
	execute(depositAccount)
}

func depositAccount(i int) {
	httpClient := NewHttpClient()
	header := NewHttpHeader()
	method := "POST"

	transaction_id := idgen.GenID()
	localRequestCreateUrl := &account.BalanceChangeRequest{
		Ac:  int32(i),
		Tx:  int32(transaction_id),
		Am:  100,
		Act: 0,
	}

	dataJson, _ := json.Marshal(localRequestCreateUrl)

	requestString := string(dataJson)
	fmt.Println("INTERNAL >> ", requestString)
	data := strings.NewReader(requestString)
	req, _ := http.NewRequest(method, bettingUrl, data)
	req.Header = header
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func create1000kAccounts(i int) {
	httpClient := NewHttpClient()
	header := NewHttpHeader()
	method := "POST"

	// deep copy to fix goroutine the same time
	localRequestCreateUrl := &account.CreateAccountRequest{
		AccountId: int32(i),
	}

	dataJson, _ := json.Marshal(localRequestCreateUrl)
	requestString := string(dataJson)
	fmt.Println("INTERNAL >> ", requestString)
	data := strings.NewReader(requestString)
	req, _ := http.NewRequest(method, createAccountUrl, data)
	req.Header = header
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Printf("Failed to create account %d: %v\n", i, err)
		return
	}
	fmt.Println(string(body))
}

func initInfra() {
	conf := &config.Config{}
	err := config.LoadConfig("config", "", conf)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config, %v", err))
		return
	}

	rd, err := redis.NewInitRedisClient(&conf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
		return
	}
	redisClient = rd
	slog.Info(fmt.Sprintf("REDIS CONNECT: %v", rd != nil))

}
