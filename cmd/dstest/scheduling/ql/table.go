package ql

// Table is the qualtiy table which stores the quality of an action by state.
type Table interface {
	// GetMax returns the action with the max Q value for a given state hash.
	GetMax(state uint32) (action int, qValue float32, err error)

	// Get the Q value for the given state and action.
	Get(state uint32, action int) (float32, error)

	// Set the q value of the action taken for a given state.
	Set(state uint32, action int, value float32) error

	// Clear the table.
	Clear() error

	// Pretty print the table.
	Print()
}
