package process

import (
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"log"
)

type ProcessManager struct {
	Config    *config.Config
	Workers   []*Worker
	IdCounter int
	Log       *log.Logger
}

func (pm *ProcessManager) getWorkerId() int {
	id := pm.IdCounter
	pm.IdCounter++
	return id
}

func (pm *ProcessManager) Init(config *config.Config) {
	pm.Config = config

	// Generate worker configurations
	workerConfig := pm.generateReplicaWorkerConfig()

	// Create and init workers
	for i := 0; i < pm.Config.ProcessConfig.NumReplicas; i++ {
		worker := new(Worker)
		worker.Init(workerConfig[i])
		pm.Workers = append(pm.Workers, worker)
	}
}

func (pm *ProcessManager) generateReplicaWorkerConfig() []map[string]any {
	var config []map[string]any
	for i := 0; i < pm.Config.ProcessConfig.NumReplicas; i++ {
		conf := make(map[string]any)
		conf["runScript"] = pm.Config.ProcessConfig.ReplicaScript
		conf["workerId"] = pm.getWorkerId()
		conf["type"] = Replica
		config = append(config, conf)
	}

	return config
}

func (pm *ProcessManager) Run() {
	// Run processes
	for i := 0; i < pm.Config.ProcessConfig.NumReplicas; i++ {
		go pm.Workers[i].RunWorker()
	}

	// Wait for process fault injections

	// Inject when requested
}
