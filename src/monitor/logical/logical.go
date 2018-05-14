package logical

import (
	"time"
)

type Logical struct {
	serviceDict map[int32]*taskUnit
}

func NewLogical() *Logical {
	return &Logical{
		serviceDict: make(map[int32]*taskUnit),
	}
}

func (o *Logical) StartAll() {
	for _, task := range o.serviceDict {
		if task.getPause() {
			continue
		}

		go o.run(task)
	}
}

func (o *Logical) StopAll() {
	for _, task := range o.serviceDict {
		if task.isRunning() {
			task.setStop()
		}
	}

	for _, task := range o.serviceDict {
		if task.isRunning() {
			task.waitStop()
		}
	}
}

func (o *Logical) run(task *taskUnit) bool {
	task.setRunning()
	defer task.UnsetRunning()

	for {
		// do job if it's not paused
		if task.getPause() {
			task.setPaused(true)
		} else {
			task.setPaused(false)
			o.doJob(task)
		}

		// loop once every 5 seconds
		select {
		case <-task.stopChan:
			break

		case time.After(time.Second * 5):
			break
		}
	}
	return true
}

func (o *Logical) doJob(task *taskUnit) {

}