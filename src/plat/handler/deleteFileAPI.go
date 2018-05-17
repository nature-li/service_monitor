package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	_"os"
	_"path/filepath"
	"path/filepath"
	"os"
	"session"
	"strconv"
)

type deleteFileAPI struct {
	session session.Session
	Success bool `json:"success"`
	Desc string `json:"desc"`
}

func (o *deleteFileAPI) handle(w http.ResponseWriter, r *http.Request) {
	userEmail := o.session.Get("user_email")
	if userEmail == "" {
		o.render(w, false, "USER_EMAIL_EMPTY")
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(r.Form.Encode())

	fileId := r.Form.Get("file_id")
	if fileId == "" {
		o.render(w, false, "FILE_ID_EMPTY")
		return
	}

	db, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "OPEN_DB_FAILED")
		return
	}
	defer db.Close()
	db.Exec("PRAGMA busy_timeout=30000")

	if !o.checkModifyRight(db, fileId, userEmail) {
		o.render(w, false, "DELETE_DENIED")
		return
	}

	if o.deleteFromDisk(w, db, fileId) {
		o.deleteFromDB(w, db, fileId)
	}
}

func (o *deleteFileAPI) render(w http.ResponseWriter, success bool, desc string) {
	o.Success = success
	o.Desc = desc

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

func (o *deleteFileAPI) checkModifyRight(db *sql.DB, fileId, userEmail string) bool {
	userRight := o.session.Get("user_right")
	digitRight, err := strconv.ParseInt(userRight, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	if (digitRight & MANAGER_RIGHT) != 0 {
		return true
	}

	querySQL := "select count(1) as count from file_list where user_email = ? and id = ?"
	logger.Info(querySQL)
	rows, err := db.Query(querySQL, userEmail, fileId)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			logger.Error(err.Error())
			return false
		}
	}
	if count == 0 {
		logger.Warnf("user_email=?, file_id=? does not exist", userEmail, fileId)
		return false
	}

	return true
}

func (o *deleteFileAPI) deleteFromDisk(w http.ResponseWriter, db *sql.DB, fileId string) bool {
	querySql := "SELECT url_name FROM file_list WHERE id = ?"
	rows, err := db.Query(querySql, fileId)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "QUERY_DB_FAILED")
		return false
	}
	defer rows.Close()

	var urlName string
	for rows.Next() {
		err = rows.Scan(&urlName)
		if err != nil {
			logger.Error(err.Error())
			o.render(w, false, "SCAN_DB_FAILED")
			return false
		}
		break
	}

	if urlName == "" {
		logger.Error("fileId" + fileId + " not exist")
		o.render(w, false, "FILE_NOT_EXISTS")
		return false
	}

	fullPath := filepath.Join(config.UploadDataPath, urlName)
	logger.Info("remove file: " + fullPath)
	err = os.Remove(fullPath)
	if err != nil {
		logger.Warn(err.Error())
	}

	return true
}

func (o *deleteFileAPI) deleteFromDB(w http.ResponseWriter, db *sql.DB, fileId string) bool {
	querySql := "DELETE FROM file_list WHERE id = ?"
	stmt, err := db.Prepare(querySql)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "DB_PREPARE_FAILED")
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileId)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "DB_EXEC_FAILED")
	}

	o.render(w, true, "SUCCESS")
	return true
}