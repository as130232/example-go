package kafkaReceiver

import (
	"github.com/segmentio/kafka-go"
	"linebot-go/cmd"
	"linebot-go/common/interface/kafkaReceiver"
)

type KafkaReceiver struct {
	baseKafkaReceiver *kafkaReceiver.BaseKafkaReceiver
	matchListReader   *kafka.Reader
	matchListHandler  kafkaReceiver.IKafkaHandler
}

func InitKafkaReceiver(app *cmd.App) {
	receiver := &KafkaReceiver{baseKafkaReceiver: &kafkaReceiver.BaseKafkaReceiver{}}//matchListReader:  pkgKafka.NewDecimalKafkaReader(global.ServerConfig.DecimalCricketKafkaEnv.GroupId, kafkaTopic.MatchListTopic),
	//matchListHandler: executor.NewMatchListExecutor(utils.NewPool(30), app.MatchPoolService),

	receiver.receiveMessage()
}

func (k *KafkaReceiver) receiveMessage() {
	go k.baseKafkaReceiver.Receive(k.matchListReader, k.matchListHandler.HandleMessage)
}
