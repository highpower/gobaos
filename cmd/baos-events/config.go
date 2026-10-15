package main

import (
	"fmt"
	"os"

	"gobaos/pkg/baos"

	"gopkg.in/yaml.v3"
)

type config struct {
	Device      string             `yaml:"device"`
	Switches    []baos.DatapointID `yaml:"switches"`
	Thermostats []baos.DatapointID `yaml:"thermostats"`
}

func newConfig(name string) (config, error) {
	file, err := os.Open(name)
	if err != nil {
		return config{}, fmt.Errorf("open config %q: %w", name, err)
	}
	defer func() { _ = file.Close() }()
	result := config{}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&result); err != nil {
		return config{}, fmt.Errorf("decode config %q: %w", name, err)
	}
	return result, nil
}
