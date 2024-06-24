package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling/ql"
	"math/rand"
)

type QLScheduler struct {
	Scheduler
	agent *ql.Agent
}

func (s *QLScheduler) Init() {
	s.agent = ql.NewAgent(ql.DefaultAgentConfig, nil)
}

func (s *QLScheduler) Reset() {

}
func (s *QLScheduler) Shutdown() {

}

// Returns a random index from available messages
func (s *QLScheduler) Next(messages []*network.Message) int {
	return rand.Intn(len(messages))
}
