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
	Title      string     `json:"title"`
	Link       string     `json:"link"`
	Summary    string     `json:"summary"`
	UpdateTime *time.Time `json:"update_time"`
	PushTime   *time.Time `json:"push_time"`
	FriendId   int        `json:"friend_id"`
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

	SubscribeUrl   string `json:"subscribe_url"`
	FeedType       string `json:"feed_type"`
	SubscribeError string `json:"subscribe_error"`

	LastPubTime *time.Time `json:"last_pub_time"`
}

func NewFriends() *Friends {
	return &Friends{}
}

func (*Friends) TableName() string {
	return "friends"
}
