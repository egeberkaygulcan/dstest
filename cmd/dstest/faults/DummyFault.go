package faults

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults/behavior"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults/trigger"
)

type DummyFault struct {
	BaseFault
}

var _ Fault = (*DummyFault)(nil)

func NewDummyFault() *DummyFault {
	return &DummyFault{
		BaseFault: BaseFault{
			FaultTrigger:  &trigger.DummyFaultTrigger{},
			FaultBehavior: &behavior.DummyFaultyBehavior{},
		},
	}
}
