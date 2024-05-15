package process

import (
	"fmt"
	"os/exec"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
)

// type ProcessType int

// const (
// 	Replica ProcessType = 0
// 	Client	ProcessType = 1
// )

func RunReplicaWorker(config config.ProcessConfig, output *string) {
	fmt.Println("Replica sript: " + config.ReplicaScript)
	cmd := exec.Command("/bin/sh", config.ReplicaScript)
	// cmd := exec.Command("echo", "hello!")
	out, err := cmd.Output()

    if err != nil {
    	fmt.Printf("error %s", err)
    }
 
    *output = string(out)
}

func RunClientWorker() {

}