package global

import (
	"net/url"
	"net/http"
	"io/ioutil"
)

func SendMail(receivers string, subject, content string) bool {
	if len(receivers) == 0 {
		return false
	}

	formValue := url.Values{
		"receiver": {receivers},
		"subject": {subject},
		"content": {content},
	}

	httpUrl := "http://notice.ops.m.com/send_mail"
	resp, err := http.PostForm(httpUrl, formValue)
	if err != nil {
		Logger.Error(err.Error())
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Logger.Error(err.Error())
			return false
		}

		Logger.Error(string(body))
		return false
	}

	return true
}
