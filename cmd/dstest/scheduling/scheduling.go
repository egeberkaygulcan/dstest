package scheduling

import (
	"fmt"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type Scheduler interface {
	Init(config *config.Config)
	Reset()
	Shutdown()
	NextIteration()
	GetClientRequest() int
	Next([]*network.Message, []*faults.Fault, faults.FaultContext) SchedulerDecision
	ApplyFault(*faults.Fault) error
}

type DecisionType int

const (
	NoOp DecisionType = iota
	SendMessage
	InjectFault
)

func (dt DecisionType) String() string {
	switch dt {
	case NoOp:
		return "NoOp"
	case SendMessage:
		return "SendMessage"
	case InjectFault:
		return "InjectFault"
	default:
		return "Unknown"
	}
}

type SchedulerDecision struct {
	DecisionType DecisionType
	Index        int
}

type SchedulerType string

const (
	Random SchedulerType = "random"
	QL     SchedulerType = "ql"
	Pct    SchedulerType = "pct"
	Replay SchedulerType = "replay"
)

func NewScheduler(schedulerType SchedulerType) (Scheduler, error) {
	switch schedulerType {
	case Random:
		return new(RandomScheduler), nil
	case QL:
		return new(QLScheduler), nil
	case Pctcp:
		return new(PCT), nil
	case Replay:
		return new(ReplayScheduler), nil
	default:
		return nil, fmt.Errorf("Unknown scheduler type: %s", schedulerType)
	}
}
