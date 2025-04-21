package flag

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"os"
	"server/global"
	"server/model/elasticsearch"
	"server/model/other"
	"time"
)

// ElasticSearchExport 导出Es中的数据到json文件
func ElasticSearchExport() error {
	var response other.ESIndexResponse
	
	// 发起第一次搜索请求，设置查询条件：索引名称、滚动时间（1分钟）、每次返回的文档数量、查询条件（匹配所有文档）
	searchRes, err := global.ESClient.Search().
		Index(elasticsearch.ArticleIndex()).
		Scroll("1m").
		Size(1000).
		Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
		Do(context.TODO())
	if err != nil {
		return err
	}
	
	// 遍历第一次查询到的结果
	for _, hit := range searchRes.Hits.Hits {
		// 为每个文档创建一个Data结构体，并将其id和source（文档内容）存储
		data := other.Data{
			ID:  hit.Id_,
			Doc: hit.Source_,
		}
		response.Data = append(response.Data, data)
	}
	
	// 使用Scroll API进行滚动查询，适合处理大批量数据（如导出数据）
	scrollId := *searchRes.ScrollId_
	for {
		scrollRes, err := global.ESClient.Scroll().ScrollId(scrollId).Scroll("1m").Do(context.TODO())
		if err != nil {
			return err
		}
		
		// 如果没有更多数据，结束滚动循环
		if len(scrollRes.Hits.Hits) == 0 {
			break
		}
		
		for _, hit := range scrollRes.Hits.Hits {
			data := other.Data{
				ID:  hit.Id_,
				Doc: hit.Source_,
			}
			response.Data = append(response.Data, data)
		}
		scrollId = *scrollRes.ScrollId_
	}
	
	// 清除滚动查询，释放es上的资源
	_, err = global.ESClient.ClearScroll().ScrollId(scrollId).Do(context.TODO())
	if err != nil {
		return err
	}
	
	fileName := fmt.Sprintf("es_%s.json", time.Now().Format("20060102"))
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	byteData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	
	_, err = file.Write(byteData)
	if err != nil {
		return err
	}
	
	return nil
}
