package faults

import (
	"fmt"
)

type NodeIsolationFault struct {
	BaseFault
}

type NodeIsolationFaultParams struct {
	name string
	age  int
}

var _ Fault = (*NodeIsolationFault)(nil)

func NewNodeIsolationFault(params map[string]interface{}) (*NodeIsolationFault, error) {
	fmt.Println("Creating a new NodeIsolationFault")

	if _, ok := params["name"]; !ok {
		return nil, fmt.Errorf("name parameter is required")
	}

	if _, ok := params["age"]; !ok {
		return nil, fmt.Errorf("age parameter is required")
	}

	parsedParams := &NodeIsolationFaultParams{
		name: params["name"].(string),
		age:  params["age"].(int),
	}

	fmt.Println("Creating a new NodeIsolationFault with params: ", parsedParams)

	return &NodeIsolationFault{
		BaseFault: BaseFault{
			Precondition: &AlwaysEnabledPrecondition{},
			Behavior:     &DummyFaultyBehavior{},
		},
	}, nil
}
