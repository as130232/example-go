package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"linebot-go/cmd"
	"linebot-go/common/application/utils"
	"linebot-go/common/global"
	"linebot-go/common/infrastructure/consts/contextKey"
	"linebot-go/common/infrastructure/pkg/pyroscope"
	"linebot-go/common/infrastructure/pkg/redis"
	"linebot-go/common/infrastructure/pkg/sqldatabase"
	"linebot-go/common/interface/executor"
	"linebot-go/common/interface/telegram"
	appConfig "linebot-go/infrastructure/config"
	"linebot-go/interface/http/router"
	"linebot-go/interface/job"
	"linebot-go/interface/kafkaReceiver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var commitId string
var buildTime string

func main() {
	defer func() {
		r := recover()
		if r != nil {
			utils.LogServerPanic(r)
			utils.SendTelegramServerPanicMessage(commitId, buildTime)
			panic(r)
		}
	}()

	saveVersion()
	global.AppName = "linebot-go"
	global.ServerConfig = appConfig.NewServerConfig(global.AppName)
	global.AppConfig = appConfig.NewAppConfig()
	global.TelegramBot = telegram.NewBot(global.ServerConfig.AppEnv, nil)
	global.Redis = redis.NewRedis(global.ServerConfig)
	global.RedisReadOnly = redis.NewRedisReadOnly(global.ServerConfig)

	sqldatabase.SetupDB(global.ServerConfig)

	c := context.WithValue(context.Background(), contextKey.ActionLogs, make(map[string]any))
	utils.SetActionId(c, uuid.NewString())

	// Setup app
	app := cmd.InitApp()

	pyroscope.Init(global.AppName, global.ServerConfig)

	// Message Receiver
	kafkaReceiver.InitKafkaReceiver(app)
	//rabbitmqReceiver.InitRabbitMqReceiver(app, global.ServerConfig.RabbitMq)

	cron := job.Init(app)

	log.Printf("%s server startup, hostname:%s, git commit:%s, buildTime:%s, listen:%s", global.AppName, global.ServerConfig.HostName, commitId, buildTime, global.ServerConfig.HttpServer.Address)

	ginRouter := router.InitRouter(app)
	httpServer := &http.Server{
		Addr:    global.ServerConfig.HttpServer.Address,
		Handler: ginRouter,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	global.IsShutdown = true

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("%s server shutdown, hostname:%s, git commit:%s, buildTime:%s, listen:%s, error:%+v", global.AppName, global.ServerConfig.HostName, commitId, buildTime, global.ServerConfig.HttpServer.Address, err)
	}

	cron.Stop()

	executor.WaitForShutdown()

	global.SqlDB.Close()
	global.SqlDBReadOnly.Close()

	pyroscope.Stop()

	os.Exit(0)
}

func saveVersion() {
	var vFlag bool
	flag.BoolVar(&vFlag, "v", false, "show version")
	flag.Parse()
	global.CommitId = commitId
	global.BuildTime = buildTime
}
