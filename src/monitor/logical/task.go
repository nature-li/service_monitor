package logical

import (
	"monitor/baseServer"
	"sync"
)

type taskUnit struct {
	id     int32
	server baseServer.IServer
	pause  int32
	quit   chan int32
	wg     *sync.WaitGroup
}

func newTaskUnit(id int32, server baseServer.IServer, pause int32) *taskUnit {
	return &taskUnit{
		id:     id,
		server: server,
		pause:  pause,
		quit:   make(chan int32, 1),
		wg:     &sync.WaitGroup{},
	}
}