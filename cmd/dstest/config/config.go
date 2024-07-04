package config

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type TestConfig struct {
	Name         string
	Experiments  int
	Iterations   int
	WaitDuration int
}

type SchedulerConfig struct {
	Type   string
	Steps  int
	Seed   int
	Params map[string]string
}

type NetworkConfig struct {
	BaseReplicaPort     int
	BaseInterceptorPort int
	Protocol            string
}

type FaultConfig struct {
	Faults []struct {
		Type   string
		Params map[string]interface{}
	}
}

type ProcessConfig struct {
	NumReplicas   int
	Timeout       int
	OutputDir     string
	ReplicaScript string
	ClientScripts []string
	CleanScript   string
	ReplicaParams []string
}

type Config struct {
	TestConfig      *TestConfig
	SchedulerConfig *SchedulerConfig
	NetworkConfig   *NetworkConfig
	FaultConfig     *FaultConfig
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
	//viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configuration file type
	viper.SetConfigFile(viper.GetString("config"))
	viper.SetConfigType("yml")

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
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
