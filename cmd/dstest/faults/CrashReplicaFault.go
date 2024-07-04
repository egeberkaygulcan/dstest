package faults

import (
	"fmt"
)

type CrashReplicaFault struct {
	nodeId int
	BaseFault
}

type CrashReplicaFaultParams struct {
	node int
}

var _ Fault = (*CrashReplicaFault)(nil)

func NewCrashReplicaFault(params map[string]interface{}) (*CrashReplicaFault, error) {
	if _, ok := params["node"]; !ok {
		return nil, fmt.Errorf("node parameter is required")
	}

	parsedParams := &CrashReplicaFaultParams{
		node: params["node"].(int),
	}

	return &CrashReplicaFault{
		BaseFault: BaseFault{
			Precondition: &AlwaysEnabledPrecondition{},
			Behavior: &CrashReplicaBehavior{
				nodeId: parsedParams.node,
			},
		},
		nodeId: -1,
	}, nil
}
