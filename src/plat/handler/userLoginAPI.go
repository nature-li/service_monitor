package handler

import (
	"net/http"
	"database/sql"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"strings"
	_ "session/cookie"
	"session"
)

type userLoginAPI struct {
	session session.Session
	Success bool `json:"success"`
	Msg string `json:"message"`

	db *sql.DB
	userName string
	userRight string
}

func (o *userLoginAPI) handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	userEmail := r.Form.Get("user_email")
	password := r.Form.Get("user_password")
	cat := r.Form.Get("captcha_value")

	if userEmail == "" {
		logger.Error("userName is empty")
		o.render(w, false, "用户无效")
		return
	}

	if password == "" {
		logger.Error("password is empty")
		o.render(w, false, "密码为空")
		return
	}

	if cat == "" {
		logger.Error("captcha is empty")
		o.render(w, false, "验证码空")
		return
	}

	if !strings.EqualFold(cat, o.session.Get("secret_captcha_value")) {
		logger.Error("captcha not match")
		o.render(w, false, "验证码错误")
		return
	}

	o.db, err = sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "内部错误")
		return
	}
	defer o.db.Close()
	o.db.Exec("PRAGMA busy_timeout=30000")

	success, message := o.checkPassword(userEmail, password)
	if success == true {
		o.session.Set("is_login", "1")
		o.session.Set("user_email", userEmail)
		o.session.Set("user_name", o.userName)
		o.session.Set("user_right", o.userRight)
	}
	o.render(w, success, message)
}

func (o *userLoginAPI) render(w http.ResponseWriter, success bool, desc string) {
	o.Success = success
	o.Msg = desc

	result, err := json.Marshal(o)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = w.Write(result)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (o *userLoginAPI) checkPassword(email, password string) (success bool, message string) {
	md5Value := md5.Sum([]byte(password))
	hexMd5 := hex.EncodeToString(md5Value[:])
	querySql := "SELECT user_email,user_name,user_right,passwd FROM user_list WHERE user_email=?"
	rows, err := o.db.Query(querySql, email)
	if err != nil {
		logger.Error(err.Error())
		return false, "内部错误"
	}

	var emailInDB string
	var passwordInDB string
	var count = 0
	for rows.Next() {
		err = rows.Scan(&emailInDB, &o.userName, &o.userRight, &passwordInDB)
		if err != nil {
			logger.Error(err.Error())
			return false, "内部错误"
		}
		count++
	}

	if count == 0 {
		return false, "用户无效"
	}

	if !strings.EqualFold(hexMd5, passwordInDB) {
		return false, "密码错误"
	}

	return true, "成功"
}
