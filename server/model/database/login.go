package database

import "server/global"

// Login 登录日志表
type Login struct {
	global.Model
	UserID      uint   `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
	LoginMethod string `json:"login_method"`
	IP          string `json:"ip"`
	Address     string `json:"address"`
	OS          string `json:"os"`
	DeviceInfo  string `json:"device_info"`
	BrowserInfo string `json:"browser_info"`
	Status      int    `json:"status"`
}
