package config

import (
	"errors"
	"fmt"
	"gobaos/pkg/baos"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type loggerData struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type mqttSettings struct {
	CaCert   string `yaml:"ca-cert"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
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

type Data struct {
	Device      string               `yaml:"device"`
	Logger      loggerData           `yaml:"logger"`
	Mqtt        mqttSettings         `yaml:"mqtt"`
	Sensors     []sensorSettings     `yaml:"sensors"`
	Controllers []controllerSettings `yaml:"controllers"`
}

func (l *loggerData) Valid() error {
	switch {
	case l.Level == "":
		return fmt.Errorf("%w: empty log level", ErrInvalid)
	case l.Output == "":
		return fmt.Errorf("%w: empty log output", ErrInvalid)
	}
	return nil
}

func (s *mqttSettings) Valid() error {
	if s.Host == "" {
		return fmt.Errorf("%w: empty mqtt host", ErrInvalid)
	}
	if (s.User != "" && s.Password == "") || (s.User == "" && s.Password != "") {
		return fmt.Errorf("%w: invalid mqtt user/password", ErrInvalid)
	}
	return nil
}

func (s *sensorSettings) Valid() error {
	if len(s.ids) == 0 {
		return fmt.Errorf("%w: no sensor id found", ErrInvalid)
	}
	if s.StateTopic == "" {
		return fmt.Errorf("%w: empty state topic", ErrInvalid)
	}
	return nil
}

func (s *controllerSettings) Valid() error {
	if len(s.ids) == 0 {
		return fmt.Errorf("%w: no controller id found", ErrInvalid)
	}
	switch {
	case s.StateTopic == "":
		return fmt.Errorf("%w: empty state topic", ErrInvalid)
	case s.CommandTopic == "":
		return fmt.Errorf("%w: empty command topic", ErrInvalid)
	}
	return nil
}

func (d *Data) Valid() error {
	if d.Device == "" {
		return fmt.Errorf("%w: missing device", ErrInvalid)
	}
	if err := errors.Join(d.Logger.Valid(), d.Mqtt.Valid()); err != nil {
		return err
	}
	for _, s := range d.Sensors {
		if err := s.Valid(); err != nil {
			return err
		}
	}
	for _, c := range d.Controllers {
		if err := c.Valid(); err != nil {
			return err
		}
	}
	return nil
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
	if err := decoder.Decode(&result); err != nil {
		return Data{}, err
	}
	if err := result.parseSensorIDs(); err != nil {
		return Data{}, err
	}
	if err := result.parseControllerIDs(); err != nil {
		return Data{}, err
	}
	if err := result.Valid(); err != nil {
		return Data{}, err
	}
	return result, nil
}

func (d *Data) parseSensorIDs() error {
	for i := range d.Sensors {
		ids, err := parseIDStrings(d.Sensors[i].IDStrings)
		if err != nil {
			return err
		}
		d.Sensors[i].ids = ids
	}
	return nil
}

func (d *Data) parseControllerIDs() error {
	for i := range d.Controllers {
		ids, err := parseIDStrings(d.Controllers[i].IDStrings)
		if err != nil {
			return err
		}
		d.Controllers[i].ids = ids
	}
	return nil
}

func parseIDStrings(idStrings []string) ([]baos.DatapointID, error) {
	result := make([]baos.DatapointID, 0, len(idStrings))
	for _, str := range idStrings {
		if strings.ContainsRune(str, '-') {
			ids, err := parseDatapointIDRange(str)
			if err != nil {
				return nil, err
			}
			result = append(result, ids...)
			continue
		}
		val, err := parseDatapointID(str)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

func parseDatapointID(str string) (baos.DatapointID, error) {
	val, err := strconv.ParseUint(str, 10, 16)
	switch {
	case err != nil:
		return 0, err
	case val == 0:
		return 0, fmt.Errorf("%w: datapoint id == 0", ErrInvalid)
	}
	return baos.DatapointID(val), nil
}

func parseDatapointIDRange(str string) ([]baos.DatapointID, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id range: %s", str)
	}
	end, err := parseDatapointID(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, err
	}
	start, err := parseDatapointID(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, err
	}
	result := make([]baos.DatapointID, 0, end-start)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result, nil
}
