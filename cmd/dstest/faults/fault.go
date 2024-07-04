package faults

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
)

type FaultContext interface {
	GetConfig() *config.Config
	GetNetworkManager() *network.Manager
	GetProcessManager() *process.ProcessManager
}

type Fault interface {
	IsEnabled() (bool, error)
	ApplyBehaviorIfPreconditionMet(context *FaultContext) error
	String() string
}

type Precondition interface {
	Satisfies() (bool, error)
	String() string
}

type BaseFault struct {
	Precondition
	Behavior
}

var _ Fault = (*BaseFault)(nil)

func (f *BaseFault) ApplyBehaviorIfPreconditionMet(context *FaultContext) error {
	triggered, err := f.Satisfies()
	if err != nil {
		return err
	}
	if triggered {
		return f.Apply(context)
	}
	return nil
}

func (f *BaseFault) IsEnabled() (bool, error) {
	return f.Satisfies()
}

func (f *BaseFault) String() string {
	return fmt.Sprintf("if %s then %s", f.Precondition, f.Behavior)
}
