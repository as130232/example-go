package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"example-go/common/infrastructure/consts/errorCode"

	"example-go/common/application/dto"
	"example-go/common/infrastructure/consts/contextKey"
	"example-go/common/infrastructure/consts/errorType"
	"github.com/gin-gonic/gin"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.Set(contextKey.HasError, true)

			switch r := r.(type) {
			case *dto.HttpMsgError:
				handleHttpMsgErrorRecover(r, c)
			default:
				handleErrorRecover(r, c)
			}
			c.Abort()
			return
		}
	}()

	c.Next()
}

func handleHttpMsgErrorRecover(ei *dto.HttpMsgError, c *gin.Context) {
	if ei.ErrorData == nil {
		c.Set(contextKey.ErrorType, errorType.InternalServerError)
		c.Set(contextKey.StackTrace, string(debug.Stack()))
		c.JSON(http.StatusInternalServerError, dto.CreateErrorResponse(errorCode.ErrorCodeMap[errorType.InternalServerError], fmt.Sprintf("%+v", ei), c.GetString(contextKey.ActionId), errorType.InternalServerError))
		return
	}

	if ei.ErrorData.Type == "" {
		ei.ErrorData.Type = errorType.InternalServerError
	}
	c.Set(contextKey.ErrorType, ei.ErrorData.Type)
	c.Set(contextKey.ErrorMessage, ei.ErrorData.Message)
	ei.ErrorData.TraceId = c.GetString(contextKey.ActionId)
	if ei.StatusCode != 0 {
		c.JSON(ei.StatusCode, dto.CreateErrorResponse(errorCode.ErrorCodeMap[ei.Type], ei.Message, c.GetString(contextKey.ActionId), ei.Type))
	} else {
		c.Set(contextKey.StackTrace, string(debug.Stack()))
		c.JSON(http.StatusInternalServerError, dto.CreateErrorResponse(errorCode.ErrorCodeMap[ei.Type], ei.Message, c.GetString(contextKey.ActionId), ei.Type))
	}
}

func handleErrorRecover(r any, c *gin.Context) {
	ed := dto.ErrorData{TraceId: c.GetString(contextKey.ActionId), Type: errorType.InternalServerError, Message: fmt.Sprintf("%+v", r)}

	c.Set(contextKey.ErrorType, ed.Type)
	c.Set(contextKey.ErrorMessage, ed.Message)
	c.Set(contextKey.StackTrace, string(debug.Stack()))
	c.JSON(http.StatusInternalServerError, dto.CreateErrorResponse(errorCode.ErrorCodeMap[ed.Type], ed.Message, c.GetString(contextKey.ActionId), ed.Type))
}
