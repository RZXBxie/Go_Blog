package database

// ArticleTag 文章标签表
type ArticleTag struct {
	Tag    string `json:"tag" gorm:"primaryKey"`
	Number int    `json:"number"`
}
