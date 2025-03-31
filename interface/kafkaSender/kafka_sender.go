package kafkaSender

import (
	commonKafka "example-go/common/interface/kafkaSender"
)

type KafkaSender struct {
	baseKafkaSender commonKafka.IKafkaSender
	//deltaWriter     []*kafka.Writer
}

func NewKafkaSender() *KafkaSender {
	return &KafkaSender{
		baseKafkaSender: &commonKafka.BaseKafkaSender{},
		//deltaWriter:     pkgKafka.NewDecimalCricketKafkaWriterList(kafkaTopic.DeltaTopic),
	}
}

// SendDeltaMessage 全環境發送賽事列表，用於賽事落地
//func (k *KafkaSender) SendDeltaMessage(dto kafkaDto.DecimalCricketRaw) {
//	dtoBytes, _ := json.Marshal(dto)
//	k.baseKafkaSender.SendAllNoLog(k.deltaWriter, &kafka.Message{
//		Value: dtoBytes,
//	})
//}
