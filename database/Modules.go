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
	UpdateTime   *time.Time
	PushTime     *time.Time
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

	WebUrl      string `json:"web_url"`
	WebTitle    string `json:"web_title"`
	WebDescribe string `json:"web_describe"`

	AuthorName   string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`

	SubscribeUrl string `json:"subscribe_url"`

	LastUpdateTime *time.Time `json:"last_update_time"`
}

func NewFriends() *Friends {
	return &Friends{}
}

func (*Friends) TableName() string {
	return "friends"
}
