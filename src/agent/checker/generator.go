package checker

import "rpcs"

type Generator struct {
}

func (o *Generator) IsHealthy(req *rpcs.HealthyReq, reply *rpcs.HealthyReply) error {
	return nil
}

func (o *Generator) Stop(req *rpcs.StopReq, reply *rpcs.StopReply) error {
	return nil
}

func (o *Generator) Start(req *rpcs.StartReq, reply *rpcs.StartReply) error {
	return nil
}
