package flag

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"os"
	"server/elasticsearch"
	"server/global"
	"server/other"
	"server/service"
)

// ElasticSearchImport 从指定的json文件中导入数据到es中
func ElasticSearchImport(jsonPath string) (int, error) {
	byteData, err := os.ReadFile(jsonPath)
	if err != nil {
		return 0, err
	}

	var response other.ESIndexResponse
	err = json.Unmarshal(byteData, &response)

	// 创建ES索引
	esService := service.ServiceGroupApp.EsService
	esExists, err := esService.IndexExist(elasticsearch.ArticleIndex())
	if err != nil {
		return 0, err
	}
	if esExists {
		if err = esService.IndexDelete(elasticsearch.ArticleIndex()); err != nil {
			return 0, err
		}
	}
	err = esService.IndexCreate(elasticsearch.ArticleIndex(), elasticsearch.ArticleMapping())
	if err != nil {
		return 0, err
	}

	// 构建批量请求数据
	var request bulk.Request
	for _, data := range response.Data {
		// 为每条数据创建索引操作，创建文档的ID
		request = append(request, types.OperationContainer{Index: &types.IndexOperation{Id_: data.ID}})
		request = append(request, data.Doc)
	}

	// 使用es客户端执行批量操作
	_, err = global.ESClient.Bulk().
		Request(&request).                   // 提交请求数据
		Index(elasticsearch.ArticleIndex()). // 指定索引名称
		Refresh(refresh.True).               // 强制刷新索引以使文档立即可见
		Do(context.TODO())                   // 执行请求
	if err != nil {
		return 0, err
	}
	total := len(response.Data)
	return total, nil
}
