package logical

import (
	"monitor/baseServer"
	"sync"
	"sync/atomic"
)

type taskUnit struct {
	server      baseServer.IServer
	runningStat int32
	pauseStat   int32
	pauseFlag   int32
	stopChan    chan int32
	wg          sync.WaitGroup
}

func newTaskUnit(server baseServer.IServer) *taskUnit {
	return &taskUnit{
		server:      server,
		runningStat: 0,
		pauseStat:   0,
		pauseFlag:   0,
		stopChan:    make(chan int32, 1),
		wg:          sync.WaitGroup{},
	}
}

func (o* taskUnit)close() {
	close(o.stopChan)
}

func (o *taskUnit) setPause(pause bool) {
	if pause {
		atomic.StoreInt32(&o.pauseFlag, 1)
	} else {
		atomic.StoreInt32(&o.pauseFlag, 0)
	}
}

func (o *taskUnit) getPause() bool {
	return atomic.LoadInt32(&o.pauseFlag) == 1
}

func (o *taskUnit) setPaused(paused bool) {
	if paused {
		atomic.StoreInt32(&o.pauseStat, 1)
	} else {
		atomic.StoreInt32(&o.pauseStat, 0)
	}
}

func (o *taskUnit) isPaused() bool {
	return atomic.LoadInt32(&o.pauseStat) == 1
}

func (o *taskUnit) setRunning() {
	o.wg.Add(1)
	atomic.StoreInt32(&o.runningStat, 1)
}

func (o *taskUnit) UnsetRunning() {
	atomic.StoreInt32(&o.runningStat, 0)
	o.wg.Done()
}

func (o *taskUnit) isRunning() bool {
	return atomic.LoadInt32(&o.runningStat) == 1
}

func (o *taskUnit) setStop() {
	o.stopChan <- 1
}

func (o *taskUnit) waitStop() {
	o.wg.Wait()
}
