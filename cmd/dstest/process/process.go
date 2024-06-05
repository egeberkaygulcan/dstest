package process

import (
  "fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
)

type ProcessManager struct {
	Config *config.Config
	Workers map[int]*Worker
	Iteration int
	IdCounter int
	Log *log.Logger
	CrashedWorkers map[int]bool
	WaitGroup *sync.WaitGroup
}

func (pm *ProcessManager) getWorkerId() int {
	id := pm.IdCounter
	pm.IdCounter++
	return id
}

func (pm *ProcessManager) Init(config *config.Config, iteration int) {
	pm.Config = config
	pm.Log = log.New(os.Stdout, "", log.Default().Flags())
	pm.CrashedWorkers = make(map[int]bool)
	pm.Workers = make(map[int]*Worker)
	pm.WaitGroup = new(sync.WaitGroup)
	pm.Iteration = iteration

	if err := os.MkdirAll(pm.Config.ProcessConfig.OutputDir, os.ModePerm); err != nil {
        log.Fatalf("Error while creating output folder: %s", err)
    }

	// Generate worker configurations
	workerConfig := pm.generateReplicaWorkerConfig()

	// Create and init workers
	for i := 0; i < pm.Config.ProcessConfig.NumReplicas; i++ {
		worker := new(Worker)
		worker.Init(workerConfig[i])
		pm.Workers[(workerConfig[i]["workerId"]).(int)] = worker
	}
}

func (pm *ProcessManager) generateReplicaWorkerConfig() []map[string]any {
	var config []map[string]any
	for i := 0; i < pm.Config.ProcessConfig.NumReplicas; i++ {
		conf := make(map[string]any)
		conf["runScript"] = pm.Config.ProcessConfig.ReplicaScript
		conf["cleanScript"] = pm.Config.ProcessConfig.CleanScript
		conf["clientScripts"] = pm.Config.ProcessConfig.ClientScripts
		conf["workerId"] = pm.getWorkerId()
		conf["type"] = Replica
		conf["baseInterceptorPort"] = pm.Config.NetworkConfig.BaseInterceptorPort
		conf["numReplicas"] = pm.Config.ProcessConfig.NumReplicas
		conf["params"] = pm.Config.ProcessConfig.ReplicaParams[i]
		conf["timeout"] = pm.Config.ProcessConfig.Timeout
		basedir := filepath.Join(pm.Config.ProcessConfig.OutputDir, 
								fmt.Sprintf("%s_%s_%d", 
											pm.Config.TestConfig.Name, 
											pm.Config.SchedulerConfig.Type,
											pm.Iteration))
		if err := os.MkdirAll(basedir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		conf["basedir"] = basedir 

		stdout, err := os.OpenFile(filepath.Join(basedir, fmt.Sprintf("stdout_%d.log", conf["workerId"])), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		conf["stdout"] = stdout

		stderr, err := os.OpenFile(filepath.Join(basedir, fmt.Sprintf("stderr_%d.log", conf["workerId"])), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		conf["stderr"] = stderr

		config = append(config, conf)
	}

	return config
}

func (pm *ProcessManager) Run() {
	// Run processes
	for _, worker := range pm.Workers {
		pm.WaitGroup.Add(1)
		go func(worker *Worker, wg *sync.WaitGroup) {
			worker.RunWorker()
			wg.Done()
			delete(pm.Workers, worker.WorkerId)
		} (worker, pm.WaitGroup)
	}

	pm.WaitGroup.Wait()

	for workerId, worker := range pm.Workers {
		pm.Log.Printf("Worker %d status: %d\n", workerId, worker.Status)
	}
}

func (pm *ProcessManager) Reset() {
	// TODO - Kill workers
	for _, worker := range pm.Workers {
		worker.StopWorker()
	}

	// TODO - Reset properties
	pm.IdCounter = 0

}

func (pm *ProcessManager) CrashReplica(workerId int) bool {
	for _, w := range pm.Workers {
		if w.WorkerId == workerId {
			w.KillWorker()
			pm.WaitGroup.Done()
			pm.CrashedWorkers[workerId] = true
			return true
		}
	}

	return false
}

func (pm *ProcessManager) RestartReplica(workerId int) bool {
	for key, _ := range pm.CrashedWorkers {
		if key == workerId {
			pm.WaitGroup.Add(1)
			pm.Workers[key].RestartWorker()
			// delete(pm.Workers, key)
			return true
		}
	}

	return false
}

func (pm *ProcessManager) generateClientWorkerConfig() []map[string]any {
	// TODO
	return nil
}

func (pm *ProcessManager) RunClient() {
	// Initialize client

	// Call client worker as goroutine
}