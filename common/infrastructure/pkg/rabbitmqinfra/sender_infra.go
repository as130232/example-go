package rabbitmqinfra

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Declare(channel *amqp.Channel, exchangeNames []string, queueNames []string) {
	defer channel.Close()

	// declare exchange
	for _, exchangeName := range exchangeNames {
		err := channel.ExchangeDeclare(
			exchangeName, // name
			"fanout",     // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			panic(err)
		}
	}

	// declare queue
	for _, queueName := range queueNames {
		_, err := channel.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			panic(err)
		}
	}
}
