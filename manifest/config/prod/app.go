package prod

import (
	"example-go/common/infrastructure/config"
	"example-go/common/infrastructure/consts/profile"
)

func CreateAppConfig() config.AppConfig {
	var conf = make(map[string]string)
	conf[profile.BlockProfileRate] = "0"
	conf[profile.MutexProfileFraction] = "0"

	return config.AppConfig{
		Config: conf,
	}
}
