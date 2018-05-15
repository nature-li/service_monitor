package main

import (
	"fmt"
	"mt/mtlog"
	"flag"
	"monitor/global"
	"net/http"
	"monitor/config"
)


func main()  {
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
	global.Logger = mtlog.NewLogger(false, mtlog.DEVELOP, mtlog.INFO, global.Conf.LogPath, global.Conf.LogName, global.Conf.LogFileSize, global.Conf.LogFileCount)
	if !global.Logger.Start() {
		fmt.Println("logger.Start failed")
	}

	// start http server
	err := http.ListenAndServe(global.Conf.HttpListenPort, nil)
	if err != nil {
		global.Logger.Error(err.Error())
	}

	// stop logger
	global.Logger.Stop()
}
