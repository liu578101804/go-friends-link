package modules

type LinkItem struct {
	Url    string `json:"url"`
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Avatar string `json:"avatar"`
}

type ConfigModule struct {
	Links []*LinkItem `json:"links"`
	Token string      `json:"token"`
	Port  int         `json:"port"`
	Cron  string      `json:"cron"`
}
