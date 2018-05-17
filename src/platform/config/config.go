package config

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

type Conf struct {
	LogPath            string `yaml:"log_path"`
	LogName            string `yaml:"log_name"`
	LogFileSize        int64  `yaml:"log_file_size"`
	LogFileCount       int    `yaml:"log_file_count"`
	LogLevel           int    `yaml:"log_level"`
	HttpListenPort     string `yaml:"http_listen_port"`
	HttpTemplatePath   string `yaml:"http_template_path"`
	HttpCookieSecret   string `yaml:"http_cookie_secret"`
	HttpSessionId      string `yaml:"http_session_id"`
	HttpAccessTime     string `yaml:"http_access_time"`
	HttpSessionTimeout int64  `yaml:"http_session_timeout"`
	OauthAppId         string `yaml:"oauth_app_id"`
	OauthAppSecret     string `yaml:"oauth_app_secret"`
	OauthUserUrl       string `yaml:"oauth_user_url"`
	OauthAuthUrl       string `yaml:"oauth_auth_url"`
	OauthRedirectUrl   string `yaml:"oauth_redirect_url"`
	OauthTokenUrl      string `yaml:"oauth_token_url"`
	MysqlAddress       string `yaml:"mysql_address"`
	MysqlPort          string `yaml:"mysql_port"`
	MysqlUser          string `yaml:"mysql_user"`
	MysqlPwd           string `yaml:"mysql_pwd"`
	MysqlDbName        string `yaml:"mysql_db_name"`
	ServerLocalMode    bool   `yaml:"server_local_mode"`

	PublicTemplatePath  string
	PrivateTemplatePath string
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

	o.PublicTemplatePath = filepath.Join(o.HttpTemplatePath, "public")
	o.PrivateTemplatePath = filepath.Join(o.HttpTemplatePath, "private")
	return o
}

func (o *Conf) Show() {
	fmt.Println("log_path:", o.LogPath)
	fmt.Println("log_name:", o.LogName)
	fmt.Println("log_file_size:", o.LogFileSize)
	fmt.Println("log_file_count:", o.LogFileCount)
	fmt.Println("http_listen_port:", o.HttpListenPort)
	fmt.Println("http_template_path:", o.HttpTemplatePath)
	fmt.Println("http_cookie_secret:", o.HttpCookieSecret, len(o.HttpCookieSecret))
	fmt.Println("http_session_id:", o.HttpSessionId)
	fmt.Println("http_access_time:", o.HttpAccessTime)
	fmt.Println("http_session_timeout:", o.HttpSessionTimeout)
	fmt.Println("oauth_app_id:", o.OauthAppId)
	fmt.Println("oauth_app_secret:", o.OauthAppSecret)
	fmt.Println("oauth_redirect_url:", o.OauthRedirectUrl)
	fmt.Println("oauth_user_url:", o.OauthUserUrl)
	fmt.Println("oauth_auth_url:", o.OauthAuthUrl)
	fmt.Println("oath_token_url:", o.OauthTokenUrl)
	fmt.Println("mysql_address:", o.MysqlAddress)
	fmt.Println("mysql_port:", o.MysqlPort)
	fmt.Println("mysql_user:", o.MysqlUser)
	fmt.Println("mysql_password:", o.MysqlPwd)
	fmt.Println("mysql_db_name:", o.MysqlDbName)
	fmt.Println("server_local_mode:", o.ServerLocalMode)
}
