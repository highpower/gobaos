package config

import (
	"fmt"
	"gobaos/pkg/baos"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

type loggerData struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

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

type mqttSettings struct {
	CaCert   string `yaml:"ca-cert"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Data struct {
	Device      string               `yaml:"device"`
	Logger      loggerData           `yaml:"logger"`
	Mqtt        mqttSettings         `yaml:"mqtt"`
	Sensors     []sensorSettings     `yaml:"sensors"`
	Controllers []controllerSettings `yaml:"controllers"`
}

func Parse(name string) (Data, error) {
	file, err := os.Open(name)
	if err != nil {
		return Data{}, fmt.Errorf("cat not read config file %w", err)
	}
	defer func() { _ = file.Close() }()
	result := Data{}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err = decoder.Decode(&result); err != nil {
		return Data{}, err
	}
	if err := result.parseSensorIDs(); err != nil {
		return Data{}, err
	}
	if err := result.parseControllerIDs(); err != nil {
		return Data{}, err
	}
	return result, nil
}

func (c *Data) parseSensorIDs() error {
	for _, sensor := range c.Sensors {
		ids, err := parseIDStrings(sensor.IDStrings)
		if err != nil {
			return err
		}
		sensor.ids = ids
	}
	return nil
}

func (c *Data) parseControllerIDs() error {
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
