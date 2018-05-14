package iagent

type IAgent interface {
	LocalStart() bool
	LocalStop() bool
}
