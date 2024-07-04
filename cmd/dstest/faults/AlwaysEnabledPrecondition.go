package faults

type AlwaysEnabledPrecondition struct {
}

func (ft *AlwaysEnabledPrecondition) Satisfies() (bool, error) {
	// always enabled
	return true, nil
}

func (ft *AlwaysEnabledPrecondition) String() string {
	return "true"
}
