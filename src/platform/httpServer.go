package main

import (
	"net/http"
	"fmt"
	"mt/mtlog"
	"mt/session/cookie"
	"flag"
	"platform/global"
	"platform/handler"
	"platform/config"
)

func main() {
	var confPath = flag.String("conf", "", "config file path")
	flag.Parse()

	if *confPath == "" {
		fmt.Println("conf is empty")
		return
	}

	global.Conf = &config.Conf{}
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

	// public files in template
	publicFs := http.FileServer(http.Dir(global.Conf.PublicTemplatePath))
	http.Handle("/templates/public/", http.StripPrefix("/templates/public/", publicFs))
	// private files in template
	http.HandleFunc("/templates/private/", handler.PrivateFileHandler)
	// icon
	http.HandleFunc("/favicon.ico", handler.FaviconHandler)

	// 404 page
	http.HandleFunc("/not_found", handler.NotFoundHandler)
	// refuse page
	http.HandleFunc("/not_allowed", handler.NotAllowHandler)
	// captcha
	http.HandleFunc("/captcha", handler.CaptchaAPIHandler)

	// login
	if global.Conf.ServerLocalMode {
		// login by login page
		http.HandleFunc("/user_login", handler.UserLoginHandler)
		http.HandleFunc("/user_login_api", handler.UserLoginAPIHandler)
	} else {
		// login by OA
		http.HandleFunc("/user_login_auth", handler.UserLoginAuthHandler)
		http.HandleFunc("/user_login_auth_api", handler.UserLoginAuthAPIHandler)
	}
	// logout
	http.HandleFunc("/user_logout", handler.UserLogoutHandler)

	// home page
	http.HandleFunc("/", handler.ListUserHandler)
	// user list
	http.HandleFunc("/list_user", handler.ListUserHandler)

	// user list api
	http.HandleFunc("/list_user_api", handler.ListUserAPIHandler)
	// edit list api
	http.HandleFunc("/edit_user_api", handler.EditUserAPIHandler)
	// add user
	http.HandleFunc("/add_user_api", handler.AddUserAPIHandler)
	// delete user
	http.HandleFunc("/del_user_api", handler.DelUserAPIHandler)


	err = http.ListenAndServe(global.Conf.HttpListenPort, nil)
	if err != nil {
		global.Logger.Error(err.Error())
	}
}
