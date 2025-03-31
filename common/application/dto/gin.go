package dto

import (
	"example-go/common/infrastructure/consts/contextKey"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	TraceId   string `json:"traceId"`
	ErrorType string `json:"errorType"`
}

func CreateEmptyResponse(c *gin.Context) *Response {
	return CreateResponse(c, struct{}{})
}

func CreateResponse(c *gin.Context, data any) *Response {
	if data == nil {
		data = struct{}{}
	}
	return &Response{
		Code:      0,
		Message:   "Success",
		Data:      data,
		TraceId:   c.GetString(contextKey.ActionId),
		ErrorType: "",
	}
}

func CreateErrorResponse(code int, message string, traceId string, errorType string) *Response {
	if code == 0 { // 找不到error code
		code = -1
	}
	return &Response{
		Code:      code,
		Message:   message,
		Data:      struct{}{},
		TraceId:   traceId,
		ErrorType: errorType,
	}
}
