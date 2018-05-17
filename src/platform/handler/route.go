package handler

import (
	"net/http"
	"html/template"
	"net/url"
	"path/filepath"
	"mt/session"
	"platform/global"
	"strconv"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(filepath.Join(global.Conf.PrivateTemplatePath, "html/404.html"))
	if err != nil {
		global.Logger.Error(err.Error())
	}

	t.Execute(w, nil)
}

func NotAllowHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(filepath.Join(global.Conf.PrivateTemplatePath, "html/refuse.html"))
	if err != nil {
		global.Logger.Error(err.Error())
	}

	t.Execute(w, nil)
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	t, err := template.ParseFiles(filepath.Join(global.Conf.PublicTemplatePath, "html/user_login.html"))
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	handler := newCaptchaHandler(s)
	_, hex, err := handler.createCaptcha()
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	value := struct {
		Base64Img string
		*pageData
	} {
		Base64Img: hex,
		pageData : newPageData(w, r, s),
	}

	t.Execute(w, value)
}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	global.Manager.SessionDestroy(w, r)

	if global.Conf.ServerLocalMode {
		http.Redirect(w, r, "/user_login", 302)
	} else {
		http.Redirect(w, r, "/user_login_auth", 302)
	}
}

func UserLoginAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	handler := userLoginAPI{session:s}
	handler.handle(w, r)
}

func UserLoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	var redirectUrl = global.Conf.OauthAuthUrl
	redirectUrl += "?appid=" + global.Conf.OauthAppId
	redirectUrl += "&response_type=code"
	redirectUrl += "&redirect_uri=" + url.QueryEscape(global.Conf.OauthRedirectUrl)
	redirectUrl += "&scope=user_info"
	redirectUrl += "&state=test"
	http.Redirect(w, r, redirectUrl, 302)
}

func UserLoginAuthAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	handler := userLoginAuthAPI{session:s}
	handler.handle(w, r)
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, global.Conf.PublicTemplatePath + "/img/favicon.ico")
}

func CaptchaAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	handler := newCaptchaHandler(s)
	handler.handle(w, r)
}

func PrivateFileHandler(w http.ResponseWriter, r *http.Request) {
	// 检测登录
	s := global.Manager.SessionStart(w, r)
	if !checkLogin(s) {
		if global.Conf.ServerLocalMode {
			http.Redirect(w, r, "/user_login", 302)
		} else {
			http.Redirect(w, r, "/user_login_auth", 302)
		}
		return
	}

	dataFs := http.FileServer(http.Dir(global.Conf.PrivateTemplatePath))
	http.StripPrefix("/templates/private/", dataFs).ServeHTTP(w, r)
}

func checkLogin(s session.Session) bool {
	isLogin := s.Get("is_login")
	if isLogin == "1" {
		return true
	}

	return false
}

func checkRight(s session.Session, right int64) bool {
	var userRight string
	if userRight = s.Get("user_right"); userRight == "" {
		return false
	}
	digitRight, err := strconv.ParseInt(userRight, 10, 64)
	if err != nil {
		global.Logger.Error(err.Error())
		return false
	}
	if digitRight & right == 0 {
		return false
	}

	return true
}

func ListUserAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	if !checkLogin(s) {
		if global.Conf.ServerLocalMode {
			http.Redirect(w, r, "/user_login", 302)
		} else {
			http.Redirect(w, r, "/user_login_auth", 302)
		}
		return
	}

	if !checkRight(s, MANAGER_RIGHT) {
		http.Redirect(w, r, "/not_allowed", 302)
		return
	}

	handler := listUserAPI{session:s}
	handler.handle(w, r)
}

func EditUserAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	if !checkLogin(s) {
		if global.Conf.ServerLocalMode {
			http.Redirect(w, r, "/user_login", 302)
		} else {
			http.Redirect(w, r, "/user_login_auth", 302)
		}
		return
	}

	if !checkRight(s, MANAGER_RIGHT) {
		http.Redirect(w, r, "/not_allowed", 302)
		return
	}

	handler := editUserAPI{session:s}
	handler.handle(w, r)
}

func AddUserAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	if !checkLogin(s) {
		if global.Conf.ServerLocalMode {
			http.Redirect(w, r, "/user_login", 302)
		} else {
			http.Redirect(w, r, "/user_login_auth", 302)
		}
		return
	}

	if !checkRight(s, MANAGER_RIGHT) {
		http.Redirect(w, r, "/not_allowed", 302)
		return
	}

	handler := addUserAPI{session:s}
	handler.handle(w, r)
}

func DelUserAPIHandler(w http.ResponseWriter, r *http.Request) {
	s := global.Manager.SessionStart(w, r)
	if !checkLogin(s) {
		if global.Conf.ServerLocalMode {
			http.Redirect(w, r, "/user_login", 302)
		} else {
			http.Redirect(w, r, "/user_login_auth", 302)
		}
		return
	}

	if !checkRight(s, MANAGER_RIGHT) {
		http.Redirect(w, r, "/not_allowed", 302)
		return
	}

	handler := delUserAPI{session:s}
	handler.handle(w, r)
}