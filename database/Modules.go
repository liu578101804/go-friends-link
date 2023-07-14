package database

import (
	"time"
)

type Model struct {
	ID        int       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Articles struct {
	Model
	Title      string
	Link       string
	Summary    string
	UpdateTime *time.Time
	PushTime   *time.Time
	FriendId   int
}

func NewArticles() *Articles {
	return &Articles{}
}

func (*Articles) TableName() string {
	return "articles"
}

type Friends struct {
	Model
	SiteUrl      string `json:"site_url"`
	SiteTitle    string `json:"site_title"`
	SiteDescribe string `json:"site_describe"`
	SiteLogo     string `json:"site_logo"`

	SubscribeUrl string `json:"subscribe_url"`

	LastUpdateTime *time.Time `json:"last_update_time"`
}

func NewFriends() *Friends {
	return &Friends{}
}

func (*Friends) TableName() string {
	return "friends"
}
