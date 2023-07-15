package config

import (
	"encoding/json"
	"fmt"
	"friends-rss/helper"
	"friends-rss/modules"
	"io/ioutil"
	"log"
	"time"
)

var ConfigInstance *modules.ConfigModule

const configFilePath = "./data/config.json"

func InitConfig() {
	ConfigInstance.Port = 80 // 默认80端口
	// 随机生成 8 位字符
	token := helper.RandString(8)
	ConfigInstance.Token = token
	fmt.Println("生产随机数 token：", token)
	ConfigInstance.Cron = "30 * * * *"
	fmt.Println("默认的调度(每小时的30分) cron:", ConfigInstance.Cron)
	ConfigInstance.LastCrawlingTime = time.Now().Format("2006-01-02 15:04:05")
	syncFile()
}

func Init() {
	ConfigInstance = &modules.ConfigModule{}
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("读取文件错误：", err.Error())
		log.Println("开始新建默认配置文件")
		InitConfig()
		return
	}
	if err := json.Unmarshal(data, ConfigInstance); err != nil {
		log.Println("解析文件错误：", err.Error())
		InitConfig()
	}
}

func UpdateCrawlingTime() {
	ConfigInstance.LastCrawlingTime = time.Now().Format("2006-01-02 15:04:05")
	syncFile()
}

func syncFile() error {
	data, _ := json.Marshal(ConfigInstance)
	err := ioutil.WriteFile(configFilePath, data, 0666) //写入文件(字节数组)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
