package dev

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"linebot-go/common/infrastructure/config"
	"linebot-go/common/infrastructure/pkg/pyroscope"
	"os"
	"strconv"
	"time"
)

func CreateServerConfig() config.ServerConfig {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		panic("MYSQL_PORT required")
	}

	mysqlReadOnlyUser := os.Getenv("MYSQL_READ_ONLY_USER")
	mysqlReadOnlyPassword := os.Getenv("MYSQL_READ_ONLY_PASSWORD")
	mysqlReadOnlyHost := os.Getenv("MYSQL_READ_ONLY_HOST")
	mysqlReadOnlyPort, err := strconv.Atoi(os.Getenv("MYSQL_READ_ONLY_PORT"))
	if err != nil {
		panic("MYSQL_READ_ONLY_PORT required")
	}

	decimalCricketRedisAddress := os.Getenv("DECIMAL_CRICKET_REDIS_ADDRESS")
	decimalCricketRedisReadOnlyAddress := os.Getenv("DECIMAL_CRICKET_READ_ONLY_REDIS_ADDRESS")

	decimalCricketKafka := os.Getenv("DECIMAL_CRICKET_KAFKA_SERVERS")

	devDecimalCricketKafka := os.Getenv("DEV_DECIMAL_CRICKET_KAFKA_SERVERS")
	cqaDecimalCricketKafka := os.Getenv("CQA_DECIMAL_CRICKET_KAFKA_SERVERS")
	uatDecimalCricketKafka := os.Getenv("UAT_DECIMAL_CRICKET_KAFKA_SERVERS")
	prodDecimalCricketKafka := os.Getenv("PROD_DECIMAL_CRICKET_KAFKA_SERVERS")

	rabbitMqUser := os.Getenv("RABBIT_MQ_USER")
	rabbitMqPassword := os.Getenv("RABBIT_MQ_PASSWORD")
	rabbitMqHost := os.Getenv("RABBIT_MQ_HOST")
	rabbitMqPort, err := strconv.Atoi(os.Getenv("RABBIT_MQ_PORT"))
	if err != nil {
		panic(err)
	}

	pyroscopeServerAddress := os.Getenv("PYROSCOPE_URL")
	pyroscopeAuthToken := os.Getenv("PYROSCOPE_TOKEN")

	decimalK8sDomain := os.Getenv("DECIMAL_K8S_DOMAIN")

	return config.ServerConfig{
		Log:        &config.LogConfig{Type: "Json"},
		HttpServer: &config.HttpServerConfig{Address: ":8080", Mode: gin.ReleaseMode},
		Mysql: &config.MysqlConfig{
			Username:        mysqlUser,
			Password:        mysqlPassword,
			DbHost:          mysqlHost,
			DbPort:          mysqlPort,
			DbName:          "deci_cricket",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 60 * time.Second,
			LogMode:         logger.Info,
		},
		MysqlReadOnly: &config.MysqlConfig{
			Username:        mysqlReadOnlyUser,
			Password:        mysqlReadOnlyPassword,
			DbHost:          mysqlReadOnlyHost,
			DbPort:          mysqlReadOnlyPort,
			DbName:          "deci_cricket",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 60 * time.Second,
			LogMode:         logger.Info,
		},
		DecimalCricketRedis:         &config.RedisConfig{Address: decimalCricketRedisAddress, MinIdleConns: 10, PoolSize: 100},
		DecimalCricketRedisReadOnly: &config.RedisConfig{Address: decimalCricketRedisReadOnlyAddress, MinIdleConns: 10, PoolSize: 100},
		DecimalCricketKafka:         &config.KafkaConfig{Servers: decimalCricketKafka, GroupId: "decimal-cricket-websocket"},
		DecimalCricketKafkaEnv: &config.KafkaEnvConfig{GroupId: "decimal-cricket-websocket",
			Dev:  devDecimalCricketKafka,
			Cqa:  cqaDecimalCricketKafka,
			Uat:  uatDecimalCricketKafka,
			Prod: prodDecimalCricketKafka,
		},
		RabbitMq:   &config.RabbitMqConfig{Account: rabbitMqUser, Password: rabbitMqPassword, Host: rabbitMqHost, Port: rabbitMqPort},
		SysApiCidr: &config.CidrConfig{IpNet: "0.0.0.0/0"},
		Pyroscope: &config.PyroscopeConfig{
			Execute:       true,
			ServerAddress: pyroscopeServerAddress,
			AuthToken:     pyroscopeAuthToken,
			LogLevel:      pyroscope.LevelInfo,
			OpenGMB:       true,
		},
		DecimalK8sDomain: decimalK8sDomain,
	}
}
