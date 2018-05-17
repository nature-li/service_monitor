package handler

import (
	"net/http"
	"path/filepath"
	"os"
	"encoding/json"
	"time"
	"strconv"
	"io"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	_ "github.com/mattn/go-sqlite3"
	"session"
)


type uploadFileAPI struct {
	session session.Session
	RetCode int `json:"code"`
	Desc string `json:"desc"`
}

func (o *uploadFileAPI)handle(w http.ResponseWriter, r *http.Request) {
	// 检测文件大小
	r.Body = http.MaxBytesReader(w, r.Body, config.UploadMaxSize)
	if err := r.ParseMultipartForm(config.UploadMaxSize); err != nil {
		logger.Error(err.Error())
		o.render(w, http.StatusBadRequest, "FILE_TOO_BIG")
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Error(err.Error())
		o.render(w, http.StatusInternalServerError, "PARSE_FORM_FAILED")
		return
	}
	fileVersion := r.Form.Get("file_version")
	fileReferLink := r.Form.Get("file_refer_link")
	fileDesc := r.Form.Get("file_desc")

	// 检测数据长度
	if len([]rune(fileVersion)) > MAX_VERSION_LEN {
		o.render(w, http.StatusBadRequest, "FILE_VERSION_BIG")
		return
	}
	if len([]rune(fileReferLink)) > MAX_LINK_LEN {
		o.render(w, http.StatusBadRequest, "REFER_LINK_BIG")
		return
	}
	if len([]rune(fileDesc)) > MAX_DESC_LEN {
		o.render(w, http.StatusBadRequest, "FILE_DESC_BIG")
		return
	}

	// 打开旧文件
	srcFile, handler, err := r.FormFile("uploadFile")
	if err != nil {
		logger.Error(err.Error())
		o.render(w, http.StatusBadRequest, "INVALID_FILE")
		return
	}
	defer srcFile.Close()

	// 打开新文件
	urlName := o.getNewName()
	newPath := filepath.Join(config.UploadDataPath, urlName)
	dstFile, err := os.Create(newPath)
	if err != nil {
		logger.Error(err.Error())
		o.render(w, http.StatusInternalServerError, "CREATE_FILE_ERROR")
		return
	}
	defer dstFile.Close()

	// 计算MD5同时复制文件
	hash := md5.New()
	reader := io.TeeReader(srcFile, dstFile)
	if _, err = io.Copy(hash, reader); err != nil {
		logger.Error(err.Error())
		o.render(w, http.StatusBadRequest, "SAVE_FILE_ERROR")
		return
	}
	hexValue := hex.EncodeToString(hash.Sum(nil))

	// 插入数据库
	userName := o.session.Get("user_name")
	httpCode, desc := o.insertToDb(handler.Filename, handler.Size, urlName, fileVersion, hexValue, userName, fileDesc, fileReferLink)
	if httpCode != http.StatusOK {
		o.deleteFile(newPath)
		o.render(w, httpCode, desc)
		return
	}

	// 记录日志
	logger.Info("receive success: file_name=" + handler.Filename + ", url_name=" + urlName)

	// 返回成功
	o.render(w, http.StatusOK, "SUCCESS")
}

func (o *uploadFileAPI)deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

func (o *uploadFileAPI)render(w http.ResponseWriter, httpCode int, desc string) {
	w.WriteHeader(httpCode)

	o.RetCode = httpCode
	o.Desc = desc
	byteResult, err := json.Marshal(o)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	_, err = w.Write(byteResult)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (o *uploadFileAPI)getNewName() string {
	now := time.Now().UnixNano()
	nowStr := strconv.FormatInt(now, 10)
	return nowStr
}

func (o *uploadFileAPI)insertToDb(fileName string, fileSize int64, urlName, version, md5, userName, desc, referLink string) (int, string) {
	db, err := sql.Open("sqlite3", config.SqliteDbPath)
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, "OPEN_DB_FAILED"
	}
	defer db.Close()
	db.Exec("PRAGMA busy_timeout=30000")

	userEmail := o.session.Get("user_email")
	stmt, err := db.Prepare("INSERT INTO file_list(file_name, file_size, url_name, version, md5_value, user_email, user_name, desc, create_time, update_time, refer_link) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, "DB_PREPARE_FAILED"
	}
	defer stmt.Close()

	createTime := time.Now().Unix()
	updateTime := createTime
	_, err = stmt.Exec(fileName, fileSize, urlName, version, md5, userEmail, userName, desc, createTime, updateTime, referLink)
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, "DB_PREPARE_FAILED"
	}

	return http.StatusOK, "SUCCESS"
}