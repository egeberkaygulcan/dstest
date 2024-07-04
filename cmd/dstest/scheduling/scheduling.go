package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type Scheduler interface {
	Init()
	Reset()
	Shutdown()
	Next([]*network.Message) int
}

type SchedulerType string

const (
	Random SchedulerType = "random"
	QL     SchedulerType = "ql"
)

func NewScheduler(schedulerType SchedulerType) Scheduler {
	switch schedulerType {
	case Random:
		return new(RandomScheduler)
	case QL:
		return new(QLScheduler)
	default:
		return new(RandomScheduler)
	}
}
