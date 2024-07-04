package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type ProcessType int

const (
	Replica ProcessType = 0
	Client  ProcessType = 1
)

type ProcessStatus int

const (
	Initialized ProcessStatus = 0
	Running 	ProcessStatus = 1
	Done 		ProcessStatus = 2
	Crashed		ProcessStatus = 3
	Timeout		ProcessStatus = 4
	Exception	ProcessStatus = 5		
)

func (s ProcessStatus) String() string {
    switch s {
    case Initialized:
        return "Initialized"
    case Running:
        return "Running"
	case Done:
		return "Done"
	case Crashed:
		return "Crashed"
	case Timeout:
		return "Timeout"
	case Exception:
		return "Exception"
    default:
        return fmt.Sprintf("%d", int(s))
    }
}

type Worker struct{
	RunScript 	  string
	NumReplicas         int
	BaseInterceptorPort int
	CleanScript	  string
	WorkerId 	  int
	Type 		  ProcessType
	Params 		  string
	
	Timeout		  int
	TimeoutDelta  int
	TimeoutTimer  *time.Timer
	Status		  ProcessStatus
	Cmd    		  *exec.Cmd

	Stdout		  *os.File
	Stderr		  *os.File
	Log 		  *log.Logger
}

func (worker *Worker) Init(config map[string]any) {
	worker.Type = config["type"].(ProcessType)
	worker.WorkerId = config["workerId"].(int)
	worker.RunScript = config["runScript"].(string)
	worker.Timeout = config["timeout"].(int)

	if worker.Type == Replica {
		worker.CleanScript = config["cleanScript"].(string)
		worker.NumReplicas = config["numReplicas"].(int)
		worker.BaseInterceptorPort = config["baseInterceptorPort"].(int)
		worker.Params = config["params"].(string)

		worker.Log = log.New(os.Stdout, fmt.Sprintf("[Worker %d] ", worker.WorkerId), log.LstdFlags)
	} else {
		worker.Log = log.New(os.Stdout, fmt.Sprintf("[Client %d] ", worker.WorkerId), log.LstdFlags)
	}

	worker.TimeoutTimer = nil

	worker.Stdout = config["stdout"].(*os.File)
	worker.Stderr = config["stderr"].(*os.File)

	worker.Status = Initialized
}

func (worker *Worker) RunWorker() {
	defer worker.clean()

	worker.Log.Println("Running worker with: " + worker.RunScript + " " + worker.Params)

	worker.Cmd = exec.Command("/bin/sh", strings.Fields(worker.RunScript + " " + worker.Params)...)
	worker.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	worker.Cmd.Stdout = worker.Stdout
	worker.Cmd.Stderr = worker.Stderr

	errch := make(chan error, 1)

	err := worker.Cmd.Start()
	if err != nil {
		worker.Log.Printf("Error while starting worker. \nError: %s\n", err)
	}

	if worker.TimeoutTimer == nil {
		worker.TimeoutTimer = time.NewTimer(time.Duration(worker.Timeout) * time.Second)
	}
	worker.Status = Running

	go func() {
		errch <- worker.Cmd.Wait()
	} ()
	

	select {
	case <- worker.TimeoutTimer.C:
		worker.Log.Println("Timeout, killing process.")
		worker.Status = Timeout
		worker.KillWorker()
		return
	case err:= <- errch:
		if err != nil {
			if worker.Status != Crashed && worker.Status != Done {
				worker.Log.Printf("Error while waiting worker. \nError: %s\n", err)
				worker.Status = Exception
				// worker.KillWorker()
			}
			return
		} else {
			worker.Status = Done
		}
	}
}

func (worker *Worker) KillWorker() {
	syscall.Kill(-worker.Cmd.Process.Pid, syscall.SIGKILL)
	worker.Log.Printf("Killed worker %d\n", worker.WorkerId)
}

func (worker *Worker) CrashWorker() {
	worker.Status = Crashed
	worker.KillWorker()
}

func (worker *Worker) StopWorker() {
	data, err := os.ReadFile(worker.Stderr.Name())
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines[0]) == 0 {
		worker.Status = Done // TODO - Bug! Never deletes the output folder
	}

	worker.KillWorker()
}

func (worker *Worker) RestartWorker() {
	if worker.Status == Crashed {
		worker.RunWorker()
	}
}

func (worker *Worker) clean() {
	if len(worker.CleanScript) == 0 {
		return
	}

	worker.Log.Println("Calling the clean script.")
	cmd := exec.Command("/bin/bash", worker.CleanScript)

	err := cmd.Start()
	if err != nil {
		worker.Log.Printf("Error while cleaning up worker. \nError: %s\n", err)
	}

	err = cmd.Wait()
	if err != nil {
		worker.Log.Printf("Error while waiting worker cleanup. \nError: %s\n", err)
	}
}
