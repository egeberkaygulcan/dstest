package faults

import "github.com/egeberkaygulcan/dstest/cmd/dstest/network"

type Fault interface {
	ApplyBehaviorIfTriggered(message *network.Message) error
}

type FaultTrigger interface {
	Satisfies() (bool, error)
}

type FaultBehavior interface {
	Apply(message *network.Message) error
}

type BaseFault struct {
	FaultTrigger
	FaultBehavior
}

var _ Fault = (*BaseFault)(nil)

func (f *BaseFault) ApplyBehaviorIfTriggered(message *network.Message) error {
	triggered, err := f.Satisfies()
	if err != nil {
		return err
	}
	if triggered {
		return f.Apply(message)
	}
	return nil
}
