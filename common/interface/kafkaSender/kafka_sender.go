package kafkaSender

import (
	"context"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"
)

const retries = 30

type IKafkaSender interface {
	Send(parentContext context.Context, writer *kafka.Writer, msg *kafka.Message)
	SendWithErrorHandle(parentContext context.Context, writer *kafka.Writer, msg *kafka.Message, errorHandle func(c context.Context))
	SendAllNoLog(writerList []*kafka.Writer, msg *kafka.Message)
	SendWithErrorHandleNoLog(writer *kafka.Writer, msg *kafka.Message, errorHandle func())
}

type BaseKafkaSender struct {
}

func (b *BaseKafkaSender) Send(parentContext context.Context, writer *kafka.Writer, msg *kafka.Message) {
	b.SendWithErrorHandle(parentContext, writer, msg, nil)
}

func (b *BaseKafkaSender) SendWithErrorHandle(parentContext context.Context, writer *kafka.Writer, msg *kafka.Message, errorHandle func(c context.Context)) {
	go func() {
		c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
		actionLogs, _ := utils.GetActionLogs(c)

		now := time.Now()
		actionId := utils.GetActionId(c)

		defer func() {
			r := recover()
			if r != nil {
				log.Printf("SendWithErrorHandle error:%+v,actionId:%s", r, actionId)
				utils.HandleErrorRecover(r, actionLogs, now)
				utils.OutputLog(actionLogs)
				if errorHandle != nil {
					errorHandle(c)
				}
				return
			}
		}()

		b.logBegin(writer.Topic, actionLogs, now)
		if actionId != "" { // 優先使用串接 actionId
			actionLogs[logKey.Id] = actionId
		}

		msg.Headers = []kafka.Header{{Key: logKey.Id, Value: []byte(actionLogs[logKey.Id].(string))}}
		for retry := 1; retry <= retries; retry++ {
			err := writer.WriteMessages(c, *msg)
			if err != nil {
				log.Printf("SendWithErrorHandle error:%+v,actionId:%s,retry:%d", err, actionId, retry)
				if retry == retries {
					panic(err)
				} else {
					time.Sleep(time.Second * 1)
					continue
				}
			}
			actionLogs[logKey.KafkaRetry] = retry
			break
		}

		utils.LogEnd(actionLogs, now)
	}()
}

func (b *BaseKafkaSender) SendAllNoLog(writerList []*kafka.Writer, msg *kafka.Message) {
	for _, writer := range writerList {
		b.SendWithErrorHandleNoLog(writer, msg, nil)
	}
}

func (b *BaseKafkaSender) SendWithErrorHandleNoLog(writer *kafka.Writer, msg *kafka.Message, errorHandle func()) {
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				if errorHandle != nil {
					errorHandle()
				}
				return
			}
		}()

		msg.Headers = []kafka.Header{{Key: logKey.Id, Value: []byte(uuid.NewString())}}
		for retry := 0; retry < retries; retry++ {
			err := writer.WriteMessages(context.Background(), *msg)
			if err != nil {
				if retry == retries {
					panic(err)
				} else {
					time.Sleep(time.Second * 1)
					continue
				}
			}
			break
		}
	}()
}

func (b *BaseKafkaSender) logBegin(topicName string, actionLogs map[string]any, now time.Time) {
	actionLogs[logKey.Type] = "kafka-sender"
	actionLogs[logKey.ServerTime] = now.String()
	actionLogs[logKey.Id] = uuid.NewString()
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	actionLogs[logKey.KafkaTopicName] = topicName
}
