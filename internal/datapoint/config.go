package datapoint

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

type readOnlySettings struct {
	ids        []baos.DatapointID
	IDStrings  []string `yaml:"ids"`
	StateTopic string   `yaml:"state-topic"`
}
type updatableSettings struct {
	ids         []baos.DatapointID
	IDStrings   []string `yaml:"ids"`
	StateTopic  string   `yaml:"state-topic"`
	UpdateTopic string   `yaml:"update-topic"`
}

type Config struct {
	Device    string              `yaml:"device"`
	Logger    loggerData          `yaml:"logger"`
	ReadOnly  []readOnlySettings  `yaml:"read-only"`
	Updatable []updatableSettings `yaml:"updatable"`
	Mqtt      mqttSettings        `yaml:"mqtt"`
}

func (l *loggerData) Valid() error {
	switch {
	case l.Level == "":
		return fmt.Errorf("%w: empty log level", ErrBadConfig)
	case l.Output == "":
		return fmt.Errorf("%w: empty log output", ErrBadConfig)
	}
	return nil
}

func (s *mqttSettings) Valid() error {
	if s.Host == "" {
		return fmt.Errorf("%w: empty mqtt host", ErrBadConfig)
	}
	if (s.User != "" && s.Password == "") || (s.User == "" && s.Password != "") {
		return fmt.Errorf("%w: invalid mqtt user/password", ErrBadConfig)
	}
	return nil
}

func (s *readOnlySettings) Valid() error {
	if len(s.ids) == 0 {
		return fmt.Errorf("%w: no read-only id found", ErrBadConfig)
	}
	if s.StateTopic == "" {
		return fmt.Errorf("%w: empty state topic", ErrBadConfig)
	}
	return nil
}

func (s *updatableSettings) Valid() error {
	if len(s.ids) == 0 {
		return fmt.Errorf("%w: no updatable id found", ErrBadConfig)
	}
	switch {
	case s.StateTopic == "":
		return fmt.Errorf("%w: empty state topic", ErrBadConfig)
	case s.UpdateTopic == "":
		return fmt.Errorf("%w: empty update topic", ErrBadConfig)
	}
	return nil
}

func (d *Config) Valid() error {
	if d.Device == "" {
		return fmt.Errorf("%w: missing device", ErrBadConfig)
	}
	if err := errors.Join(d.Logger.Valid(), d.Mqtt.Valid()); err != nil {
		return err
	}
	for _, s := range d.ReadOnly {
		if err := s.Valid(); err != nil {
			return err
		}
	}
	for _, c := range d.Updatable {
		if err := c.Valid(); err != nil {
			return err
		}
	}
	return nil
}

func NewConfig(name string) (Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return Config{}, fmt.Errorf("newConfig: %w", err)
	}
	defer func() { _ = file.Close() }()
	result := Config{}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&result); err != nil {
		return Config{}, err
	}
	if err := result.parseSensorIDs(); err != nil {
		return Config{}, err
	}
	if err := result.parseControllerIDs(); err != nil {
		return Config{}, err
	}
	if err := result.Valid(); err != nil {
		return Config{}, err
	}
	return result, nil
}

func (d *Config) parseSensorIDs() error {
	for i := range d.ReadOnly {
		ids, err := parseIDStrings(d.ReadOnly[i].IDStrings)
		if err != nil {
			return err
		}
		d.ReadOnly[i].ids = ids
	}
	return nil
}

func (d *Config) parseControllerIDs() error {
	for i := range d.Updatable {
		ids, err := parseIDStrings(d.Updatable[i].IDStrings)
		if err != nil {
			return err
		}
		d.Updatable[i].ids = ids
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
		return 0, fmt.Errorf("%w: datapoint id == 0", ErrBadConfig)
	}
	return baos.DatapointID(val), nil
}

func parseDatapointIDRange(str string) ([]baos.DatapointID, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: invalid id range %s", ErrBadConfig, str)
	}
	end, err := parseDatapointID(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, err
	}
	start, err := parseDatapointID(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, err
	}
	if start > end {
		return nil, fmt.Errorf("%w: invalid id range %s", ErrBadConfig, str)
	}
	result := make([]baos.DatapointID, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result, nil
}
