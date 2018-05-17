package logical

import (
	"time"
	"monitor/global"
	"fmt"
)

type Logical struct {
	serviceDict map[int32]*TaskUnit
}

func NewLogical() *Logical {
	return &Logical{
		serviceDict: make(map[int32]*TaskUnit),
	}
}

func (o *Logical) SetServiceDict(services []*TaskUnit) {
	for _, task := range services {
		o.serviceDict[task.id] = task
	}
}

func (o *Logical) Start() {
	for _, task := range o.serviceDict {
		if task.pause != 1 {
			task.wg.Add(1)
			go o.run(task)
		}
	}
}

func (o *Logical) Stop() {
	for _, task := range o.serviceDict {
		if task.pause != 1 {
			task.quit <- 1
		}
	}

	for _, task := range o.serviceDict {
		if task.pause != 1 {
			task.wg.Wait()
		}
		close(task.quit)
	}
}

func (o *Logical) run(task *TaskUnit) bool {
	defer task.wg.Done()

	quit := false
	timer := time.NewTimer(time.Second * 5)
	for !quit{
		// do the job
		o.doJob(task)

		// loop once every 5 seconds
		select {
		case <-task.quit:
			quit = true
			break

		case <-timer.C:
			timer.Reset(time.Second * 1)
		}
	}
	return true
}

func (o *Logical) doJob(task *TaskUnit) {
	service := task.server

	// init
	if !service.Init() {
		global.Logger.Errorf("service[%s] init failed", service.GetDetail())
		return
	}
	defer service.Stop()

	// check service healthy
	health, err := service.IsHealthy()

	// check service health failed
	if err != nil {
		content := fmt.Sprintf("service[%s] check healthy failed", service.GetDetail())
		global.Logger.Errorf(content)
		global.SendMail(global.MailReceivers, "Service Monitor", content)
		return
	}

	// service is healthy
	if health {
		global.Logger.Infof("service[%s] is healthy", service.GetDetail())
		return
	}

	// record message
	content := fmt.Sprintf("service[%s] is down", service.GetDetail())
	global.Logger.Infof(content)
	global.SendMail(global.MailReceivers, "Service Monitor", content)

	// check start condition
	if !service.CheckStartCondition() {
		content = fmt.Sprintf("service[%s] check condition failed", service.GetDetail())
		global.Logger.Info(content)
		global.SendMail(global.MailReceivers, "Service Monitor", content)
		return
	}

	// start service
	if !service.RemoteStart() {
		content = fmt.Sprintf("start service[%s] failed", service.GetDetail())
		global.Logger.Infof(content)
		global.SendMail(global.MailReceivers, "Service Monitor", content)
		return
	}

	// record success message
	content = fmt.Sprintf("start service[%s] success", service.GetDetail())
	global.Logger.Infof(content)
	global.SendMail(global.MailReceivers, "Service Monitor", content)
}
