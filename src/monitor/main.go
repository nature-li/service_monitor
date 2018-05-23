package main

import (
	"fmt"
	"mt/mtlog"
	"flag"
	"monitor/global"
	"net/http"
	"monitor/config"
	"monitor/task"
	"monitor/handler"
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

	// run logic task
	m := task.NewMonitor()
	m.Start()
	defer m.Stop()

	// http listening
	global.Logger.Infof("monitor is listening...")
	http.HandleFunc("/", handler.MainHandler)
	err := http.ListenAndServe(global.Conf.HttpListenPort, nil)
	if err != nil {
		global.Logger.Error(err.Error())
	}
}
