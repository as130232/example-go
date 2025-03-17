package kafkaReceiver

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"linebot-go/common/infrastructure/consts/contextKey"
	"time"

	"github.com/google/uuid"
	"linebot-go/common/application/utils"
	"linebot-go/common/global"
	"linebot-go/common/infrastructure/consts/logKey"
)

type IKafkaReceiver interface {
	Receive(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msg *kafka.Message))
}

type IKafkaHandler interface {
	HandleMessage(c context.Context, reader *kafka.Reader, msg *kafka.Message)
}

type BaseKafkaReceiver struct {
}

func (b *BaseKafkaReceiver) receive(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msg *kafka.Message)) {
	logMessage := map[string]any{}
	now := time.Now()
	var message kafka.Message

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, logMessage, now)
			logMessage[logKey.KafkaMessage] = string(message.Value)
			utils.OutputLog(logMessage)
			return
		}
	}()

	b.logBegin(logMessage, now)
	c := context.WithValue(context.Background(), contextKey.ActionLogs, logMessage)

	message, err := reader.FetchMessage(c)
	now = time.Now() // 在FetchMessage才是真的時間
	if len(message.Headers) > 0 {
		logMessage[logKey.RefId] = logMessage[logKey.Id].(string)
		logMessage[logKey.Id] = string(message.Headers[0].Value)
	}

	if err != nil {
		logMessage[logKey.ErrorMessage] = fmt.Sprintf("%+v", err)
		utils.LogEnd(logMessage, now)
		return
	}
	utils.LogKafkaInfo(logMessage, reader.Config().GroupID, &message)
	logMessage[logKey.KafkaLag] = reader.Stats().Lag

	callback(c, reader, &message)

	utils.LogEnd(logMessage, now)
}

func (b *BaseKafkaReceiver) Receive(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msg *kafka.Message)) {
	for !global.IsShutdown {
		b.receive(reader, callback)
	}
}

func (b *BaseKafkaReceiver) logBegin(logMessage map[string]any, now time.Time) {
	logMessage[logKey.Type] = "kafka-receiver"
	logMessage[logKey.ServerTime] = now.String()
	logMessage[logKey.Id] = uuid.NewString()
	logMessage[logKey.ServiceName] = global.AppName
	logMessage[logKey.HostName] = global.ServerConfig.HostName
}
