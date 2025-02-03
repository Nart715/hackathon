package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server" json:"server,omitempty"`
	Redis      RedisConfig      `mapstructure:"redis" json:"redis,omitempty"`
	GrpcClient GrpcConfigClient `mapstructure:"grpcClient" json:"grpc_client,omitempty"`
	Kafka      KafkaConfig      `mapstructure:"kafka" json:"kafka,omitempty"`
	KafkaTopic KafkaTopic       `mapstructure:"kafkaTopics" json:"kafka_topics,omitempty"`
}

type ServerConfig struct {
	Http ServerInfo `mapstructure:"http" json:"http,omitempty"`
	Grpc ServerInfo `mapstructure:"grpc" json:"grpc,omitempty"`
}

type ServerInfo struct {
	Host           string `mapstructure:"host" json:"host,omitempty"`
	Port           int    `mapstructure:"port" json:"port,omitempty"`
	ConnectTimeOut int    `mapstructure:"connectTimeOut" json:"connect_time_out,omitempty"`
}

type GrpcConfigClient struct {
	Host         string
	Port         int
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

type RedisConfig struct {
	Password      string        `mapstructure:"password" json:"password,omitempty"`
	DB            int           `mapstructure:"db" json:"db,omitempty"`
	DialTimeout   time.Duration `mapstructure:"dialTimeout" json:"dial_timeout,omitempty"`
	ReadTimeout   time.Duration `mapstructure:"readTimeout" json:"read_timeout,omitempty"`
	WriteTimeout  time.Duration `mapstructure:"writeTimeout" json:"write_timeout,omitempty"`
	MaxIdle       int           `mapstructure:"maxIdle" json:"max_idle,omitempty"`
	ReadOnly      bool          `mapstructure:"readOnly" json:"read_only,omitempty"`
	RouteRandomly bool          `mapstructure:"routeRandomly" json:"route_randomly,omitempty"`
	MaxRedirects  int           `mapstructure:"maxRedirects" json:"max_redirects,omitempty"`
	PoolSize      int           `mapstructure:"poolSize" json:"pool_size,omitempty"`
	MinIdleConns  int           `mapstructure:"minIdleConns" json:"min_idle_conns,omitempty"`
	Channel       string        `mapstructure:"channel" json:"channel,omitempty"`
	Clusters      []string      `mapstructure:"clusters" json:"clusters,omitempty"`
}

type KafkaConfig struct {
	ClientId         string         `mapstructure:"clientId" json:"client_id,omitempty"`
	DialTimeout      time.Duration  `mapstructure:"dialTimeout" json:"dial_timeout,omitempty"`
	ReadTimeout      time.Duration  `mapstructure:"readTimeout" json:"read_timeout,omitempty"`
	WriteTimeout     time.Duration  `mapstructure:"writeTimeout" json:"write_timeout,omitempty"`
	MaxRetry         int            `mapstructure:"maxRetry" json:"max_retry,omitempty"`
	RetryBackoff     time.Duration  `mapstructure:"retryBackoff" json:"retry_backoff,omitempty"`
	RefreshFrequency time.Duration  `mapstructure:"refreshFrequency" json:"refresh_frequency,omitempty"`
	Brokers          []string       `mapstructure:"brokers" json:"brokers,omitempty"`
	Consumer         ConsumerConfig `mapstructure:"consumer" json:"consumer,omitempty"`
	Producer         ProducerConfig `mapstructure:"producer" json:"producer,omitempty"`
	EnableTLS        bool           `mapstructure:"enableTLS" json:"enable_tls,omitempty"`
	EnableSASL       bool           `mapstructure:"enableSASL" json:"enable_sasl,omitempty"`
}

type ConsumerConfig struct {
	GroupId           string        `mapstructure:"groupId" json:"group_id,omitempty"`
	MaxProcessingTime time.Duration `mapstructure:"maxProcessingTime" json:"max_processing_time,omitempty"`
	FetchMin          int32         `mapstructure:"fetchMin" json:"fetch_min,omitempty"`
	FetchMax          int32         `mapstructure:"fetchMax" json:"fetch_max,omitempty"`
	RetryBackoff      time.Duration `mapstructure:"retryBackoff" json:"retry_backoff,omitempty"`
}

type ProducerConfig struct {
	MaxMessageBytes int           `mapstructure:"maxMessageBytes" json:"max_message_bytes,omitempty"`
	Compression     string        `mapstructure:"compression" json:"compression,omitempty"`
	FlushFrequency  time.Duration `mapstructure:"flushFrequency" json:"flush_frequency,omitempty"`
	FlushBytes      int           `mapstructure:"flushBytes" json:"flush_bytes,omitempty"`
	RetryMax        int           `mapstructure:"retryMax" json:"retry_max,omitempty"`
	RetryBackoff    time.Duration `mapstructure:"retryBackoff" json:"retry_backoff,omitempty"`
}

type KafkaTopic struct {
	BalanceChange string `mapstructure:"balanceChange" json:"balance_change,omitempty"`
}

func loadConfig(configfile string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(configfile)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	if err := v.ReadInConfig(); err != nil {
		log.Println("config file not found. Using exists environment variables")
		return nil, err
	}
	overrideconfig(v)
	return v, nil
}

func overrideconfig(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		envKey := "APP_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		envValue := os.Getenv(envKey)
		if envValue != "" {
			v.Set(key, envValue)
		}

	}
}

func LoadConfig(pathToFile string, env string, config any) error {
	configPath := pathToFile + "/" + "application"
	if len(env) > 0 {
		configPath = configPath + "-" + env
	}

	pwd, err := os.Getwd()
	if err == nil && len(pwd) > 0 {
		configPath = pwd + "/" + configPath
	}

	confFile := fmt.Sprintf("%s.yaml", configPath)
	slog.Info(fmt.Sprintf("Config file path: %s", confFile))
	v, err := loadConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	return nil
}
