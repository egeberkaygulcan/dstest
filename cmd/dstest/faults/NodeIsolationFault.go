package faults

import (
	"fmt"
)

type NodeIsolationFault struct {
	BaseFault
}

type NodeIsolationFaultParams struct {
	nodeId int
}

var _ Fault = (*NodeIsolationFault)(nil)

func NewNodeIsolationFault(params map[string]interface{}) (*NodeIsolationFault, error) {
	fmt.Println("Creating a new NodeIsolationFault")

	if _, ok := params["nodeId"]; !ok {
		return nil, fmt.Errorf("nodeId parameter is required")
	}

	parsedParams := &NodeIsolationFaultParams{
		nodeId: params["nodeId"].(int),
	}

	fmt.Println("Creating a new NodeIsolationFault with params: ", parsedParams)

	return &NodeIsolationFault{
		BaseFault: BaseFault{
			Precondition: &AlwaysEnabledPrecondition{},
			Behavior: &IsolateNodeBehavior{
				nodeId: parsedParams.nodeId,
			},
		},
	}, nil
}
