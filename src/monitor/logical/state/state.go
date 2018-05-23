package state

import "time"

type state struct {
	healthy bool
	changeTime   time.Time
	lastCheck    time.Time
}

func newState(healthy bool) *state {
	now := time.Now()
	return &state{
		healthy:healthy,
		changeTime:now,
		lastCheck:now,
	}
}

func (o *state) setState(healthy bool) {
	if healthy == o.healthy {
		o.lastCheck = time.Now()
		return
	}

	o.healthy = healthy
	o.changeTime = time.Now()
	o.lastCheck = o.changeTime
}