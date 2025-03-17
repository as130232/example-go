package handler

import (
	"net/http"

	"linebot-go/common/application/utils"

	"linebot-go/common/global"

	"github.com/gin-gonic/gin"
)

func GetAppConfig(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	c.JSON(http.StatusOK, global.AppConfig)
}
