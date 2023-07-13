package modules

type ResponseArticleModule struct {
	Floor   int    `json:"floor"`
	Title   string `json:"title"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Link    string `json:"link"`
	Author  string `json:"author"`
	Avatar  string `json:"avatar"`
}

type ResponseStatisticalDataModule struct {
	FriendsNum      int    `json:"friends_num"`
	ActiveNum       int    `json:"active_num"`
	ErrorNum        int    `json:"error_num"`
	ArticleNum      int    `json:"article_num"`
	LastUpdatedTime string `json:"last_updated_time"`
}

// ResponseModule 响应体
type ResponseModule struct {
	StatisticalData *ResponseStatisticalDataModule `json:"statistical_data"`
	ArticleData     []*ResponseArticleModule       `json:"article_data"`
	Page            int                            `json:"page"`
}
