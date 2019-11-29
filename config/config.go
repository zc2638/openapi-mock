package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/**
 * Created by zc on 2019-10-24.
 */
// config 配置项
type configure struct {
	Name    string  `yaml:"name"`
	Host    string  `yaml:"host"`
	Port    string  `yaml:"port"`
	OpenApi OpenApi `yaml:"openApi"`
}

type OpenApi struct {
	Host string `yaml:"host"`
}

var Cfg = &configure{}

func init() {
	// 初始化配置文件
	var err error
	b, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, Cfg); err != nil {
		panic(err)
	}
}
