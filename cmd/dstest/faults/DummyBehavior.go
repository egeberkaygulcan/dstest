package faults

type DummyBehavior struct {
}

func (fb *DummyBehavior) Apply(context FaultContext) error {
	// do nothing
	return nil
}

func (fb *DummyBehavior) String() string {
	return "do nothing"
}
