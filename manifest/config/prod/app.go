package prod

import (
	"linebot-go/common/infrastructure/config"
	"linebot-go/common/infrastructure/consts/profile"
)

func CreateAppConfig() config.AppConfig {
	var conf = make(map[string]string)
	conf[profile.BlockProfileRate] = "0"
	conf[profile.MutexProfileFraction] = "0"

	return config.AppConfig{
		Config: conf,
	}
}
