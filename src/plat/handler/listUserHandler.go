package handler

import (
	"net/http"
	"path/filepath"
	"html/template"
	"plat/global"
)

func ListUserHandler(w http.ResponseWriter, r *http.Request)  {
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

	t, err := template.ParseFiles(filepath.Join(global.Conf.PrivateTemplatePath, "html/list_user.html"))
	if err != nil {
		global.Logger.Error(err.Error())
	}

	data := newPageData(w, r, s)
	t.Execute(w, data)
}
