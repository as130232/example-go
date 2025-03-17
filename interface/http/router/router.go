package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"linebot-go/cmd"
	"linebot-go/common/global"
	commonMiddleware "linebot-go/common/interface/http/middleware"
	"linebot-go/common/interface/http/router"
	"linebot-go/docs"
)

func InitRouter(app *cmd.App) *gin.Engine {
	ginRouter := router.InitRouter()

	docs.SwaggerInfo.Title = "decimal-cricket-websocket API"
	docs.SwaggerInfo.Description = "板球 websocket 服務"
	docs.SwaggerInfo.Version = "X.0"
	if "local" == global.ServerConfig.AppEnv {
		docs.SwaggerInfo.Host = "localhost" + global.ServerConfig.HttpServer.Address
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else {
		docs.SwaggerInfo.Host = "deci-api.ljb" + global.ServerConfig.AppEnv + ".site"
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
	docs.SwaggerInfo.BasePath = "/decimal-cricket-websocket"

	docs.SwaggerInfo.LeftDelim = "{{"
	docs.SwaggerInfo.RightDelim = "}}"

	serviceRouter := ginRouter.Group("/" + global.AppName)

	// Swagger Router & API
	swaggerRouter := serviceRouter.Group("/swagger")
	swaggerRouter.Use(commonMiddleware.CheckInternalIp) // IP check
	swaggerRouter.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return ginRouter
}
