package global

import (
	"mt/mtlog"
	"mt/session"
	"platform/config"
)

var Logger *mtlog.Logger = nil
var Manager session.Manager = nil
var Conf *config.Conf = nil
