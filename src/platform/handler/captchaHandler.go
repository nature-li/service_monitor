package handler

import (
	"image/color"
	"github.com/afocus/captcha"
	"bytes"
	"image/jpeg"
	"encoding/base64"
	"net/http"
	"mt/session"
	"platform/global"
	"encoding/json"
	"path/filepath"
)

type captchaHandler struct {
	Success bool `json:"success"`
	Value string `json:"value"`

	session session.Session
	cap *captcha.Captcha
}

func newCaptchaHandler(s session.Session) *captchaHandler {
	cat := captcha.New()
	err := cat.AddFont(filepath.Join(global.Conf.PublicTemplatePath, "fonts/comic.ttf"))
	if err != nil {
		global.Logger.Error(err.Error())
		return nil
	}

	cat.SetSize(100, 40)
	cat.SetDisturbance(captcha.MEDIUM)
	cat.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cat.SetBkgColor(color.RGBA{0, 153, 0, 255})

	handler := &captchaHandler{cap: cat, session:s}
	return handler
}

func (o *captchaHandler) createCaptcha() (value, hex string, err error) {
	image, value := o.cap.Create(4, captcha.ALL)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, nil)
	if err != nil {
		global.Logger.Error(err.Error())
		return "", "", err
	}

	hex = base64.StdEncoding.EncodeToString(buf.Bytes())
	o.session.Set("secret_captcha_value", value)

	return value, hex, nil
}

func (o *captchaHandler) handle(w http.ResponseWriter, r *http.Request) {
	value, imgValue, err := o.createCaptcha()
	if err != nil {
		o.render(w, false, "CREATE_CAPTCHA_ERROR", "CREATE_CAPTCHA_ERROR")
		return
	}

	o.render(w, true, value, imgValue)
}

func (o *captchaHandler) render(w http.ResponseWriter, success bool, value, imgValue string) {
	o.Success = success
	o.Value = imgValue

	jsonValue, err := json.Marshal(&o)
	if err != nil {
		return
	}

	w.Write(jsonValue)
}