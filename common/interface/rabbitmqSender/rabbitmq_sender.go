package rabbitmqSender

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"log"
	"sync"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/pkg/rabbitmqinfra"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/config"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitMqSender interface {
	SendExchangeAny(parentContext context.Context, any any, exchangeName string)
	SendExchangeAnyWithErrorHandle(parentContext context.Context, any any, exchangeName string, errorHandle func())
	SendExchangeWithErrorHandle(parentContext context.Context, msg string, exchangeName string, errorHandle func())
	SendQueueAny(parentContext context.Context, any any, queueName string)
	SendQueue(parentContext context.Context, msg string, queueName string)
}

type BaseRabbitMqSender struct {
	rabbitMqConfig *config.RabbitMqConfig
	mux            sync.RWMutex
	connection     *amqp.Connection
}

func NewBaseRabbitMqSender(rabbitMqConfig *config.RabbitMqConfig) *BaseRabbitMqSender {
	return &BaseRabbitMqSender{rabbitMqConfig: rabbitMqConfig}
}

func (b *BaseRabbitMqSender) Connect() {
	conn, err := rabbitmqinfra.CreateConnection(b.rabbitMqConfig)
	if err != nil {
		panic(err)
	}
	b.connection = conn
	go b.reconnect()
}

// reconnect reconnects to server if the connection or a channel
// is closed unexpectedly. Normal shutdown is ignored. It tries
// maximum of 7200 times and sleeps half a second in between
// each try which equals to 1 hour.
func (b *BaseRabbitMqSender) reconnect() {
WATCH:

	conErr := <-b.connection.NotifyClose(make(chan *amqp.Error))
	if conErr != nil {
		utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Sender CRITICAL: Connection dropped, reconnecting\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
		log.Println("RabbitMq Sender CRITICAL: Connection dropped, reconnecting")

		var err error
		for i := 1; i <= 172800; i++ {
			b.mux.RLock()
			b.connection, err = rabbitmqinfra.CreateConnection(b.rabbitMqConfig)
			b.mux.RUnlock()

			if err == nil {
				utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Sender INFO: Reconnected\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
				log.Println("RabbitMq Sender INFO: Reconnected")
				goto WATCH
			}

			time.Sleep(500 * time.Millisecond)
		}

		utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Sender CRITICAL: Failed to reconnect\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
		log.Println(errors.New("RabbitMq Sender CRITICAL: Failed to reconnect"))
	} else {
		utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Sender INFO: Connection dropped normally, will not reconnect\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
		log.Println("RabbitMq Sender INFO: Connection dropped normally, will not reconnect")
	}
}

func (b *BaseRabbitMqSender) channel() *amqp.Channel {
	if b.connection == nil {
		connection, err := rabbitmqinfra.CreateConnection(b.rabbitMqConfig)
		if err != nil {
			panic(err)
		}
		b.connection = connection
	}

	channel, err := b.connection.Channel()
	if err != nil {
		panic(err)
	}
	return channel
}

func (b *BaseRabbitMqSender) SendExchangeAny(parentContext context.Context, any any, exchangeName string) {
	msg, _ := json.Marshal(any)
	b.SendExchangeWithErrorHandle(parentContext, string(msg), exchangeName, nil)
}

func (b *BaseRabbitMqSender) SendExchangeAnyWithErrorHandle(parentContext context.Context, any any, exchangeName string, errorHandle func()) {
	msg, _ := json.Marshal(any)
	b.SendExchangeWithErrorHandle(parentContext, string(msg), exchangeName, errorHandle)
}

func (b *BaseRabbitMqSender) SendExchangeWithErrorHandle(parentContext context.Context, msg string, exchangeName string, errorHandle func()) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			if errorHandle != nil {
				errorHandle()
			}
			return
		}
	}()

	b.logBegin(exchangeName, "", actionLogs, now, utils.GetActionId(c))

	actionLogs[logKey.Msg] = msg

	channel := b.channel()
	defer channel.Close()
	err := channel.PublishWithContext(c,
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:     amqp.Table{logKey.Id: actionLogs[logKey.Id].(string)},
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		panic(err)
	}

	utils.LogEnd(actionLogs, now)
}

func (b *BaseRabbitMqSender) SendQueueAny(parentContext context.Context, any any, queueName string) {
	msg, _ := json.Marshal(any)
	b.SendQueue(parentContext, string(msg), queueName)
}

func (b *BaseRabbitMqSender) SendQueue(parentContext context.Context, msg string, queueName string) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.logBegin("", queueName, actionLogs, now, utils.GetActionId(c))

	actionLogs[logKey.Msg] = msg

	channel := b.channel()
	defer channel.Close()
	err := channel.PublishWithContext(c,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			Headers:     amqp.Table{logKey.Id: actionLogs[logKey.Id].(string)},
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		panic(err)
	}

	utils.LogEnd(actionLogs, now)
}

func (b *BaseRabbitMqSender) logBegin(exchangeName string, queueName string, actionLogs map[string]any, now time.Time, actionId string) {
	actionLogs[logKey.Type] = "rabbitmq-sender"
	actionLogs[logKey.ServerTime] = now.String()
	if actionId == "" {
		actionId = uuid.NewString()
	}
	actionLogs[logKey.Id] = actionId
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	if exchangeName != "" {
		actionLogs[logKey.RabbitmqExchangeName] = exchangeName
	}
	if queueName != "" {
		actionLogs[logKey.RabbitmqQueueName] = queueName
	}
}
