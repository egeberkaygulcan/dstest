package scheduling

import (
	"encoding/json"
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/actions"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"os"
	"path/filepath"
	"strings"
)

type ReplayScheduler struct {
	Scheduler
	Config  *config.Config
	actions []actions.Action
	index   int
}

// assert ReplayScheduler implements the Scheduler interface
var _ Scheduler = &ReplayScheduler{}

type UnparsedAction struct {
	ActionType actions.ActionType
	Action     map[string]interface{}
}

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
	scheduleActions, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		fmt.Printf("Error reading actions file: %s\n", err)
		return
	}

	// read actions from file (one per line)
	actionStrings := strings.Split(strings.TrimSpace(string(scheduleActions)), "\n")

	// convert each string into an action
	for _, actionStr := range actionStrings {
		//fmt.Printf("ActionStr: %s\n", actionStr)
		var actionEntry UnparsedAction
		err := json.Unmarshal([]byte(actionStr), &actionEntry)
		if err != nil {
			fmt.Printf("Error unmarshalling action: %s\n", err)
			continue
		}

		parsedAction := actions.NewAction(actionEntry.ActionType, actionEntry.Action)
		if parsedAction == nil {
			fmt.Printf("Error parsing action\n")
			continue
		}

		// print parsed action
		//fmt.Printf("parsedAction: %s\n", parsedAction)

		s.actions = append(s.actions, parsedAction)
	}

	s.index = 0

	panic("at the disco")
}

func (s *ReplayScheduler) NextIteration() {}
func (s *ReplayScheduler) Reset()         {}
func (s *ReplayScheduler) Shutdown()      {}

// Returns a random index from available messages
func (s *ReplayScheduler) Next(messages []*network.Message, faults []*faults.Fault, context faults.FaultContext) SchedulerDecision {
	// if no more actions, return NoOp
	if s.index >= len(s.actions) {
		fmt.Printf("No more actions to schedule\n")
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}

	nextAction := s.actions[s.index]

	// Next action is a client request
	// FIXME: should this go into GetClientRequest?
	if nextAction.GetType() == actions.ClientRequest {
		// check if client requests are available
		// FIXME: is this correct?!
		if s.Config.SchedulerConfig.ClientRequests > 0 {
			s.Config.SchedulerConfig.ClientRequests--
			s.index++
			return SchedulerDecision{
				DecisionType: SendMessage,
				Index:        0,
			}
		}
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}

	// Next action is a message
	if nextAction.GetType() == actions.SendMessage {
		// search for the action with same sender and receiver
		for i, message := range messages {
			if message.Sender == nextAction.(*actions.DeliverMessageAction).Sender &&
				message.Receiver == nextAction.(*actions.DeliverMessageAction).Receiver &&
				message.Name == nextAction.(*actions.DeliverMessageAction).Name {
				s.index++
				return SchedulerDecision{
					DecisionType: SendMessage,
					Index:        i,
				}
			}
		}
		// if not found, return NoOp… maybe next time?
		return SchedulerDecision{
			DecisionType: NoOp,
		}
	}

	actionStr := s.actions[s.index]
	fmt.Printf("Selecting action: %s\n", actionStr)

	// parse action
	// check if action is available
	s.index++
	return SchedulerDecision{
		DecisionType: SendMessage,
		Index:        0,
	}

	// if not, return NoOp
	return SchedulerDecision{
		DecisionType: NoOp,
	}
}

func (s *ReplayScheduler) GetClientRequest() int {
	return -1
}
