package baseServer

import (
	"net/rpc"
	"monitor/global"
	"rpcs"
)

type BaseServer struct {
	*rpcs.ServiceInfo
	*rpcs.ServiceType
	Client       *rpc.Client
	StartRely    []IServer
	StopRely     []IServer
}

func NewBaseServer(agentAddress string) *BaseServer {
	server :=  &BaseServer{
		Client:       nil,
		StartRely:    make([]IServer, 0),
		StopRely:     make([]IServer, 0),
	}
	server.AgentAddress = agentAddress
	return server
}

func (o *BaseServer) AddStartReply(server IServer) {
	o.StartRely = append(o.StartRely, server)
}

func (o *BaseServer) AddStopReply(server IServer) {
	o.StopRely = append(o.StopRely, server)
}

func (o *BaseServer) CheckStartCondition() bool {
	for _, s := range o.StartRely {
		health, err := s.IsHealthy()
		if err != nil {
			global.Logger.Errorf("check service[%s] for failed", s.GetDetail())
			continue
		}

		if !health {
			return false
		}
	}

	return true
}

func (o *BaseServer) CheckStopCondition() bool {
	for _, s := range o.StopRely {
		health, err := s.IsHealthy()
		if err != nil {
			global.Logger.Errorf("check service[%s] for failed", s.GetDetail())
			continue
		}

		if health {
			return false
		}
	}

	return true
}
