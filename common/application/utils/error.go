package utils

import (
	"linebot-go/common/application/dto"
)

func GenErrorMsg(statusCode int, msgType, message string) error {
	return &dto.HttpMsgError{StatusCode: statusCode, ErrorData: &dto.ErrorData{Type: msgType, Message: message}}
}
