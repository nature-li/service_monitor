package global

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"crypto/x509"
	"net"
	"errors"
	"encoding/pem"
)

func SendMail(receivers string, subject, content string) bool {
	if len(receivers) == 0 {
		return false
	}

	if Conf.SendMail == 0 {
		Logger.Infof("receivers=[%s], subject=[%s], content=[%s]", receivers, subject, content)
		return true
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

func decrypt(buffer []byte, pwd []byte) ([]byte, error) {
	block, rest := pem.Decode(buffer)
	if len(rest) > 0 {
		return nil, errors.New("extra data included in key")
	}
	der, err := x509.DecryptPEMBlock(block, pwd)
	if err != nil {
		return nil, err
	}
	return der, nil
}

func publicKeyFile(keyFile, keyPwd string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(keyFile)
	if err != nil {
		Logger.Error(err.Error())
		return nil
	}

	var signer ssh.Signer
	if keyPwd == "" {
		signer, err = ssh.ParsePrivateKey(buffer)
		if err != nil {
			Logger.Error(err.Error())
			return nil
		}
	} else {
		der, err := decrypt(buffer, []byte(keyPwd))
		if err != nil {
			return nil
		}

		key, err := x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			Logger.Error(err.Error())
			return nil
		}

		signer, err = ssh.NewSignerFromKey(key)
		if err != nil {
			Logger.Error(err.Error())
			return nil
		}
	}

	return ssh.PublicKeys(signer)
}

func GetSSHConfig(loginUser string) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: loginUser,
		Auth: []ssh.AuthMethod{
			publicKeyFile(Conf.KeyFile, Conf.KeyPwd),
		},
		HostKeyCallback:func(host string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return config
}

