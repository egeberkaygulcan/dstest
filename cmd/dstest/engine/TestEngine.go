package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/actions"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/faults"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/scheduling"
)

// FIXME: Repetition of FaultManager interface to avoid cyclic import
// how to avoid this?
type FaultManager interface {
	Init(config *config.Config) error
	GetFaults() []*faults.Fault
	GetEnabledFaults() []*faults.Fault
	PrintFaults()
}

type TestEngine struct {
	Config         *config.Config
	Scheduler      scheduling.Scheduler
	NetworkManager *network.Manager
	ProcessManager *process.ProcessManager
	FaultManager   FaultManager
	Log            *log.Logger

	Experiments   int
	Iterations    int
	Steps         int
	SleepDuration time.Duration
	ReplicaIds    []int
}

func (te *TestEngine) Init(config *config.Config) error {
	te.Config = config
	te.Experiments = config.TestConfig.Experiments
	te.Iterations = config.TestConfig.Iterations
	te.Steps = config.SchedulerConfig.Steps
	te.SleepDuration = time.Duration(config.TestConfig.WaitDuration) * time.Millisecond
	te.ReplicaIds = make([]int, te.Config.ProcessConfig.NumReplicas)
	for i := 0; i < te.Config.ProcessConfig.NumReplicas; i++ {
		te.ReplicaIds[i] = i
	}

	// Initialize scheduler
	scheduler, err := scheduling.NewScheduler(scheduling.SchedulerType(config.SchedulerConfig.Type))
	if err != nil {
		return fmt.Errorf("Error initializing scheduler: %s", err.Error())
	}
	te.Scheduler = scheduler

	te.NetworkManager = new(network.Manager)
	te.ProcessManager = new(process.ProcessManager)

	if scheduling.SchedulerType(config.SchedulerConfig.Type) == scheduling.Pctcp {
		config.SchedulerConfig.Params["network_manager"] = te.NetworkManager
	}
	// te.ProcessManager.Init(config, te.Iterations)
	te.FaultManager = new(faults.FaultManager)

	if err := te.FaultManager.Init(config); err != nil {
		return fmt.Errorf("Error initializing FaultManager: %s", err.Error())
	}

	te.Log = log.New(os.Stdout, "[TestEngine] ", log.LstdFlags)

	return nil
}

func (te *TestEngine) Run() error {
	for i := 0; i < te.Experiments; i++ {
		te.Log.Printf("Starting experiment %d...\n", i+1)

		te.Scheduler.Init(te.Config)
		for j := 0; j < te.Iterations; j++ {
			te.Log.Printf("Starting iteration %d\n", j+1)

			// Initialize NetworkManager
			err := te.NetworkManager.Init(te.Config, te.ReplicaIds)
			if err != nil {
				return fmt.Errorf("Error initializing NetworkManager: %s", err.Error())
			}

			// Initialize FaultManager
			err = te.FaultManager.Init(te.Config)
			if err != nil {
				return fmt.Errorf("Error initializing FaultManager: %s", err.Error())
			}
			te.FaultManager.PrintFaults()

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

			time.Sleep(time.Duration(te.Config.TestConfig.StartupDuration) * time.Second)

			schedule := make([]actions.Action, 0)
			for s := 0; s < te.Steps; {
				if te.ProcessManager.BugCandidate {
					break
				}
				enabledMessages := te.NetworkManager.GetActions()
				sc := te.Scheduler.GetClientRequest()
				if sc >= 0 {
					te.ProcessManager.RunClient(sc)
					action := &actions.ClientRequestAction{
						ClientId: sc,
					}
					schedule = append(schedule, action)
				}
				// TODO - Get fault from scheduler
				var faultContext faults.FaultContext = NewEngineFaultContext(te)
				decision := te.Scheduler.Next(enabledMessages, te.FaultManager.GetFaults(), faultContext)

				fmt.Printf("decision: %+v\n", decision)

				if decision.DecisionType == scheduling.SendMessage {
					action := decision.Index
					fmt.Println("action: ", action)
					te.NetworkManager.SendMessage(enabledMessages[action].MessageId)
					schedule = append(schedule, &actions.DeliverMessageAction{
						Sender:   enabledMessages[action].Sender,
						Receiver: enabledMessages[action].Receiver,
						Name:     enabledMessages[action].Name,
					})
					s++
				}

				if decision.DecisionType == scheduling.InjectFault {
					fault := te.FaultManager.GetFaults()[decision.Index]
					te.Log.Printf("Applying fault: %+v\n", fault)
					err := (*fault).ApplyBehaviorIfPreconditionMet(&faultContext)
					if err != nil {
						te.Log.Printf("Error applying fault: %s\n", err)
					}
					schedule = append(schedule, &actions.InjectFaultAction{
						Fault: *fault,
					})
				}

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
			if true { //te.ProcessManager.BugCandidate {
				outputFile, err := os.OpenFile(filepath.Join(te.ProcessManager.Basedir, "schedule.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					te.Log.Printf("Could not create schedule file.\n Err: %s\n", err)
				}

				for _, action := range schedule {
					/*fmt.Fprintf(outputFile, "%+#v\n", struct {
						ActionType ActionType
						Action     *Action
					}{
						ActionType: action.GetType(),
						Action:     &action,
					})*/
					actionJson, _ := json.Marshal(struct {
						ActionType actions.ActionType
						Action     *actions.Action
					}{
						ActionType: action.GetType(),
						Action:     &action,
					})
					fmt.Fprintf(outputFile, "%s\n", (actionJson))
				}
				outputFile.Close()
			}
			te.Log.Println("Iteration complete.")
			te.Scheduler.NextIteration()
		}
		te.Scheduler.Reset()
	}
	return nil
}

type EngineFaultContext struct {
	engine *TestEngine
}

func NewEngineFaultContext(engine *TestEngine) *EngineFaultContext {
	return &EngineFaultContext{engine: engine}
}

func (efc *EngineFaultContext) GetConfig() *config.Config {
	return efc.engine.Config
}

func (efc *EngineFaultContext) GetNetworkManager() *network.Manager {
	return efc.engine.NetworkManager
}

func (efc *EngineFaultContext) GetProcessManager() *process.ProcessManager {
	return efc.engine.ProcessManager
}

// confirm that EngineFaultContext implements FaultContext
var _ faults.FaultContext = (*EngineFaultContext)(nil)
