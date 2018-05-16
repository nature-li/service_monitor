package checker

import (
	"rpcs"
	"agent/global"
	"reflect"
)

type IndexServer struct {
}

func (o *IndexServer) IsHealthy(req *rpcs.HealthyReq, reply *rpcs.HealthyReply) error {
	global.Logger.Debugf("checking %s healthy", reflect.TypeOf(*o))
	reply.Health = true
	reply.Msg = "check success"
	return nil
}

func (o *IndexServer) Stop(req *rpcs.StopReq, reply *rpcs.StopReply) error {
	return nil
}

func (o *IndexServer) Start(req *rpcs.StartReq, reply *rpcs.StartReply) error {
	return nil
}


