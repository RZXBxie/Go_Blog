package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func initServer(addr string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           addr,             // 设置服务器监听的地址
		Handler:        router,           // 设置路由
		ReadTimeout:    10 * time.Minute, // 设置请求的读取超时时间
		WriteTimeout:   10 * time.Minute, // 设置请求的写入超时时间
		MaxHeaderBytes: 1 << 20,          // 设置最大请求头的大小
	}
}
