package utils

import (
	"linebot-go/common/global"
	"strings"
)

func GetKafkaAddress() []string {
	return strings.Split(global.ServerConfig.DecimalCricketKafka.Servers, ",")
}

func GetKafkaAddressByArg(kafkaAddress string) []string {
	return strings.Split(kafkaAddress, ",")
}
