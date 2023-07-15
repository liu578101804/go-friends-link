package main

import (
	"fmt"
	"friends-rss/config"
	"friends-rss/crawling"
	"friends-rss/database"
	"friends-rss/helper"
	"friends-rss/server/handler"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"strings"
)

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

// Cron 全局调度器
var Cron *cron.Cron

// startCron 定时任务
func startCron() {
	cronStr := config.ConfigInstance.Cron
	if cronStr == "" {
		log.Println("调度表达式异常")
		return
	}
	log.Println("调度表达式为：", cronStr)
	Cron = cron.New()
	var entryId cron.EntryID
	entryId, _ = Cron.AddFunc(cronStr, func() {
		// 定时任务执行开始标识
		log.Println("+++ 开始执行定时任务...")
		// 获取数据
		crawling.Crawling()
		// 定时任务执行结束标识
		log.Println("+++ 执行了一次定时任务...")
		// 打印下次执行时间
		printNextRunTime(entryId)
	})
	// 启动，自动启分支
	Cron.Start()
	// 打印下次执行时间
	printNextRunTime(entryId)
}
func printNextRunTime(id cron.EntryID) {
	nextTime := Cron.Entry(id).Next.Format("2006-01-02 15:04:05")
	log.Println("调度器下次执行时间为：", nextTime)
}

// checkNeedDir 检查必须的文件夹是否存在
func checkNeedDir() {
	log.Println("检查所需要的文件夹是否存在")
	paths := []string{"./data"}
	for _, path := range paths {
		log.Println("检查文件夹：", path)
		if has, _ := helper.PathExists(path); !has {
			log.Println("开始创建文件夹：", path)
			if err := helper.MkdirAllDir(path); err != nil {
				log.Println("文件夹创建失败：", path)
				panic(err)
			}
		}
	}
}

func main() {
	// 检查必须的文件夹是否存在
	checkNeedDir()
	// 初始化配置文件
	config.Init()
	// 初始化数据库
	database.Init()
	// 启动定时任务
	startCron()
	// 引擎
	r := gin.Default()
	// 注册中间件
	r.Use(Cors())
	// 添加后台
	r.Static("/ui", "./ui")
	// 注册API
	handler.RegisterAPI(r)
	// 启动
	r.Run(fmt.Sprintf(":%d", config.ConfigInstance.Port))
}
