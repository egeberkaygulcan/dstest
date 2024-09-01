package scheduling

import (
	"math/rand"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"sort"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type PCT struct {
	Scheduler
	Config                   *config.Config
	RequestQuota             int
	NumClientTypes           int
	ClientRequestProbability float64
	NetworkManager           *network.Manager
	Rand                     *rand.Rand
	Depth                    int
	InitialPriorities        []int
	PriorityChangePoints     []int
	NumPriorityChange        int //this is used to iterate through the InitialPriorities.
	NumPriorityChanged       int //this is used to iterate thtrough the priority change points.
	Step                     int
}

// assert PCT implements the Scheduler interface
var _ Scheduler = &PCT{}

func (s *PCT) Init(config *config.Config) {
	s.Config = config
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)
	s.Depth = config.SchedulerConfig.Params["d"].(int)
	s.NumPriorityChange = 0
	s.NumPriorityChanged = 0
	s.Step = 0
	s.PriorityChangePoints = make([]int, 0)
	s.InitialPriorities = make([]int, 0)
	for i := 1; i < s.Depth; i++ { 
		s.PriorityChangePoints = append(s.PriorityChangePoints, s.distinctRandomInteger(s.Config.SchedulerConfig.Steps, s.PriorityChangePoints))
	}
	//PriorityChangePoints is iterated and compared to the current step number (increasing sequence) - so the list should be sorted.
	sort.Ints(s.PriorityChangePoints)

	for i := 0; i < s.Config.ProcessConfig.NumReplicas; i++ {
		s.InitialPriorities = append(s.InitialPriorities, s.distinctRandomInteger(s.Config.ProcessConfig.NumReplicas, s.InitialPriorities))
	}
}

func (s *PCT) distinctRandomInteger(max int, arr []int) int {
	for {
		i := s.Rand.Intn(max)
		if !contains(i, arr) {
			return i
		}
	}
}

func contains(i int, s []int) bool {
	if len(s) == 0 {
		return false
	}

	for _, val := range s {
		if val == i {
			return true
		}
	}
	return false
}

func (s *PCT) NextIteration() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.NumPriorityChange = 0
	s.NumPriorityChanged = 0
	s.Step = 0
	s.PriorityChangePoints = make([]int, 0)
	s.InitialPriorities = make([]int, 0)
	for i := 1; i < s.Depth; i++ { 
		s.PriorityChangePoints = append(s.PriorityChangePoints, s.distinctRandomInteger(s.Config.SchedulerConfig.Steps, s.PriorityChangePoints))
	}
	sort.Ints(s.PriorityChangePoints)

	for i := 0; i < s.Config.ProcessConfig.NumReplicas; i++ {
		s.InitialPriorities = append(s.InitialPriorities, s.distinctRandomInteger(s.Config.ProcessConfig.NumReplicas, s.InitialPriorities))
	}
}

func (s *PCT) Reset() {
	s.RequestQuota = s.Config.SchedulerConfig.ClientRequests
	s.Rand = rand.New(rand.NewSource(int64(s.Config.SchedulerConfig.Seed)))
}
func (s *PCT) Shutdown() {

}

// Returns a random index from available messages
func (s *PCT) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	if s.NumPriorityChanged < (s.Depth-1) { 
		if s.Step == s.PriorityChangePoints[s.NumPriorityChanged] {
			s.NumPriorityChanged++
			s.NumPriorityChange  =  (s.NumPriorityChange +1 )  % ( s.Config.ProcessConfig.NumReplicas -1 )
		}
	}

	if len(messages) > 0 {
		s.Step++
		decision := -1
		for i, message := range messages {
			if message.Sender == s.InitialPriorities[s.NumPriorityChange] {
				decision = i
			}
		}

		if decision < 0 {
			//s.Step = s.Step - 1 //if step is not decremented, the step number does not match with the step number the test engine counts.
			return SchedulerDecision{
				DecisionType: NoOp,
			}
		}

		return SchedulerDecision{
			DecisionType: SendMessage,
			Index:        decision,
		}
	} else {
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}
}

func (s *PCT) GetClientRequest() int {
	if s.RequestQuota > 0 {
		r := s.Rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {
			s.RequestQuota--
			return s.Rand.Intn(s.NumClientTypes)
		}
	}
	return -1
}
