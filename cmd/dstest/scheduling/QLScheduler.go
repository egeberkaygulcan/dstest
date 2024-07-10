package scheduling

import (
	"encoding/json"
	"fmt"
	agentv1 "github.com/aunum/gold/pkg/v1/agent"
	"github.com/aunum/gold/pkg/v1/common"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling/ql"
	"github.com/segmentio/fasthash/fnv1a"
	"gorgonia.org/tensor"
	"math/rand"
)

type QLScheduler struct {
	Scheduler
	RequestQuota             int
	NumClientTypes           int
	agent                    *ql.Agent
	Rand                     *rand.Rand
	ClientRequestProbability float64
}

type QLSchedulerConfig struct {
	client_request_probability float64
	RequestQuota               int
	Epsilon                    float64
	Gamma                      float64
	Alpha                      float64
}

// assert QLScheduler implements the Scheduler interface
var _ Scheduler = &QLScheduler{}

var DefaultAlpha float32 = 0.3
var DefaultGamma float32 = 0.7
var DefaultEpsilon float32 = 0.1

func (s *QLScheduler) Init(config *config.Config) {
	s.Rand = rand.New(rand.NewSource(int64(config.SchedulerConfig.Seed)))
	s.RequestQuota = config.SchedulerConfig.ClientRequests
	s.NumClientTypes = len(config.ProcessConfig.ClientScripts)
	s.ClientRequestProbability = config.SchedulerConfig.Params["client_request_probability"].(float64)

	Epsilon := DefaultEpsilon
	if config.SchedulerConfig.Params["Epsilon"] != nil {
		Epsilon = config.SchedulerConfig.Params["Epsilon"].(float32)
	}

	Gamma := DefaultGamma
	if config.SchedulerConfig.Params["Gamma"] != nil {
		Gamma = config.SchedulerConfig.Params["Gamma"].(float32)
	}

	Alpha := DefaultAlpha
	if config.SchedulerConfig.Params["Alpha"] != nil {
		Alpha = config.SchedulerConfig.Params["Alpha"].(float32)
	}

	hyperparameters := &ql.AgentConfig{
		Hyperparameters: &ql.Hyperparameters{
			Epsilon: common.NewConstantSchedule(Epsilon),
			Gamma:   Gamma,
			Alpha:   Alpha,
		},
		Base: agentv1.NewBase("Q"),
	}

	fmt.Printf("QLScheduler: Epsilon: %f, Gamma: %f, Alpha: %f\n", Epsilon, Gamma, Alpha)

	s.agent = ql.NewAgent(hyperparameters, nil)
}

func (s *QLScheduler) Reset() {

}
func (s *QLScheduler) Shutdown() {

}

func StateHash(messages []*network.Message) uint32 {
	// convert each message to a hash
	hashes := make([]uint32, len(messages))
	for i, m := range messages {
		msgJson, _ := json.Marshal(m)
		hashes[i] = fnv1a.HashBytes32(msgJson)
	}

	// combine the hashes
	finalHash := fnv1a.HashUint32(hashes[0])
	for i := 1; i < len(hashes); i++ {
		finalHash = fnv1a.HashUint32(finalHash ^ hashes[i])
	}

	return finalHash
}

// Returns a random index from available messages
func (s *QLScheduler) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	// visualize the agent ??
	s.agent.Visualize()

	// create a tensor.Dense from the list of Messages
	values := make([]uint32, len(messages))
	for i, message := range messages {
		messageJson, _ := json.Marshal(*message)
		values[i] = fnv1a.HashBytes32(messageJson)
	}
	state := tensor.New(tensor.WithBacking(messages))

	// use the agent to get an action
	//state := StateHash(messages)
	action, err := s.agent.Action(state, messages)
	if err != nil {
		fmt.Println("ERROR GETTING ACTION:", err)
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}

	fmt.Println("QL chose action:", action)
	return SchedulerDecision{
		DecisionType: SendMessage,
		Index:        action,
	}
}

func (s *QLScheduler) GetClientRequest() int {
	if s.RequestQuota > 0 {
		r := s.Rand.Float64()
		if r <= s.ClientRequestProbability || s.ClientRequestProbability == 1.0 {
			s.RequestQuota--
			return s.Rand.Intn(s.NumClientTypes)
		}
	}
	return -1
}
