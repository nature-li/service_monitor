package global

import (
	"mt/mtlog"
	"mt/session"
	"plat/config"
)

var Logger *mtlog.Logger = nil
var Manager session.Manager = nil
var Conf *config.Conf = nil
