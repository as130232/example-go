package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"example-go/common/application/utils"
	"example-go/common/global"
	"example-go/common/infrastructure/consts/contextKey"
	"example-go/common/infrastructure/consts/customHeaderKey"
	"example-go/common/infrastructure/consts/logKey"
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/google/uuid"
)

const MaxPrintBodyLen = 204800

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuffer *bytes.Buffer
}

func (blw bodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	blw.bodyBuffer.Write(b)
	return blw.ResponseWriter.Write(b)
}

func Log(c *gin.Context) {
	logBegin(c)

	blw := &bodyLogWriter{bodyBuffer: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	logEnd(c, blw)
}

func logBegin(c *gin.Context) {
	c.Set(contextKey.ActionId, uuid.NewString())

	actionLogs := make(map[string]any)
	c.Set(contextKey.ActionLogs, actionLogs)

	c.Set(contextKey.RequestTime, time.Now())

	c.Set(contextKey.RequestHeader, fmt.Sprintf("%+v", c.Request.Header))

	headerByte := []byte(c.GetString(contextKey.RequestHeader))
	c.Set(contextKey.RequestHeaderLength, len(headerByte))

	requestMethod := c.Request.Method
	if http.MethodPost == requestMethod || http.MethodPut == requestMethod {
		getRawDataNow := time.Now()
		body, _ := c.GetRawData()
		c.Set(contextKey.GetRawDataNowTime, time.Since(getRawDataNow))

		requestBodyBufferNow := time.Now()
		requestBodyBuffer := bytes.NewBuffer(body)
		c.Set(contextKey.RequestBodyBufferTime, time.Since(requestBodyBufferNow))

		nopCloserNow := time.Now()
		c.Request.Body = io.NopCloser(requestBodyBuffer)
		c.Set(contextKey.NopCloserTime, time.Since(nopCloserNow))

		c.Set(contextKey.RequestBody, string(body))
		c.Set(contextKey.RequestBodySize, len(body))
	}

	c.Set(contextKey.SqlTrace, c.Request.Header.Get(customHeaderKey.Trace))
	c.Set(contextKey.DeviceType, c.Request.Header.Get(customHeaderKey.DeviceType))
}

func logEnd(c *gin.Context, blw *bodyLogWriter) {
	actionId := c.GetString(contextKey.ActionId)
	requestTime := c.GetTime(contextKey.RequestTime)
	requestMethod := c.Request.Method
	contentType := c.Request.Header.Get(headers.ContentType)

	responseBody := strings.Trim(blw.bodyBuffer.String(), "\n")
	if len(responseBody) > MaxPrintBodyLen {
		responseBody = responseBody[:(MaxPrintBodyLen - 1)]
	}

	actionLogs := make(map[string]any)
	actionLogs[logKey.GetRawDataNowTime] = c.GetInt64(contextKey.GetRawDataNowTime)
	actionLogs[logKey.RequestBodyBufferTime] = c.GetInt64(contextKey.RequestBodyBufferTime)
	actionLogs[logKey.NopCloserTime] = c.GetInt64(contextKey.NopCloserTime)
	actionLogs[logKey.Type] = "http"
	actionLogs[logKey.ServerTime] = time.Now().String()
	actionLogs[logKey.Id] = actionId
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	actionLogs[logKey.Method] = requestMethod
	actionLogs[logKey.RequestUri] = c.Request.RequestURI
	actionLogs[logKey.UserAgent] = c.Request.UserAgent()
	actionLogs[logKey.ContentType] = contentType
	actionLogs[logKey.ClientIp] = utils.GetClientIp(c)
	//actionLogs[logKey.RequestHeader] = c.GetString(contextKey.RequestHeader)
	actionLogs[logKey.RequestHeaderLength] = c.GetInt(contextKey.RequestHeaderLength)
	actionLogs[logKey.RequestTime] = requestTime.String()
	actionLogs[logKey.ContentLength] = c.Request.ContentLength
	actionLogs[logKey.ResponseStatus] = c.Writer.Status()
	timeUsed := time.Since(requestTime)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	if timeUsed.Seconds() > 5 {
		actionLogs[logKey.SlowHttp] = true
	}

	processActionLogs(actionLogs, c, responseBody)
	utils.OutputLog(actionLogs)
}

func processActionLogs(actionLogs map[string]any, c *gin.Context, responseBody string) {
	if c.GetBool(contextKey.HasError) {
		actionLogs[logKey.RequestHeader] = c.GetString(contextKey.RequestHeader)
		actionLogs[logKey.RequestHeaderLength] = c.GetInt(contextKey.RequestHeaderLength)

		if http.MethodPost == actionLogs[logKey.Method] || http.MethodPut == actionLogs[logKey.Method] {
			actionLogs[logKey.RequestBody] = c.GetString(contextKey.RequestBody)
		}

		if responseBody != "" {
			actionLogs[logKey.ResponseBody] = responseBody
		}

		actionLogs[logKey.ErrorType] = c.GetString(contextKey.ErrorType)
		actionLogs[logKey.ErrorMessage] = c.GetString(contextKey.ErrorMessage)

		if http.StatusInternalServerError == c.Writer.Status() {
			actionLogs[logKey.Trace] = c.GetString(contextKey.StackTrace)
		}
	} else if c.GetBool(contextKey.LogRequestBody) {
		actionLogs[logKey.RequestBody] = c.GetString(contextKey.RequestBody)
	}

	if http.MethodPost == actionLogs[logKey.Method] || http.MethodPut == actionLogs[logKey.Method] {
		actionLogs[logKey.RequestBodySize] = c.GetInt(contextKey.RequestBodySize)
	}

	for k, v := range c.GetStringMap(contextKey.ActionLogs) {
		actionLogs[k] = v
	}
}
