package database

import (
	"gorm.io/gorm"
	"time"
)

type Articles struct {
	gorm.Model
	Title        string
	Link         string
	Summary      string
	Updated      string
	AuthorName   string
	AuthorAvatar string
}

func NewArticles() *Articles {
	return &Articles{}
}

func (*Articles) TableName() string {
	return "articles"
}

type Friends struct {
	gorm.Model
	Url            string
	WebTitle       string
	AuthorName     string
	IsSubscribe    int
	SubscribeType  string
	LastUpdateTime time.Time
}

func NewFriends() *Friends {
	return &Friends{}
}

func (*Friends) TableName() string {
	return "friends"
}
