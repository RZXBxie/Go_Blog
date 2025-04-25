package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"server/global"
	"strings"
	"time"
)

// GinLogger 是一个 Gin 中间件，用于记录请求日志。
// 该中间件会在每次请求结束后，使用 Zap 日志记录请求信息。
// 通过此中间件，可以方便地追踪每个请求的情况以及性能。
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		startTime := time.Now()

		// 获取请求的路径和查询参数
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 继续执行其他的中间件
		c.Next()
		cost := time.Since(startTime)

		// 使用zap记录请求日志
		global.Log.Info(path,
			// 记录响应状态码
			zap.Int("status", c.Writer.Status()),
			// 记录请求方法
			zap.String("method", c.Request.Method),
			// 记录请求路径
			zap.String("path", path),
			// 记录查询参数
			zap.String("query", query),
			// 记录客户端ip
			zap.String("ip", c.ClientIP()),
			// 记录user-agent信息
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 是一个 Gin 中间件，用于捕获和处理请求中的 panic 错误。
// 该中间件的主要作用是确保服务在遇到未处理的异常时不会崩溃，并通过日志系统提供详细的错误追踪。
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				var e error
				if tempErr, ok := err.(error); ok {
					e = tempErr
				} else {
					e = fmt.Errorf("%v", err) // 转换为 error 处理
				}

				var netErr *net.OpError
				if errors.As(e, &netErr) {
					var syscallErr *os.SyscallError
					if errors.As(netErr.Err, &syscallErr) {
						errMsg := strings.ToLower(syscallErr.Error())
						if strings.Contains(errMsg, "broken pipe") || strings.Contains(errMsg, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求体信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					global.Log.Error(c.Request.URL.Path,
						zap.Any("error", e),
						zap.String("request", string(httpRequest)),
					)
					_ = c.Error(e)
					c.Abort()
					return
				}

				if stack {
					global.Log.Error("[Recovery from panic]",
						zap.Any("error", e),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					global.Log.Error("[Recovery from panic]",
						zap.Any("error", e),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
