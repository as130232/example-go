package rabbitmqinfra

import amqp "github.com/rabbitmq/amqp091-go"

func CreateExchangeDeliveryChan(channel *amqp.Channel, exchangeName string) <-chan amqp.Delivery {
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

	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	deliveryChan, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}

	return deliveryChan
}

func CreateQueueDeliveryChan(channel *amqp.Channel, queueName string) <-chan amqp.Delivery {
	queue, err := channel.QueueDeclare(
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

	deliveryChan, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}

	return deliveryChan
}
