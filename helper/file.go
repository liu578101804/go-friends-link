package helper

import (
	"os"
	"strings"
)

// GetWorkDir 获取当前工作目录
func GetWorkDir() string {
	wd, _ := os.Getwd()
	return strings.Replace(wd, "\\", "/", -1)
}

// PathExists 判断 文件夹 或者 文件 是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MkdirAllDir 创建文件夹
func MkdirAllDir(dir string) error {
	// 创建文件夹
	err := os.MkdirAll(dir, os.ModePerm)
	return err
}
