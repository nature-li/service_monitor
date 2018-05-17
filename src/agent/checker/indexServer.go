package checker

import (
	"rpcs"
	"agent/global"
	"reflect"
)

type IndexServer struct {
}

func (o *IndexServer) Stop(req *rpcs.StopReq, reply *rpcs.StopReply) error {
	return nil
}

func (o *IndexServer) Start(req *rpcs.StartReq, reply *rpcs.StartReply) error {
	return nil
}


