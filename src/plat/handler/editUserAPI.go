package handler

import (
	"session"
	"net/http"
	"encoding/json"
	"database/sql"
	"strconv"
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
		logger.Error(err.Error())
		return
	}

	userId := r.Form.Get("user_id")
	downloadRight := r.Form.Get("download_right")
	uploadRight := r.Form.Get("upload_right")
	managerRight := r.Form.Get("manager_right")

	var userRight int64 = 0
	if downloadRight == "true" {
		userRight |= DOWNLOAD_RIGHT
	}

	if uploadRight == "true" {
		userRight |= UPLOAD_RIGHT
	}

	if managerRight == "true" {
		userRight |= MANAGER_RIGHT
	}

	db, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, false, "OPEN_DB_ERROR", nil)
	}
	defer db.Close()
	db.Exec("PRAGMA busy_timeout=30000")

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
		logger.Error(err.Error())
		return
	}

	w.Write(result)
}

func (o *editUserAPI) updateRight(db *sql.DB, userId string, userRight int64) bool {
	querySQL := "UPDATE user_list SET user_right = ? WHERE id = ?"
	result, err := db.Exec(querySQL, userRight, userId)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	if rows == 1 {
		return true
	}

	return false
}

func (o *editUserAPI) queryUser(db *sql.DB, userId string) *jsonListUserAPI {
	querySQL := "SELECT id,user_email,user_right,create_time FROM user_list WHERE id=?"
	rows, err := db.Query(querySQL, userId)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var userEmail string
		var userRight string
		var when sql.NullInt64

		err = rows.Scan(&id, &userEmail, &userRight, &when)
		if err != nil {
			logger.Error(err.Error())
			return nil
		}

		digitRight, err := strconv.ParseInt(userRight, 10, 64)
		if err != nil {
			logger.Error(err.Error())
			return nil
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

		user := &jsonListUserAPI{
			Id:            id,
			UserEmail:     userEmail,
			DownloadRight: downloadRight,
			UploadRight:   uploadRight,
			ManagerRight:  managerRight,
		}
		user.setCreateTime(when)

		return user
	}

	return nil
}
