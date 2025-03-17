package rabbitmqinfra

import (
	"fmt"
	"log"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateChannel(rabbitMqConfig *config.RabbitMqConfig) *amqp.Channel {
	conn, err := CreateConnection(rabbitMqConfig)
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return channel
}

func CreateConnection(rabbitMqConfig *config.RabbitMqConfig) (*amqp.Connection, error) {
	mqAddress := buildRabbitMqAddress(rabbitMqConfig)
	var connection *amqp.Connection
	var err error
	for retry := 1; retry <= 3; retry++ {
		connection, err = amqp.Dial(mqAddress)
		if err == nil {
			break
		}
		log.Printf("CreateConnection mqAddress:%s,error:%+v,retry:%d", mqAddress, err, retry)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func buildRabbitMqAddress(rabbitMqConfig *config.RabbitMqConfig) string {
	mqAddress := fmt.Sprintf("amqp://%v:%v@%v:%v/", rabbitMqConfig.Account, rabbitMqConfig.Password, rabbitMqConfig.Host, rabbitMqConfig.Port)
	return mqAddress
}
