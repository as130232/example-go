package appConfig

import (
	"example-go/common/application/utils"
	"example-go/common/infrastructure/config"
	"example-go/manifest/config/dev"
	"example-go/manifest/config/local"
	"example-go/manifest/config/prod"
	"os"
)

func NewServerConfig(appName string) *config.ServerConfig {
	appEnv := os.Getenv("APP_ENV")
	var serverConfig config.ServerConfig

	switch appEnv {
	case "local":
		serverConfig = local.CreateServerConfig()
	case "dev":
		serverConfig = dev.CreateServerConfig()
	case "prod":
		serverConfig = prod.CreateServerConfig()
	default:
		panic("APP_ENV must be local|dev|cqa|uat|prod")
	}

	serverConfig.AppEnv = appEnv

	if serverConfig.HttpServer != nil && appName != "" {
		serverConfig.HttpServer.ServerName = appName
	}

	if serverConfig.Log == nil {
		serverConfig.Log = &config.LogConfig{Type: "Console"}
	}

	serverConfig.InterfaceNameAndIpMap = utils.GetAllInterfaceNameAndIp()
	serverConfig.HostName = utils.GetHostName(serverConfig.InterfaceNameAndIpMap)

	if serverConfig.Pyroscope != nil {
		serverConfig.Pyroscope.PodName = serverConfig.HostName
	}

	utils.LogServerConfig(&serverConfig)

	return &serverConfig
}
