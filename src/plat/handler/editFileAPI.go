package handler

import (
	"database/sql"
	"session"
	"net/http"
	"encoding/json"
	"strconv"
)

type editFileAPI struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	session session.Session
}

func (o *editFileAPI) queryFileList(queryId string) *tableRow {
	db, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer db.Close()
	db.Exec("PRAGMA busy_timeout=30000")

	querySql := "SELECT id,file_name,file_size,url_name,version,md5_value,user_email,user_name,desc,create_time,update_time,refer_link FROM file_list WHERE id = ?"
	rows, err := db.Query(querySql, queryId)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer rows.Close()

	row := &tableRow{}
	for rows.Next() {
		err = rows.Scan(&row.Id, &row.FileName, &row.FileSize, &row.UrlName, &row.Version, &row.Md5, &row.UserEmail, &row.UserName, &row.Desc, &row.createTime, &row.updateTime, &row.ReferLink)
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
		break
	}
	row.format()

	return row
}

func (o *editFileAPI) editFile(w http.ResponseWriter, r *http.Request) {
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
	fileVersion := r.Form.Get("file_version")
	fileReferLink := r.Form.Get("file_refer_link")
	fileDesc := r.Form.Get("file_desc")

	// 检测数据长度
	if len([]rune(fileVersion)) > MAX_VERSION_LEN {
		o.render(w, false, "FILE_VERSION_BIG")
		return
	}
	if len([]rune(fileReferLink)) > MAX_LINK_LEN {
		o.render(w, false, "REFER_LINK_BIG")
		return
	}
	if len([]rune(fileDesc)) > MAX_DESC_LEN {
		o.render(w, false, "FILE_DESC_BIG")
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
		o.render(w, false, "EDIT_DENIED")
		return
	}

	if o.editDB(db, fileId, fileVersion, fileDesc, fileReferLink) {
		o.render(w, true, "SUCCESS")
	} else {
		o.render(w, false, "EDIT_DB_FAILED")
	}
}

func (o *editFileAPI) checkModifyRight(db *sql.DB, fileId, userEmail string) bool {
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

func (o *editFileAPI) editDB(db *sql.DB, fileId, fileVersion, fileDesc, referLink string) bool  {
	querySql := "update file_list set version=?, desc=?, refer_link=? where id=?"
	logger.Info(querySql)
	rows, err := db.Exec(querySql, fileVersion, fileDesc, referLink, fileId)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	count, err := rows.RowsAffected()
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Infof("affected rows: %v", count)
	return true
}

func (o *editFileAPI) render(w http.ResponseWriter, success bool, msg string) {
	o.Success = success
	o.Msg = msg

	result, err := json.Marshal(o)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	w.Write(result)
}
