package database

import (
	"server/global"
	"server/model/apptypes"
)

// Image 图片表
type Image struct {
	global.Model
	Name     string            `json:"name"`                        // 名字
	Path     string            `json:"path" gorm:"size:255;unique"` // 路径
	Category apptypes.Category `json:"category"`                    // 图片类别
	Storage  apptypes.Storage  `json:"storage"`                     // 图片存储方式
}
