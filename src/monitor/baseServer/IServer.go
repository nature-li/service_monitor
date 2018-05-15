package baseServer

type IServer interface {
	IsHealthy() bool
	CheckStartCondition() bool
	CheckStopCondition() bool
	RemoteStart() bool
	RemoteStop() bool
	GetDescription() string
}
