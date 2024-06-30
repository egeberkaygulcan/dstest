package faults

import (
	"container/list"
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
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
	case "crash":
		return NewCrashReplicaFault(Params)
	case "restart":
		return NewRestartReplicaFault(Params)
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
			return fmt.Errorf("error creating %s fault: %v", f.Type, err)
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
func (fm *FaultManager) ApplyFaults(context *FaultContext) error {
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(Fault)
		err := f.ApplyBehaviorIfPreconditionMet(context)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFaults returns all the faults in the list of faults
func (fm *FaultManager) GetFaults() []*Fault {
	faults := make([]*Fault, 0)
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(Fault)
		faults = append(faults, &f)
	}
	return faults
}

// GetEnabledFaults returns all the enabled faults in the list of faults
func (fm *FaultManager) GetEnabledFaults() []*Fault {
	faults := make([]*Fault, 0)
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(Fault)
		enabled, err := f.IsEnabled()
		if err != nil {
			fmt.Println("Error getting enabled status of fault: ", err)
		}
		if enabled {
			faults = append(faults, &f)
		}
	}
	return faults
}

// PrintFaults prints all the faults in the list of faults
func (fm *FaultManager) PrintFaults() {
	fmt.Println("Faults:")
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(Fault)
		fmt.Printf("\t- %+v\n", f)
	}
}
