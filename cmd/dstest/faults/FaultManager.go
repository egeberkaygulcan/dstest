package faults

import (
	"container/list"
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
	"strings"
)

type FaultManager struct {
	// list of Faults
	Faults list.List
}

func NewFault(Name string, Params map[string]interface{}) (Fault, error) {
	switch strings.ToLower(Name) {
	case "dummy":
		return NewDummyFault(Params)
	default:
		return nil, fmt.Errorf("Fault '%s' not found", Name)
	}
}

func (fm *FaultManager) Init(config *config.Config) error {
	fm.Faults.Init()

	// Add faults to the list
	faults := config.FaultConfig.Faults

	for _, f := range faults {
		fmt.Println("Going to create a fault: ", f.Type)
		fault, err := NewFault(f.Type, f.Params)
		if err != nil {
			return fmt.Errorf("error creating fault: %v", err)
		}
		fm.AddFault(fault)
	}

	return nil
}

// AddFault adds a fault to the list of faults
func (fm *FaultManager) AddFault(f Fault) {
	fm.Faults.PushBack(f)
}

// ApplyFaults applies all the faults in the list of faults
func (fm *FaultManager) ApplyFaults(context FaultContext) error {
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(Fault)
		err := f.ApplyBehaviorIfTriggered(context)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintFaults prints all the faults in the list of faults
func (fm *FaultManager) PrintFaults() {
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(Fault)
		fmt.Println("Fault: ", f)
	}
}

type Context interface {
	GetConfig() *config.Config
	GetNetworkManager() *network.Manager
	GetProcessManager() *process.ProcessManager
}
