package profile

import (
	"linebot-go/common/infrastructure/config"
	"linebot-go/common/infrastructure/consts/profile"
	"runtime"
	"strconv"
)

func InitProfile(appConfig *config.AppConfig) {
	blockProfileRateValue := appConfig.Config[profile.BlockProfileRate]
	blockProfileRate, err := strconv.Atoi(blockProfileRateValue)
	if err != nil {
		panic(err)
	}

	mutexProfileFractionValue := appConfig.Config[profile.MutexProfileFraction]
	mutexProfileFraction, err := strconv.Atoi(mutexProfileFractionValue)
	if err != nil {
		panic(err)
	}

	runtime.SetBlockProfileRate(blockProfileRate)         //sampling probability = 1%(1/rate)
	runtime.SetMutexProfileFraction(mutexProfileFraction) //sampling probability = 1%(1/rate)
}
