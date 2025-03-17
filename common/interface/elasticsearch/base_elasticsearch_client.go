package elasticsearchClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/contextKey"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esutil"

	"github.com/elastic/go-elasticsearch/v7/esapi"

	"git-new.okkia.site/crk/decimal-cricket-common/application/utils"
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/logKey"
	"github.com/google/uuid"
)

type IElasticClient interface {
	CreateOrUpdate(parentContext context.Context, index string, documentId string, document any) *esapi.Response
	BulkCreateOrUpdate(parentContext context.Context, index string, elasticDocuments []ElasticDocument)
	Search(parentContext context.Context, indexArr []string, query string) *esapi.Response
	SearchPage(parentContext context.Context, indexArr []string, query string, from *int, size *int) *esapi.Response
	SearchPageAndSort(parentContext context.Context, indexArr []string, query string, from *int, size *int, sort []string) *esapi.Response
	Count(parentContext context.Context, indexArr []string, query string) *esapi.Response
}

type BaseElasticClient struct {
	client *elasticsearch.Client
}

type ElasticDocument struct {
	DocumentId string
	Document   any
}

func NewBaseElasticClient(client *elasticsearch.Client) *BaseElasticClient {
	return &BaseElasticClient{client: client}
}

func (b *BaseElasticClient) CreateOrUpdate(parentContext context.Context, index string, documentId string, document any) *esapi.Response {
	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.logBegin([]string{index}, document, actionLogs, now, utils.GetActionId(c))

	data, err := json.Marshal(document)
	if err != nil {
		panic(err)
	}

	request := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       bytes.NewReader(data),
	}

	response, err := request.Do(c, b.client)
	if err != nil {
		panic(err)
	}
	if response.IsError() {
		actionLogs[logKey.ElasticResponseBody] = fmt.Sprintf("%+v", response)
	}

	b.logEnd(response, actionLogs, now)

	return response
}

func (b *BaseElasticClient) BulkCreateOrUpdate(parentContext context.Context, index string, elasticDocuments []ElasticDocument) {
	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.logBegin([]string{index}, elasticDocuments, actionLogs, now, utils.GetActionId(c))

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  index,
		Client: b.client,
		OnError: func(ctx context.Context, err error) {
			log.Printf("BulkIndexer error: %+v, index: %s, actionId: %s", err, index, utils.GetActionId(c))
		},
	})
	if err != nil {
		panic(err)
	}

	elasticDocumentLen := len(elasticDocuments)
	actionLogs[logKey.ElasticBulkCount] = elasticDocumentLen
	var countSuccessful uint64
	for index := 0; index < elasticDocumentLen; index++ {
		elasticDocument := elasticDocuments[index]

		data, err := json.Marshal(elasticDocument.Document)
		if err != nil {
			actionLogs[logKey.ElasticBulkMarshalFail+elasticDocument.DocumentId] = true
			continue
		}

		err = bulkIndexer.Add(c, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: elasticDocument.DocumentId,
			Body:       bytes.NewReader(data),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&countSuccessful, 1)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				log.Printf("BulkIndexer add onFailure Item:%+v, res:%+v, error:%s, index: %d, actionId: %s", item, res, err, index, utils.GetActionId(c))
			},
		})
		if err != nil {
			actionLogs[logKey.ElasticBulkAddFail+elasticDocument.DocumentId] = true
		}
	}

	if err = bulkIndexer.Close(c); err != nil {
		panic(err)
	}

	stats := bulkIndexer.Stats()
	actionLogs[logKey.ElasticBulkAddedNum] = stats.NumAdded
	actionLogs[logKey.ElasticBulkIndexedNum] = stats.NumIndexed
	actionLogs[logKey.ElasticBulkFailedNum] = stats.NumFailed

	if countSuccessful != uint64(elasticDocumentLen) {
		actionLogs[logKey.ElasticBulkResult] = "fail"
	} else {
		actionLogs[logKey.ElasticBulkResult] = "success"
	}

	b.logEnd(nil, actionLogs, now)
}

func (b *BaseElasticClient) Search(parentContext context.Context, indexArr []string, query string) *esapi.Response {
	from := new(int)
	*from = 0
	size := new(int)
	*size = 10000
	return b.SearchPage(parentContext, indexArr, query, from, size)
}

func (b *BaseElasticClient) SearchPage(parentContext context.Context, indexArr []string, query string, from *int, size *int) *esapi.Response {
	return b.SearchPageAndSort(parentContext, indexArr, query, from, size, []string{})
}

func (b *BaseElasticClient) SearchPageAndSort(parentContext context.Context, indexArr []string, query string, from *int, size *int, sort []string) *esapi.Response {
	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	if from == nil {
		from = new(int)
		*from = 0
	}

	if size == nil {
		size = new(int)
		*size = 10
	}

	if *size > 10000 {
		*size = 10000
	}

	if sort == nil {
		sort = []string{}
	}

	b.logBegin(indexArr, query, actionLogs, now, utils.GetActionId(c))

	request := esapi.SearchRequest{
		Index: indexArr,
		Body:  strings.NewReader(query),
		From:  from,
		Size:  size,
		Sort:  sort,
	}

	response, err := request.Do(c, b.client)
	if err != nil {
		panic(err)
	}
	if response.IsError() {
		actionLogs[logKey.ElasticResponseBody] = fmt.Sprintf("%+v", response)
	}

	b.logEnd(response, actionLogs, now)

	return response
}

func (b *BaseElasticClient) Count(parentContext context.Context, indexArr []string, query string) *esapi.Response {
	c := context.WithValue(parentContext, contextKey.ActionLogs, make(map[string]any))
	actionLogs, _ := utils.GetActionLogs(c)

	now := time.Now()

	defer func() {
		r := recover()
		if r != nil {
			utils.HandleErrorRecover(r, actionLogs, now)
			utils.OutputLog(actionLogs)
			return
		}
	}()

	b.logBegin(indexArr, query, actionLogs, now, utils.GetActionId(c))

	request := esapi.CountRequest{
		Index: indexArr,
		Body:  strings.NewReader(query),
	}

	response, err := request.Do(c, b.client)
	if err != nil {
		panic(err)
	}
	if response.IsError() {
		actionLogs[logKey.ElasticResponseBody] = fmt.Sprintf("%+v", response)
	}

	b.logEnd(response, actionLogs, now)

	return response
}

func (b *BaseElasticClient) logBegin(indexArr []string, document any, actionLogs map[string]any, now time.Time, actionId string) {
	actionLogs[logKey.Type] = "elastic-client"
	actionLogs[logKey.ServerTime] = now.String()
	if actionId == "" {
		actionId = uuid.NewString()
	}
	actionLogs[logKey.Id] = actionId
	actionLogs[logKey.ServiceName] = global.AppName
	actionLogs[logKey.HostName] = global.ServerConfig.HostName
	actionLogs[logKey.ElasticIndexName] = fmt.Sprintf("%+v", indexArr)
	documentStr := fmt.Sprintf("%+v", document)
	actionLogs[logKey.ElasticRequestBody] = strings.ReplaceAll(strings.ReplaceAll(documentStr, "\n", ""), "\t", "")
}

func (b *BaseElasticClient) logEnd(response *esapi.Response, actionLogs map[string]any, now time.Time) {
	if response != nil {
		actionLogs[logKey.ElasticResponseStatus] = fmt.Sprintf("%+v", response.Status())
	}
	timeUsed := time.Since(now)
	actionLogs[logKey.TimeUsed] = timeUsed.String()
	actionLogs[logKey.TimeUsedNano] = timeUsed
	utils.OutputLog(actionLogs)
}
