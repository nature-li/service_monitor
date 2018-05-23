package logical

import "sync"
import (
	"errors"
	"sync/atomic"
	"monitor/global"
	"monitor/logical/state"
	"monitor/logical/job"
	"fmt"
)

type JobRunner struct {
	jobDict map[int32]*job.Job
	wg      sync.WaitGroup
	lock    *sync.Mutex
	stop    int32

	manager *state.Manager
}

func NewJobRunner() *JobRunner {
	return &JobRunner{
		jobDict: make(map[int32]*job.Job),
		lock:    new(sync.Mutex),
		stop:    0,
		manager: state.NewManager(),
	}
}

func (o *JobRunner) AddJob(job *job.Job) (bool, error) {
	if atomic.LoadInt32(&o.stop) == 1 {
		return false, errors.New("JobRunner is stopped")
	}

	if o.isJobRunning(job.Id()) {
		return false, errors.New("last job is still running")
	}

	o.jobDict[job.Id()] = job
	o.wg.Add(1)
	go o.run(job)
	return true, nil
}

func (o *JobRunner) isJobRunning(id int32) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	if _, ok := o.jobDict[id]; ok {
		return true
	}
	return false
}

func (o *JobRunner) run(job *job.Job) {
	defer o.wg.Done()

	o.doJob(job)

	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.jobDict, job.Id())
}

func (o *JobRunner) doJob(job *job.Job) {
	// check all policy
	healthy, err := job.CheckAll()
	if err != nil {
		content := job.ServiceName() + ": " + err.Error()
		global.Logger.Error(content)
		global.SendMail(job.MailReceiver(), global.AD_TECH_MONITOR, content)
		return
	}

	// set service's health state
	changeTime := o.manager.SetState(job.Id(), healthy)

	// output report log
	if healthy {
		content := job.ServiceName() + ": is healthy"
		global.Logger.Report(content)
		return
	}

	// send waning mail
	if !job.IsAutoRecover() {
		content := job.ServiceName() + ": is not healthy but auto recover is disable"
		global.SendMail(job.MailReceiver(), global.AD_TECH_MONITOR, content)
		return
	}

	// try to restart unhealthy service
	content := job.ServiceName() + ": is not healthy, so trying to restart it"
	global.SendMail(job.MailReceiver(), global.AD_TECH_MONITOR, content)

	// check if all rely services are healthy
	rely := job.Rely()
	for _, jobId := range rely {
		good, when, err := o.manager.GetState(jobId)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}

		if when.Before(changeTime) {
			global.Logger.Warn(fmt.Sprintf("%s is waiting jobId[%s] for the lastest check", job.ServiceName(), jobId))
			return
		}

		if !good {
			global.Logger.Warn(fmt.Sprintf("%s is waiting jobId[%s] to be healthy", job.ServiceName(), jobId))
			return
		}
	}

	err = job.Restart()
	if err != nil {
		content := job.ServiceName() + ":" + err.Error()
		global.Logger.Error(content)
		global.SendMail(job.MailReceiver(), global.AD_TECH_MONITOR, content)
	}
}

func (o *JobRunner) Stop() {
	atomic.StoreInt32(&o.stop, 1)
	o.wg.Wait()
}
