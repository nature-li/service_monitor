package service

import (
	"monitor/service/indexServer"
	"monitor/logical"
)

func LoadServices() ([]*logical.TaskUnit, error) {
	taskList := make([]*logical.TaskUnit, 0)

	idxServ := indexServer.NewIndexServer("127.0.0.1:1234")
	task := logical.NewTaskUnit(1, idxServ, 0)
	taskList = append(taskList, task)
	return taskList, nil
}