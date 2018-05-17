package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"session"
)

type listFileAPI struct {
	session session.Session
	Success   string            `json:"success"`
	ItemCount int               `json:"item_count"`
	Content   []jsonListFileAPI `json:"content"`
}

func (o *listFileAPI) handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Error(err.Error())
		o.render(w, "false", 0, nil)
		return
	}

	fileName := r.Form.Get("file_name")
	offset := r.Form.Get("off_set")
	limit := r.Form.Get("limit")

	logger.Info(r.Form.Encode())

	db, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, "false", 0, nil)
	}
	defer db.Close()
	db.Exec("PRAGMA busy_timeout=30000")

	code, totalCount := o.countDB(db, fileName)
	if code != http.StatusOK {
		o.render(w, "false", 0, nil)
		return
	}

	code, _, rows := o.queryDB(db, fileName, limit, offset)
	if code != http.StatusOK {
		o.render(w, "false", 0, nil)
		return
	}

	o.render(w, "true", totalCount, rows)
}

func (o *listFileAPI) render(w http.ResponseWriter, success string, itemCount int, content []jsonListFileAPI) {
	w.WriteHeader(http.StatusOK)

	o.Success = success
	o.ItemCount = itemCount
	o.Content = content
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

func (o *listFileAPI) countDB(db *sql.DB, fileName string) (code int, totalCount int) {
	var err error
	var querySql = "SELECT COUNT(1) AS COUNT FROM file_list"
	var rows *sql.Rows
	if len(fileName) != 0 {
		querySql += " WHERE file_name like ?"
		logger.Info(querySql)
		rows, err = db.Query(querySql, "%"+fileName+"%")
	} else {
		logger.Info(querySql)
		rows, err = db.Query(querySql, fileName)
	}
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, 0
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&totalCount)
	}

	return http.StatusOK, totalCount
}

func (o *listFileAPI) queryDB(db *sql.DB, fileName, limit, offset string) (int, string, []jsonListFileAPI) {
	var err error
	var dataSql = "SELECT id,file_name,file_size,url_name,version,md5_value,create_time FROM file_list"
	if len(fileName) != 0 {
		dataSql += " WHERE file_name like ? order by create_time desc limit ? offset ?"
	} else {
		dataSql += " order by create_time desc limit ? offset ?"
	}

	var rows *sql.Rows
	if len(fileName) != 0 {
		logger.Info(dataSql)
		rows, err = db.Query(dataSql, "%"+fileName+"%", limit, offset)
	} else {
		logger.Info(dataSql)
		rows, err = db.Query(dataSql, limit, offset)
	}
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
	}
	defer rows.Close()

	var rowList []jsonListFileAPI
	for rows.Next() {
		row := jsonListFileAPI{}
		err = rows.Scan(&row.Id, &row.FileName, &row.FileSize, &row.urlName, &row.Version, &row.Md5, &row.createTime)
		if err != nil {
			logger.Error(err.Error())
			return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
		}

		row.format()
		rowList = append(rowList, row)
	}

	return http.StatusOK, "SUCCESS", rowList
}
