package elasticsearch

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

// Article 文章表
type Article struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Cover    string   `json:"cover"`   // 封面
	Title    string   `json:"title"`   // 标题
	Keyword  string   `json:"keyword"` // 文章标题-关键字
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Abstract string   `json:"abstract"` // 简介
	Content  string   `json:"content"`
	Views    int      `json:"views"`    // 浏览量
	Comments int      `json:"comments"` // 评论数
	Likes    int      `json:"likes"`    // 收藏量
}

// ArticleIndex 文章es索引
func ArticleIndex() string {
	return "article_index"
}

// ArticleMapping 文章Mapping映射
func ArticleMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"updated_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"cover":      types.TextProperty{},
			"title":      types.TextProperty{},    // 全文搜索/模糊匹配
			"keyword":    types.KeywordProperty{}, // 精确匹配
			"category":   types.KeywordProperty{},
			"tags":       []types.KeywordProperty{},
			"abstract":   types.TextProperty{},
			"content":    types.TextProperty{},
			"views":      types.IntegerNumberProperty{},
			"comments":   types.IntegerNumberProperty{},
			"likes":      types.IntegerNumberProperty{},
		},
	}
}
