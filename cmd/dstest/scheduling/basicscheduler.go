package scheduling

import (
	"fmt"
)

type BasicScheduler struct {
	// initialization code goes here
}

// Check if BasicScheduler implements Scheduler interface

func (bs *BasicScheduler) Init() {
	fmt.Println("BasicScheduler initialized")
}

func (bs *BasicScheduler) OnStartup() {
	// Initialization code goes here
	fmt.Println("BasicScheduler Startup!")
}

func (bs *BasicScheduler) OnQueuedMessage(m *any) {
	// Initialization code goes here
	fmt.Println("BasicScheduler initialized")
}

func (bs *BasicScheduler) OnShutdown() {
	// Initialization code goes here
	fmt.Println("BasicScheduler Shutdown!")
}
