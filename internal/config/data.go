package config

import (
	"bytes"
	"errors"
	"fmt"
	"gobaos/internal/pattern"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Mqtt      mqttData        `yaml:"mqtt"`
	Baos      baosData        `yaml:"baos"`
	Logger    loggerData      `yaml:"logger"`
	ReadOnly  []datapointData `yaml:"read-only"`
	Updatable []datapointData `yaml:"updatable"`
	readOnly  []idData
	updatable []idData
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

func (d *Data) VisitReadOnly(f func(id int, pattern pattern.Parsed) error) error {
	for _, ip := range d.readOnly {
		if err := f(ip.id, ip.pat); err != nil {
			return fmt.Errorf("data.visitReadOnly: %w", err)
		}
	}
	return nil
}

func (d *Data) VisitUpdatable(f func(id int, pattern pattern.Parsed) error) error {
	for _, ip := range d.updatable {
		if err := f(ip.id, ip.pat); err != nil {
			return fmt.Errorf("data.visitUpdatable: %w", err)
		}
	}
	return nil
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
	for _, dp := range d.ReadOnly {
		if err := dp.valid(); err != nil {
			return err
		}
	}
	for _, dp := range d.Updatable {
		if err := dp.valid(); err != nil {
			return err
		}
	}
	return errors.Join(d.Logger.valid(), d.Mqtt.valid(), d.Baos.valid())
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
	return &result, nil
}
