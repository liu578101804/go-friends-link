package config

import (
	"encoding/json"
	"fmt"
	"friends-rss/modules"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

var ConfigInstance *modules.ConfigModule

func InitConfig() {
	ConfigInstance.Port = 80 // 默认80端口
	// 随机生成
	token := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100000000))
	ConfigInstance.Token = token
	fmt.Println("生产随机数 token：", token)
	ConfigInstance.Cron = "*/1 * * * *"
	fmt.Println("默认的调度(每分钟一次) cron:", ConfigInstance.Cron)
	syncFile()
}

func init() {
	ConfigInstance = &modules.ConfigModule{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("读取文件错误：", err.Error())
		InitConfig()
		return
	}
	if err := json.Unmarshal(data, ConfigInstance); err != nil {
		fmt.Println("解析文件错误：", err.Error())
		InitConfig()
	}
}

func UpdateCrawlingTime() {
	ConfigInstance.LastCrawlingTime = time.Now().Format("2006-01-02 15:04:05")
	syncFile()
}

func syncFile() error {
	data, _ := json.Marshal(ConfigInstance)
	err := ioutil.WriteFile("./config.json", data, 0666) //写入文件(字节数组)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
