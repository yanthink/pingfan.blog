package bootstrap

import (
	"blog/app"
	"blog/app/models"
	"blog/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
	"time"
)

func SetupDatabase() *gorm.DB {
	app.DB = ConnectDB()

	sqlDB, _ := app.DB.DB()
	// 设置最大连接数
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnections)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)
	// 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeSeconds) * time.Second)

	if config.Database.AutoMigrate {
		migration(app.DB)
	}

	return app.DB
}

func ConnectDB() *gorm.DB {
	var logLevel logger.LogLevel

	if config.App.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	// 连接 MySQL 数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Pass,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zapgorm2.New(app.Logger).LogMode(logLevel),
		// AutoMigrate 默认会自动创建数据库外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能
		SkipDefaultTransaction: true,
		// 禁用嵌套事务，true 嵌套中的任何一个事务发生回滚整个事务都会一起回滚
		DisableNestedTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	if config.App.Env == "local" && config.Database.AutoMigrate {
		migration(db)
	}

	return db
}

func migration(db *gorm.DB) {
	_ = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Tag{},
		&models.Like{},
		&models.Favorite{},
		&models.Comment{},
		&models.Upvote{},
		&models.Notification{},
	)
}
