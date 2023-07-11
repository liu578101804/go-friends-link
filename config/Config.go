package config

import (
	"encoding/json"
	"errors"
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
	ConfigInstance.Links = make([]*modules.LinkItem, 0)
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

func getLinkItem(url string) (item *modules.LinkItem, index int, err error) {
	for i := 0; i < len(ConfigInstance.Links); i++ {
		item = ConfigInstance.Links[i]
		if item.Url == url {
			return item, i, nil
		}
	}
	return nil, 0, errors.New("没找到想要的友链")
}

func AddLinks(m *modules.LinkItem) error {
	_, _, err := getLinkItem(m.Url)
	if err == nil {
		return errors.New("友链已经存在")
	}
	ConfigInstance.Links = append(ConfigInstance.Links, m)
	syncFile()
	return nil
}

func UpdateLinks(m *modules.LinkItem) error {
	item, _, err := getLinkItem(m.Url)
	if err == nil { //存在友链
		// 添加
		AddLinks(m)
		// 删除旧的
		DelLinks(item.Url)
	}
	ConfigInstance.Links = append(ConfigInstance.Links, m)
	syncFile()
	return nil
}

func DelLinks(url string) error {
	_, index, err := getLinkItem(url)
	if err != nil {
		return err
	}
	ConfigInstance.Links = append(ConfigInstance.Links[:index], ConfigInstance.Links[index+1:]...)
	syncFile()
	return nil
}

func syncFile() error {
	data, _ := json.Marshal(ConfigInstance)
	err := ioutil.WriteFile("./config.json", data, 0666) //写入文件(字节数组)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
