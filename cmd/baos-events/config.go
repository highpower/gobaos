package main

import (
	"gobaos/pkg/baos"
	"os"

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
		return config{}, err
	}
	defer func() { _ = file.Close() }()
	result := config{}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&result); err != nil {
		return config{}, err
	}
	return result, nil
}
