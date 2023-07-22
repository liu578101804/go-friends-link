package crawling

import (
	"friends-rss/config"
	"friends-rss/database"
	"github.com/mmcdole/gofeed"
	"log"
	"time"
)

// Crawling 去获取消息
func Crawling() {

	// 获取全部友链
	friends, err := database.GetAllFriends()
	if err != nil {
		log.Println("查询获取友链失败", err)
		return
	}
	// 判断是否没有友链
	if len(friends) == 0 {
		log.Println("没有友链存在，本次结束!")
		return
	}
	// 遍历所有朋友
	for i := 0; i < len(friends); i++ {
		friend := friends[i]
		// 打印日志
		log.Println("准备获取：", friend.SubscribeUrl, friend.SubscribeUrl)
		// 判断是否开启订阅
		if friend.SubscribeUrl == "" {
			log.Println("本次获取的站点信息：", friend.SubscribeUrl, friend.SiteTitle)
			log.Println("此朋友没开启订阅")
			break
		}

		// 这个朋友的最后的更新时间
		if friend.LastPubTime == nil || friend.LastPubTime.IsZero() {
			nowTime := time.Now().Add(-time.Hour * 7200) //往前推300天
			friend.LastPubTime = &nowTime
			log.Println("此站点最后更新时间为空，给他默认的时间：", nowTime)
		}
		// 缓存站点的最后更新时间
		siteLastPubTime := friend.LastPubTime
		log.Println("此站点最后更新时间为：", siteLastPubTime)

		// 使用三方库解析
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(friend.SubscribeUrl)
		if err != nil {
			log.Println("获取订阅异常：", err)
			// 写入最新的更新时间
			if err = database.D.Model(&friend).
				Update("last_pub_time", friend.LastPubTime).
				Update("subscribe_error", err.Error()).Error; err != nil {
				log.Println("更新朋友最后更新时间异常:", err)
			}
			continue //跳出这次获取
		}
		log.Println(friend.SubscribeUrl, "获取完成，数据量：", len(feed.Items))
		if len(feed.Items) > 1 {
			for _, item := range feed.Items {
				if siteLastPubTime.Unix() < item.PublishedParsed.Unix() {
					// 存入数据
					article := database.NewArticles()
					article.Link = item.Link
					article.PushTime = item.PublishedParsed
					article.UpdateTime = item.UpdatedParsed
					article.Title = item.Title
					article.Summary = item.Description
					article.FriendId = friend.ID
					database.D.Save(article)
					// 更新最新文章时间
					if article.PushTime.Unix() > friend.LastPubTime.Unix() {
						friend.LastPubTime = article.PushTime
					}
				}
			}
		} else {
			log.Println("本次获取数量异常")
		}

		// 写入最新的更新时间
		if err = database.D.Model(&friend).
			Update("last_pub_time", friend.LastPubTime).
			Update("feed_type", feed.FeedType).
			Update("subscribe_error", "").Error; err != nil {
			log.Println("更新朋友最后更新时间异常:", err)
		}
	}

	// 更新下最后更新时间
	config.UpdateCrawlingTime()
}
