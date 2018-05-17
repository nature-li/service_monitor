package main

import (
	"net/http"
	"fmt"
	"mt/mtlog"
	"mt/session/cookie"
	"flag"
	"plat/global"
	"plat/handler"
)

func main() {
	var confPath = flag.String("conf", "", "config file path")
	flag.Parse()

	if *confPath == "" {
		fmt.Println("conf is empty")
		return
	}

	if global.Conf.GetConf(*confPath) == nil {
		fmt.Println("parse config file error")
		return
	}
	global.Conf.Show()

	global.Logger = mtlog.NewLogger(false, mtlog.DEVELOP, mtlog.Level(global.Conf.LogLevel), global.Conf.LogPath, global.Conf.LogName, global.Conf.LogFileSize, global.Conf.LogFileCount)
	if !global.Logger.Start() {
		fmt.Println("logger.Start failed")
	}

	var err error
	global.Manager, err = cookie.NewManager(global.Conf.HttpCookieSecret, global.Conf.HttpSessionId, global.Conf.HttpAccessTime, global.Conf.HttpSessionTimeout)
	if err != nil {
		global.Logger.Error("NewManager failed")
		return
	}
	defer global.Logger.Stop()

	// 公开模板文件
	publicFs := http.FileServer(http.Dir(global.Conf.PublicTemplatePath))
	http.Handle("/templates/public/", http.StripPrefix("/templates/public/", publicFs))
	// 私有模板文件
	http.HandleFunc("/templates/private/", handler.PrivateFileHandler)
	// 图标
	http.HandleFunc("/favicon.ico", handler.FaviconHandler)

	// 404页面
	http.HandleFunc("/not_found", handler.NotFoundHandler)
	// 拒绝访问页面
	http.HandleFunc("/not_allowed", handler.NotAllowHandler)
	// 验证码
	http.HandleFunc("/captcha", handler.CaptchaAPIHandler)

	// 登录页面
	if global.Conf.ServerLocalMode {
		http.HandleFunc("/user_login", handler.UserLoginHandler)
		http.HandleFunc("/user_login_api", handler.UserLoginAPIHandler)
	} else {
		// OA登录
		http.HandleFunc("/user_login_auth", handler.UserLoginAuthHandler)
		http.HandleFunc("/user_login_auth_api", handler.UserLoginAuthAPIHandler)
	}
	// 退出登录
	http.HandleFunc("/user_logout", handler.UserLogoutHandler)

	// 首页
	http.HandleFunc("/", handler.ListUserHandler)
	// 用户列表
	http.HandleFunc("/list_user", handler.ListUserHandler)

	// 用户列表
	http.HandleFunc("/list_user_api", handler.ListUserAPIHandler)
	// 编辑用户
	http.HandleFunc("/edit_user_api", handler.EditUserAPIHandler)
	// 添加用户
	http.HandleFunc("/add_user_api", handler.AddUserAPIHandler)
	// 删除用户
	http.HandleFunc("/del_user_api", handler.DelUserAPIHandler)


	err = http.ListenAndServe(global.Conf.HttpListenPort, nil)
	if err != nil {
		global.Logger.Error(err.Error())
	}
}
