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
	Database   DataSourceConfig `mapstructure:"database" json:"database,omitempty"`
	Redis      RedisConfig      `mapstructure:"redis" json:"redis,omitempty"`
	Nats       NatServerConfig  `mapstructure:"nats" json:"nats,omitempty"`
	GrpcClient GrpcConfigClient `mapstructure:"grpcClient" json:"grpc_client,omitempty"`
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

type DataSourceConfig struct {
	DriverName         string        `mapstructure:"driverName" json:"driver_name,omitempty"`
	Host               string        `mapstructure:"host" json:"host,omitempty"`
	Port               int           `mapstructure:"port" json:"port,omitempty"`
	UserName           string        `mapstructure:"userName" json:"user_name,omitempty"`
	Password           string        `mapstructure:"password" json:"password,omitempty"`
	DBName             string        `mapstructure:"dbName" json:"db_name,omitempty"`
	SSLMode            string        `mapstructure:"sslMode" json:"ssl_mode,omitempty"`
	MaxOpenConnections int           `mapstructure:"maxOpenConnections" json:"max_open_connections,omitempty"`
	MaxIdleConnections int           `mapstructure:"maxIdleConnections" json:"max_idle_connections,omitempty"`
	MaxConnLifetime    time.Duration `mapstructure:"maxConnLifetime" json:"max_conn_lifetime,omitempty"`
	MaxConnIdleTime    time.Duration `mapstructure:"maxConnIdleTime" json:"max_conn_idle_time,omitempty"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host" json:"host,omitempty"`
	Port         int           `mapstructure:"port" json:"port,omitempty"`
	Password     string        `mapstructure:"password" json:"password,omitempty"`
	DB           int           `mapstructure:"db" json:"db,omitempty"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout" json:"dial_timeout,omitempty"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"read_timeout,omitempty"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"write_timeout,omitempty"`
	MaxIdle      int           `mapstructure:"maxIdle" json:"max_idle,omitempty"`
}

type NatServerConfig struct {
	Host   string         `mapstructure:"host" json:"host,omitempty"`
	Port   int            `mapstructure:"port" json:"port,omitempty"`
	Topics NatTopicConfig `mapstructure:"topics" json:"topics,omitempty"`
}

type NatTopicConfig struct {
	BettingTopic string `mapstructure:"bettingTopic" json:"bettingTopic,omitempty"`
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

func (r *DataSourceConfig) BuildDnsMYSQL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", r.UserName, r.Password, r.Host, r.Port, r.DBName)
}

func (r *RedisConfig) BuildRedisConnectionString() string {
	return fmt.Sprintf("redis://%s:%s@%s:%d", "", "", r.Host, r.Port)
}
