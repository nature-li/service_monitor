package handler

import (
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"session"
	"time"
	"strconv"
)

type jsonListUserAPI struct {
	Id            int    `json:"id"`
	UserEmail     string `json:"user_email"`
	DownloadRight bool   `json:"download_right"`
	UploadRight   bool   `json:"upload_right"`
	ManagerRight  bool   `json:"manager_right"`
	CreateTime    string `json:"create_time"`
}

func (o *jsonListUserAPI) setCreateTime(createTime sql.NullInt64) {
	var when time.Time
	if createTime.Valid {
		when = time.Unix(createTime.Int64, 0)
	} else {
		when = time.Unix(0, 0)
	}

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
		logger.Error(err.Error())
		o.render(w, "false", 0, nil)
		return
	}

	userEmail := r.Form.Get("user_email")
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
		logger.Error(err.Error())
		return
	}

	_, err = w.Write(result)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (o *listUserAPI) countDB(db *sql.DB, userEmail string) (code int, totalCount int) {
	var err error
	var querySql = "SELECT COUNT(1) AS COUNT FROM user_list"
	var rows *sql.Rows
	if len(userEmail) != 0 {
		querySql += " WHERE user_email like ?"
		logger.Info(querySql)
		rows, err = db.Query(querySql, "%"+userEmail+"%")
	} else {
		logger.Info(querySql)
		rows, err = db.Query(querySql, userEmail)
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

func (o *listUserAPI) queryDB(db *sql.DB, userEmail, limit, offset string) (int, string, []jsonListUserAPI) {
	var err error
	var dataSql = "SELECT id,user_email,user_right,create_time FROM user_list"
	if len(userEmail) != 0 {
		dataSql += " WHERE user_email like ? order by create_time desc limit ? offset ?"
	} else {
		dataSql += " order by create_time desc limit ? offset ?"
	}

	var rows *sql.Rows
	if len(userEmail) != 0 {
		logger.Info(dataSql)
		rows, err = db.Query(dataSql, "%"+userEmail+"%", limit, offset)
	} else {
		logger.Info(dataSql)
		rows, err = db.Query(dataSql, limit, offset)
	}
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
	}
	defer rows.Close()

	var rowList []jsonListUserAPI
	for rows.Next() {
		var id int
		var userEmail string
		var userRight string
		var when sql.NullInt64
		err = rows.Scan(&id, &userEmail, &userRight, &when)
		if err != nil {
			logger.Error(err.Error())
			return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
		}

		digitRight, err := strconv.ParseInt(userRight, 10, 64)
		if err != nil {
			logger.Error(err.Error())
			return http.StatusInternalServerError, "QUERY_DB_FAILED", nil
		}

		downloadRight := false
		if (digitRight & DOWNLOAD_RIGHT) != 0 {
			downloadRight = true
		}

		uploadRight := false
		if (digitRight & UPLOAD_RIGHT) != 0 {
			uploadRight = true
		}

		managerRight := false
		if (digitRight & MANAGER_RIGHT) != 0 {
			managerRight = true
		}

		row := jsonListUserAPI{
			Id:            id,
			UserEmail:     userEmail,
			DownloadRight: downloadRight,
			UploadRight:   uploadRight,
			ManagerRight:  managerRight,
		}
		row.setCreateTime(when)
		rowList = append(rowList, row)
	}

	return http.StatusOK, "SUCCESS", rowList
}
