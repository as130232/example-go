package middleware

import (
	"example-go/common/application/utils"
	"github.com/gin-gonic/gin"
)

func CheckInternalIp(c *gin.Context) {
	ipAddress := utils.GetRemoteIp(c)
	utils.CheckIpAddressInCIDR(c, ipAddress)

	c.Next()
}
