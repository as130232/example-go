package service

import (
	"example-go/common/global"
	"example-go/common/infrastructure/consts/env"
	"os"
	"strings"
)

// crawlerSwitch 預設開啟
var crawlerSwitch *bool

func GetCrawlerSwitch() bool {
	//初始化時根據環境判斷是否要觸發排程，後續可透過api更新開關切換
	if crawlerSwitch == nil {
		//若是gateway服務，排程只有dev、prod環境執行
		if strings.Contains(global.AppName, "gateway") {
			s := false
			if env.Prod == os.Getenv("APP_ENV") || env.Dev == os.Getenv("APP_ENV") {
				s = true
			}
			crawlerSwitch = &s
		} else {
			//若是一般服務，則所有環境排程預設都是打開
			s := true
			crawlerSwitch = &s
		}
	}
	return *crawlerSwitch
}

func UpdateCrawlerSwitch(isOpen bool) {
	s := isOpen
	crawlerSwitch = &s
}
