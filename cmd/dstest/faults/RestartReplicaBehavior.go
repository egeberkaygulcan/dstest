package faults

import "fmt"

type RestartReplicaBehavior struct {
	nodeId int
}

var _ Behavior = (*RestartReplicaBehavior)(nil)

func NewRestartReplicaBehavior(nodeId int) *RestartReplicaBehavior {
	return &RestartReplicaBehavior{
		nodeId,
	}
}

func (fb *RestartReplicaBehavior) Apply(context FaultContext) error {
	context.GetProcessManager().RestartReplica(fb.nodeId)
	return nil
}

func (fb *RestartReplicaBehavior) String() string {
	return fmt.Sprintf("restart replica %d", fb.nodeId)
}
