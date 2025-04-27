package middleware

import (
	"github.com/gin-gonic/gin"
	"server/model/apptypes"
	"server/model/response"
	"server/utils"
)

// AdminAuth 管理员权限认证
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := utils.GetRoleID(c)
		if roleID != apptypes.Admin {
			response.Forbidden("Access denied. Admin privileges are required", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
