package modules

// AddFriendRequestM 添加友链请求体
type AddFriendRequestM struct {
	SubscribeUrl string `json:"subscribe_url"`
	SiteLogo     string `json:"site_logo"`
}
