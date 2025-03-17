package config

import (
	"net"
	"time"

	"gorm.io/gorm/logger"
)

type LogConfig struct {
	Type string
}

type HttpServerConfig struct {
	Address    string
	ServerName string
	Mode       string
}

type MysqlConfig struct {
	Regin           string
	Username        string
	Password        string `json:"-"`
	DbHost          string
	DbPort          int
	DbName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogMode         logger.LogLevel
}

type RedisConfig struct {
	Address      string
	PoolSize     int
	MinIdleConns int
}

type KafkaConfig struct {
	Servers string
	GroupId string
}

type KafkaEnvConfig struct {
	Dev     string
	Cqa     string
	Uat     string
	Prod    string
	GroupId string
}

type RabbitMqConfig struct {
	Account  string
	Password string `json:"-"`
	Host     string
	Port     int
}

type CidrConfig struct {
	IpNet string
}

type GrpcServerConfig struct {
	Address string
}

type ElasticsearchConfig struct {
	Addresses []string
}

type PyroscopeConfig struct {
	Execute       bool
	PodName       string
	ServerAddress string
	AuthToken     string
	LogLevel      int
	OpenGMB       bool
}

type ServerConfig struct {
	AppEnv                      string
	HostName                    string
	Log                         *LogConfig
	HttpServer                  *HttpServerConfig
	GrpcServer                  *GrpcServerConfig
	Mysql                       *MysqlConfig
	MysqlReadOnly               *MysqlConfig
	DecimalCricketRedis         *RedisConfig
	DecimalCricketRedisReadOnly *RedisConfig
	LockRedis                   *RedisConfig
	DecimalCricketKafka         *KafkaConfig
	DecimalCricketKafkaEnv      *KafkaEnvConfig
	RabbitMq                    *RabbitMqConfig
	ElasticsearchServer         *ElasticsearchConfig
	LogElasticsearchServer      *ElasticsearchConfig
	SysApiCidr                  *CidrConfig
	DecimalK8sDomain            string
	DecimalApiDomain            string
	InterfaceNameAndIpMap       map[string][]net.IP
	Pyroscope                   *PyroscopeConfig
}

func BuildDefaultConfig() *ServerConfig {
	return new(ServerConfig)
}
