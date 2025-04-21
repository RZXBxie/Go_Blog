package flag

import (
	"os"
	"server/global"
	"strings"
)

// SQLImport 导入SQL数据
func SQLImport(sqlPath string) (errs []error) {
	byteData, err := os.ReadFile(sqlPath)
	if err != nil {
		errs = append(errs, err)
	}
	sqlList := strings.Split(string(byteData), ";")
	for _, sql := range sqlList {
		// 去除字符串开头和结尾的空白符
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		err = global.DB.Exec(sql).Error
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
