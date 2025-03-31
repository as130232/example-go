package utils

import (
	"example-go/common/global"
	"fmt"
)

func SendTelegramMessage(message string) {
	if global.TelegramBot != nil {
		if len(message) > 4096 {
			message = message[0:4096]
		}
		global.TelegramBot.SendMessage(message)
	}
}

func SendTelegramServerPanicMessage(commitId string, buildTime string) {
	if global.TelegramBot != nil && global.ServerConfig != nil {
		global.TelegramBot.SendMessage(fmt.Sprintf("=====\n%s %s server panic, \nhostname:%s, \ngit commit:%s, \nbuildTime:%s\n=====", global.ServerConfig.AppEnv, global.AppName, global.ServerConfig.HostName, commitId, buildTime))
	}
}
