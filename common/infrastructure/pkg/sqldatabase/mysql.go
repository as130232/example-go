package sqldatabase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"linebot-go/common/global"
	"log"
	"os"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"linebot-go/common/infrastructure/config"
)

func SetupDB(config *config.ServerConfig) {
	NewDB(config)
	NewDBReadOnly(config)
}

func setupAWSAuthProvider(mysqlConfig *config.MysqlConfig) *sql.DB {
	authProvider := NewAWSAuthProvider(mysqlConfig.Regin, fmt.Sprintf("%s:%d", mysqlConfig.DbHost, mysqlConfig.DbPort), mysqlConfig.Username, mysqlConfig.DbName)

	authProvider.Register()
	dsn := authProvider.DataSourceName()

	cfg, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}
	cfg.CheckConnLiveness = true
	beforeConnect := mysqlDriver.BeforeConnect(func(ctx context.Context, cfg *mysqlDriver.Config) error {
		cfg.Passwd = authProvider.AuthToken()
		log.Println("[BeforeConnect] AuthToken:", cfg.Passwd)
		return nil
	})
	err = cfg.Apply(beforeConnect)
	if err != nil {
		panic(err)
	}

	connector, err := mysqlDriver.NewConnector(cfg)
	if err != nil {
		panic(err)
	}

	return sql.OpenDB(connector)
}

func openDB(mysqlConfig *config.MysqlConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var errorMessage string

	for retry := 1; retry <= 3; retry++ {
		db, err = gorm.Open(mysql.Open(buildDbLink(mysqlConfig)), &gorm.Config{
			Logger: NewCustomLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             5 * time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			}).LogMode(mysqlConfig.LogMode),
		})
		if err == nil {
			break
		}
		errorMessage = fmt.Sprintf("%+v", err)
		if strings.Contains(errorMessage, "bad connection") {
			mysqlConfigJson, _ := json.Marshal(mysqlConfig)
			log.Printf("openDB has bad connection mysqlConfig:%+v,retry:%d", mysqlConfigJson, retry)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	return db, err
}

func NewDB(config *config.ServerConfig) *gorm.DB {
	mysqlConfig := config.Mysql

	var db *gorm.DB
	var sqlDB *sql.DB
	var err error
	if mysqlConfig.Regin == "" {
		db, err = openDB(mysqlConfig)
	} else {
		sqlDB = setupAWSAuthProvider(mysqlConfig)
		db, err = gorm.Open(&mysql.Dialector{Config: &mysql.Config{Conn: sqlDB}}, &gorm.Config{
			Logger: NewCustomLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             5 * time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			}).LogMode(mysqlConfig.LogMode),
		})
		if err != nil {
			db, err = openDB(mysqlConfig)
		}
	}
	if err != nil {
		panic(err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns) //设置数据库连接池最大连接数
	sqlDB.SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime)

	global.DB = db
	global.SqlDB = sqlDB
	return db
}

func buildDbLink(c *config.MysqlConfig) string {
	if c == nil {
		return "root:2uh7x8T5@tcp(faerun-dev-instance-1.cc0kglfryxap.ap-northeast-1.rds.amazonaws.com:3306)/faerun_game?charset=utf8mb4&parseTime=True&loc=Asia%2fTaipei"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s", c.Username, c.Password, c.DbHost, c.DbPort, c.DbName, "Asia%2fTaipei")
}

func openDBReadOnly(mysqlConfig *config.MysqlConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var errorMessage string

	for retry := 1; retry <= 3; retry++ {
		db, err = gorm.Open(mysql.Open(buildDbLink(mysqlConfig)), &gorm.Config{
			Logger: NewCustomLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             5 * time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			}).LogMode(mysqlConfig.LogMode),
			//Logger: logger.Default.LogMode(mysqlConfig.LogMode),
		})
		if err == nil {
			break
		}
		errorMessage = fmt.Sprintf("%+v", err)
		if strings.Contains(errorMessage, "bad connection") {
			log.Printf("openDBReadOnly has bad connection mysqlConfig:%+v,retry:%d", mysqlConfig, retry)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	return db, err
}

func NewDBReadOnly(config *config.ServerConfig) *gorm.DB {
	mysqlConfig := config.MysqlReadOnly

	var db *gorm.DB
	var sqlDB *sql.DB
	var err error
	if mysqlConfig.Regin == "" {
		db, err = openDBReadOnly(mysqlConfig)
	} else {
		sqlDB = setupAWSAuthProvider(mysqlConfig)
		db, err = gorm.Open(&mysql.Dialector{Config: &mysql.Config{Conn: sqlDB}}, &gorm.Config{
			Logger: logger.Default.LogMode(mysqlConfig.LogMode),
		})
		if err != nil {
			db, err = openDBReadOnly(mysqlConfig)
		}
	}
	if err != nil {
		panic(err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns) //设置数据库连接池最大连接数
	sqlDB.SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime)

	global.DBReadOnly = db
	global.SqlDBReadOnly = sqlDB
	return db
}
