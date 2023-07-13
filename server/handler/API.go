package handler

import (
	"errors"
	"friends-rss/config"
	"friends-rss/database"
	"friends-rss/modules"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// SaveFriendFunc 添加朋友
func SaveFriendFunc(c *gin.Context) {
	// 校验token
	token := c.DefaultQuery("token", "")
	if token != config.ConfigInstance.Token {
		RsError(errors.New("token error"), c)
		return
	}
	// 解析请求数据
	requestM := new(modules.AddFriendRequestM)
	err := c.BindJSON(requestM)
	if err != nil {
		RsError(err, c)
		return
	}

	// 查询是否存在
	friend := database.NewFriends()
	// 保存到数据库
	friend.WebUrl = requestM.WebUrl

	// 最后一次更新时间
	lastUpdateTime := time.Now().Add(-time.Hour * 4800) //往前推200天

	// 如果没找到
	if err = database.D.Where("web_url=?", requestM.WebUrl).First(friend).Error; err == gorm.ErrRecordNotFound {
		log.Println("新增友链")
		friend.SubscribeUrl = requestM.SubscribeUrl
		friend.WebTitle = requestM.WebTitle
		friend.WebDescribe = requestM.WebDescribe
		friend.AuthorName = requestM.AuthorName
		friend.AuthorAvatar = requestM.AuthorAvatar
		friend.LastUpdateTime = &lastUpdateTime
		err = database.D.Save(friend).Error
	} else {
		log.Println("更新友链")
		friend.SubscribeUrl = requestM.SubscribeUrl
		friend.WebTitle = requestM.WebTitle
		friend.WebDescribe = requestM.WebDescribe
		friend.AuthorName = requestM.AuthorName
		friend.AuthorAvatar = requestM.AuthorAvatar
		friend.LastUpdateTime = &lastUpdateTime
		err = database.D.Updates(friend).Error
	}
	if err != nil {
		RsError(err, c)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    friend,
	})
}

// DelFriendFunc 删除朋友
func DelFriendFunc(c *gin.Context) {
	// 校验token
	token := c.DefaultQuery("token", "")
	if token != config.ConfigInstance.Token {
		RsError(errors.New("token error"), c)
		return
	}
	// 获取请求参数
	id := c.DefaultQuery("id", "")
	if id == "" {
		RsError(errors.New("参数不能为空"), c)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		RsError(errors.New("非法参数"), c)
		return
	}
	// 查询是否存在
	friend := database.NewFriends()
	if err := database.D.Where("id=?", idInt).First(friend).Error; err == gorm.ErrRecordNotFound {
		RsError(errors.New("没找到友链"), c)
		return
	}
	// 删除
	if err := database.D.Model(friend).Unscoped().Delete("id=?", friend.ID).Error; err != nil {
		log.Println(err.Error())
		RsError(errors.New("删除友链失败"), c)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}

// DontGetConfigFunc 禁止访问配置文件
func DontGetConfigFunc(c *gin.Context) {
	c.JSON(403, gin.H{
		"message": "Forbidden",
	})
}

// GetAllFunc 获取数据
func GetAllFunc(c *gin.Context) {
	// 获取页码
	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		RsError(err, c)
		return
	}
	if pageInt-1 <= 0 {
		pageInt = 1
	}

	rsArticles := make([]*modules.ResponseArticleModule, 0)

	// 获取十条
	articles := make([]*database.Articles, 0)
	if err := database.D.Order("push_time desc").Offset((pageInt - 1) * 10).Limit(10).Find(&articles).Error; err != nil {
		log.Println("获取异常：", err.Error())
		RsError(errors.New("获取异常"), c)
		return
	}

	for _, article := range articles {

		art := &modules.ResponseArticleModule{
			Floor:   0,
			Title:   article.Title,
			Created: article.PushTime.Format("2006-01-02 15:04:05"),
			Link:    article.Link,
			Author:  article.AuthorName,
			Avatar:  article.AuthorAvatar,
		}

		if article.UpdateTime != nil && !article.UpdateTime.IsZero() {
			art.Updated = article.UpdateTime.Format("2006-01-02 15:04:05")
		}

		rsArticles = append(rsArticles, art)
	}

	// 文章数量
	var articleCount64 int64
	if err := database.D.Model(database.NewArticles()).Count(&articleCount64).Error; err != nil {
		log.Println("获取异常：", err.Error())
		RsError(errors.New("获取异常"), c)
		return
	}

	// 获取友链情况
	friends, err := database.GetAllFriends()
	if err != nil {
		log.Println("获取异常：", err.Error())
		RsError(errors.New("获取异常"), c)
		return
	}

	rst := modules.ResponseModule{
		StatisticalData: &modules.ResponseStatisticalDataModule{
			FriendsNum:      len(friends),
			ActiveNum:       len(friends),
			ErrorNum:        0,
			ArticleNum:      int(articleCount64),
			LastUpdatedTime: config.ConfigInstance.LastCrawlingTime,
		},
		ArticleData: rsArticles,
		Page:        pageInt,
	}

	c.JSON(200, rst)
}

// GetAllFriends 获取全部朋友
func GetAllFriends(c *gin.Context) {
	friends, err := database.GetAllFriends()
	if err != nil {
		log.Println("获取异常：", err.Error())
		RsError(errors.New("获取异常"), c)
		return
	}
	c.JSON(200, gin.H{
		"message": "加载成功",
		"data":    friends,
	})
}

// RsError 返回错误信息
func RsError(err error, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": err.Error(),
	})
}

func RegisterAPI(r *gin.Engine) {
	// 获取全部朋友
	r.GET("/friend", GetAllFriends)
	// 获取数据
	r.GET("/all", GetAllFunc)
	// 禁止访问配置文件
	r.GET("/config.json", DontGetConfigFunc)
	// 删除友链
	r.DELETE("/friend", DelFriendFunc)
	// 添加友链
	r.POST("/friend", SaveFriendFunc)
}
