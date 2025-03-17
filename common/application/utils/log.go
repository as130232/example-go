package utils

import (
	"encoding/json"
	"fmt"
	"linebot-go/common/infrastructure/config"
	"linebot-go/common/infrastructure/consts/logKey"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
)

func BuildTag(str string) string {
	return "[" + str + "] "
}

func BuildFunctionTag() string {
	counter, _, _, ok := runtime.Caller(1)
	if !ok {
		return "[Unknown]"
	}
	fullPathFunctionName := runtime.FuncForPC(counter).Name()
	return "[" + buildPackageFunctionName(fullPathFunctionName) + "] "
}

func buildPackageFunctionName(fullName string) string {
	split := strings.Split(fullName, "/")
	return split[len(split)-1]
}

func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("Err %s: %s\n", msg, err)
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("Err %s: %s\n", msg, err)
	}
}

func LogTimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s, took: %s, tookLongTime:%v", name, elapsed, elapsed.Milliseconds() > 1000)
}

func LogTimeTrackIfOver(start time.Time, name string, thresholdMilli int64) {
	elapsed := time.Since(start)
	if elapsed.Milliseconds() > thresholdMilli {
		log.Printf("%s, took %s", name, elapsed)
	}
}

func LogServerPanic(r any) {
	logMessage := map[string]any{}
	logMessage[logKey.Id] = uuid.NewString()
	logMessage[logKey.ServerPanic] = true
	logMessage[logKey.ErrorMessage] = fmt.Sprintf("%+v", r)
	logMessage[logKey.StackTrace] = string(debug.Stack())
	b, _ := json.Marshal(logMessage)
	msg := string(b)
	log.Println(msg)
}

func LogServerConfig(serverConfig *config.ServerConfig) {
	b, _ := json.Marshal(serverConfig)
	log.Printf("LogServerConfig:%+v", string(b))
}
