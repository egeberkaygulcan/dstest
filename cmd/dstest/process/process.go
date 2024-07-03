package process

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
)

type ProcessManager struct {
	Config *config.Config
	Workers map[int]*Worker
	ClientWorkers map[int]*Worker
	WorkerIds []int
	Iteration int
	Log *log.Logger
	CrashedWorkers map[int]bool
	WaitGroup *sync.WaitGroup
	Basedir string
	BugCandidate bool
	ClientIdCounter atomic.Uint32
}

func (pm *ProcessManager) Init(config *config.Config, workerIds []int, iteration int) {
	pm.Config = config
	pm.Log = log.New(os.Stdout, "[ProcessManager]", log.Default().Flags())
	pm.CrashedWorkers = make(map[int]bool)
	pm.Workers = make(map[int]*Worker)
	pm.ClientWorkers = make(map[int]*Worker)
	pm.WorkerIds = workerIds
	pm.WaitGroup = new(sync.WaitGroup)
	pm.Iteration = iteration
	pm.BugCandidate = false

	if err := os.MkdirAll(pm.Config.ProcessConfig.OutputDir, os.ModePerm); err != nil {
        pm.Log.Printf("Could not create output directory.\n Err: %s\n", err)
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
		conf["workerId"] = pm.WorkerIds[i]
		conf["type"] = Replica
		conf["baseInterceptorPort"] = pm.Config.NetworkConfig.BaseInterceptorPort
		conf["numReplicas"] = pm.Config.ProcessConfig.NumReplicas
		conf["params"] = pm.Config.ProcessConfig.ReplicaParams[i]
		conf["timeout"] = pm.Config.ProcessConfig.Timeout
		pm.Basedir = filepath.Join(pm.Config.ProcessConfig.OutputDir, 
								fmt.Sprintf("%s_%s_%d", 
											pm.Config.TestConfig.Name, 
											pm.Config.SchedulerConfig.Type,
											pm.Iteration))
		if err := os.MkdirAll(pm.Basedir, os.ModePerm); err != nil {
			pm.Log.Printf("Could not create iteration directory.\n Err: %s\n", err)
		}
		conf["basedir"] = pm.Basedir 

		stdout, err := os.OpenFile(filepath.Join(pm.Basedir, fmt.Sprintf("stdout_%d.log", conf["workerId"])), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			pm.Log.Printf("Could not create worker stdout.\n Err: %s\n", err)
		}
		conf["stdout"] = stdout

		stderr, err := os.OpenFile(filepath.Join(pm.Basedir, fmt.Sprintf("stderr_%d.log", conf["workerId"])), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			pm.Log.Printf("Could not create worker stderr.\n Err: %s\n", err)
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
		} (worker, pm.WaitGroup)
	}

	pm.WaitGroup.Wait()

	var bug bool = false
	for workerId, worker := range pm.Workers {
		if worker.Status != Done {
			bug = true
		}
		pm.Log.Printf("Worker %d status: %s\n", workerId, worker.Status.String())
	}

	if !bug {
		pm.deleteDir()
	} else {
		pm.Log.Printf("Found bug candidate at iteration %d\n", pm.Iteration)
		pm.BugCandidate = true
	}
}

func (pm *ProcessManager) Shutdown() {
	for _, worker := range pm.Workers {
		if worker.Status != Exception && worker.Status != Timeout {
			worker.StopWorker()
		}
	}

	for _, worker := range pm.ClientWorkers {
		if worker.Status != Exception && worker.Status != Timeout && worker.Status != Done {
			worker.StopWorker()
		}
	}
}

func (pm *ProcessManager) CrashReplica(workerId int) bool {
	for _, w := range pm.Workers {
		if w.WorkerId == workerId {
			w.CrashWorker()
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

func (pm *ProcessManager) generateClientWorkerConfig(clientType int) map[string]any {
	conf := make(map[string]any)
	conf["runScript"] = pm.Config.ProcessConfig.ClientScripts[clientType]
	conf["type"] = Client
	conf["workerId"] = int(pm.ClientIdCounter.Add(1))
	conf["timeout"] = pm.Config.ProcessConfig.Timeout
	pm.Basedir = filepath.Join(pm.Config.ProcessConfig.OutputDir, 
							fmt.Sprintf("%s_%s_%d", 
										pm.Config.TestConfig.Name, 
										pm.Config.SchedulerConfig.Type,
										pm.Iteration))
	conf["basedir"] = pm.Basedir 

	stdout, err := os.OpenFile(filepath.Join(pm.Basedir, fmt.Sprintf("client_stdout_%d.log", conf["workerId"])), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		pm.Log.Printf("Could not create worker stdout.\n Err: %s\n", err)
	}
	conf["stdout"] = stdout

	stderr, err := os.OpenFile(filepath.Join(pm.Basedir, fmt.Sprintf("client_stderr_%d.log", conf["workerId"])), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		pm.Log.Printf("Could not create worker stderr.\n Err: %s\n", err)
	}
	conf["stderr"] = stderr
	return conf
}

func (pm *ProcessManager) RunClient(clientType int) {
	// Initialize client
	config := pm.generateClientWorkerConfig(clientType)
	clientWorker := new(Worker)
	clientWorker.Init(config)
	pm.ClientWorkers[config["workerId"].(int)] = clientWorker

	// Call client worker as goroutine
	pm.WaitGroup.Add(1)
	go func(worker *Worker, wg *sync.WaitGroup) {
		worker.RunWorker()
		wg.Done()
	} (clientWorker, pm.WaitGroup)
}

func (pm *ProcessManager) deleteDir() {
	os.RemoveAll(pm.Basedir)
}