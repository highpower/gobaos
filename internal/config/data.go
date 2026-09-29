package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Mqtt       mqttData       `yaml:"mqtt"`
	Logger     loggerData     `yaml:"logger"`
	Datapoints datapointsData `yaml:"datapoints"`
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
	return d.Datapoints.Pattern
}

func (d *Data) UpdatablePattern() string {
	return d.Datapoints.UpdatablePattern
}

func (d *Data) Every() time.Duration {
	return d.Datapoints.Every
}

func (d *Data) VisitAll(f func(pattern int, updatable int) error) error {
	for _, m := range d.Datapoints.IDMap {
		if err := f(m.ReadOnly, m.Updatable); err != nil {
			return err
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
	return errors.Join(d.Logger.valid(), d.Mqtt.valid(), d.Datapoints.valid())
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
