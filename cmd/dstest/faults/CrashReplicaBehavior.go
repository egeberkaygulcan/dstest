package faults

type CrashReplicaBehavior struct {
	nodeId int
}

var _ FaultBehavior = (*CrashReplicaBehavior)(nil)

func NewCrashReplicaBehavior(nodeId int) *CrashReplicaBehavior {
	return &CrashReplicaBehavior{
		nodeId,
	}
}

func (fb *CrashReplicaBehavior) Apply(context FaultContext) error {
	// do nothing
	return nil
}
