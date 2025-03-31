package kafkaReceiver

import (
	"example-go/cmd"
	"example-go/common/interface/kafkaReceiver"
	"github.com/segmentio/kafka-go"
)

type KafkaReceiver struct {
	baseKafkaReceiver *kafkaReceiver.BaseKafkaReceiver
	matchListReader   *kafka.Reader
	matchListHandler  kafkaReceiver.IKafkaHandler
}

func InitKafkaReceiver(app *cmd.App) {
	receiver := &KafkaReceiver{baseKafkaReceiver: &kafkaReceiver.BaseKafkaReceiver{}} //matchListReader:  pkgKafka.NewDecimalKafkaReader(global.ServerConfig.DecimalCricketKafkaEnv.GroupId, kafkaTopic.MatchListTopic),
	//matchListHandler: executor.NewMatchListExecutor(utils.NewPool(30), app.MatchPoolService),

	receiver.receiveMessage()
}

func (k *KafkaReceiver) receiveMessage() {
	go k.baseKafkaReceiver.Receive(k.matchListReader, k.matchListHandler.HandleMessage)
}
