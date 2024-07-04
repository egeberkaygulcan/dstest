package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"math/rand"
)

type RandomScheduler struct {
	Scheduler
}

// confirm it satisfies the interface Scheduler
var _ Scheduler = (*RandomScheduler)(nil)

func (s *RandomScheduler) Init() {

}

func (s *RandomScheduler) Reset() {

}
func (s *RandomScheduler) Shutdown() {

}

// Next returns a random index from available messages
func (s *RandomScheduler) Next(messages []*network.Message, faults []*faults.Fault, ctx faults.FaultContext) SchedulerDecision {
	// Apply faults with a probability of 1%
	if rand.Float64() < 0.01 {
		return SchedulerDecision{
			DecisionType: InjectFault,
			Index:        rand.Intn(len(faults)),
		}
	}

	return SchedulerDecision{
		DecisionType: SendMessage,
		Index:        rand.Intn(len(messages)),
	}
}
