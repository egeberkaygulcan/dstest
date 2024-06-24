package behavior

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type DummyFaultyBehavior struct {
}

//var _ faults.FaultBehavior = (*DummyFaultyBehavior)(nil)

func (fb *DummyFaultyBehavior) Apply(message *network.Message) error {
	// do nothing
	return nil
}
