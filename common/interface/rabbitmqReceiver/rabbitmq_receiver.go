package rabbitmqReceiver

import (
	"context"
	"errors"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"log"
	"sync"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/config"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/pkg/rabbitmqinfra"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BaseRabbitMqReceiver struct {
	rabbitMqConfig      *config.RabbitMqConfig
	mux                 sync.RWMutex
	connection          *amqp.Connection
	channel             *amqp.Channel
	exchangeCallbackMap map[string]func(c context.Context, deliver amqp.Delivery)
	queueCallbackMap    map[string]func(c context.Context, deliver amqp.Delivery)
}

func NewBaseRabbitMqReceiver(rabbitMqConfig *config.RabbitMqConfig) *BaseRabbitMqReceiver {
	return &BaseRabbitMqReceiver{
		rabbitMqConfig:      rabbitMqConfig,
		exchangeCallbackMap: make(map[string]func(c context.Context, deliver amqp.Delivery)),
		queueCallbackMap:    make(map[string]func(c context.Context, deliver amqp.Delivery)),
	}
}

func (b *BaseRabbitMqReceiver) AddExchangeCallback(topicName string, callback func(c context.Context, deliver amqp.Delivery)) {
	b.exchangeCallbackMap[topicName] = callback
}

func (b *BaseRabbitMqReceiver) AddQueueCallback(topicName string, callback func(c context.Context, deliver amqp.Delivery)) {
	b.queueCallbackMap[topicName] = callback
}

func (b *BaseRabbitMqReceiver) bindCallback() {
	for topic, fun := range b.exchangeCallbackMap {
		go b.ExchangeReceive(topic, rabbitmqinfra.CreateExchangeDeliveryChan(b.channel, topic), fun)
	}
	for topic, fun := range b.queueCallbackMap {
		go b.QueueReceive(topic, rabbitmqinfra.CreateQueueDeliveryChan(b.channel, topic), fun)
	}
}

func (b *BaseRabbitMqReceiver) Connect() *amqp.Channel {
	conn, err := rabbitmqinfra.CreateConnection(b.rabbitMqConfig)
	if err != nil {
		panic(err)
	}
	b.connection = conn

	channel, err := b.connection.Channel()
	if err != nil {
		panic(err)
	}
	b.channel = channel

	b.bindCallback()
	go b.reconnect()

	return channel
}

// reconnect reconnects to server if the connection or a channel
// is closed unexpectedly. Normal shutdown is ignored. It tries
// maximum of 7200 times and sleeps half a second in between
// each try which equals to 1 hour.
func (b *BaseRabbitMqReceiver) reconnect() {
WATCH:

	conErr := <-b.connection.NotifyClose(make(chan *amqp.Error))
	if conErr != nil {
		utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Receiver CRITICAL: Connection dropped, reconnecting\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
		log.Println("RabbitMq Receiver CRITICAL: Connection dropped, reconnecting")

		b.channel.Close()

		var err error
		for i := 1; i <= 172800; i++ {
			b.mux.RLock()
			b.connection, err = rabbitmqinfra.CreateConnection(b.rabbitMqConfig)
			b.mux.RUnlock()

			if err == nil {
				b.channel, err = b.connection.Channel()
				if err == nil {
					utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Receiver INFO: Reconnected\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
					log.Println("RabbitMq Receiver INFO: Reconnected")
					b.bindCallback()
					goto WATCH
				}
			}

			time.Sleep(500 * time.Millisecond)
		}

		utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Receiver CRITICAL: Failed to reconnect\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
		log.Println(errors.New("RabbitMq Receiver CRITICAL: Failed to reconnect"))
	} else {
		utils.SendTelegramMessage(fmt.Sprintf("=====\n%s %s\nhostname:%s RabbitMq Receiver INFO: Connection dropped normally, will not reconnect\n=====", global.ServerConfig.AppEnv, global.AppName, fmt.Sprintf(global.ServerConfig.HostName)))
		log.Println("RabbitMq Receiver INFO: Connection dropped normally, will not reconnect")
	}
}

func (b *BaseRabbitMqReceiver) exchangeReceive(exchangeName string, delivery amqp.Delivery, callback func(c context.Context, deliver amqp.Delivery)) {
	actionLogs := make(map[string]any)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			actionLogs[logKey.RabbitmqDeliveryBody] = string(delivery.Body)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.logBegin(delivery.Headers, exchangeName, "", actionLogs, now)
	c := context.WithValue(context.Background(), contextKey.ActionLogs, actionLogs)

	callback(c, delivery)

	utils.LogEnd(actionLogs, now)
}

func (b *BaseRabbitMqReceiver) ExchangeReceive(exchangeName string, deliveryChan <-chan amqp.Delivery, callback func(c context.Context, deliver amqp.Delivery)) {
	for delivery := range deliveryChan {
		b.exchangeReceive(exchangeName, delivery, callback)
	}
}

func (b *BaseRabbitMqReceiver) queueReceive(queueName string, delivery amqp.Delivery, callback func(c context.Context, deliveryChan amqp.Delivery)) {
	actionLogs := make(map[string]any)
	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.logBegin(delivery.Headers, "", queueName, actionLogs, now)
	c := context.WithValue(context.Background(), contextKey.ActionLogs, actionLogs)

	callback(c, delivery)

	utils.LogEnd(actionLogs, now)
}

func (b *BaseRabbitMqReceiver) QueueReceive(queueName string, deliveryChan <-chan amqp.Delivery, callback func(c context.Context, deliveryChan amqp.Delivery)) {
	for delivery := range deliveryChan {
		b.queueReceive(queueName, delivery, callback)
	}
}

func (b *BaseRabbitMqReceiver) logBegin(headers amqp.Table, exchangeName string, queueName string, actionLogs map[string]any, now time.Time) {
	actionLogs[logKey.Type] = "rabbitmq-receiver"
	actionLogs[logKey.ServerTime] = now.String()
	actionLogs[logKey.Id] = uuid.NewString()
	if len(headers) > 0 {
		actionLogs[logKey.RefId] = actionLogs[logKey.Id].(string)
		actionLogs[logKey.Id] = headers[logKey.Id].(string)
	}
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	if exchangeName != "" {
		actionLogs[logKey.RabbitmqExchangeName] = exchangeName
	}
	if queueName != "" {
		actionLogs[logKey.RabbitmqQueueName] = queueName
	}
}
