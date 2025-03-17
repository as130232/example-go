// Package pyroscope Pyroscope性能分析工具封装
package pyroscope

import (
	"github.com/grafana/pyroscope-go"
	"log"

	"linebot-go/common/infrastructure/config"
)

var profiler *pyroscope.Profiler

// Config 配置
type Config struct {
	ApplicationName string            //当前应用的名称
	Tags            map[string]string // 標籤
	ServerAddress   string            //服务地址
	AuthToken       string            //授权token
	LogLevel        int               //日志等级
	OpenGMB         bool              //default false,开启会造成过大的性能消耗
}

func Init(appName string, serverConfig *config.ServerConfig) {
	if serverConfig.Pyroscope == nil || !serverConfig.Pyroscope.Execute {
		return
	}

	cfg := Config{
		ApplicationName: appName,
		Tags: map[string]string{
			"PodName": serverConfig.HostName,
		},
		ServerAddress: serverConfig.Pyroscope.ServerAddress,
		AuthToken:     serverConfig.Pyroscope.AuthToken,
		LogLevel:      serverConfig.Pyroscope.LogLevel,
		OpenGMB:       serverConfig.Pyroscope.OpenGMB,
	}

	start(&cfg)
}

// Start 使用传入配置启动pyroscope
func start(cfg *Config) {
	if cfg.LogLevel == 0 {
		cfg.LogLevel = LevelDebug
	}

	var profileTypes = []pyroscope.ProfileType{
		pyroscope.ProfileCPU,
		pyroscope.ProfileAllocObjects,
		pyroscope.ProfileAllocSpace,
		pyroscope.ProfileInuseObjects,
		pyroscope.ProfileInuseSpace,
	}
	if cfg.OpenGMB {
		profileTypes = append(profileTypes,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		)
	}

	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.ApplicationName,
		Tags:            cfg.Tags,
		ServerAddress:   cfg.ServerAddress,
		AuthToken:       cfg.AuthToken,
		SampleRate:      5,
		Logger:          newLogger(cfg.LogLevel), //pyroscope.StandardLogger
		ProfileTypes:    profileTypes,
		DisableGCRuns:   false,
	})

	if err != nil {
		panic(err)
	}
}

// Stop 停止pyroscope
func Stop() {
	if profiler == nil {
		return
	}

	err := profiler.Stop()
	if err != nil {
		log.Printf("pyroscope stop error:%+v", err)
	}
}
