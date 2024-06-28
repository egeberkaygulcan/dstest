package faults

import "fmt"

type CrashReplicaBehavior struct {
	nodeId int
}

var _ Behavior = (*CrashReplicaBehavior)(nil)

func NewCrashReplicaBehavior(nodeId int) *CrashReplicaBehavior {
	return &CrashReplicaBehavior{
		nodeId,
	}
}

func (fb *CrashReplicaBehavior) Apply(context FaultContext) error {
	context.GetProcessManager().CrashReplica(fb.nodeId)
	return nil
}

func (fb *CrashReplicaBehavior) String() string {
	return fmt.Sprintf("crash replica %d", fb.nodeId)
}
