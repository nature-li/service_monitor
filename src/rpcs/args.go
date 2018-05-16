package rpcs

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
