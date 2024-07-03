package config

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultConfiguration []byte

type TestConfig struct {
	Name        string
	Experiments int
	Iterations  int
	WaitDuration int
}

type SchedulerConfig struct {
	Type   string
	Steps  int
	ClientRequests int
	Seed   int
	Params map[string]any
}

type NetworkConfig struct {
	BaseReplicaPort     int
	BaseInterceptorPort int
}

type ProcessConfig struct {
	NumReplicas   int
	Timeout 	  int
	OutputDir	  string
	ReplicaScript string
	ClientScripts []string
	CleanScript	  string
	ReplicaParams []string
}

type Config struct {
	TestConfig      *TestConfig
	SchedulerConfig *SchedulerConfig
	NetworkConfig   *NetworkConfig
	ProcessConfig   *ProcessConfig
}

func ModifyFilepath(config *Config) {
	wd, _ := os.Getwd()
	wd = filepath.Clean(filepath.Join(wd, "../.."))

	config.ProcessConfig.OutputDir = filepath.Join(wd, config.ProcessConfig.OutputDir)
	config.ProcessConfig.ReplicaScript = filepath.Join(wd, config.ProcessConfig.ReplicaScript)
	if len(config.ProcessConfig.CleanScript) > 0 {
		config.ProcessConfig.CleanScript = filepath.Join(wd, config.ProcessConfig.CleanScript)
	}
}

func Read() (*Config, error) {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configuration file type
	viper.SetConfigType("yml")

	// Read configuration
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		return nil, err
	}

	// Unmarshal the configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	ModifyFilepath(&config)
	return &config, nil
}
