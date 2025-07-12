package esclient

import (
	"context"
	"encoding/json"
	testutil2 "example.com/m/test/testutil"
	"example.com/m/util/copier"
	"example.com/m/util/printhelper"
	"example.com/m/util/stringutil"
	"example.com/m/util/timeutil"
	"example.com/m/util/utilerror"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"time"
)

type EsClient struct {
	esClient *elastic.Client
}

func NewEsClient() *EsClient {
	ctx := context.Background()
	host := testutil2.GetStringConfig(ctx, testutil2.ModuleEs, testutil2.ConfigESHost)
	esClient, err := elastic.DialContext(ctx,
		elastic.SetSniff(false),
		elastic.SetURL([]string{host}...))
	if err != nil {
		panic(err)
	}
	return &EsClient{
		esClient: esClient,
	}
}

func GetQueryEsTimeoutMilliSeconds(ctx context.Context) int64 {
	return 5000
}

func (m *EsClient) EsAdd(ctx context.Context, index string, id string, body interface{}) *utilerror.UtilError {
	printhelper.Printf("es add %s : %s", index, stringutil.Object2String(body))
	start := timeutil.GetCurrentTimestamp()
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(GetQueryEsTimeoutMilliSeconds(ctx))*time.Millisecond)
	defer cancel()
	resp, err := m.esClient.Index().
		Index(index).
		Id(id).
		Refresh("true").
		BodyJson(body).
		Do(requestCtx)
	printhelper.Printf("EsAdd cost %d ms", timeutil.GetCurrentTimestamp()-start)
	if err != nil {
		return utilerror.NewError(err.Error())
	}
	if resp == nil {
		return utilerror.NewError("empty rsp from ES")
	}
	if resp.Id != id {
		return utilerror.NewError("empty rsp from ES")
	}
	return nil
}

func (m *EsClient) EsBatchCreate(ctx context.Context, data interface{}, f func(list interface{}) []*elastic.BulkIndexRequest) *utilerror.UtilError {
	start := timeutil.GetCurrentTimestamp()
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(GetQueryEsTimeoutMilliSeconds(ctx))*time.Millisecond)
	defer cancel()
	bulkReq := m.esClient.Bulk()
	docs := f(data)
	for _, doc := range docs {
		bulkReq.Add(doc)
	}
	printhelper.Printf("es batchCreate %s", stringutil.Object2String(bulkReq))
	if len(docs) == 0 {
		return nil
	}
	_, esErr := bulkReq.Do(requestCtx)
	printhelper.Printf("EsBatchCreate cost %d ms", timeutil.GetCurrentTimestamp()-start)
	if esErr != nil {
		return utilerror.NewError(esErr.Error())
	}
	return nil
}

func (m *EsClient) EsBatchGet(ctx context.Context, index string, param string, ids []string, resp interface{}) error {
	start := timeutil.GetCurrentTimestamp()
	if len(ids) == 0 {
		return nil
	}
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(GetQueryEsTimeoutMilliSeconds(ctx))*time.Millisecond)
	defer cancel()
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewTermsQuery(param, ToSliceAny(ids)...))
	src, _ := boolQ.Source()
	data, _ := json.Marshal(src)
	msg := fmt.Sprintf("es query idx %s  cond: %v", index, string(data))
	printhelper.Printf(msg)
	fmt.Println(msg)
	searchResult, esErr := m.esClient.
		Search(index).
		Query(boolQ).
		TrackTotalHits(true).
		Size(500).Do(requestCtx)
	printhelper.Printf("EsBatchGet cost %d ms", timeutil.GetCurrentTimestamp()-start)
	if esErr != nil {
		return esErr
	}
	var EsRespList []interface{}
	if searchResult != nil && searchResult.Hits != nil && len(searchResult.Hits.Hits) > 0 {
		for _, hit := range searchResult.Hits.Hits {
			EsRespList = append(EsRespList, hit.Source)
		}
		cErr := copier.JsonCopy(EsRespList, resp)
		if cErr != nil {
			return cErr
		}
	}
	return nil
}

func ToSliceAny(data interface{}) []interface{} {
	d := reflect.ValueOf(data)
	if d.Kind() != reflect.Slice {
		return nil
	}

	result := make([]interface{}, d.Len())
	for i := 0; i < d.Len(); i++ {
		result[i] = d.Index(i).Interface()
	}

	return result
}
