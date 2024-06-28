package faults

type DummyFaultyBehavior struct {
}

func (fb *DummyFaultyBehavior) Apply(context FaultContext) error {
	// do nothing
	return nil
}

func (fb *DummyFaultyBehavior) String() string {
	return "do nothing"
}
