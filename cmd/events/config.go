package main

import (
	"gobaos/pkg/baos"
	"gopkg.in/yaml.v3"
	"os"
)

type datapoint struct {
	ID   []baos.DatapointID `yaml:"ids"`
	Type string             `yaml:"type"`
}

type config struct {
	Device     string      `yaml:"device"`
	Datapoints []datapoint `yaml:"datapoints"`
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
