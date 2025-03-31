package utils

import (
	"context"
	"encoding/json"
	"example-go/common/global"
	"example-go/common/infrastructure/consts/contextKey"
	"example-go/common/infrastructure/consts/errorType"
	"example-go/common/infrastructure/consts/logKey"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

func TimingLog(logMessage map[string]any, key string, examine func()) {
	now := time.Now()
	examine()
	logMessage[key] = time.Since(now)
}

func GetOrGenActionId(actionLogs map[string]any) string {
	var ret string
	if actionId, ok := actionLogs[logKey.Id].(string); ok {
		ret = actionId
	}
	return ret
}

func SetActionLog(c context.Context, key string, value any) {
	actionLogs := c.Value(contextKey.ActionLogs)
	actionLogs.(map[string]any)[key] = value
}

func GetActionLogs(ctx context.Context) (map[string]any, bool) {
	actionLogs := ctx.Value(contextKey.ActionLogs)
	if actionLogs == nil {
		return nil, false
	}

	return actionLogs.(map[string]any), true
}

func SetActionId(c context.Context, actionId string) {
	actionLogs := c.Value(contextKey.ActionLogs)
	actionLogs.(map[string]any)[logKey.Id] = actionId
}

func GetActionId(c context.Context) string {
	actionLogs := c.Value(contextKey.ActionLogs)
	if actionLogs == nil {
		return ""
	}

	actionId, exists := actionLogs.(map[string]any)[logKey.Id]
	if exists {
		return actionId.(string)
	} else {
		return ""
	}
}

func KafkaLog(actionLogs map[string]any) {
	actionLogsByte, _ := json.Marshal(actionLogs)
	log.Println(string(actionLogsByte))
}

func ConsoleLog(actionLogs map[string]any) {
	var logValue strings.Builder
	for k, v := range actionLogs {
		logValue.WriteString("\n")
		logValue.WriteString(k)
		logValue.WriteString("=")
		vByte, _ := json.Marshal(v)
		logValue.WriteString(string(vByte))
	}

	log.Println(logValue.String())
}

func HandleErrorRecover(r any, actionLogs map[string]any, now time.Time) {
	actionLogs[logKey.ErrorType] = errorType.InternalServerError
	actionLogs[logKey.ErrorMessage] = fmt.Sprintf("%+v", r)
	actionLogs[logKey.StackTrace] = string(debug.Stack())
	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	if timeUsed.Seconds() > 5 {
		actionLogs[logKey.SlowProcess] = true
	}
}

func LogKafkaInfo(actionLogs map[string]any, groupId string, msg *kafka.Message) {
	actionLogs[logKey.KafkaGroupId] = groupId
	actionLogs[logKey.KafkaTopicName] = msg.Topic
	actionLogs[logKey.KafkaPartition] = msg.Partition
	actionLogs[logKey.KafkaOffset] = msg.Offset
}

func LogEnd(actionLogs map[string]any, now time.Time) {
	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	if timeUsed.Seconds() > 5 {
		actionLogs[logKey.SlowProcess] = true
	}
	OutputLog(actionLogs)
}

func OutputLog(actionLogs map[string]any) {
	//processedactionLogs := processactionLogsKeyDotToUnderscore(actionLogs)
	if global.ServerConfig.Log.Type == "Json" {
		KafkaLog(actionLogs)
	} else {
		ConsoleLog(actionLogs)
	}
}

// processactionLogsKeyDotToUnderscore
// log message key 全部使用 _ 取代 .
// 運維表示 Elasticsearch using dot key 會有建 index 階層的 issue
// https://stackoverflow.com/questions/59472323/elasticsearch-dynamic-field-mapping-and-json-dot-notation
func processActionLogsKeyDotToUnderscore(actionLogs map[string]any) map[string]any {
	processedActionLogs := make(map[string]any, len(actionLogs))
	for k, v := range actionLogs {
		if strings.Contains(k, ".") {
			processedActionLogs[strings.Replace(k, ".", "_", -1)] = v
		} else {
			processedActionLogs[k] = v
		}
	}
	return processedActionLogs
}
