package faults

type AlwaysEnabledPrecondition struct {
}

//var _ faults.FaultTrigger = (*DummyFaultTrigger)(nil)

func (ft *AlwaysEnabledPrecondition) Satisfies() (bool, error) {
	// always enabled
	return true, nil
}
