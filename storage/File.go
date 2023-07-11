package storage

import (
	"encoding/json"
	"fmt"
	"friends-rss/config"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type ArticleData struct {
	Floor   int    `json:"floor"`
	Title   string `json:"title"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Link    string `json:"link"`
	Author  string `json:"author"`
	Avatar  string `json:"avatar"`
}

type StatisticalData struct {
	FriendsNum      int    `json:"friends_num"`
	ActiveNum       int    `json:"active_num"`
	ErrorNum        int    `json:"error_num"`
	ArticleNum      int    `json:"article_num"`
	LastUpdatedTime string `json:"last_updated_time"`
}

type StorageModule struct {
	StatisticalData StatisticalData `json:"statistical_data"`
	ArticleData     []*ArticleData  `json:"article_data"`
}

var StorageInstance *StorageModule

const DataFileDir = "./tmp"

var DataFilePath = fmt.Sprintf("%s%s", DataFileDir, "/data.json")

func init() {

	StorageInstance = new(StorageModule)

	// 判断tmp文件是否存在
	_, err := os.Stat(DataFileDir)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建
			err := os.MkdirAll(DataFileDir, os.ModePerm)
			if err != nil {
				fmt.Println("创建数据文件失败", err)
				return
			}
			SaveArticleData(make([]*ArticleData, 0))
		}
	}

	data, err := ioutil.ReadFile(DataFilePath)
	if err != nil {
		fmt.Println("读取数据文件错误：", err)
		StorageInstance = new(StorageModule)
		SaveArticleData(make([]*ArticleData, 0))
	} else {
		err = json.Unmarshal(data, StorageInstance)
		if err != nil {
			fmt.Println("序列化数据文件错误：", err)
			SaveArticleData(make([]*ArticleData, 0))
		}
	}
}

func SaveArticleData(articles []*ArticleData) error {
	StorageInstance.ArticleData = articles
	StorageInstance.StatisticalData.FriendsNum = len(config.ConfigInstance.Links)
	StorageInstance.StatisticalData.ArticleNum = len(articles)
	StorageInstance.StatisticalData.LastUpdatedTime = time.Now().Format("2006-01-02 15:04:05")
	data, _ := json.Marshal(StorageInstance)
	err := ioutil.WriteFile(DataFilePath, data, 0666) //写入文件(字节数组)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
