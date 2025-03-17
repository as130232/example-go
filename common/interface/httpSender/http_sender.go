package httpSender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-http-utils/headers"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

const (
	ContentTypeNone = ""
	ContentTypeJson = binding.MIMEJSON
)

var (
	transport *http.Transport
	client    *http.Client
)

func init() {
	transport = http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 10
	client = &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}
}

// SendRequest
// actionId: give "" if no actionId
// contentType: give "" if no contentType
// header: give nil if no header
// reqBody: give nil if no request body
func SendRequest(actionId string, httpMethod string, uri string, contentType string, header map[string]string, reqBody any, resBodyLog bool) (resBytes []byte, err error) {
	start := time.Now()
	logMessage := make(map[string]any)
	logBegin(logMessage, start)
	defer logEnd(logMessage, start)

	// header
	if header == nil {
		header = make(map[string]string)
	}
	if actionId == "" {
		actionId = uuid.NewString()
	}
	header[logKey.Id] = actionId
	header["Content-Type"] = contentType
	logMessage[logKey.Id] = actionId

	// body
	reqBodyBytes := []byte("")
	if reqBody != nil {
		reqBodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			logMessage[logKey.ErrorMessage] = err.Error()
			return
		}
	}

	resp, err := do(logMessage, httpMethod, uri, header, reqBodyBytes)
	if err != nil {
		logMessage[logKey.HttpClientResponseError] = err.Error()
		return
	}

	if resp.StatusCode != http.StatusOK {
		logMessage[logKey.HttpClientResponseError] = resp.Status
	}

	// read body
	var bodyReader = resp.Body
	defer resp.Body.Close()
	// gzip decompression if needed
	encoding := resp.Header.Get(headers.ContentEncoding)
	if len(encoding) > 0 && encoding == "gzip" {
		bodyReader, err = gzip.NewReader(resp.Body)
		logMessage[logKey.HttpClientResponseError] = err.Error()
		if err != nil {
			return
		}
		defer bodyReader.Close()
	}

	resBytes, err = io.ReadAll(bodyReader)
	if err != nil {
		logMessage[logKey.HttpClientResponseError] = err.Error()
		return
	}

	if resBodyLog {
		logMessage[logKey.HttpClientResponseBody] = string(resBytes)
	}
	return resBytes, nil
}

func do(logMessage map[string]any, httpMethod string, uri string, header map[string]string, reqBody []byte) (*http.Response, error) {
	logMessage[logKey.HttpClientRequestMethod] = httpMethod
	logMessage[logKey.HttpClientRequestUri] = uri
	logMessage[logKey.HttpClientRequestHeader] = fmt.Sprintf("%+v", header)
	if len(reqBody) > 0 {
		logMessage[logKey.HttpClientRequestBody] = string(reqBody)
	}

	req, err := http.NewRequest(httpMethod, uri, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	logMessage[logKey.HttpClientRequestHost] = req.URL.Host
	logMessage[logKey.HttpClientRequestPort] = req.URL.Port()
	logMessage[logKey.HttpClientRequestPath] = req.URL.Path

	resp, err := client.Do(req)
	if resp != nil {
		logMessage[logKey.HttpClientResponseHeader] = fmt.Sprintf("%+v", resp.Header)
		logMessage[logKey.HttpClientResponseContentLength] = resp.ContentLength
		logMessage[logKey.HttpClientResponseStatus] = resp.Status
		logMessage[logKey.HttpClientResponseStatusCode] = resp.StatusCode
	}
	return resp, err
}

func logBegin(logMessage map[string]any, start time.Time) {
	logMessage[logKey.Type] = "http-client"
	logMessage[logKey.ServerTime] = start.String()
	logMessage[logKey.ServiceName] = global.AppName
}

func logEnd(logMessage map[string]any, start time.Time) {
	executeDuration := time.Since(start)
	logMessage[logKey.TimeUsed] = executeDuration.String()
	logMessage[logKey.TimeUsedNano] = executeDuration
	if executeDuration.Seconds() > 5 {
		logMessage[logKey.SlowHttp] = true
	}

	if reqBody, ok := logMessage[logKey.HttpClientRequestBody]; ok {
		reqBodyLog := map[string]any{
			logKey.Type:  "httpSender",
			logKey.Id:    logMessage[logKey.Id],
			"z-req-body": reqBody,
		}
		utils.OutputLog(reqBodyLog)
	}
	delete(logMessage, logKey.HttpClientRequestBody)
	utils.OutputLog(logMessage)
}
