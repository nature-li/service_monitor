package baseServer

type IServer interface {
	Init() bool
	Stop() bool
	IsHealthy() (bool, error)
	CheckStartCondition() bool
	CheckStopCondition() bool
	RemoteStart() bool
	RemoteStop() bool
	GetDetail() string
}
