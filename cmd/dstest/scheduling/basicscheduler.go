package scheduling

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
)

type BasicScheduler struct {
	// initialization code goes here
}

// Check if BasicScheduler implements Scheduler interface
var _ Scheduler = (*BasicScheduler)(nil)

func (bs *BasicScheduler) Init() {
	fmt.Println("BasicScheduler initialized")
}

func (bs *BasicScheduler) OnStartup() {
	// Initialization code goes here
	fmt.Println("BasicScheduler Startup!")
}

func (bs *BasicScheduler) OnQueuedMessage(m *network.Message) {
	// Initialization code goes here
	fmt.Println("BasicScheduler initialized")
}

func (bs *BasicScheduler) OnShutdown() {
	// Initialization code goes here
	fmt.Println("BasicScheduler Shutdown!")
}
