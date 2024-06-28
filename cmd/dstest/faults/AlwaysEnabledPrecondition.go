package faults

type AlwaysEnabledPrecondition struct {
}

//var _ faults.Precondition = (*DummyFaultTrigger)(nil)

func (ft *AlwaysEnabledPrecondition) Satisfies() (bool, error) {
	// always enabled
	return true, nil
}

func (ft *AlwaysEnabledPrecondition) String() string {
	return "true"
}
