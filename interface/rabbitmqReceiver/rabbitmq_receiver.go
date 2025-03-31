package rabbitmqReceiver

import (
	"example-go/cmd"
	"example-go/common/infrastructure/config"
	"example-go/common/interface/rabbitmqReceiver"
)

type RabbitMqReceiver struct {
	app *cmd.App
}

func InitRabbitMqReceiver(app *cmd.App, rabbitMqConfig *config.RabbitMqConfig) {
	//receiver := RabbitMqReceiver{app}

	baseRabbitMqReceiver := rabbitmqReceiver.NewBaseRabbitMqReceiver(rabbitMqConfig)

	baseRabbitMqReceiver.Connect()
}
