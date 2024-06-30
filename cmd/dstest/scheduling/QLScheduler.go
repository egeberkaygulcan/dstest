package scheduling

import (
	"encoding/json"
	"fmt"
	agentv1 "github.com/aunum/gold/pkg/v1/agent"
	"github.com/aunum/gold/pkg/v1/common"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling/ql"
	"github.com/segmentio/fasthash/fnv1a"
	"gorgonia.org/tensor"
)

type QLScheduler struct {
	Scheduler
	agent *ql.Agent
}

var DefaultQLSchedulerConfig = &ql.AgentConfig{
	Hyperparameters: &ql.Hyperparameters{
		Epsilon: common.NewConstantSchedule(0.1),
		Gamma:   0.7,
		Alpha:   0.3,
	},
	Base: agentv1.NewBase("Q"),
}

func (s *QLScheduler) Init() {
	s.agent = ql.NewAgent(DefaultQLSchedulerConfig, nil)
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
func (s *QLScheduler) Next(messages []*network.Message) int {
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
		return 0
	}

	fmt.Println("QL chose action:", action)
	return action
}
