package main

import (
	"fmt"
	"log"
	"time"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/process"
)

func main() {
	fmt.Println("Starting dstest")
	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Name: " + cfg.TestConfig.Name)

	output := ""

	go process.RunReplicaWorker(*cfg.ProcessConfig, &output)

	time.Sleep(1 * time.Second)
	fmt.Println("Output: " + output)
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