package state

import (
	"sync"
	"time"
	"sync/atomic"
	"errors"
	"fmt"
)

type Manager struct {
	stat map[int32]*state
	lock *sync.Mutex

	stop int32
	wg   sync.WaitGroup
}

func NewManager() *Manager {
	return &Manager{
		stat: make(map[int32]*state),
		lock: new(sync.Mutex),
	}
}

func (o *Manager) Start() {
	o.wg.Add(1)
	go o.clean()
}

func (o *Manager) Stop() {
	atomic.StoreInt32(&o.stop, 1)
	o.wg.Wait()
}

func (o *Manager) clean() {
	defer o.wg.Done()

	interval := 0
	for {
		time.Sleep(1)
		if atomic.LoadInt32(&o.stop) == 1 {
			break
		}

		// clean once every 60 seconds
		interval++
		if interval < 60 {
			continue
		}
		interval = 0
		o.cleanOnce()
	}
}

func (o *Manager) cleanOnce() {
	o.lock.Lock()
	defer o.lock.Unlock()

	now := time.Now()
	for id, s := range o.stat {
		duration := now.Sub(s.lastCheck)
		if duration.Seconds() > 24 * 3600 {
			delete(o.stat, id)
		}
	}
}

func (o *Manager) SetState(id int32, healthy bool) (changeTime time.Time) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if s, ok := o.stat[id]; ok {
		s.setState(healthy)
		changeTime = s.changeTime
		return
	}

	s := newState(healthy)
	o.stat[id] = s
	changeTime = s.changeTime
	return
}

func (o *Manager) GetState(id int32) (healthy bool, lastCheck time.Time, err error) {
	if s, ok := o.stat[id]; ok {
		healthy = s.healthy
		lastCheck = s.lastCheck
		return
	}

	err = errors.New(fmt.Sprintf("id:%d not exist", id))
	return
}

func (o *Manager) GetChangeTime(id int32) (changeTime time.Time, err error) {
	if s, ok := o.stat[id]; ok {
		changeTime = s.changeTime
		return
	}

	err = errors.New(fmt.Sprintf("id:%d not exist", id))
	return
}