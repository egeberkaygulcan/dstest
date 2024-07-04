package scheduling

import (
	"math/rand"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type PCTCP struct {
	Scheduler
	Config                   *config.Config
	RequestQuota             int
	NumClientTypes           int
	ClientRequestProbability float64
	NetworkManager *network.Manager
	Rand *rand.Rand
	Depth int
	PriorityChangePoints []int
}

func (s *PCTCP) Init(config *config.Config) {
	s.Config = config
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.NetworkManager = config.SchedulerConfig.Params["network_manager"].(*network.Manager)
	s.Depth = config.SchedulerConfig.Params["d"].(int)
	s.PriorityChangePoints = make([]int, s.Depth-1)
	for i := 0; i < s.Depth - 1; i++ {
		s.PriorityChangePoints[i] = s.DistinctRandomInteger(s.Config.SchedulerConfig.Steps)
	}
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}

func (s *PCTCP) DistinctRandomInteger(max int) int {
	for {
		i := s.Rand.Intn(max)
		if !contains(i, s.PriorityChangePoints) {
			return i
		}
	}
}

func contains(i int, s []int) bool {
	b := false

	for _, val := range s {
		if val == i {
			b = true
		}
	}
	return b
}

func (s *PCTCP) NextIteration() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
}

func (s *PCTCP) Reset() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}
func (s *PCTCP) Shutdown() {

}

// Returns a random index from available messages
func (s *PCTCP) Next(messages []*network.Message, iteration int) int {
	return s.Rand.Intn(len(messages))
}

func (s *PCTCP) GetClientRequest() int {
	if s.RequestQuota > 0 {
		r := s.Rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {
			s.RequestQuota--
			return s.Rand.Intn(s.NumClientTypes)
		}
	}
	return -1
}
