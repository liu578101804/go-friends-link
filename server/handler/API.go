package handler

import (
	"errors"
	"friends-rss/config"
	"friends-rss/modules"
	"friends-rss/storage"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

// 添加友链
func AddLinkFunc(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token != config.ConfigInstance.Token {
		RsError(errors.New("token error"), c)
		return
	}

	requestM := new(modules.LinkItem)
	err := c.BindJSON(requestM)
	if err != nil {
		RsError(err, c)
		return
	}
	err = config.AddLinks(requestM)
	if err != nil {
		RsError(err, c)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}

// 添加友链
func UpdateLinkFunc(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token != config.ConfigInstance.Token {
		RsError(errors.New("token error"), c)
		return
	}

	requestM := new(modules.LinkItem)
	err := c.BindJSON(requestM)
	if err != nil {
		RsError(err, c)
		return
	}
	err = config.UpdateLinks(requestM)
	if err != nil {
		RsError(err, c)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}

// 删除友链
func DelLinkFunc(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token != config.ConfigInstance.Token {
		RsError(errors.New("token error"), c)
		return
	}

	url := c.DefaultQuery("url", "")
	if url == "" {
		RsError(errors.New("参数不能为空"), c)
		return
	}
	err := config.DelLinks(url)
	if err != nil {
		RsError(err, c)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}

// 禁止访问配置文件
func DontGetConfigFunc(c *gin.Context) {
	c.JSON(403, gin.H{
		"message": "Forbidden",
	})
}

// 获取数据
func GetAllFunc(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		RsError(err, c)
		return
	}
	if pageInt-1 <= 0 {
		pageInt = 1
	}

	// 按更新时间排序
	sort.Sort(storage.StorageInstance.ArticleData)
	articles := make([]*storage.ArticleData, 0)
	for i := (pageInt - 1) * 10; i < len(storage.StorageInstance.ArticleData) && i < pageInt*10; i++ {
		articles = append(articles, storage.StorageInstance.ArticleData[i])
	}

	rst := storage.StorageModule{
		StatisticalData: storage.StorageInstance.StatisticalData,
		ArticleData:     articles,
		Page:            pageInt,
	}

	c.JSON(200, rst)
}

// 返回错误信息
func RsError(err error, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": err.Error(),
	})
}

func RegisterAPI(r *gin.Engine) {
	// 获取数据
	r.GET("/all", GetAllFunc)
	// 禁止访问配置文件
	r.GET("/config.json", DontGetConfigFunc)
	// 删除友链
	r.DELETE("/link", DelLinkFunc)
	// 修改友链
	r.PUT("/link", UpdateLinkFunc)
	// 添加友链
	r.POST("/link", AddLinkFunc)
}
