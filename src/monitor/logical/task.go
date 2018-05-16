package logical

import (
	"monitor/baseServer"
	"sync"
)

type TaskUnit struct {
	id     int32
	server baseServer.IServer
	pause  int32
	quit   chan int32
	wg     *sync.WaitGroup
}

func NewTaskUnit(id int32, server baseServer.IServer, pause int32) *TaskUnit {
	return &TaskUnit{
		id:     id,
		server: server,
		pause:  pause,
		quit:   make(chan int32, 1),
		wg:     &sync.WaitGroup{},
	}
}