package handler

import (
	"net/http"
	"database/sql"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"strings"
	_ "mt/session/cookie"
	"mt/session"
	"platform/global"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type userLoginAPI struct {
	session session.Session
	Success bool   `json:"success"`
	Msg     string `json:"message"`

	db        *sql.DB
	userName  sql.NullString
	userRight sql.NullString
}

func (o *userLoginAPI) handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	userEmail := r.Form.Get("user_email")
	password := r.Form.Get("user_password")
	cat := r.Form.Get("captcha_value")

	if userEmail == "" {
		global.Logger.Error("userName is empty")
		o.render(w, false, "用户无效")
		return
	}

	if password == "" {
		global.Logger.Error("password is empty")
		o.render(w, false, "密码为空")
		return
	}

	if cat == "" {
		global.Logger.Error("captcha is empty")
		o.render(w, false, "验证码空")
		return
	}

	if !strings.EqualFold(cat, o.session.Get("secret_captcha_value")) {
		global.Logger.Error("captcha not match")
		o.render(w, false, "验证码错误")
		return
	}

	connectStr := fmt.Sprintf("%s:%s@/%s", global.Conf.MysqlUser, global.Conf.MysqlPwd, global.Conf.MysqlDbName)
	o.db, err = sql.Open("mysql", connectStr)
	if err != nil {
		global.Logger.Error(err.Error())
		o.render(w, false, "内部错误")
		return
	}
	defer o.db.Close()

	success, message := o.checkPassword(userEmail, password)
	if success == true {
		o.session.Set("is_login", "1")
		o.session.Set("user_email", userEmail)
		o.session.Set("user_name", o.userName.String)
		o.session.Set("user_right", o.userRight.String)
	}
	o.render(w, success, message)
}

func (o *userLoginAPI) render(w http.ResponseWriter, success bool, desc string) {
	o.Success = success
	o.Msg = desc

	result, err := json.Marshal(o)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	_, err = w.Write(result)
	if err != nil {
		global.Logger.Error(err.Error())
	}
}

func (o *userLoginAPI) checkPassword(email, password string) (success bool, message string) {
	md5Value := md5.Sum([]byte(password))
	hexMd5 := hex.EncodeToString(md5Value[:])
	querySql := "SELECT user_name,user_email,user_pwd,user_right FROM user_list WHERE user_email=? and user_type=1"
	rows, err := o.db.Query(querySql, email)
	if err != nil {
		global.Logger.Error(err.Error())
		return false, "内部错误"
	}

	var emailInDB sql.NullString
	var passwordInDB sql.NullString
	var count = 0
	for rows.Next() {
		err = rows.Scan(&o.userName, &emailInDB, &passwordInDB, &o.userRight)
		if err != nil {
			global.Logger.Error(err.Error())
			return false, "内部错误"
		}
		count++
	}

	if count == 0 {
		return false, "用户无效"
	}

	if !strings.EqualFold(hexMd5, passwordInDB.String) {
		return false, "密码错误"
	}

	return true, "成功"
}
