package executor

import (
	"context"
	"example-go/common/infrastructure/consts/contextKey"
	"github.com/panjf2000/ants/v2"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"

	"example-go/common/application/utils"
	"example-go/common/global"
	"example-go/common/infrastructure/consts/logKey"
	"github.com/google/uuid"
)

type IExecutor interface {
	Run(parentContext context.Context, callback func(c context.Context))
	RunWithKafka(parentContext context.Context, callback func(c context.Context), groupId string, msg *kafka.Message)
}

type BaseExecutor struct {
	Name string
	Pool *ants.Pool
}

var executorMap = new(sync.Map)

func (e *BaseExecutor) RunWithKafka(parentContext context.Context, callback func(c context.Context), groupId string, msg *kafka.Message) {
	e.Pool.Submit(func() {
		executorId := uuid.NewString()
		executorMap.Store(executorId, true)

		c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
		actionLogs, _ := utils.GetActionLogs(c)

		now := time.Now()

		defer func() {
			r := recover()
			if r != nil {
				utils.HandleErrorRecover(r, actionLogs, now)
				actionLogs[logKey.ExecutorKafkaMessage] = string(msg.Value)
				utils.OutputLog(actionLogs)
				executorMap.Delete(executorId)
				return
			}
		}()
		utils.LogKafkaInfo(actionLogs, groupId, msg)

		headers := msg.Headers
		id := ""
		if len(headers) > 0 {
			id = string(headers[0].Value)
		} else {
			id = utils.GetActionId(c)
		}

		e.logBegin(actionLogs, now, id)

		callback(c)

		e.logEnd(actionLogs, now)

		executorMap.Delete(executorId)
	})
}

func (e *BaseExecutor) RunWithKafkaNoLog(parentContext context.Context, callback func(c context.Context), groupId string, msg *kafka.Message) {
	e.Pool.Submit(func() {
		executorId := uuid.NewString()
		executorMap.Store(executorId, true)

		c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
		actionLogs, _ := utils.GetActionLogs(c)

		now := time.Now()

		defer func() {
			r := recover()
			if r != nil {
				utils.HandleErrorRecover(r, actionLogs, now)
				actionLogs[logKey.ExecutorKafkaMessage] = string(msg.Value)
				utils.OutputLog(actionLogs)
				executorMap.Delete(executorId)
				return
			}
		}()
		//utils.LogKafkaInfo(actionLogs, groupId, msg)

		//headers := msg.Headers
		//id := ""
		//if len(headers) > 0 {
		//	id = string(headers[0].Value)
		//} else {
		//	id = utils.GetActionId(c)
		//}

		//e.logBegin(actionLogs, now, id)

		callback(c)

		//e.logEnd(actionLogs, now)

		executorMap.Delete(executorId)
	})
}
func (e *BaseExecutor) Run(parentContext context.Context, callback func(c context.Context)) {
	e.Pool.Submit(func() {
		executorId := uuid.NewString()
		executorMap.Store(executorId, true)

		c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
		actionLogs, _ := utils.GetActionLogs(c)

		now := time.Now()

		defer func() {
			r := recover()
			if r != nil {
				utils.HandleErrorRecover(r, actionLogs, now)
				utils.OutputLog(actionLogs)
				executorMap.Delete(executorId)
				return
			}
		}()

		e.logBegin(actionLogs, now, utils.GetActionId(c))

		callback(c)

		e.logEnd(actionLogs, now)

		executorMap.Delete(executorId)
	})
}

func (e *BaseExecutor) logBegin(actionLogs map[string]any, now time.Time, id string) {
	actionLogs[logKey.Type] = "goroutine"
	actionLogs[logKey.ServerTime] = now.String()
	if id != "" {
		actionLogs[logKey.Id] = id
	} else {
		actionLogs[logKey.Id] = uuid.NewString()
	}
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	actionLogs[logKey.ExecutorName] = e.Name
	actionLogs[logKey.ExecutorCapacity] = e.Pool.Cap()
	actionLogs[logKey.ExecutorUsage] = e.Pool.Running()
}

func (e *BaseExecutor) logEnd(actionLogs map[string]any, now time.Time) {
	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	if timeUsed.Seconds() > 5 {
		actionLogs[logKey.SlowProcess] = true
	}
	utils.OutputLog(actionLogs)
}

func WaitForShutdown() {
	for waitCount := 1; waitCount < 120; waitCount++ {
		executorMapSize := utils.LenSyncMap(executorMap)
		if executorMapSize > 0 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
}
