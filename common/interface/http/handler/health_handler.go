package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"linebot-go/common/application/utils"
)

func GetStatus(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	c.Status(http.StatusOK)
}
