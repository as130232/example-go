package handler

import (
	"net/http"

	"linebot-go/common/global"

	"github.com/gin-gonic/gin"
	"linebot-go/common/application/utils"
)

func GetServerConfig(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	c.JSON(http.StatusOK, global.ServerConfig)
}
