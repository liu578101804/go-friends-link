package database

// GetAllFriends 获取所有的友链
func GetAllFriends() ([]*Friends, error) {
	friends := make([]*Friends, 0)
	err := D.Find(&friends).Error
	return friends, err
}

type ArticleUniFriendM struct {
	Friends
	Articles
}

func (*ArticleUniFriendM) TableName() string {
	return "articles"
}

// GetAllArticlesUniFriend 获取所有的文章联合朋友
func GetAllArticlesUniFriend(start, count int) ([]*ArticleUniFriendM, error) {
	articles := make([]*ArticleUniFriendM, 0)
	db := D.
		Select("friends.site_title,friends.site_logo,articles.*").
		Joins("left join friends on articles.friend_id = friends.id")
	if start != -1 {
		if count > 10 || count < 1 {
			count = 10
		}
		db = db.Offset(start).Limit(count)
	}
	err := db.Find(&articles).Error
	return articles, err
}
