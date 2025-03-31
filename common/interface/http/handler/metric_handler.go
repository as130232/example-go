package handler

import (
	"example-go/common/application/utils"
	"github.com/gin-gonic/gin"
)

func GetMetrics(c *gin.Context) {
	utils.CheckIpAddressInCIDR(c, utils.GetRemoteIp(c))

	//metricObject := global.MetricService.GetMetric()
	//utils.SetActionLog(c, logKey.CpuTotalPercent, metricObject.GoPsUtilMetric.CpuTotalPercent)
	//utils.SetActionLog(c, logKey.MemTotalPercent, metricObject.GoPsUtilMetric.MemUsedPercent)
	//utils.SetActionLog(c, logKey.NumGC, metricObject.GcMetric.NumGC)
	//utils.SetActionLog(c, logKey.IncrNumGC, metricObject.GcMetric.IncrNumGC)
	//utils.SetActionLog(c, logKey.PauseTotalNs, metricObject.GcMetric.PauseTotalNs)
	//utils.SetActionLog(c, logKey.IncrPauseNs, metricObject.GcMetric.IncrPauseNs)

	//c.JSON(http.StatusOK, metricObject)
}
