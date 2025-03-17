package httpClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	httpHeaders "github.com/go-http-utils/headers"
	"github.com/google/uuid"
)

const (
	ContentTypeNone = ""
	ContentTypeJson = "application/json"
)

type BaseHttpClient struct {
	httpClient *http.Client
}

var httpClientIdMap = new(sync.Map)

func NewBaseHttpClient(httpClient *http.Client) *BaseHttpClient {
	return &BaseHttpClient{httpClient: httpClient}
}

func (b *BaseHttpClient) Get(url string, headerMap map[string]string, id string) (*http.Response, error) {
	if headerMap == nil {
		headerMap = make(map[string]string)
	}

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodGet, url, headerMap, nil, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) GetJson(url string, id string) (*http.Response, error) {
	headerMap := make(map[string]string)
	headerMap[httpHeaders.ContentType] = ContentTypeJson

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodGet, url, headerMap, nil, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) Post(url string, headerMap map[string]string, body any, id string) (*http.Response, error) {
	if headerMap == nil {
		headerMap = make(map[string]string)
	}

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodPost, url, headerMap, body, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) PostJson(url string, body any, id string) (*http.Response, error) {
	headerMap := make(map[string]string)
	headerMap[httpHeaders.ContentType] = ContentTypeJson

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodPost, url, headerMap, body, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) Put(url string, headerMap map[string]string, body any, id string) (*http.Response, error) {
	if headerMap == nil {
		headerMap = make(map[string]string)
	}

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodPut, url, headerMap, body, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) PutJson(url string, body any, id string) (*http.Response, error) {
	headerMap := make(map[string]string)
	headerMap[httpHeaders.ContentType] = ContentTypeJson

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodPut, url, headerMap, body, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) Delete(url string, headerMap map[string]string, id string) (*http.Response, error) {
	if headerMap == nil {
		headerMap = make(map[string]string)
	}

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodDelete, url, headerMap, nil, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) DeleteJson(url string, id string) (*http.Response, error) {
	headerMap := make(map[string]string)
	headerMap[httpHeaders.ContentType] = ContentTypeJson

	var response *http.Response
	var err error
	for retryCount := 0; retryCount <= 10; retryCount++ {
		response, err = b.do(http.MethodDelete, url, headerMap, nil, retryCount, id)
		if err == nil && http.StatusServiceUnavailable == response.StatusCode {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return response, err
}

func (b *BaseHttpClient) do(method string, url string, headerMap map[string]string, body any, retryCount int, id string) (response *http.Response, err error) {
	httpClientId := uuid.NewString()
	httpClientIdMap.Store(httpClientId, true)

	logMessage := map[string]any{}

	now := time.Now()

	defer func() {
		httpClientIdMap.Delete(httpClientId)

		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, logMessage, now)

			if body != nil {
				bodyByte, _ := json.Marshal(body)
				logMessage[logKey.HttpClientRequestBody] = string(bodyByte)
			}

			if response != nil {
				bodyByte, _ := io.ReadAll(response.Body)
				logMessage[logKey.HttpClientResponseBody] = string(bodyByte)
				logMessage[logKey.HttpClientResponseHeader] = fmt.Sprintf("%+v", response.Header)
				logMessage[logKey.HttpClientResponseContentLength] = response.ContentLength
				logMessage[logKey.HttpClientResponseStatus] = response.Status
				logMessage[logKey.HttpClientResponseStatusCode] = response.StatusCode
			}

			utils.OutputLog(logMessage)
			return
		}
	}()

	b.logBegin(logMessage, now, retryCount, id)

	var request *http.Request
	if body != nil {
		bodyByte, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		request, err = http.NewRequest(method, url, bytes.NewReader(bodyByte))
	} else {
		request, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		panic(err)
	}

	if headerMap != nil {
		for k, v := range headerMap {
			request.Header.Set(k, v)
		}
	}

	if id != "" {
		request.Header.Set(logKey.Id, id)
	}
	logMessage[logKey.HttpClientRequestMethod] = request.Method
	logMessage[logKey.HttpClientRequestHeader] = fmt.Sprintf("%+v", request.Header)
	logMessage[logKey.HttpClientRequestHost] = request.URL.Host
	logMessage[logKey.HttpClientRequestPort] = request.URL.Port()
	logMessage[logKey.HttpClientRequestPath] = request.URL.Path
	logMessage[logKey.HttpClientUsed] = utils.LenSyncMap(httpClientIdMap)
	if request.RequestURI != "" {
		logMessage[logKey.HttpClientRequestUri] = request.RequestURI
	}

	response, err = b.httpClient.Do(request)
	if err != nil {
		if body != nil {
			bodyByte, _ := json.Marshal(body)
			logMessage[logKey.HttpClientRequestBody] = string(bodyByte)
		}
		if response != nil {
			bodyByte, _ := io.ReadAll(response.Body)
			logMessage[logKey.HttpClientResponseBody] = string(bodyByte)
		}
		logMessage[logKey.HttpClientResponseError] = fmt.Sprintf("%+v", err)
	}
	if response != nil {
		logMessage[logKey.HttpClientResponseHeader] = fmt.Sprintf("%+v", response.Header)
		logMessage[logKey.HttpClientResponseContentLength] = response.ContentLength
		logMessage[logKey.HttpClientResponseStatus] = response.Status
		logMessage[logKey.HttpClientResponseStatusCode] = response.StatusCode
	}
	utils.LogEnd(logMessage, now)

	return response, err
}

func (b *BaseHttpClient) logBegin(logMessage map[string]any, now time.Time, retryCount int, id string) {
	logMessage[logKey.Type] = "http-client"
	logMessage[logKey.ServerTime] = now.String()
	if id != "" {
		logMessage[logKey.Id] = id
	} else {
		logMessage[logKey.Id] = uuid.NewString()
	}
	logMessage[logKey.HttpClientRetryCount] = retryCount
	logMessage[logKey.ServiceName] = global.AppName
	logMessage[logKey.HostName] = global.ServerConfig.HostName
}
