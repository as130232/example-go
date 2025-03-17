package middleware

import (
	"github.com/gin-gonic/gin"
	"linebot-go/common/application/utils"
)

func CheckInternalIp(c *gin.Context) {
	ipAddress := utils.GetRemoteIp(c)
	utils.CheckIpAddressInCIDR(c, ipAddress)

	c.Next()
}
