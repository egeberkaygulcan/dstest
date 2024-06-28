package faults

type DummyFaultTrigger struct {
}

//var _ faults.Precondition = (*DummyFaultTrigger)(nil)

func (ft *DummyFaultTrigger) Satisfies() (bool, error) {
	// never triggered
	return false, nil
}

func (ft *DummyFaultTrigger) String() string {
	return "false"
}
