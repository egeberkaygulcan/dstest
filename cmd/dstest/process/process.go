package process

import "github.com/egeberkaygulcan/dstest/cmd/dstest/config"

type ProcessManager struct {
	ProcessConfig *config.ProcessConfig
}

func (pm ProcessManager) init(config config.ProcessConfig) {
	// Generate worker configurations

	// Create workers

	// Init workers
}

func (pm ProcessManager) run() {
	// Run processes

	// Wait for process fault injections

	// Inject when requested
}