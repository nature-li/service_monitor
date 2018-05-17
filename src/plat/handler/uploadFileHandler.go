package handler

import (
	"net/http"
	"html/template"
	"path/filepath"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 检测是否登录
	s := manager.SessionStart(w, r)
	if !checkLogin(s) {
		if config.ServerLocalMode {
			http.Redirect(w, r, "/user_login", 302)
		} else {
			http.Redirect(w, r, "/user_login_auth", 302)
		}
		return
	}

	// 检测上传权限
	if !checkRight(s, UPLOAD_RIGHT) {
		http.Redirect(w, r, "/not_allowed", 302)
		return
	}

	t, err := template.ParseFiles(filepath.Join(config.privateTemplatePath, "html/upload_file.html"))
	if err != nil {
		logger.Error(err.Error())
	}

	data := newPageData(w, r, s)
	t.Execute(w, data)
}