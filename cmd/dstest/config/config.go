package config

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultConfiguration []byte

type TestConfig struct {
	Name        string
	Experiments int
	Iterations  int
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
}

type ProcessConfig struct {
	NumReplicas   int
	ReplicaScript string
	ClientScripts []string
}

type Config struct {
	TestConfig      *TestConfig
	SchedulerConfig *SchedulerConfig
	NetworkConfig   *NetworkConfig
	ProcessConfig   *ProcessConfig
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
	return &config, nil
}
