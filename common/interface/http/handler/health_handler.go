package handler

import (
	"net/http"

	"example-go/common/application/utils"
	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	c.Status(http.StatusOK)
}
