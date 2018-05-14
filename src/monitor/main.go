package main

import (
	"fmt"
	"mt/mtlog"
	"flag"
	"monitor/config"
	"net/http"
)

var conf config.Conf
var logger *mtlog.Logger

func main()  {
	// parse config
	var confPath = flag.String("conf", "", "config file path")
	flag.Parse()
	if *confPath == "" {
		fmt.Println("conf is empty")
		return
	}
	if conf.GetConf(*confPath) == nil {
		fmt.Println("parse config file error")
		return
	}
	conf.Show()

	// init logger
	logger = mtlog.NewLogger(false, mtlog.DEVELOP, mtlog.INFO, conf.LogPath, conf.LogName, conf.LogFileSize, conf.LogFileCount)
	if !logger.Start() {
		fmt.Println("logger.Start failed")
	}

	// start http server
	err := http.ListenAndServe(conf.HttpListenPort, nil)
	if err != nil {
		logger.Error(err.Error())
	}

	// stop logger
	logger.Stop()
}
