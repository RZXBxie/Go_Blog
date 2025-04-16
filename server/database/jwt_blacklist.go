package database

import "server/global"

// JWTBlacklist JWT黑名单表
type JWTBlacklist struct {
	global.Model
	Jwt string `json:"jwt" gorm:"type:text"`
}
