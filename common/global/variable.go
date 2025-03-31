package global

import (
	"database/sql"
	"example-go/common/infrastructure/config"
	"example-go/common/interface/telegram"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	AppName       string
	ServerConfig  *config.ServerConfig
	AppConfig     *config.AppConfig
	DB            *gorm.DB
	DBReadOnly    *gorm.DB
	SqlDB         *sql.DB
	SqlDBReadOnly *sql.DB
	Redis         *redis.Client
	RedisReadOnly *redis.Client
	LockRedis     *redis.Client
	TelegramBot   *telegram.Bot
	CommitId      string
	BuildTime     string
	IsShutdown    bool
)
