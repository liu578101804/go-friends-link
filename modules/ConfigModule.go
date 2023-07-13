package modules

type ConfigModule struct {
	Token            string `json:"token"`
	Port             int    `json:"port"`
	Cron             string `json:"cron"`
	LastCrawlingTime string `json:"last_crawling_time"`
}
