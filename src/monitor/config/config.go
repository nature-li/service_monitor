package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	LogPath            string `yaml:"log_path"`
	LogName            string `yaml:"log_name"`
	LogFileSize        int64  `yaml:"log_file_size"`
	LogFileCount       int    `yaml:"log_file_count"`
	HttpListenPort     string `yaml:"http_listen_port"`
	HttpTemplatePath   string `yaml:"http_template_path"`
	HttpCookieSecret   string `yaml:"http_cookie_secret"`
	HttpSessionId      string `yaml:"http_session_id"`
	HttpAccessTime     string `yaml:"http_access_time"`
	HttpSessionTimeout int64  `yaml:"http_session_timeout"`
	UploadDataPath     string `yaml:"upload_data_path"`
	UploadMaxSize      int64  `yaml:"upload_max_size"`
	SqliteDbPath       string `yaml:"sqlite_db_path"`
	OauthAppId         string `yaml:"oauth_app_id"`
	OauthAppSecret     string `yaml:"oauth_app_secret"`
	OauthUserUrl       string `yaml:"oauth_user_url"`
	OauthAuthUrl       string `yaml:"oauth_auth_url"`
	OauthRedirectUrl   string `yaml:"oauth_redirect_url"`
	OauthTokenUrl      string `yaml:"oauth_token_url"`
	ServerLocalMode    bool   `yaml:"server_local_mode"`

	maxUploadSizeStr    string
	publicTemplatePath  string
	privateTemplatePath string
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

	o.maxUploadSizeStr = strconv.FormatFloat(float64(o.UploadMaxSize)/1024/1024, 'f', 2, 64)
	o.publicTemplatePath = filepath.Join(o.HttpTemplatePath, "public")
	o.privateTemplatePath = filepath.Join(o.HttpTemplatePath, "private")
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
	fmt.Println("upload_data_path:", o.UploadDataPath)
	fmt.Println("upload_max_size:", o.UploadMaxSize)
	fmt.Println("sqlite_db_path:", o.SqliteDbPath)
	fmt.Println("oauth_app_id:", o.OauthAppId)
	fmt.Println("oauth_app_secret:", o.OauthAppSecret)
	fmt.Println("oauth_redirect_url:", o.OauthRedirectUrl)
	fmt.Println("oauth_user_url:", o.OauthUserUrl)
	fmt.Println("oauth_auth_url:", o.OauthAuthUrl)
	fmt.Println("oath_token_url:", o.OauthTokenUrl)
	fmt.Println("server_local_mode:", o.ServerLocalMode)
}
