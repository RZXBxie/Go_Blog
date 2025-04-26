package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"
	"strconv"
)

var jwtService = service.ServiceGroupApp.JwtService

// JWTAuth 是一个中间件函数，验证请求中的JWT token是否合法
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := utils.GetAccessToken(c)
		refreshToken := utils.GetRefreshToken(c)

		if jwtService.IsInBlacklist(refreshToken) {
			utils.ClearRefreshToken(c)
			response.NoAuth("Account logged in from another location or token is invalid", c)
			c.Abort()
			return
		}

		j := utils.NewJWT()

		claims, err := j.ParseAccessToken(accessToken)

		if err != nil {
			// accessToken为空或者已过期
			if accessToken == "" || errors.Is(err, utils.TokenExpired) {

				// 尝试解析refreshToken
				refreshClaims, err := j.ParseRefreshToken(accessToken)
				if err != nil {
					utils.ClearRefreshToken(c)
					response.NoAuth("Refresh Token expired or invalid", c)
					c.Abort()
					return
				}

				var user database.User
				if err := global.DB.Select("uuid", "role_id").Take(&user, refreshClaims.UserID).Error; err != nil {
					utils.ClearRefreshToken(c)
					response.NoAuth("User not exist", c)
					c.Abort()
					return
				}

				newAccessClaims := j.CreateAccessClaims(
					request.BaseClaims{
						UserID: refreshClaims.UserID,
						RoleID: user.RoleID,
						UUID:   user.UUID,
					})
				newAccessToken, err := j.CreateAccessToken(newAccessClaims)
				if err != nil {
					utils.ClearRefreshToken(c)
					response.NoAuth("Create access token fail", c)
					c.Abort()
					return
				}

				// 将新的accessToken和过期时间添加到响应头中

				c.Header("new-access-token", newAccessToken)
				c.Header("new-access-expires-at", strconv.FormatInt(newAccessClaims.ExpiresAt.Unix(), 10))

				c.Set("claims", newAccessClaims)
				c.Next()
				return
			}
			utils.ClearRefreshToken(c)
			response.NoAuth("Invalid access token", c)
			c.Abort()
			return
		}

		// accessToken合法，就将其存到Context中
		c.Set("claims", claims)
		c.Next()
	}
}
