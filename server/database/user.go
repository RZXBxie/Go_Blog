package database

import (
	"github.com/gofrs/uuid"
	"server/global"
	"server/model/apptypes"
)

// User 用户表
type User struct {
	global.Model
	UUID      uuid.UUID         `json:"uuid" gorm:"type:char(36);unique;not null"`
	Username  string            `json:"username"`
	Password  string            `json:"-"`
	Email     string            `json:"email"`
	OpenID    string            `json:"openid"`
	Avatar    string            `json:"avatar" gorm:"size:255"` // 头像
	Address   string            `json:"address"`
	Signature string            `json:"signature" gorm:"default:'这位用户有点低调，未设置任何签名'"`
	RoleID    apptypes.RoleID   `json:"role_id"`
	Register  apptypes.Register `json:"register"`
	Freeze    bool              `json:"freeze"` //账户是否被冻结
}
