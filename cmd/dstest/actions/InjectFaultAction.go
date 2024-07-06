package actions

import "github.com/egeberkaygulcan/dstest/cmd/dstest/faults"

type InjectFaultAction struct {
	Fault faults.Fault
}

// make sure InjectFaultAction implements the Action interface
var _ Action = (*InjectFaultAction)(nil)

func (ifa *InjectFaultAction) GetType() ActionType {
	return InjectFault
}
