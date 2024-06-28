package faults

import (
	"fmt"
)

type CrashReplicaFault struct {
	nodeId int
	BaseFault
}

type CrashReplicaFaultParams struct {
	nodeId int
}

var _ Fault = (*CrashReplicaFault)(nil)

func NewCrashReplicaFault(context FaultContext, params map[string]interface{}) (*CrashReplicaFault, error) {
	if _, ok := params["nodeId"]; !ok {
		return nil, fmt.Errorf("nodeId parameter is required")
	}

	parsedParams := &CrashReplicaFaultParams{
		nodeId: params["nodeId"].(int),
	}

	return &CrashReplicaFault{
		BaseFault: BaseFault{
			Precondition: &AlwaysEnabledPrecondition{},
			Behavior: &CrashReplicaBehavior{
				nodeId: parsedParams.nodeId,
			},
		},
		nodeId: -1,
	}, nil
}
