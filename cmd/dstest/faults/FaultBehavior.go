package faults

type FaultBehavior interface {
	Apply(message FaultContext) error
}
