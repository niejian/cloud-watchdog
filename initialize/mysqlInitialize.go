// initialize doc

package initialize

import (
	"cloud-watchdog/config/parser"
	"cloud-watchdog/global"
	"cloud-watchdog/zapLog"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

//MysqlInitialize doc
//@Description: mysql初始化
//@Author niejian
//@Date 2021-05-08 14:38:17
func MysqlInitialize()  {
	config, err := parser.SysConfigParser()
	if err != nil {
		zapLog.LOGGER().Error("数据库连接失败，请检查配置", zap.String("err", err.Error()))
		panic("数据库连接失败，请检查配置")
		os.Exit(0)
	}
	admin := config.Mysql
	if "" == admin.Username {
		zapLog.LOGGER().Error("请配置数据库账号", zap.String("err", err.Error()))
		os.Exit(0)
	}

	if "" == admin.Password {
		zapLog.LOGGER().Error("请配置数据库密码", zap.String("err", err.Error()))
		os.Exit(0)
	}
	if "" == admin.Dbname {
		zapLog.LOGGER().Error("请配置库名", zap.String("err", err.Error()))
		os.Exit(0)
	}
	mysqlConfig := mysql.Config{
		DSN:                       admin.Username + ":" + admin.Password + "@(" + admin.Path + ")/" + admin.Dbname + "?" + admin.Config, // DSN data source name
		DefaultStringSize:         191,                                                                                                  // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                                 // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                                 // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                                 // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                                // 根据版本自动配置
	}
	var gormConfig *gorm.Config
	gormConfig = &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		zapLog.LOGGER().Error("MySQL启动异常", zap.Any("err", err))
		os.Exit(0)
	} else {
		global.GLOBAL_DB = db
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(admin.MaxIdleConns)
		sqlDB.SetMaxOpenConns(admin.MaxOpenConns)
	}

}
