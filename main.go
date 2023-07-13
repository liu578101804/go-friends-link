package main

import (
	"encoding/xml"
	"fmt"
	"friends-rss/config"
	_ "friends-rss/database"
	"friends-rss/modules"
	"friends-rss/server/handler"
	"friends-rss/storage"
	"friends-rss/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/noaway/dateparse"
	"net/http"
	"strings"
	"time"
)

// 去爬取消息
func crawling() {

	links := config.ConfigInstance.Links

	if len(links) == 0 {
		fmt.Println("没有友链存在，本次结束!")
		return
	}

	articleData := make([]*storage.ArticleData, 0)

	for i := 0; i < len(links); i++ {
		fmt.Println("开始获取：", links[i].Url, links[i].Kind)

		rstData, err := utils.GetBody(links[i].Url)
		if err != nil {
			fmt.Println(err)
		}
		if links[i].Kind == "feed" {

			fmt.Println(links[i].Url, "feed 获取完成")
			rss := modules.RssModule{}
			err := xml.Unmarshal(rstData, &rss)
			if err != nil {
				fmt.Printf("Error Decode: %v\n", err)
				return
			}

			for _, item := range rss.Channel.Items {

				// 解析发布时间
				t, err := dateparse.ParseAny(item.PubDate)
				if err != nil {
					fmt.Println("转换时间格式错误")
					panic(err.Error())
				}
				articleData = append(articleData, &storage.ArticleData{
					Floor:   0,
					Title:   item.Title,
					Updated: t.Format("2006-01-02 15:04:05"),
					Link:    item.Link,
					Author:  links[i].Author,
					Avatar:  links[i].Avatar,
				})
			}

		} else if links[i].Kind == "atom" {

			fmt.Println(links[i].Url, "atom 获取完成")
			atom := modules.AtomModule{}
			err := xml.Unmarshal(rstData, &atom)
			if err != nil {
				fmt.Printf("Error Decode: %v\n", err)
				return
			}

			for _, item := range atom.Entry {
				// 解析发布时间
				t, err := dateparse.ParseAny(item.Updated)
				if err != nil {
					fmt.Println("转换时间格式错误")
					panic(err.Error())
				}
				articleData = append(articleData, &storage.ArticleData{
					Floor:   0,
					Title:   item.Title,
					Updated: t.Format("2006-01-02 15:04:05"),
					Link:    item.Link,
					Author:  links[i].Author,
					Avatar:  links[i].Avatar,
				})
			}

		}

	}
	// 保存到本地
	storage.SaveArticleData(articleData)

	fmt.Println("本次获取数据完成")
}

// Cors 开放所有接口的OPTIONS方法
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			origin := c.Request.Header.Get("Origin")
			c.Header("Access-Control-Allow-Origin", origin)                                    // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

/*定时任务*/
func StartCron() {
	s := gocron.NewScheduler(time.UTC)

	if config.ConfigInstance.Cron == "" {
		fmt.Println("调度表达式异常")
		return
	}
	fmt.Println("调度表达式为：", config.ConfigInstance.Cron)
	s.Cron(config.ConfigInstance.Cron).Tag("crawling").Do(func() {
		// 获取数据
		crawling()
		fmt.Println("执行了一次定时任务...")
		_, t := s.NextRun()
		fmt.Println("下一次执行时间：", t.Format("2006-01-02 15:04:05"))
	})

	// 启动
	s.StartBlocking()
}

func main() {

	//启动定时任务
	go StartCron()

	r := gin.Default()
	r.Use(Cors())

	// 注册API
	handler.RegisterAPI(r)

	r.Run(fmt.Sprintf(":%d", config.ConfigInstance.Port))
}
