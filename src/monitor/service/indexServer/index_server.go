package indexServer

import (
)

type IndexServer struct {
}

func NewIndexServer() *IndexServer {
	return &IndexServer{
	}
}

func (o *IndexServer) IsHealthy() bool {
	return true
}

func (o *IndexServer) CheckStartCondition() bool {
	return true
}

func (o *IndexServer) CheckStopCondition() bool {
	return true
}

func (o *IndexServer) RemoteStart() bool {
	return true
}

func (o *IndexServer) RemoteStop() bool {
	return true
}
