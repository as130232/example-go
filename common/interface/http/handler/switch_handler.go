package handler

import (
	"example-go/common/application/dto"
	"example-go/common/application/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary 排程開關
// @Tags	開關
// @Version 1.0
// @Produce	json
// @Success	200 {object} object{crawlerSwitch=boolean}
// @Router	/crawler/switch [get]
func GetCrawlerSwitch(c *gin.Context) {
	c.JSON(http.StatusOK, dto.CreateResponse(c, gin.H{"crawlerSwitch": service.GetCrawlerSwitch()}))
}

// @Summary 排程開關
// @Tags	開關
// @Version 1.0
// @Produce	json
// @Param isOpen path bool true "是否開啟開關"
// @Success	200 {object} object{crawlerSwitch=boolean}
// @Router	/crawler/switch [put]
func UpdateCrawlerSwitch(c *gin.Context) {
	isOpenStr := c.Param("isOpen")
	isOpen, _ := strconv.ParseBool(isOpenStr)
	service.UpdateCrawlerSwitch(isOpen)
	GetCrawlerSwitch(c)
}
