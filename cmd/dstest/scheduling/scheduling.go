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
)

func NewScheduler(schedulerType SchedulerType) Scheduler {
	switch schedulerType {
	case Random:
		return new(RandomScheduler)
	default:
		return new(RandomScheduler)
	}
}
