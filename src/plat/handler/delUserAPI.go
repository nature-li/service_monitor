package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"fmt"
	"session"
	"strings"
)

type delUserAPI struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`

	session session.Session
}

func (o *delUserAPI) handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "PARSE_FORM_ERROR")
		return
	}

	userIdList := r.Form.Get("user_id_list")
	if userIdList == "" {
		o.render(w, false, "USER_ID_LIST_EMPTY")
		return
	}

	db, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "OPEN_DB_FAILED")
	}
	defer db.Close()
	db.Exec("PRAGMA busy_timeout=30000")

	if !o.delUser(db, userIdList) {
		o.render(w, false, "DEL_USER_FAILED")
	}

	o.render(w, true, "SUCCESS")
}

func (o *delUserAPI) render(w http.ResponseWriter, success bool, msg string) {
	o.Success = success
	o.Msg = msg

	result, err := json.Marshal(o)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	w.Write(result)
}

func (o *delUserAPI) delUser(db *sql.DB, userIdList string) bool {
	idList := strings.Split(userIdList, ",")
	sqlPart := strings.Join(idList, ",")
	if len(sqlPart) == 0 {
		logger.Error("user_id_list is empty")
		return false
	}
	if len(sqlPart) > 0 {
		sqlPart = sqlPart[0:len(sqlPart) - 1]
	}

	querySQL := "DELETE FROM user_list WHERE id in (" + sqlPart + ")"
	logger.Info(querySQL)
	results, err := db.Exec(querySQL)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	affectRows, err := results.RowsAffected()
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	logger.Info(fmt.Sprintf("delete rows: %v", affectRows))
	return true
}
