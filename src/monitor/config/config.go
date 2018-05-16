package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	LogPath        string `yaml:"log_path"`
	LogName        string `yaml:"log_name"`
	LogLevel       int    `yaml:"log_level"`
	LogFileSize    int64  `yaml:"log_file_size"`
	LogFileCount   int    `yaml:"log_file_count"`
	HttpListenPort string `yaml:"http_listen_port"`
}

func (o *Conf) GetConf(path string) *Conf {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	err = yaml.Unmarshal(yamlFile, o)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return o
}

func (o *Conf) Show() {
	fmt.Println("log_path:", o.LogPath)
	fmt.Println("log_name:", o.LogName)
	fmt.Println("log_level:", o.LogLevel)
	fmt.Println("log_file_size:", o.LogFileSize)
	fmt.Println("log_file_count:", o.LogFileCount)
	fmt.Println("http_listen_port:", o.HttpListenPort)
}
