package handler

import (
	"net/http"

	"example-go/common/global"

	"example-go/common/application/utils"
	"github.com/gin-gonic/gin"
)

func GetServerConfig(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	c.JSON(http.StatusOK, global.ServerConfig)
}
