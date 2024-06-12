package main

import (
	"fmt"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/network"
	"log"
	"sync"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	// "github.com/egeberkaygulcan/dstest/cmd/dstest/network"
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
	pm.Init(cfg, 0)
	// out, err := process.RunReplicaWorker(*cfg.ProcessConfig)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// fmt.Println(out)
	// Init scheduler

	// Init network
	nm := new(network.Manager)
	nm.Init(cfg)

	// Init processes

	// Run network
	go nm.Run()

	// Run scheduler

	// Spawn processes
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pm.Run()
		wg.Done()
	}()

	// Later wrap this process around an experiment module

	wg.Wait()
}
