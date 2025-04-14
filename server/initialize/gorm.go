package initialize

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"server/global"
)

// InitGorm 通过gorm连接到mysql
func InitGorm() *gorm.DB {
	mysqlConfig := global.Config.Mysql
	db, err := gorm.Open(mysql.Open(mysqlConfig.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlConfig.LogLevel()),
	})
	if err != nil {
		global.Log.Error("Failed to connect to mysql:", zap.Error(err))
		os.Exit(1)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxOpenConns)
	return db
}
