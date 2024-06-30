package faults

import (
	"fmt"
)

type RestartReplicaFault struct {
	nodeId int
	BaseFault
}

type RestartReplicaFaultParams struct {
	node int
}

var _ Fault = (*RestartReplicaFault)(nil)

func NewRestartReplicaFault(params map[string]interface{}) (*RestartReplicaFault, error) {
	// print params
	fmt.Println("Creating a new RestartReplicaFault: ", params)

	if _, ok := params["node"]; !ok {
		return nil, fmt.Errorf("node parameter is required")
	}

	parsedParams := &RestartReplicaFaultParams{
		node: params["node"].(int),
	}

	return &RestartReplicaFault{
		BaseFault: BaseFault{
			Precondition: &AlwaysEnabledPrecondition{},
			Behavior: &RestartReplicaBehavior{
				nodeId: parsedParams.node,
			},
		},
		nodeId: -1,
	}, nil
}
