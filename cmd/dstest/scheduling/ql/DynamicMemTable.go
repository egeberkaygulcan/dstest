package ql

import "fmt"

type DynamicMemTable struct {
	table map[uint32]map[int]float32
}

// confirm DynamicMemTable implements Table interface
var _ Table = (*DynamicMemTable)(nil)

func NewDynamicMemTable() *DynamicMemTable {
	return &DynamicMemTable{
		table: map[uint32]map[int]float32{},
	}
}

func (m *DynamicMemTable) GetMax(state uint32) (action int, qValue float32, err error) {
	qv, ok := m.table[state]
	if !ok {
		return 0, 0.0, nil
	}
	maxQValue := float32(-1)
	for action, value := range qv {
		if value > maxQValue {
			maxQValue = value
			action = action
		}
	}
	return action, maxQValue, nil
}

func (m *DynamicMemTable) Get(state uint32, action int) (float32, error) {
	qv, ok := m.table[state]
	if !ok {
		return 0.0, nil
	}
	qValue, ok := qv[action]
	if !ok {
		return 0.0, fmt.Errorf("action %s does not exist in state %d", action, state)
	}
	return qValue, nil
}

func (m *DynamicMemTable) Set(state uint32, action int, qValue float32) error {
	qv, ok := m.table[state]
	if !ok {
		qv = map[int]float32{}
	}
	qv[action] = qValue
	m.table[state] = qv
	return nil
}

func (m *DynamicMemTable) Clear() error {
	m.table = map[uint32]map[int]float32{}
	return nil
}

func (m *DynamicMemTable) Print() {
	for state, values := range m.table {
		fmt.Println("-----")
		fmt.Println("State: ", state)
		for action, qValue := range values {
			fmt.Println(action, ":", qValue)
		}
		fmt.Println("-----")
	}
}
