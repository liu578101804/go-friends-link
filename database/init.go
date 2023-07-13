package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var D *gorm.DB

func init() {
	var err error
	D, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动同步
	D.AutoMigrate(
		NewArticles(),
		NewFriends(),
	)
}
