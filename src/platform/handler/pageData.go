package handler

import (
	"net/http"
	"mt/session"
	"strconv"
	"platform/global"
)

type pageData struct {
	LoginName            string
	WrapperClass         string
	PinLock              string
	HiddenClass          string
	UploadMaxFileSize    int64
	UploadMaxFileSizeStr string
	UserRight            bool
}

func newPageData(w http.ResponseWriter, r *http.Request, s session.Session) *pageData {
	// 是否展开侧边栏
	wrapperClass := ""
	hiddenClass := ""
	if cookie, ok := r.Cookie("pin_nav"); ok == nil {
		if cookie.Value == "1" {
			wrapperClass = "toggled"
			hiddenClass = "hidden-self"
		}
	}

	// 是否锁住浮动锁
	pinLock := "glyphicon-pushpin"
	if cookie, ok := r.Cookie("pin_lock"); ok == nil {
		if cookie.Value == "1" {
			pinLock = "glyphicon-lock"
		}
	}

	// 登录相关
	var loginName = ""
	var userRight = ""
	if s != nil {
		if s.Get("is_login") == "1" {
			loginName = s.Get("user_name")
			userRight = s.Get("user_right")
		}
	}

	var digitRight int64 = 0
	var err error = nil
	if len(userRight) != 0 {
		digitRight, err = strconv.ParseInt(userRight, 10, 64)
		if err != nil {
			global.Logger.Error(err.Error())
			digitRight = 0
		}
	}

	managerRight := false
	if (digitRight & MANAGER_RIGHT) != 0 {
		managerRight = true
	}

	return &pageData{
		LoginName:            loginName,
		WrapperClass:         wrapperClass,
		PinLock:              pinLock,
		HiddenClass:          hiddenClass,
		UserRight:            managerRight}
}
