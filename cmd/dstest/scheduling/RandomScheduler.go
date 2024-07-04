package scheduling

import (
	"math/rand"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type RandomScheduler struct {
	Scheduler
	Config                   *config.Config
	RequestQuota             int
	NumClientTypes           int
	ClientRequestProbability float64
	Rand *rand.Rand
}

func (s *RandomScheduler) Init(config *config.Config) {
	s.Config = config
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}

func (s *RandomScheduler) NextIteration() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
}

func (s *RandomScheduler) Reset() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}
func (s *RandomScheduler) Shutdown() {

}

// Returns a random index from available messages
func (s *RandomScheduler) Next(messages []*network.Message) int {
	return s.Rand.Intn(len(messages))
}

func (s *RandomScheduler) GetClientRequest() int {
	if s.RequestQuota > 0 {
		r := s.Rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {
			s.RequestQuota--
			return s.Rand.Intn(s.NumClientTypes)
		}
	}
	return -1
}
