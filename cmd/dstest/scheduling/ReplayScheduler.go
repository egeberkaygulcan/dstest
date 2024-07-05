package scheduling

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"os"
	"strings"
)

type ReplayScheduler struct {
	Scheduler
	Config  *config.Config
	actions []string
	index   int
}

// assert ReplayScheduler implements the Scheduler interface
var _ Scheduler = &ReplayScheduler{}

func (s *ReplayScheduler) Init(config *config.Config) {
	s.Config = config

	// read "filename" from config
	filename := config.SchedulerConfig.Params["filename"].(string)
	// if nil, abort
	if filename == "" {
		fmt.Printf("Error: filename not provided!\n")
		return
	}

	// check if file exists
	actions, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading actions file: %s\n", err)
		return
	}

	// read actions from file (one per line)
	s.actions = strings.Split(string(actions), "\n")
	s.index = 0

	// print actions
	fmt.Printf("Actions:\n%v\n", s.actions)
}

func (s *ReplayScheduler) NextIteration() {}
func (s *ReplayScheduler) Reset()         {}
func (s *ReplayScheduler) Shutdown()      {}

// Returns a random index from available messages
func (s *ReplayScheduler) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	if s.index < len(s.actions) {
		actionStr := s.actions[s.index]
		fmt.Printf("Selecting action: %s\n", actionStr)
		s.index++
		return SchedulerDecision{
			DecisionType: SendMessage,
			Index:        0,
		}
	} else {
		fmt.Printf("No more actions to schedule\n")
		return SchedulerDecision{
			DecisionType: NoOp,
		}

	}
}

func (s *ReplayScheduler) GetClientRequest() int {
	return -1
}
