package engine

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling"
)

type Action struct {
	Sender int
	Receiver int
	Name string
}

type TestEngine struct {
	Config *config.Config
	Scheduler scheduling.Scheduler
	NetworkManager *network.Manager
	ProcessManager *process.ProcessManager
	Log *log.Logger

	Experiments int
	Iterations  int
	Steps	    int
	SleepDuration time.Duration
	ReplicaIds	[]int
}

func (te *TestEngine) Init(config *config.Config) {
	te.Config = config
	te.Experiments = config.TestConfig.Experiments
	te.Iterations = config.TestConfig.Iterations
	te.Steps = config.SchedulerConfig.Steps
	te.SleepDuration = time.Duration(config.TestConfig.WaitDuration) * time.Millisecond
	te.ReplicaIds = make([]int, te.Config.ProcessConfig.NumReplicas)
	for i := 0; i < te.Config.ProcessConfig.NumReplicas; i++ {
		te.ReplicaIds[i] = i
	}


	te.Scheduler = scheduling.NewScheduler(scheduling.SchedulerType(config.SchedulerConfig.Type))
	te.NetworkManager = new(network.Manager)
	// te.NetworkManager.Init(config)
	te.ProcessManager = new(process.ProcessManager)
	// te.ProcessManager.Init(config, te.Iterations)

	te.Log = log.New(os.Stdout, "[TestEngine] ", log.LstdFlags)
}

func (te *TestEngine) Run() {
	for i := 0; i < te.Experiments; i++ {
		te.Log.Printf("Starting experiment %d...\n", i)

		te.Scheduler.Init(te.Config)
		for j := 0; j < te.Iterations; j++ {
			te.Log.Printf("Starting iteration %d\n", j+1)
			te.NetworkManager.Init(te.Config, te.ReplicaIds)
			te.ProcessManager.Init(te.Config, te.ReplicaIds, j)
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				te.NetworkManager.Run()
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				te.ProcessManager.Run()
				wg.Done()
			}()

			time.Sleep(1 * time.Second)

			schedule := make([]Action, 0)
			for s := 0; s < te.Steps; s++ {
				actions := te.NetworkManager.GetActions()
				// TODO - Schedule client
				sc := te.Scheduler.GetClientRequest()
				if sc >= 0 {
					te.ProcessManager.RunClient(sc)
				}
				// TODO - Get fault from scheduler
				action := te.Scheduler.Next(actions)
				if action != 0 {
					te.NetworkManager.SendMessage(actions[action].MessageId)
					schedule = append(schedule, Action{
						Sender: actions[action].Sender,
						Receiver: actions[action].Receiver,
						Name: actions[action].Name,
					})
				}
				
				// TODO - Execute fault and append to schedule

				time.Sleep(te.SleepDuration)
			}
			// te.Schedules = append(te.Schedules, schedule)
			te.Log.Println("Shutting down ProcessManager...")
			te.ProcessManager.Shutdown()
			te.Log.Println("Shutting down NetworkManager...")
			te.NetworkManager.Shutdown()
			te.Log.Println("Shutdown complete.")
			wg.Wait()

			te.Log.Println("Checking for bugs...")
			if te.ProcessManager.BugCandidate {
				outputFile, err := os.OpenFile(filepath.Join(te.ProcessManager.Basedir, "schedule.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					te.Log.Printf("Could not create schedule file.\n Err: %s\n", err)
				}

				for _, action := range schedule {
					fmt.Fprintln(outputFile, action)
				}
				outputFile.Close()
			}
			te.Log.Println("Iteration complete.")
			te.Scheduler.NextIteration()
		}
		te.Scheduler.Reset()
	}
}