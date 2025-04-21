package database

import "server/global"

// FriendLink 友链表
type FriendLink struct {
	global.Model
	Logo        string `json:"logo" gorm:"size:255"`
	Image       Image  `json:"-" gorm:"foreignKey:Logo;references:Path"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Description string `json:"description"` // 描述
}
