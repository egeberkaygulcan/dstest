package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type Scheduler interface {
	Init()
	Reset()
	Shutdown()
	Next([]*network.Message, []*faults.Fault, faults.FaultContext) SchedulerDecision
	ApplyFault(*faults.Fault) error
}

type DecisionType int

const (
	NoOp DecisionType = iota
	SendMessage
	InjectFault
)

type SchedulerDecision struct {
	DecisionType DecisionType
	Index        int
}

type SchedulerType string

const (
	Random SchedulerType = "random"
)

func NewScheduler(schedulerType SchedulerType) Scheduler {
	switch schedulerType {
	case Random:
		return new(RandomScheduler)
	default:
		return new(RandomScheduler)
	}
}
