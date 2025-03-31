package handler

import (
	"net/http"

	"example-go/common/global"

	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"commitId": global.CommitId, "buildTime": global.BuildTime})
}
