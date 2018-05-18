package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"time"
	"fmt"
	"mt/session"
	"platform/global"
	"net/url"
)

type addUserAPI struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`

	session session.Session
}

func (o *addUserAPI) handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		global.Logger.Error(err.Error())
		o.render(w, false, "PARSE_FORM_ERROR")
		return
	}

	userEmail := r.Form.Get("user_email")
	managerRight := r.Form.Get("manager_right")

	if userEmail == "" {
		global.Logger.Error("user_email is empty")
		o.render(w, false, "user_email is empty")
		return
	}

	var userRight int64 = 0
	if managerRight == "true" {
		userRight |= MANAGER_RIGHT
	}

	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true",
		global.Conf.MysqlUser,
		global.Conf.MysqlPwd,
		global.Conf.MysqlAddress,
		global.Conf.MysqlPort,
		global.Conf.MysqlDbName,
		url.QueryEscape("Asia/Shanghai"))
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		global.Logger.Error(err.Error())
		o.render(w, false, "OPEN_DB_FAILED")
	}
	defer db.Close()

	if !o.addUser(db, userEmail, userRight) {
		o.render(w, false, "ADD_USER_FAILED")
	}

	o.render(w, true, "SUCCESS")
}

func (o *addUserAPI) render(w http.ResponseWriter, success bool, msg string) {
	o.Success = success
	o.Msg = msg

	result, err := json.Marshal(o)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	w.Write(result)
}

func (o *addUserAPI) addUser(db *sql.DB, userEmail string, userRight int64) bool {
	now := time.Now().Local()
	querySQL := "INSERT INTO users(user_email, user_right, create_time) VALUES (?,?,?)"
	results, err := db.Exec(querySQL, userEmail, userRight, now)
	if err != nil {
		global.Logger.Error(err.Error())
		return false
	}

	affectRows, err := results.RowsAffected()
	if err != nil {
		global.Logger.Error(err.Error())
		return false
	}

	if affectRows != 1 {
		global.Logger.Error(fmt.Sprintf("rows_affected is %v", affectRows))
		return false
	}

	return true
}
