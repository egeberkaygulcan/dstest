package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type ProcessType int

const (
	Replica ProcessType = 0
	Client	ProcessType = 1
)

type Worker struct{
	RunScript string
	WorkerId int
	Type ProcessType
	StdOut *string
	Log *log.Logger
}

func (worker *Worker) Init(config map[string]any) {
	worker.RunScript = config["runScript"].(string)
	worker.WorkerId = config["workerId"].(int)
	worker.Type = config["type"].(ProcessType)
	worker.StdOut = new(string)
	worker.Log = log.New(os.Stdout, fmt.Sprintf("[Worker %d] ", worker.WorkerId), log.LstdFlags)
}

func (worker *Worker) RunWorker() {
	worker.Log.Println("Run sript: " + worker.RunScript)
	cmd := exec.Command("/bin/sh", worker.RunScript)

	out, err := cmd.Output()

    if err != nil {
    	worker.Log.Printf("error %s, out %s\n", err, string(out))
    }
 
	worker.Log.Println("Output: " + string(out))

    *worker.StdOut = string(out)
}