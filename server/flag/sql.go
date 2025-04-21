package flag

import (
	"server/database"
	"server/global"
)

// SQL 表结构迁移，如果表不存在，他会创建新表；如果表存在，他会根据结构更新表
func SQL() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&database.Advertisement{},
		&database.ArticleLike{},
		&database.ArticleCategory{},
		&database.Comment{},
		&database.Feedback{},
		&database.FooterLink{},
		&database.FriendLink{},
		&database.Image{},
		&database.JWTBlacklist{},
		&database.Login{},
		&database.User{},
	)
}
