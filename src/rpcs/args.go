package rpcs

type ServiceInfo struct {
	Id           int32
	IpAddress    string
	DomainName   string
	InstallPath  string
	LogPath      string
	RunUser      string
	ListenPort   int32
	PauseFlag    int32
	ZkAddress    string
	ZkNode       string
	AgentAddress string
}

type ServiceType struct {
	ServiceType int32
	ServiceName string
	StartCmd    string
	StopCmd     string
	RestartCmd  string
}

type HealthyReq struct {
}

type HealthyReply struct {
	Health bool
	Msg    string
}

type StopReq struct {
}

type StopReply struct {
}

type StartReq struct {
}

type StartReply struct {
}
