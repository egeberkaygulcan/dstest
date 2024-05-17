package main

import (
	"fmt"
	"log"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
)

func main() {

	// ------ DO NOT CHANGE -------
	fmt.Println("Starting dstest")
	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// -----------------------------

	fmt.Println("Name: " + cfg.TestConfig.Name)

	// Spawn goroutine
	pm := new(process.ProcessManager)
	pm.Init(cfg)
	pm.Run()
	// out, err := process.RunReplicaWorker(*cfg.ProcessConfig)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// fmt.Println(out)
	// Init scheduler

	// Init network

	// Init processes

	// Run network

	// Run scheduler

	// Spawn processes

	// Later wrap this process around an experiment module
}
