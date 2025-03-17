package job

import (
	"context"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"
)

type IJob interface {
	Execute(callback func(c context.Context))
	ExecuteNoLog(callback func())
}

type BaseJob struct {
	Name string
}

func (b *BaseJob) Execute(callback func(c context.Context)) {
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

	b.logBegin(actionLogs, now)
	c := context.WithValue(context.Background(), contextKey.ActionLogs, actionLogs)

	callback(c)

	utils.LogEnd(actionLogs, now)
}

func (b *BaseJob) ExecuteNoLog(callback func()) {
	defer func() {
		r := recover()
		if r != nil {
			return
		}
	}()

	callback()
}

func (b *BaseJob) logBegin(actionLogs map[string]any, now time.Time) {
	actionLogs[logKey.Type] = "job"
	actionLogs[logKey.ServerTime] = now.String()
	actionLogs[logKey.Id] = uuid.NewString()
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	actionLogs[logKey.JobName] = b.Name
}
