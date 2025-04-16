package database

import (
	"server/global"
)

// Advertisement 广告表
type Advertisement struct {
	global.Model
	AdImage string `json:"ad_image" gorm:"size:255"` // 图片
	Image   Image  `json:"-" gorm:"foreignKey:AdImage;references:Path"`
	Link    string `json:"link"` // 链接
	Title   string `json:"title"`
	Content string `json:"content"`
}
