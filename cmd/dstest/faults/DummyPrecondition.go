package faults

type DummyFaultTrigger struct {
}

//var _ faults.FaultTrigger = (*DummyFaultTrigger)(nil)

func (ft *DummyFaultTrigger) Satisfies() (bool, error) {
	// never triggered
	return false, nil
}
