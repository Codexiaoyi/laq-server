package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	//github
	AccessToken string
	Owner       string
	Repo        string
)

func init() {
	file, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("Load config file error!", err)
	}
	LoadSource(file)
}

//加载数据资源配置
func LoadSource(file *ini.File) {
	AccessToken = file.Section("github").Key("AccessToken").MustString("")
	Owner = file.Section("github").Key("Owner").MustString("")
	Repo = file.Section("github").Key("Repo").MustString("")
}
