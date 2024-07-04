package faults

import "fmt"

type IsolateNodeBehavior struct {
	nodeId int
}

func (fb *IsolateNodeBehavior) Apply(context *FaultContext) error {
	(*context).GetNetworkManager().Router.IsolateNode(fb.nodeId)
	return nil
}

func (fb *IsolateNodeBehavior) String() string {
	return fmt.Sprintf("isolate node %d", fb.nodeId)
}
