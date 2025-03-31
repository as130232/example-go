package handler

import (
	"net/http"

	"example-go/common/application/utils"

	"example-go/common/global"

	"github.com/gin-gonic/gin"
)

func GetAppConfig(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	c.JSON(http.StatusOK, global.AppConfig)
}
