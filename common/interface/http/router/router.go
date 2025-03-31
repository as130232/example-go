package router

import (
	"example-go/common/global"
	"example-go/common/infrastructure/consts/errorCode"
	"example-go/common/interface/http/handler"
	"example-go/common/interface/http/middleware"
	"github.com/Depado/ginprom"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if gin.DebugMode == global.ServerConfig.HttpServer.Mode {
		gin.SetMode(gin.DebugMode)
	} else if gin.TestMode == global.ServerConfig.HttpServer.Mode {
		gin.SetMode(gin.TestMode)
	}

	router := gin.New()

	//Define ginprom for monitoring
	prom := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(prom.Instrument(), middleware.Log, middleware.Recover)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthRouter := router.Group("/health")
	{
		healthRouter.GET("/", handler.GetStatus)
	}

	serviceRouter := router.Group("/" + global.AppName)
	serviceRouter.GET("/health", handler.GetStatus)
	serviceRouter.GET("/server-config", handler.GetServerConfig)
	serviceRouter.GET("/app-config", handler.GetAppConfig)
	serviceRouter.GET("/metrics", handler.GetMetrics)
	serviceRouter.GET("/version", handler.GetVersion)

	// pprof 性能分析介紹：https://geektutu.com/post/hpg-pprof.html
	pprof.Register(router, "/"+global.AppName+"/pprof") // 性能

	errorCode.CheckDuplicateErrorCode()

	return router
}

// SetSwitchRouter 設置開關路由
func SetSwitchRouter(serviceRouter *gin.RouterGroup) {
	switchRouter := serviceRouter.Group("/switch")
	switchRouter.GET("/crawler", handler.GetCrawlerSwitch)
	switchRouter.PUT("/crawler/:isOpen", handler.UpdateCrawlerSwitch)
}
