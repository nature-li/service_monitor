package logical

import "sync"
import (
	"monitor/global"
	"errors"
)

type RunThread struct {
	jobDict map[int]*Job
	wg      sync.WaitGroup
	lock    *sync.Mutex
	stop    bool
}

func NewRunThread() *RunThread {
	return &RunThread{
		jobDict: make(map[int]*Job),
		lock:    new(sync.Mutex),
		stop:    false,
	}
}

func (o *RunThread) AddJob(job *Job) (bool, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	// if stopped do nothing
	if o.stop {
		return false, errors.New("RunThread is stopped")
	}

	// if job is running do nothing
	if _, ok := o.jobDict[job.id]; ok {
		return false, errors.New("last job is still running")
	}

	// add job to job dict and launch it
	o.jobDict[job.id] = job
	o.wg.Add(1)
	go o.run(job)

	return true, nil
}

func (o *RunThread) run(job *Job) {
	// do this when go-routine is ending
	defer o.wg.Done()

	// do job
	o.doJob(job)

	// remove job from job dict
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.jobDict, job.id)
}

func (o *RunThread) doJob(job *Job) {
	// check all policy
	healthy, err := job.checkAll()
	if err != nil {
		content := job.serviceName + ": " + err.Error()
		global.Logger.Error(content)
		global.SendMail(job.mailReceiver, global.AD_TECH_MONITOR, content)
		return
	}

	// output report log
	if healthy {
		content := job.serviceName + ": is healthy"
		global.Logger.Report(content)
		return
	}

	// send waning mail
	if !job.autoRecover {
		content := job.serviceName + ": is not healthy but auto recover is disable"
		global.SendMail(job.mailReceiver, global.AD_TECH_MONITOR, content)
		return
	}

	// try to restart unhealthy service
	content := job.serviceName + ": is not healthy, so trying to restart it"
	global.SendMail(job.mailReceiver, global.AD_TECH_MONITOR, content)
	err = job.restart()
	if err != nil {
		content := job.serviceName + ":" + err.Error()
		global.Logger.Error(content)
		global.SendMail(job.mailReceiver, global.AD_TECH_MONITOR, content)
	}
}

func (o *RunThread) Stop() {
	// set stop flag
	o.lock.Lock()
	o.stop = true
	o.lock.Unlock()

	// wait until all job is quited
	o.wg.Wait()
}
