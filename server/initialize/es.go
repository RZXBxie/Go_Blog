package initialize

import (
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"os"
	"server/global"
)

// ConnectES 初始化并返回好一个配置好的ElasticSearch客户端
func ConnectES() *elasticsearch.TypedClient {
	esConfig := global.Config.ES
	cfg := elasticsearch.Config{
		Addresses: []string{esConfig.URL},
		Username:  esConfig.Username,
		Password:  esConfig.Password,
	}
	// 如果配置中指定了需要打印日志到控制台，则启用日志打印
	if esConfig.IsConsolePrint {
		cfg.Logger = &elastictransport.ColorLogger{
			Output:             os.Stderr, // 设置日志输出到标准输出
			EnableRequestBody:  true,      // 启用请求体打印
			EnableResponseBody: true,      // 启用响应体打印
		}
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		global.Log.Error("Failed to connect to elasticsearch", zap.Error(err))
		os.Exit(1)
	}

	return client
}
