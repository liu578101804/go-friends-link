package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var D *gorm.DB

const DbPath = "./data/database.db"

func Init() {
	var err error
	D, err = gorm.Open(sqlite.Open(DbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动同步
	D.AutoMigrate(
		NewArticles(),
		NewFriends(),
	)
}
