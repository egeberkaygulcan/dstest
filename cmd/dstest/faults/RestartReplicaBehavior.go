package faults

type RestartReplicaBehavior struct {
	nodeId int
}

var _ FaultBehavior = (*RestartReplicaBehavior)(nil)

func NewRestartReplicaBehavior(nodeId int) *RestartReplicaBehavior {
	return &RestartReplicaBehavior{
		nodeId,
	}
}

func (fb *RestartReplicaBehavior) Apply(context FaultContext) error {
	// do nothing
	return nil
}
