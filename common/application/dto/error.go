package dto

import (
	"encoding/json"
)

type HttpMsgError struct {
	StatusCode int
	*ErrorData
}

type ErrorData struct {
	TraceId string `json:"traceId"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *HttpMsgError) Error() string {
	marshal, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return string(marshal)
}
