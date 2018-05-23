package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	LogPath         string `yaml:"log_path"`
	LogName         string `yaml:"log_name"`
	LogLevel        int    `yaml:"log_level"`
	LogFileSize     int64  `yaml:"log_file_size"`
	LogFileCount    int    `yaml:"log_file_count"`
	HttpListenPort  string `yaml:"http_listen_port"`
	KeyFile         string `yaml:"key_file"`
	KeyPwd          string `yaml:"key_pwd"`
	MysqlUser       string `yaml:"mysql_user"`
	MysqlPwd        string `yaml:"mysql_pwd"`
	MysqlAddress    string `yaml:"mysql_address"`
	MysqlPort       string `yaml:"mysql_port"`
	MysqlDbName     string `yaml:"mysql_db_name"`
	MonitorInterval int    `yaml:"monitor_interval"`
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
	fmt.Println("Key_file:", o.KeyFile)
	fmt.Println("Key_pwd:", o.KeyPwd)
	fmt.Println("yysql_user:", o.MysqlUser)
	fmt.Println("yysql_pwd:", o.MysqlPwd)
	fmt.Println("mysql_address:", o.MysqlAddress)
	fmt.Println("mysql_port:", o.MysqlPort)
	fmt.Println("mysql_db_name:", o.MysqlDbName)
	fmt.Println("monitor_interval:", o.MonitorInterval)
}
