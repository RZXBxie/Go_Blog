package database

import "server/global"

// ArticleLike 文章收藏表
type ArticleLike struct {
	global.Model
	ArticleID uint `json:"article_id"`
	UserID    uint `json:"user_id"`
	User      User `json:"-" gorm:"foreignKey:UserID"`
}
