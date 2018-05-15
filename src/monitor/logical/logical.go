package logical

import (
	"time"
	"monitor/global"
	"fmt"
)

type Logical struct {
	serviceDict map[int32]*taskUnit
}

func NewLogical() *Logical {
	return &Logical{
		serviceDict: make(map[int32]*taskUnit),
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

func (o *Logical) run(task *taskUnit) bool {
	defer task.wg.Done()

	for {
		// do the job
		o.doJob(task)

		// loop once every 5 seconds
		select {
		case <-task.quit:
			break

		case time.After(time.Second * 5):
			break
		}
	}
	return true
}

func (o *Logical) doJob(task *taskUnit) {
	service := task.server

	// check service healthy
	if service.IsHealthy() {
		global.Logger.Infof("service[%s] is healthy", service.GetDescription())
		return
	}

	// record message
	content := fmt.Sprintf("service[%s] is down", service.GetDescription())
	global.Logger.Infof(content)
	global.SendMail(global.MailReceivers, "Service Monitor", content)

	// check start condition
	if !service.CheckStartCondition() {
		content = fmt.Sprintf("service[%s] check condition failed", service.GetDescription())
		global.Logger.Info(content)
		global.SendMail(global.MailReceivers, "Service Monitor", content)
		return
	}

	// start service
	if !service.RemoteStart() {
		content = fmt.Sprintf("start service[%s] failed", service.GetDescription())
		global.Logger.Infof(content)
		global.SendMail(global.MailReceivers, "Service Monitor", content)
		return
	}

	// record success message
	content = fmt.Sprintf("start service[%s] success", service.GetDescription())
	global.Logger.Infof(content)
	global.SendMail(global.MailReceivers, "Service Monitor", content)
}
