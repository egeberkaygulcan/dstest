package faults

type DummyPrecondition struct {
}

//var _ faults.Precondition = (*DummyPrecondition)(nil)

func (ft *DummyPrecondition) Satisfies() (bool, error) {
	// never triggered
	return false, nil
}

func (ft *DummyPrecondition) String() string {
	return "false"
}
