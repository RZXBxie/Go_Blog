package flag

import (
	"server/global"
	database2 "server/model/database"
)

// SQL 表结构迁移，如果表不存在，他会创建新表；如果表存在，他会根据结构更新表
func SQL() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&database2.Advertisement{},
		&database2.ArticleLike{},
		&database2.ArticleCategory{},
		&database2.Comment{},
		&database2.Feedback{},
		&database2.FooterLink{},
		&database2.FriendLink{},
		&database2.Image{},
		&database2.JWTBlacklist{},
		&database2.Login{},
		&database2.User{},
	)
}
