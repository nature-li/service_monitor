package task

import (
	"fmt"
	"net/url"
	"database/sql"
	"monitor/global"
	"monitor/logical/job"
	"sync"
	"sync/atomic"
	"time"
	"monitor/logical"
)

type Monitor struct {
	wg   *sync.WaitGroup
	stop int32

	runner *logical.JobRunner
}

func NewMonitor() *Monitor {
	return &Monitor{
		wg:     new(sync.WaitGroup),
		stop:   0,
		runner: logical.NewJobRunner(),
	}
}

func (o *Monitor) Start() {
	o.wg.Add(1)
	go o.run()
}

func (o *Monitor) Stop() {
	atomic.StoreInt32(&o.stop, 1)
	o.runner.Stop()
	o.wg.Wait()
}

func (o *Monitor) run() {
	o.wg.Done()
	for {
		time.Sleep(time.Second * time.Duration(global.Conf.MonitorInterval))
		if atomic.LoadInt32(&o.stop) == 1 {
			break
		}

		o.monitorOnce()
	}
}

func (o *Monitor) monitorOnce() {
	serviceDict, err := o.loadJobFromDB()
	if err != nil {
		return
	}

	for _, aJob := range serviceDict {
		o.runner.AddJob(aJob)
	}
}

func (o *Monitor) loadJobFromDB() (map[int32]*job.Job, error) {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true",
		global.Conf.MysqlUser,
		global.Conf.MysqlPwd,
		global.Conf.MysqlAddress,
		global.Conf.MysqlPort,
		global.Conf.MysqlDbName,
		url.QueryEscape("Asia/Shanghai"))
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	defer db.Close()

	serviceDict, err := o.queryService(db)
	if err != nil {
		return nil, err
	}

	err = o.queryCheckCmd(db, serviceDict)
	if err != nil {
		return nil, err
	}

	err = o.queryRely(db, serviceDict)
	if err != nil {
		return nil, err
	}
	return serviceDict, nil
}

func (o *Monitor) queryService(db *sql.DB) (map[int32]*job.Job, error) {
	querySql := "SELECT id,service_name,ssh_user,ssh_ip,ssh_port,start_cmd,stop_cmd,auto_recover,mail_receiver " +
		"FROM services WHERE activate = 1"
	rows, err := db.Query(querySql)
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	serviceDict := make(map[int32]*job.Job)
	for rows.Next() {
		var id int32
		var serviceName string
		var sshUser string
		var sshIP string
		var sshPort string
		var startCmd string
		var stopCmd string
		var autoRecover int
		var mailReceiver string
		err = rows.Scan(&id, &serviceName, &sshUser, &sshIP, &sshPort, &startCmd, &stopCmd, &autoRecover, &mailReceiver)
		if err != nil {
			global.Logger.Error(err.Error())
			return nil, err
		}

		aJob := job.NewJob(id, serviceName, sshUser, sshIP, sshPort, startCmd, stopCmd, autoRecover, mailReceiver)
		serviceDict[aJob.Id()] = aJob
	}

	return serviceDict, nil
}

func (o *Monitor) queryCheckCmd(db *sql.DB, serviceDict map[int32]*job.Job) error {
	querySql := "SELECT id,service_id,local_check,check_shell,operator,check_value,good_match FROM check_cmd"
	rows, err := db.Query(querySql)
	if err != nil {
		global.Logger.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		var serviceId int32
		var localCheck int32
		var checkShell string
		var operator string
		var checkValue string
		var goodMatch int
		err = rows.Scan(&id, &serviceId, &localCheck, &checkShell, &operator, &checkValue, &goodMatch)
		if err != nil {
			global.Logger.Error(err.Error())
			return err
		}

		if s, ok := serviceDict[serviceId]; ok {
			checker := job.NewCheckCmd(id, serviceId, localCheck, checkShell, operator, checkValue, goodMatch)
			if localCheck == 1 {
				s.AddLocalCmd(checker)
			} else {
				s.AddRemoteCmd(checker)
			}
		}
	}

	return nil
}

func (o *Monitor) queryRely(db *sql.DB, serviceDict map[int32]*job.Job) error {
	querySql := "SELECT service_id,rely_id FROM service_rely"
	rows, err := db.Query(querySql)
	if err != nil {
		global.Logger.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var serviceId int32
		var relyId int32
		err = rows.Scan(&serviceId, relyId)
		if err != nil {
			global.Logger.Error(err.Error())
			return err
		}

		if s, ok := serviceDict[serviceId]; ok {
			s.AddRely(relyId)
		}
	}

	return nil
}
