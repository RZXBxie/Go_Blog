package utils

import (
	"server/global"
	"server/model/other"
)

// MySQLPagination MySQL分页查询，参数T是要查询的数据库表结构体，option里面有查询参数
func MySQLPagination[T any](model *T, option other.MySQLOption) (list []T, total int64, err error) {
	if option.Page < 1 {
		option.Page = 1
	}
	if option.PageSize < 1 {
		option.PageSize = 10
	}
	if option.Order == "" {
		option.Order = "id desc"
	}

	query := global.DB.Model(model)

	if option.Where != nil {
		query = query.Where(option.Where)
	}

	// SELECT count(*) FROM 表 WHERE xxx条件
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 预加载关联模型.如查询文章的时候，把作者一并查出来
	// SELECT * FROM user LEFT JOIN role ON user.role_id = role.id
	for _, preload := range option.Preload {
		query = query.Preload(preload)
	}

	err = query.Order(option.Order).
		Limit(option.PageSize).
		Offset(option.PageSize * (option.Page - 1)).
		Find(&list).Error

	return list, total, err
}
