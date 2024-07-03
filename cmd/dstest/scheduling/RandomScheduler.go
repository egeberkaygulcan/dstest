package scheduling

import (
	"math/rand"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type RandomScheduler struct {
	Scheduler
	
	RequestQuota int
	NumClientTypes int
	ClientRequestProbability float64
}

func (s *RandomScheduler) Init(config *config.Config) {
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
}

func (s *RandomScheduler) Reset() {

}
func (s *RandomScheduler) Shutdown() {

}

// Returns a random index from available messages
func (s *RandomScheduler) Next(messages []*network.Message) int {
	return rand.Intn(len(messages))
}

func (s *RandomScheduler) GetClientRequest() int {
	if s.RequestQuota > 0 {
		r := rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {
			s.RequestQuota--
			return rand.Intn(s.NumClientTypes)
		}
	}
	return -1
}