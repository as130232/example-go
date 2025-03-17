package local

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"linebot-go/common/infrastructure/config"
	"time"
)

func CreateServerConfig() config.ServerConfig {
	return config.ServerConfig{
		Log:        &config.LogConfig{Type: "Console"},
		HttpServer: &config.HttpServerConfig{Address: ":8083", Mode: gin.ReleaseMode},
		Mysql: &config.MysqlConfig{
			Username:        "dev_deci_rd_use",
			Password:        "22ZwdnyWNrcB",
			DbHost:          "mysql-deci.ljbdev.site",
			DbPort:          3306,
			DbName:          "deci_cricket",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 60 * time.Second,
			LogMode:         logger.Info,
		},
		MysqlReadOnly: &config.MysqlConfig{
			Username:        "dev_deci_rd_use",
			Password:        "22ZwdnyWNrcB",
			DbHost:          "mysql-deci.ljbdev.site",
			DbPort:          3306,
			DbName:          "deci_cricket",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 60 * time.Second,
			LogMode:         logger.Info,
		},
		DecimalCricketRedis:         &config.RedisConfig{Address: "redis-deci.ljbdev.site:6379", MinIdleConns: 10, PoolSize: 100},
		DecimalCricketRedisReadOnly: &config.RedisConfig{Address: "redis-read-deci.ljbdev.site:6379", MinIdleConns: 10, PoolSize: 100},
		DecimalCricketKafka: &config.KafkaConfig{GroupId: "decimal-cricket-websocket-local",
			Servers: "b-2.devdeci.bn8v79.c1.kafka.us-west-2.amazonaws.com:9092,b-1.devdeci.bn8v79.c1.kafka.us-west-2.amazonaws.com:9092,b-3.devdeci.bn8v79.c1.kafka.us-west-2.amazonaws.com:9092",
		},
		DecimalCricketKafkaEnv: &config.KafkaEnvConfig{GroupId: "decimal-cricket-websocket-local",
			Dev:  "b-2.devdeci.bn8v79.c1.kafka.us-west-2.amazonaws.com:9092,b-1.devdeci.bn8v79.c1.kafka.us-west-2.amazonaws.com:9092,b-3.devdeci.bn8v79.c1.kafka.us-west-2.amazonaws.com:9092",
			Cqa:  "",
			Uat:  "",
			Prod: "",
		}, RabbitMq: &config.RabbitMqConfig{Account: "decimal", Password: "HpPd8drP", Host: "deci-rabbitmq.ljbdev.site", Port: 5672},
		SysApiCidr:       &config.CidrConfig{IpNet: "0.0.0.0/0"},
		DecimalApiDomain: "https://deci-api.ljbdev.site",
	}
}
