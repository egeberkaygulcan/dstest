package faults

import (
	"fmt"
)

/**
The behavior is to call the CrashReplica method in process/process.go
Also make another fault for restarting a replica!! ;-)
*/

type CrashReplicaFault struct {
	BaseFault
}

type CrashReplicaFaultParams struct {
	name string
	age  int
}

var _ Fault = (*CrashReplicaFault)(nil)

func NewCrashReplicaFault(params map[string]interface{}) (*CrashReplicaFault, error) {
	fmt.Println("Creating a new CrashReplicaFault")

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

	fmt.Println("Creating a new CrashReplica with params: ", parsedParams)

	return &CrashReplicaFault{
		BaseFault: BaseFault{
			FaultTrigger:  &DummyFaultTrigger{},
			FaultBehavior: &DummyFaultyBehavior{},
		},
	}, nil
}
