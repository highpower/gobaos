package main

import (
	"fmt"
	"gobaos/pkg/baos"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

type sensorSettings struct {
	ids        []baos.DatapointID
	IDStrings  []string `yaml:"ids"`
	StateTopic string   `yaml:"state-topic"`
}

type controllerSettings struct {
	ids          []baos.DatapointID
	IDStrings    []string `yaml:"ids"`
	StateTopic   string   `yaml:"state-topic"`
	CommandTopic string   `yaml:"command-topic"`
}

type config struct {
	Device      string               `yaml:"device"`
	MqttHost    string               `yaml:"mqtt-host"`
	MqttUser    string               `yaml:"user"`
	MqttPass    string               `yaml:"password"`
	CaCert      string               `yaml:"ca-cert"`
	Sensors     []sensorSettings     `yaml:"sensors"`
	Controllers []controllerSettings `yaml:"controllers"`
}

func newConfig(path string) (config, error) {
	file, err := os.Open(path)
	if err != nil {
		return config{}, err
	}
	defer func() { _ = file.Close() }()
	result := config{}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err = decoder.Decode(&result); err != nil {
		return config{}, err
	}
	if err := result.parseSensorIDs(); err != nil {
		return config{}, err
	}
	if err := result.parseControllerIDs(); err != nil {
		return config{}, err
	}
	return result, nil
}

func (c *config) parseSensorIDs() error {
	for _, sensor := range c.Sensors {
		ids, err := parseIDStrings(sensor.IDStrings)
		if err != nil {
			return err
		}
		sensor.ids = ids
	}
	return nil
}

func (c *config) parseControllerIDs() error {
	for _, controller := range c.Controllers {
		ids, err := parseIDStrings(controller.IDStrings)
		if err != nil {
			return err
		}
		controller.ids = ids
	}
	return nil
}

func parseIDStrings(idStrings []string) ([]baos.DatapointID, error) {
	result := make([]baos.DatapointID, 0, len(idStrings))
	for _, str := range idStrings {
		if strings.ContainsRune(str, '-') {
			ids, err := parseIDRange(str)
			if err != nil {
				return nil, err
			}
			result = append(result, ids...)
			continue
		}
		val, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return nil, err
		}
		result = append(result, baos.DatapointID(val))
	}
	return result, nil
}

func parseIDRange(str string) ([]baos.DatapointID, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id range: %s", str)
	}
	start, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return nil, err
	}
	end, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return nil, err
	}
	result := make([]baos.DatapointID, 0, end-start)
	for i := start; i <= end; i++ {
		result = append(result, baos.DatapointID(i))
	}
	return result, nil
}
