package grpcServer

import (
	"context"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"time"

	"google.golang.org/grpc/metadata"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"

	"google.golang.org/grpc"
)

func UnaryInterceptor(parentContext context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	md, _ := metadata.FromIncomingContext(c)
	logBegin(actionLogs, now, md, info.FullMethod, req)

	m, err := handler(c, req)

	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	utils.LogEnd(actionLogs, now)

	return m, err
}

func logBegin(actionLogs map[string]any, now time.Time, md metadata.MD, method string, req interface{}) {
	actionLogs[logKey.Type] = "grpc-server"
	actionLogs[logKey.ServerTime] = now.String()
	actionLogs[logKey.Id] = uuid.NewString()
	refIds := md.Get(logKey.Id)
	if len(refIds) > 0 {
		actionLogs[logKey.Id] = refIds[0]
	}
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	actionLogs[logKey.Method] = method
	actionLogs[logKey.GrpcRequest] = fmt.Sprintf("%+v", req)
}
