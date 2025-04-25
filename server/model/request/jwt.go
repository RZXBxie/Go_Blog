package request

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"server/model/apptypes"
)

type BaseClaims struct {
	UserID uint
	UUID   uuid.UUID       // 唯一标识用户
	RoleID apptypes.RoleID // 标识用户权限级别
}

// JwtCustomClaims 结构体用于存储JWT的自定义Claims，继承自BaseClaims，并包含标准的JWT注册信息
type JwtCustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

// JwtCustomRefreshClaims 结构体用于存储刷新Token的自定义Claims，包含用户ID和标准的JWT注册信息
type JwtCustomRefreshClaims struct {
	UserID uint
	jwt.RegisteredClaims
}
