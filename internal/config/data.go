package config

import (
	"bytes"
	"errors"
	"fmt"
	"gobaos/pkg/baos"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Mqtt      mqttData      `yaml:"mqtt"`
	Baos      baosData      `yaml:"baos"`
	Logger    loggerData    `yaml:"logger"`
	ReadOnly  datapointData `yaml:"read-only"`
	Updatable datapointData `yaml:"updatable"`
	readOnly  map[int]bool
	updatable map[int]bool
}

func (d *Data) Host() string {
	return d.Mqtt.Host
}

func (d *Data) User() string {
	return d.Mqtt.User
}

func (d *Data) Password() string {
	return d.Mqtt.Password
}

func (d *Data) CaCert() string {
	return d.Mqtt.CaCert
}

func (d *Data) Level() string {
	return d.Logger.Level
}

func (d *Data) Output() string {
	return d.Logger.Output
}

func (d *Data) Pattern() string {
	return d.ReadOnly.Pattern
}

func (d *Data) UpdatablePattern() string {
	return d.Updatable.Pattern
}

func (d *Data) ReadOnlyContains(id baos.DatapointID) bool {
	_, ok := d.readOnly[int(id)]
	return ok
}

func (d *Data) UpdatableContains(id baos.DatapointID) bool {
	_, ok := d.updatable[int(id)]
	return ok
}

func (d *Data) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("config[")
	if err := d.dump(&buffer, ",", false); err != nil {
		panic(err)
	}
	buffer.WriteString("]")
	return buffer.String()
}

func (d *Data) dump(writer io.Writer, sep string, always bool) error {
	parts := []fmt.Stringer{&d.Logger, &d.Mqtt}
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

func (d *Data) valid() error {
	return errors.Join(d.Logger.valid(), d.Mqtt.valid(), d.Baos.valid(),
		d.ReadOnly.valid(), d.Updatable.valid())
}

func Parse(name string) (*Data, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("config.parse: %w", err)
	}
	defer func() { _ = file.Close() }()
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	result := Data{}
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
