package rabbitmqReceiver

import (
	"linebot-go/cmd"
	"linebot-go/common/infrastructure/config"
	"linebot-go/common/interface/rabbitmqReceiver"
)

type RabbitMqReceiver struct {
	app *cmd.App
}

func InitRabbitMqReceiver(app *cmd.App, rabbitMqConfig *config.RabbitMqConfig) {
	//receiver := RabbitMqReceiver{app}

	baseRabbitMqReceiver := rabbitmqReceiver.NewBaseRabbitMqReceiver(rabbitMqConfig)

	baseRabbitMqReceiver.Connect()
}
