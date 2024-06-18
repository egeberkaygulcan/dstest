package scheduling

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"math/rand"
)

type RandomScheduler struct {
	Scheduler
}

func (s *RandomScheduler) Init() {

}

func (s *RandomScheduler) Reset() {

}
func (s *RandomScheduler) Shutdown() {

}

// Returns a random index from available messages
func (s *RandomScheduler) Next(messages []*network.Message) int {
	return rand.Intn(len(messages))
}