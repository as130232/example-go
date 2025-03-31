package appConfig

import (
	"example-go/common/application/utils/profile"
	"example-go/common/infrastructure/config"
	"example-go/manifest/config/dev"
	"example-go/manifest/config/local"
	"example-go/manifest/config/prod"
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
