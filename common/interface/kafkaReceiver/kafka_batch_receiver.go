package kafkaReceiver

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"linebot-go/common/application/utils"
	"linebot-go/common/global"
	"linebot-go/common/infrastructure/consts/contextKey"
	"linebot-go/common/infrastructure/consts/logKey"
	"sync"
	"time"
)

type IKafkaBatchReceiver interface {
	Receive(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message))
	ConsumeMessageBufferByTimeOut(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message))
}

type BaseKafkaBatchReceiver struct {
	BufferSize    int
	MsgBuffer     []kafka.Message
	consumeSignal chan bool
	BatchTimeOut  time.Duration
	mux           *sync.Mutex
}

func NewBaseKafkaBatchReceiver(bufferSize int, timeOut time.Duration) *BaseKafkaBatchReceiver {
	return &BaseKafkaBatchReceiver{
		BufferSize:    bufferSize,
		MsgBuffer:     make([]kafka.Message, 0, bufferSize),
		consumeSignal: make(chan bool, 1),
		BatchTimeOut:  timeOut,
		mux:           &sync.Mutex{},
	}
}

func (b *BaseKafkaBatchReceiver) receive(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message)) {
	actionLogs := make(map[string]any)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.receiveLogBegin(actionLogs)
	c := context.WithValue(context.Background(), contextKey.ActionLogs, actionLogs)

	m, err := reader.FetchMessage(c)
	now = time.Now() // 在FetchMessage才是真的時間

	actionLogs[logKey.ServerTime] = now.String()
	if len(m.Headers) > 0 {
		actionLogs[logKey.RefId] = actionLogs[logKey.Id].(string)
		actionLogs[logKey.Id] = string(m.Headers[0].Value)
	}
	if err != nil {
		actionLogs[logKey.ErrorMessage] = fmt.Sprintf("%+v", err)
		b.receiveLogEnd(actionLogs, now)
		return
	}
	utils.LogKafkaInfo(actionLogs, reader.Config().GroupID, &m)
	actionLogs[logKey.KafkaLag] = reader.Stats().Lag

	b.mutexReceive(c, m, reader, callback)

	b.receiveLogEnd(actionLogs, now)
}

func (b *BaseKafkaBatchReceiver) mutexReceive(c context.Context, m kafka.Message, reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message)) {
	defer func() {
		b.mux.Unlock()
		r := recover()
		if r != nil {
			panic(r)
		}
	}()
	b.mux.Lock()
	b.MsgBuffer = append(b.MsgBuffer, m)
	if len(b.MsgBuffer) == b.BufferSize {
		callback(c, reader, b.MsgBuffer)
		b.MsgBuffer = b.MsgBuffer[:0] // flush buffer, keep allocated memory
	}
}

func (b *BaseKafkaBatchReceiver) Receive(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message)) {
	for !global.IsShutdown {
		b.receive(reader, callback)
	}
}

func (b *BaseKafkaBatchReceiver) consumeMessageBufferByTimeOut(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message)) {
	actionLogs := make(map[string]any)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.consumeLogBegin(actionLogs)
	c := context.WithValue(context.Background(), contextKey.ActionLogs, actionLogs)

	select {
	case <-time.After(b.BatchTimeOut):
		actionLogs[logKey.ServerTime] = now.String()

		b.mutexConsumeMessageBufferByTimeOut(c, reader, callback)

		b.consumeLogEnd(actionLogs, now)
	}
}

func (b *BaseKafkaBatchReceiver) mutexConsumeMessageBufferByTimeOut(c context.Context, reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message)) {
	defer func() {
		b.mux.Unlock()
		r := recover()
		if r != nil {
			panic(r)
		}
	}()
	b.mux.Lock()
	utils.SetActionLog(c, logKey.KafkaTopicName, reader.Config().Topic)
	if len(b.MsgBuffer) == 0 || len(b.MsgBuffer) == b.BufferSize {
		//log.Println("no message in buffer")
		return
	}
	callback(c, reader, b.MsgBuffer)
	b.MsgBuffer = b.MsgBuffer[:0] // flush buffer, keep allocated memory
}

func (b *BaseKafkaBatchReceiver) ConsumeMessageBufferByTimeOut(reader *kafka.Reader, callback func(c context.Context, reader *kafka.Reader, msgs []kafka.Message)) {
	for !global.IsShutdown {
		b.consumeMessageBufferByTimeOut(reader, callback)
	}
}

func (b *BaseKafkaBatchReceiver) receiveLogBegin(actionLogs map[string]any) {
	actionLogs[logKey.Type] = "kafka-receiver"
	actionLogs[logKey.Id] = uuid.NewString()
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
}

func (b *BaseKafkaBatchReceiver) receiveLogEnd(actionLogs map[string]any, now time.Time) {
	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	utils.OutputLog(actionLogs)
}

func (b *BaseKafkaBatchReceiver) consumeLogBegin(actionLogs map[string]any) {
	actionLogs[logKey.Type] = "kafka-receiver"
	actionLogs[logKey.Id] = uuid.NewString()
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
}

func (b *BaseKafkaBatchReceiver) consumeLogEnd(actionLogs map[string]any, now time.Time) {
	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	utils.OutputLog(actionLogs)
}
