package handler

import (
	"mt/session"
	"net/http"
	"encoding/json"
	"database/sql"
	"strconv"
	"platform/global"
	"fmt"
	"time"
	"net/url"
)

type editUserAPI struct {
	Success bool              `json:"success"`
	Msg     string            `json:"msg"`
	Content *jsonListUserAPI `json:"content"`

	session session.Session
}

func (o *editUserAPI) handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	userId := r.Form.Get("user_id")
	managerRight := r.Form.Get("manager_right")

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
		o.render(w, false, "OPEN_DB_ERROR", nil)
	}
	defer db.Close()

	if !o.updateRight(db, userId, int64(userRight)) {
		o.render(w, false, "UPDATE_DB_FAILED", nil)
		return
	}

	user := o.queryUser(db, userId)
	if user != nil {
		o.render(w, true, "ok", user)
	} else {
		o.render(w, false, "QUERY_DB_FAILED", nil)
	}
}

func (o *editUserAPI) render(w http.ResponseWriter, success bool, msg string, content *jsonListUserAPI) {
	o.Success = success
	o.Msg = msg
	o.Content = content

	result, err := json.Marshal(o)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	w.Write(result)
}

func (o *editUserAPI) updateRight(db *sql.DB, userId string, userRight int64) bool {
	querySQL := "UPDATE users SET user_right = ? WHERE id = ?"
	global.Logger.Info(querySQL)
	result, err := db.Exec(querySQL, userRight, userId)
	if err != nil {
		global.Logger.Error(err.Error())
		return false
	}

	_, err = result.RowsAffected()
	if err != nil {
		global.Logger.Error(err.Error())
		return false
	}

	return true
}

func (o *editUserAPI) queryUser(db *sql.DB, userId string) *jsonListUserAPI {
	querySQL := "SELECT id,user_email,user_right,create_time FROM users WHERE id=?"
	rows, err := db.Query(querySQL, userId)
	if err != nil {
		global.Logger.Error(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var userEmail string
		var userRight string
		var when time.Time

		err = rows.Scan(&id, &userEmail, &userRight, &when)
		if err != nil {
			global.Logger.Error(err.Error())
			return nil
		}

		digitRight, err := strconv.ParseInt(userRight, 10, 64)
		if err != nil {
			global.Logger.Error(err.Error())
			return nil
		}

		managerRight := false
		if (digitRight & MANAGER_RIGHT) != 0 {
			managerRight = true
		}

		user := &jsonListUserAPI{
			Id:            id,
			UserEmail:     userEmail,
			ManagerRight:  managerRight,
		}
		user.setCreateTime(when)

		return user
	}

	return nil
}
