package dev

import (
	"linebot-go/common/infrastructure/config"
	"linebot-go/common/infrastructure/consts/profile"
)

func CreateAppConfig() config.AppConfig {
	var conf = make(map[string]string)
	conf[profile.BlockProfileRate] = "1"
	conf[profile.MutexProfileFraction] = "1"

	return config.AppConfig{
		Config: conf,
	}
}
