package middleware

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/service"
)

func LoginRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		C.Next()
		// 异步记录日志
		go func() {
			gaodeService := service.ServiceGroupApp.GaodeService
			var userID uint
			var address string

			ip := c.ClientIP()
			loginMethod := c.DefaultQuery("flag", "email") //若为传递flag参数，则默认为email
			userAgent := c.Request.UserAgent()

			if value, ok := c.Get("user_id"); ok {
				if id, ok := value.(uint); ok {
					userID = id
				}
			}

			// 获取用户ip的地理位置
			address = getAddressFromIP(ip, gaodeService)

			// 解析用户的浏览器、操作系统和设备信息
			os, deviceInfo, browserInfo := parseUserAgent(userAgent)

			login := database.Login{
				UserID:      userID,
				LoginMethod: loginMethod,
				IP:          ip,
				Address:     address,
				OS:          os,
				DeviceInfo:  deviceInfo,
				BrowserInfo: browserInfo,
				Status:      c.Writer.Status(),
			}

			// 将登录信息存储到数据库
			if err := global.DB.Create(&login).Error; err != nil {
				global.Log.Error("Failed to record login", zap.Error(err))
			}
		}()
	}
}

// getAddressFromIP 获取IP地址对应的地理位置信息（省市信息）
func getAddressFromIP(ip string, gaodeService service.GaodeService) string {
	resp, err := gaodeService.GetLocationByIP(ip)
	if err != nil || resp.Province == "" {
		return "未知"
	}
	if resp.City != "" && resp.Province != resp.City {
		return resp.Province + "-" + resp.City
	}
	return resp.Province
}

// parseUserAgent 解析用户代理（User-Agent）字符串，提取操作系统、设备信息和浏览器信息
func parseUserAgent(userAgent string) (os, deviceInfo, browserInfo string) {
	os = userAgent
	deviceInfo = userAgent
	browserInfo = userAgent

	parser := uaparser.NewFromSaved()
	cli := parser.Parse(userAgent)
	os = cli.Os.Family
	deviceInfo = cli.Device.Family
	browserInfo = cli.UserAgent.Family

	return
}
