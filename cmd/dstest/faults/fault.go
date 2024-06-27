package faults

import (
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
	ApplyBehaviorIfTriggered(context FaultContext) error
}

type FaultTrigger interface {
	Satisfies() (bool, error)
}

type BaseFault struct {
	FaultTrigger
	FaultBehavior
}

var _ Fault = (*BaseFault)(nil)

func (f *BaseFault) ApplyBehaviorIfTriggered(context FaultContext) error {
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
