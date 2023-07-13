package modules

// AddFriendRequestM 添加友链请求体
type AddFriendRequestM struct {
	AuthorName   string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`

	WebUrl       string `json:"web_url"`
	WebTitle     string `json:"web_title"`
	SubscribeUrl string `json:"subscribe_url"`
	WebDescribe  string `json:"web_describe"`
}
