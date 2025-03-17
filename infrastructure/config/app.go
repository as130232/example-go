package appConfig

import (
	"linebot-go/common/application/utils/profile"
	"linebot-go/common/infrastructure/config"
	"linebot-go/manifest/config/dev"
	"linebot-go/manifest/config/local"
	"linebot-go/manifest/config/prod"
	"os"
)

func NewAppConfig() *config.AppConfig {
	appEnv := os.Getenv("APP_ENV")
	var appConfig config.AppConfig

	switch appEnv {
	case "local":
		appConfig = local.CreateAppConfig()
	case "dev":
		appConfig = dev.CreateAppConfig()
	case "cqa":
		appConfig = dev.CreateAppConfig()
	case "uat":
		appConfig = prod.CreateAppConfig()
	case "prod":
		appConfig = prod.CreateAppConfig()
	default:
		panic("APP_ENV must be local|dev|cqa|uat|prod")
	}

	profile.InitProfile(&appConfig)

	return &appConfig
}
