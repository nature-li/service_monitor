package handler

import (
	"net/http"
	"session"
	"strconv"
)

type pageData struct {
	LoginName            string
	WrapperClass         string
	PinLock              string
	HiddenClass          string
	UploadMaxFileSize    int64
	UploadMaxFileSizeStr string
	DownloadRight        bool
	UploadRight          bool
	ModifyRight          bool
	UserRight            bool
}

func (o *pageData)reCalcModifyRight(loginEmail, uploaderEmail string) {
	if (loginEmail != "") && (loginEmail == uploaderEmail) {
		o.ModifyRight = true
	}
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

	digitRight, err := strconv.ParseInt(userRight, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		digitRight = 0
	}

	downloadRight := false
	if (digitRight & DOWNLOAD_RIGHT) != 0 {
		downloadRight = true
	}

	uploadRight := false
	if (digitRight & UPLOAD_RIGHT) != 0 {
		uploadRight = true
	}

	managerRight := false
	modifyRight := false
	if (digitRight & MANAGER_RIGHT) != 0 {
		managerRight = true
		modifyRight = true
	}

	return &pageData{
		LoginName:            loginName,
		WrapperClass:         wrapperClass,
		PinLock:              pinLock,
		HiddenClass:          hiddenClass,
		UploadMaxFileSize:    config.UploadMaxSize,
		UploadMaxFileSizeStr: config.maxUploadSizeStr,
		DownloadRight:        downloadRight,
		UploadRight:          uploadRight,
		UserRight:            managerRight,
		ModifyRight:          modifyRight}
}
