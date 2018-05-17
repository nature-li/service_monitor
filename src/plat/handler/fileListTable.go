package handler

import (
	"time"
	"strconv"
)

type tableRow struct {
	Id            int    `json:"id"`
	FileName      string `json:"file_name"`
	RFileSize     string `json:"r_file_size"`
	FileSize      int64  `json:"file_size"`
	UrlName       string `json:"url_name"`
	FileUrl       string `json:"file_url"`
	Version       string `json:"version"`
	Md5           string `json:"md5_value"`
	UserEmail     string `json:"user_email"`
	UserName      string `json:"user_name"`
	CreateTimeFmt string `json:"create_time"`
	UpdateTimeFmt string `json:"update_time"`
	Desc          string `json:"all_desc"`
	ReferLink     string `json:"refer_link"`

	*pageData
	createTime string
	updateTime string
}

func (o *tableRow) format() {
	o.FileUrl = "/data/" + o.UrlName

	if o.FileSize > 1073741824 {
		size := float64(o.FileSize) / 1073741824
		o.RFileSize = strconv.FormatFloat(size, 'f', 2, 64) + "G"
	} else if o.FileSize > 1048576 {
		size := float64(o.FileSize) / 1048576
		o.RFileSize = strconv.FormatFloat(size, 'f', 2, 64) + "M"
	} else if o.FileSize > 1024 {
		size := float64(o.FileSize) / 1024
		o.RFileSize = strconv.FormatFloat(size, 'f', 2, 64) + "K"
	} else {
		o.RFileSize = strconv.FormatInt(o.FileSize, 10)
	}

	createTime, err := strconv.ParseInt(o.createTime, 10, 64)
	if err != nil {
		return
	}

	when := time.Unix(createTime, 0)
	o.CreateTimeFmt = when.Format("2006-01-02 15:04:05")

	updateTime, err := strconv.ParseInt(o.updateTime, 10, 64)
	if err != nil {
		return
	}

	when = time.Unix(updateTime, 0)
	o.UpdateTimeFmt = when.Format("2006-01-02 15:04:05")
}

type jsonListFileAPI struct {
	Id            int    `json:"id"`
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	FileUrl       string `json:"file_url"`
	Version       string `json:"version"`
	Md5           string `json:"md5_value"`
	CreateTimeFmt string `json:"create_time"`

	urlName    string
	createTime string
}

func (o *jsonListFileAPI) format() {
	o.FileUrl = "/data/" + o.urlName

	createTime, err := strconv.ParseInt(o.createTime, 10, 64)
	if err != nil {
		return
	}

	when := time.Unix(createTime, 0)
	o.CreateTimeFmt = when.Format("2006-01-02 15:04:05")
}
