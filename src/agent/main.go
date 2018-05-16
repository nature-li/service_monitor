package main

import (
	"net/rpc"
	"net"
	"agent/config"
	"agent/global"
	"flag"
	"fmt"
	"mt/mtlog"
	"net/http"
	"agent/checker"
)

func main() {
	// parse config
	var confPath = flag.String("conf", "", "config file path")
	flag.Parse()
	if *confPath == "" {
		fmt.Println("conf is empty")
		return
	}

	global.Conf = &config.Conf{}
	if global.Conf.GetConf(*confPath) == nil {
		fmt.Println("parse config file error")
		return
	}
	global.Conf.Show()

	// init logger
	global.Logger = mtlog.NewLogger(false, mtlog.DEVELOP, mtlog.Level(global.Conf.LogLevel), global.Conf.LogPath, global.Conf.LogName, global.Conf.LogFileSize, global.Conf.LogFileCount)
	if !global.Logger.Start() {
		fmt.Println("logger.Start failed")
	}
	defer global.Logger.Stop()

	// register function
	rpc.Register(new(checker.IndexServer))
	rpc.Register(new(checker.Generator))
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", global.Conf.HttpListenPort)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	global.Logger.Info("agent server is starting...")
	http.Serve(listener, nil)
}
