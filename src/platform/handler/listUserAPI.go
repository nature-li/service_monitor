package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"mt/session"
	"time"
	"strconv"
	"platform/global"
	"fmt"
	"net/url"
)

type jsonListUserAPI struct {
	Id            int    `json:"id"`
	UserEmail     string `json:"user_email"`
	ManagerRight  bool   `json:"manager_right"`
	CreateTime    string `json:"create_time"`
}

func (o *jsonListUserAPI) setCreateTime(when time.Time) {
	o.CreateTime = when.Format("2006-01-02 15:04:05")
}

type listUserAPI struct {
	session   session.Session
	Success   string            `json:"success"`
	ItemCount int               `json:"item_count"`
	Content   []jsonListUserAPI `json:"content"`
}

func (o *listUserAPI) handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		global.Logger.Error(err.Error())
		o.render(w, "false", 0, nil)
		return
	}

	userEmail := r.Form.Get("user_email")
	offset := r.Form.Get("off_set")
	limit := r.Form.Get("limit")

	global.Logger.Info(r.Form.Encode())

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
		o.render(w, "false", 0, nil)
	}
	defer db.Close()

	code, totalCount := o.countDB(db, userEmail)
	if code != http.StatusOK {
		o.render(w, "false", 0, nil)
		return
	}

	code, _, rows := o.queryDB(db, userEmail, limit, offset)
	if code != http.StatusOK {
		o.render(w, "false", 0, nil)
		return
	}

	o.render(w, "true", totalCount, rows)
}

func (o *listUserAPI) render(w http.ResponseWriter, success string, itemCount int, content []jsonListUserAPI) {
	w.WriteHeader(http.StatusOK)

	o.Success = success
	o.ItemCount = itemCount
	o.Content = content
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

func (o *listUserAPI) countDB(db *sql.DB, userEmail string) (code int, totalCount int) {
	var err error
	var querySql = "SELECT COUNT(1) AS COUNT FROM users"
	var rows *sql.Rows
	if len(userEmail) != 0 {
		querySql += " WHERE user_email like ?"
		global.Logger.Info(querySql)
		rows, err = db.Query(querySql, "%"+userEmail+"%")
	} else {
		global.Logger.Info(querySql)
		rows, err = db.Query(querySql)
	}
	if err != nil {
		global.Logger.Error(err.Error())
		return http.StatusInternalServerError, 0
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&totalCount)
	}

	return http.StatusOK, totalCount
}

func (o *listUserAPI) queryDB(db *sql.DB, userEmail, limit, offset string) (int, string, []jsonListUserAPI) {
	var err error
	var dataSql = "SELECT id,user_email,user_right,create_time FROM users"
	if len(userEmail) != 0 {
		dataSql += " WHERE user_email like ?"
	}
	dataSql += " order by create_time desc limit ? offset ?"

	var rows *sql.Rows
	if len(userEmail) != 0 {
		global.Logger.Info(dataSql)
		rows, err = db.Query(dataSql, "%"+userEmail+"%", limit, offset)
	} else {
		global.Logger.Info(dataSql)
		rows, err = db.Query(dataSql, limit, offset)
	}
	if err != nil {
		global.Logger.Error(err.Error())
		return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
	}
	defer rows.Close()

	var rowList []jsonListUserAPI
	for rows.Next() {
		var id int
		var userEmail sql.NullString
		var userRight sql.NullString
		var when time.Time
		err = rows.Scan(&id, &userEmail, &userRight, &when)
		if err != nil {
			global.Logger.Error(err.Error())
			return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
		}

		digitRight, err := strconv.ParseInt(userRight.String, 10, 64)
		if err != nil {
			global.Logger.Error(err.Error())
			return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
		}

		managerRight := false
		if (digitRight & MANAGER_RIGHT) != 0 {
			managerRight = true
		}

		row := jsonListUserAPI{
			Id:           id,
			UserEmail:    userEmail.String,
			ManagerRight: managerRight,
		}
		row.setCreateTime(when)
		rowList = append(rowList, row)
	}

	return http.StatusOK, "SUCCESS", rowList
}
