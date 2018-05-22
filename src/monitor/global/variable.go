package global

import (
	"monitor/config"
	"mt/mtlog"
	"golang.org/x/crypto/ssh"
)

var Conf *config.Conf = nil
var Logger *mtlog.Logger = nil
var MailReceivers = "lyg@meitu.com"

