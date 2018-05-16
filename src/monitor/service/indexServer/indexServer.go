package indexServer

import (
	"monitor/baseServer"
	"net/rpc"
	"monitor/global"
	"rpcs"
)

type IndexServer struct {
	*baseServer.BaseServer
}

func NewIndexServer(agentAddress string) *IndexServer {
	return &IndexServer{
		baseServer.NewBaseServer(agentAddress),
	}
}

func (o *IndexServer) Init() bool {
	global.Logger.Debugf("rpc remote address: %s", o.AgentAddress)
	var err error
	o.Client, err = rpc.DialHTTP("tcp", o.AgentAddress)
	if err != nil {
		global.Logger.Error(err.Error())
		return false
	}
	return true
}

func (o *IndexServer) Stop() bool {
	if o.Client != nil {
		o.Client.Close()
	}
	return true
}

func (o *IndexServer) IsHealthy() (bool, error) {
	req := new(rpcs.HealthyReq)
	reply := new(rpcs.HealthyReply)
	err := o.Client.Call("IndexServer.IsHealthy", req, reply)
	if err != nil {
		global.Logger.Error(err.Error())
		return false, err
	}

	if !reply.Health {
		global.Logger.Errorf("Not healthy: %s", reply.Msg)
		return false, nil
	}

	return true, nil
}

func (o *IndexServer) RemoteStart() bool {
	return true
}

func (o *IndexServer) RemoteStop() bool {
	return true
}

func (o *IndexServer) GetDetail() string {
	return ""
}
