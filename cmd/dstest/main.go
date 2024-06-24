package main

import (
	"fmt"
	"log"

	"github.com/egeberkaygulcan/dstest/cmd/dstest/cmd"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/config"
	"github.com/egeberkaygulcan/dstest/cmd/dstest/engine"
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

	te := new(engine.TestEngine)
	te.Init(cfg)

	err = te.Run()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Execute()
}
