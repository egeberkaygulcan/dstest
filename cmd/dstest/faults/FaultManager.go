package faults

import (
	"container/list"
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type FaultManager struct {
	// list of Faults
	Faults list.List
}

// AddFault adds a fault to the list of faults
func (fm *FaultManager) AddFault(f BaseFault) {
	fm.Faults.PushBack(f)
}

// ApplyFaults applies all the faults in the list of faults
func (fm *FaultManager) ApplyFaults(message *network.Message) error {
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(BaseFault)
		err := f.ApplyBehaviorIfTriggered(message)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintFaults prints all the faults in the list of faults
func (fm *FaultManager) PrintFaults() {
	for e := fm.Faults.Front(); e != nil; e = e.Next() {
		f := e.Value.(BaseFault)
		fmt.Println("BaseFault: ", f)
	}
}
