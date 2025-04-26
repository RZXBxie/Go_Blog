package service

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/utils"
)

type JwtService struct {
}

// SetJwtToRedis 将JWT设置到redis中
func (jwtService *JwtService) SetJwtToRedis(jwt string, uuid uuid.UUID) error {
	duration, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	return global.Redis.Set(uuid.String(), jwt, duration).Err()
}

func (jwtService *JwtService) GetJwtFromRedis(uuid uuid.UUID) (string, error) {
	return global.Redis.Get(uuid.String()).Result()
}

func (jwtService *JwtService) InsertIntoBlacklist(jwtBlacklist database.JWTBlacklist) error {
	if err := global.DB.Create(&jwtBlacklist).Error; err != nil {
		return err
	}
	global.BlackCache.SetDefault(jwtBlacklist.Jwt, struct{}{})
	return nil
}

// IsInBlacklist 检查JWT是否在黑名单中
func (jwtService *JwtService) IsInBlacklist(jwt string) bool {
	// 从黑名单缓存中检查JWT是否存在
	_, ok := global.BlackCache.Get(jwt)
	return ok
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
func LoadAll() {
	var data []string
	// 从数据库中获取所有的黑名单JWT
	if err := global.DB.Model(&database.JWTBlacklist{}).Pluck("jwt", &data).Error; err != nil {
		// 如果获取失败，记录错误日志
		global.Log.Error("Failed to load JWT blacklist from the database", zap.Error(err))
		return
	}
	// 将所有JWT添加到BlackCache缓存中
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}
}
