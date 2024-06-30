package faults

type Behavior interface {
	Apply(message *FaultContext) error
	String() string
}
