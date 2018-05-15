package global

import (
	"monitor/config"
	"mt/mtlog"
)

var Conf *config.Conf = nil
var Logger *mtlog.Logger = nil
var MailReceivers = "lyg@meitu.com"
