package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"
)

type UserApi struct {
}

// Register 注册
func (userApi *UserApi) Register(ctx *gin.Context) {
	var req request.Register
	// 请求体自动解析 JSON 到结构体
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	session := sessions.Default(ctx)
	savedEmail := session.Get("email")
	// 两次邮箱一致性判断
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", ctx)
		return
	}

	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", ctx)
		return
	}

	savedTime := session.Get("expiration_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", ctx)
		return
	}

	user, err := userService.Register(req)

	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", ctx)
		return
	}

	userApi.tokenNext(ctx, user)
}

// Login 登录接口，根据不同的登录方式调用不同的登录方法
func (userApi *UserApi) Login(ctx *gin.Context) {
	switch ctx.Query("flag") {
	case "email":
		userApi.EmailLogin(ctx)
	case "qq":
		userApi.QQLogin(ctx)
	default:
		userApi.EmailLogin(ctx)
	}
}

// EmailLogin 邮箱登录
func (userApi *UserApi) EmailLogin(ctx *gin.Context) {
	var req request.Login

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	if store.Verify(req.CaptchaID, req.Captcha, true) {
		user, err := userService.EmailLogin(req)
		if err != nil {
			global.Log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage("Failed to login", ctx)
			return
		}
		userApi.tokenNext(ctx, user)
		return
	}

	response.FailWithMessage("Incorrect verification code", ctx)
}

// QQLogin TODO 接入QQ登录的功能
func (userApi *UserApi) QQLogin(c *gin.Context) {}

// tokenNext 为用户生成双token
func (userApi *UserApi) tokenNext(ctx *gin.Context, user database.User) {
	// 检查账户是否被冻结
	if user.Freeze {
		response.FailWithMessage("the user is frozen, please contact the administrator", ctx)
		return
	}

	baseClaims := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()

	// 创建accessToken
	accessClaims := j.CreateAccessClaims(baseClaims)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to create access token:", zap.Error(err))
		response.FailWithMessage("Failed to create access token", ctx)
		return
	}

	// 创建refreshToken
	refreshClaims := j.CreateRefreshClaims(baseClaims)
	refreshToken, err := j.CreateRefreshToken(refreshClaims)
	if err != nil {
		global.Log.Error("Failed to create refresh token:", zap.Error(err))
		response.FailWithMessage("Failed to create refresh token", ctx)
		return
	}

	// 是否开启了多地点登录拦截.
	// 如果开启了拦截，那么老的token立即失效，旧的登录会被强制下线；否则只会设置新的token，旧的登录不会下线
	if !global.Config.System.UseMultipoint {
		utils.SetRefreshToken(ctx, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		// 登录日志中间件需要记录用户登录信息，那里需要用到user_id，所以此处要存入user_id
		ctx.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:        user,
			AccessToken: accessToken,
			// 将秒转为毫秒
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "successfully login", ctx)
		return
	}
	// 检查redis中是否已存在该用户的JWT
	if jwtStr, err := jwtService.GetJwtFromRedis(user.UUID); errors.Is(err, redis.Nil) {
		// 不存在就设置新的
		if err := jwtService.SetJwtToRedis(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", ctx)
			return
		}
		utils.SetRefreshToken(ctx, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		ctx.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "successfully login", ctx)
	} else if err != nil {
		global.Log.Error("Failed to set login status:", zap.Error(err))
		response.FailWithMessage("Failed to set login status", ctx)
	} else {
		// 将旧的JWT加入黑名单，并将新的JWT设置到redis
		blackList := database.JWTBlacklist{
			Jwt: jwtStr,
		}
		if err := jwtService.InsertIntoBlacklist(blackList); err != nil {
			global.Log.Error("Failed to invalidate jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt", ctx)
			return
		}

		// 设置新的jwt到redis
		if err := jwtService.SetJwtToRedis(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", ctx)
			return
		}

		utils.SetRefreshToken(ctx, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		ctx.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "successfully login", ctx)
	}
}

// ForgotPassword 找回密码
func (userApi *UserApi) ForgotPassword(ctx *gin.Context) {
	var req request.ForgotPassword
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	session := sessions.Default(ctx)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", ctx)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", ctx)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expiration_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", ctx)
		return
	}

	err = userService.ForgotPassword(req)
	if err != nil {
		global.Log.Error("Failed to retrieve the password", zap.Error(err))
		response.FailWithMessage("Failed to retrieve the password", ctx)
		return
	}
	response.OkWithMessage("Successfully retrieved the password", ctx)
}

// UserCard 获取用户卡片信息
func (userApi *UserApi) UserCard(ctx *gin.Context) {
	var req request.UserCard
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to get user card", zap.Error(err))
		response.FailWithMessage("Failed to get user card", ctx)
		return
	}
	response.OkWithData(userCard, ctx)
}

// Logout 登出
func (userApi *UserApi) Logout(ctx *gin.Context) {
	userService.Logout(ctx)
	response.OkWithMessage("Successfully logout", ctx)
}

// UserResetPassword 重置密码
func (userApi *UserApi) UserResetPassword(ctx *gin.Context) {
	var req request.UserResetPassword
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	req.UserID = utils.GetUserID(ctx)
	err = userService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("Failed to reset user password", zap.Error(err))
		response.FailWithMessage("Failed to reset user password", ctx)
		return
	}
	response.OkWithMessage("Successfully reset user password, please login again", ctx)
}

// UserInfo 获取用户信息
func (userApi *UserApi) UserInfo(ctx *gin.Context) {
	userID := utils.GetUserID(ctx)
	user, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user information", zap.Error(err))
		response.FailWithMessage("Failed to get user information", ctx)
		return
	}
	response.OkWithData(user, ctx)
}

// UserChangeInfo 更改用户信息
func (userApi *UserApi) UserChangeInfo(ctx *gin.Context) {
	var req request.UserChangeInfo
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	req.UserID = utils.GetUserID(ctx)
	err = userService.UserChangeInfo(req)
	if err != nil {
		global.Log.Error("Failed to change user info", zap.Error(err))
		response.FailWithMessage("Failed to change user info", ctx)
		return
	}
	response.OkWithMessage("Successfully changed user info", ctx)
}

// UserWeather 获取用户所在地区的天气信息
func (userApi *UserApi) UserWeather(ctx *gin.Context) {
	ip := ctx.ClientIP()
	weather, err := userService.UserWeather(ip)
	if err != nil {
		global.Log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", ctx)
		return
	}
	response.OkWithData(weather, ctx)
}

// UserChart 获取用户图表数据，登录和注册人数
func (userApi *UserApi) UserChart(ctx *gin.Context) {
	var req request.UserChart
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	resp, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart", zap.Error(err))
		response.FailWithMessage("Failed to get user chart", ctx)
		return
	}
	response.OkWithData(resp, ctx)
}

func (userApi *UserApi) UserList(ctx *gin.Context) {
	var req request.UserList
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	list, total, err := userService.UserList(req)
	if err != nil {
		global.Log.Error("Failed to get user list", zap.Error(err))
		response.FailWithMessage("Failed to get user list", ctx)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, ctx)
}
