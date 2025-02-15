package main

import (
	"component-master/config"
	"component-master/infra/redis"
	"component-master/proto/account"
	"component-master/shared/idgen"
	"component-master/util"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	totalRequests = 1_000_000 // 1M RPS target
	batchSize     = 100_000   // chia 10 batch
	// numWorkers     = 700       // Adjust based on system capability
	requestTimeout = 2 * time.Second
	apiURL         = "http://localhost:8081/api/v1/player/balance-change"
)

var httpClient = &fasthttp.Client{
	MaxConnsPerHost: 10000, // Tăng số lượng kết nối tối đa
	ReadTimeout:     2 * time.Second,
	WriteTimeout:    2 * time.Second,
}

var (
	redisClient *redis.RedisClient

	createAccountUrl = "http://localhost:8081/api/v1/player/created-account"
	depositUrl       = "http://localhost:8081/api/v1/player/balance-change"
	bettingUrl       = "http://localhost:8081/api/v1/player/balance-change"
	requestCreateUrl = &account.CreateAccountRequest{}
	depositRequest   = &account.BalanceChangeRequest{}

	load_test = "a"
)

var requestPool = sync.Pool{
	New: func() interface{} {
		return &account.BalanceChangeRequest{}
	},
}

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
	flag.StringVar(&load_test, "load_test", "a", "Load test")
	flag.Parse()
}

func main() {
	if load_test == "a" {
		execute(create1000kAccounts)
		execute(depositAccount)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())

		for i := 0; i < totalRequests; i += batchSize {
			fmt.Printf("Running batch %d to %d...\n", i, i+batchSize)
			numWorkers := runtime.NumCPU() * 20
			loadTest(batchSize, numWorkers, sendRequest)
			time.Sleep(2 * time.Second)
		}
		// loadTest(500000, depositAccountRandom)
	}
}

func depositAccount(i int) {
	httpClient := NewHttpClient()
	header := NewHttpHeader()
	method := "POST"

	transaction_id := idgen.GenID()
	localRequestCreateUrl := &account.BalanceChangeRequest{
		Ac:  int32(i),
		Tx:  int64(transaction_id),
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

func depositAccountRandom(i int) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Sử dụng một nguồn ngẫu nhiên cục bộ
	httpClient := NewHttpClient()
	header := NewHttpHeader()
	method := "POST"

	transaction_id := idgen.GenID()
	localRequestCreateUrl := &account.BalanceChangeRequest{
		Ac:  int32(rng.Intn(100_000) + 1), // dựa theo số lượng account id
		Tx:  int64(transaction_id),
		Am:  1,
		Act: int32(rng.Intn(2)),
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

func sendRequest() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	transactionID := idgen.GenID()

	// Lấy object từ pool
	localRequest := requestPool.Get().(*account.BalanceChangeRequest)
	localRequest.Ac = int32(rng.Intn(100_000) + 1)
	localRequest.Tx = int64(transactionID)
	localRequest.Am = 1
	localRequest.Act = []int32{1, -1}[rng.Intn(2)]

	dataJson, _ := json.Marshal(localRequest)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(apiURL)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(dataJson)

	err := httpClient.DoTimeout(req, resp, requestTimeout)
	if err != nil {
		fmt.Println("Request failed:", err)
	}

	// Đưa object trở lại pool
	requestPool.Put(localRequest)
}
