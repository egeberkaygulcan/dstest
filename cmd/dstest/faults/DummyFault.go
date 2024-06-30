package faults

import (
	"fmt"
)

type DummyFault struct {
	BaseFault
}

type DummyFaultParams struct {
}

var _ Fault = (*DummyFault)(nil)

func NewDummyFault(params map[string]interface{}) (*DummyFault, error) {

	parsedParams := &DummyFaultParams{}

	fmt.Println("Creating a new DummyFault with params: ", parsedParams)

	return &DummyFault{
		BaseFault: BaseFault{
			Precondition: &DummyPrecondition{},
			Behavior:     &DummyBehavior{},
		},
	}, nil
}
