package database

import (
	"github.com/gofrs/uuid"
	"server/global"
)

// Comment 评论表
type Comment struct {
	global.Model
	ArticleID uint      `json:"article_id"`                                                 // 文章id
	PID       *uint     `json:"p_id"`                                                       // 父评论id
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID;references:ID"`                      // 父评论，如果是一级评论则为nil
	Children  []Comment `json:"children" gorm:"foreignKey:PID;constraint:OnDelete:CASCADE"` // 子评论数组
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`
	User      User      `json:"-" gorm:"foreignKey:UserUUID;references:UUID"`
	Content   string    `json:"content"`
}

// TODO 创建和删除评论时需要更新评论数
