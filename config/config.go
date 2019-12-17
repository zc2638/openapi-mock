package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

/**
 * Created by zc on 2019-10-24.
 */
// config 配置项
type configure struct {
	Name     string   `yaml:"name"`
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	OpenApi  OpenApi  `yaml:"openApi"`
	RabbitMQ RabbitMQ `yaml:"rabbitMQ"`
}

type OpenApi struct {
	Host string `yaml:"host"`
}

type RabbitMQ struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Use      bool   `yaml:"use"`
}

var Cfg = &configure{}

func init() {
	// 初始化配置文件
	var err error
	b, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Fatal("配置加载失败：", err)
	}

	if err := yaml.Unmarshal(b, Cfg); err != nil {
		log.Fatal("配置解析失败：", err)
	}

	if Cfg.Port == "" {
		Cfg.Port = ServerPort
	}
}
