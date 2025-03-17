package grpcClient

import (
	"context"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"
)

func UnaryInterceptor(parentContext context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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

	logBegin(actionLogs, now, method, req)

	md := metadata.Pairs(logKey.Id, actionLogs[logKey.Id].(string))
	ctx := metadata.NewOutgoingContext(c, md)
	err := invoker(ctx, method, req, reply, cc, opts...)

	actionLogs[logKey.Method] = method
	actionLogs[logKey.GrpcReply] = fmt.Sprintf("%+v", reply)
	utils.LogEnd(actionLogs, now)

	return err
}

func logBegin(logMessage map[string]any, now time.Time, method string, req interface{}) {
	logMessage[logKey.Type] = "grpc-client"
	logMessage[logKey.ServerTime] = now.String()
	logMessage[logKey.Id] = uuid.NewString()
	logMessage[logKey.ServiceName] = global.AppName
	logMessage[logKey.HostName] = global.ServerConfig.HostName
	logMessage[logKey.Method] = method
	logMessage[logKey.GrpcRequest] = fmt.Sprintf("%+v", req)
}
