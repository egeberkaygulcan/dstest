package scheduling

import "github.com/egeberkaygulcan/dstest/cmd/dstest/network"

type Scheduler interface {
	// Init initializes the scheduler
	Init()

	OnQueuedMessage(m *network.Message)

	OnStartup()

	OnShutdown()
}
