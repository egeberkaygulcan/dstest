package faults

type DummyFault struct {
	BaseFault
}

type DummyFaultParams struct {
}

var _ Fault = (*DummyFault)(nil)

func NewDummyFault(params map[string]interface{}) (*DummyFault, error) {

	_ = &DummyFaultParams{}

	return &DummyFault{
		BaseFault: BaseFault{
			Precondition: &DummyPrecondition{},
			Behavior:     &DummyBehavior{},
		},
	}, nil
}
