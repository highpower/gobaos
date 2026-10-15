package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"gobaos/pkg/baos"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Mqtt      mqttSettings      `yaml:"mqtt"`
	Baos      baosSettings      `yaml:"baos"`
	Logger    loggerSettings    `yaml:"logger"`
	ReadOnly  datapointSettings `yaml:"read-only"`
	Updatable datapointSettings `yaml:"updatable"`
	readOnly  map[int]bool
	updatable map[int]bool
}

func (s *Settings) Host() string {
	return s.Mqtt.Host
}

func (s *Settings) User() string {
	return s.Mqtt.User
}

func (s *Settings) Password() string {
	return s.Mqtt.Password
}

func (s *Settings) CaCert() string {
	return s.Mqtt.CaCert
}

func (s *Settings) Level() string {
	return s.Logger.Level
}

func (s *Settings) Output() string {
	return s.Logger.Output
}

func (s *Settings) Pattern() string {
	return s.ReadOnly.Pattern
}

func (s *Settings) UpdatablePattern() string {
	return s.Updatable.Pattern
}

func (s *Settings) ReadOnlyContains(id baos.DatapointID) bool {
	_, ok := s.readOnly[int(id)]
	return ok
}

func (s *Settings) UpdatableContains(id baos.DatapointID) bool {
	_, ok := s.updatable[int(id)]
	return ok
}

func (s *Settings) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("config[")
	if err := s.dump(&buffer, ",", false); err != nil {
		panic(err)
	}
	buffer.WriteString("]")
	return buffer.String()
}

func (s *Settings) dump(writer io.Writer, sep string, always bool) error {
	parts := []fmt.Stringer{&s.Logger, &s.Mqtt}
	for i, s := range parts {
		if _, err := io.WriteString(writer, s.String()); err != nil {
			return err
		}
		if i == len(parts)-1 && !always {
			break
		}
		if _, err := io.WriteString(writer, sep); err != nil {
			return err
		}
	}
	return nil
}

func (s *Settings) valid() error {
	return errors.Join(s.Logger.valid(), s.Mqtt.valid(), s.Baos.valid(),
		s.ReadOnly.valid(), s.Updatable.valid())
}

func Parse(name string) (*Settings, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	result := Settings{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}
	if err := result.valid(); err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}
	readOnly, err := constructMap(result.ReadOnly.IDs)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}
	result.readOnly = readOnly
	updatable, err := constructMap(result.Updatable.IDs)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}
	result.updatable = updatable
	return &result, nil
}

func constructMap(list []string) (map[int]bool, error) {
	m := make(map[int]bool)
	for _, ids := range list {
		idList, err := parseIDs(ids)
		if err != nil {
			return nil, fmt.Errorf("config.parse: %w", err)
		}
		for _, id := range idList {
			m[id] = true
		}
	}
	return m, nil
}
