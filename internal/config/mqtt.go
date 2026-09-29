package config

import (
	"errors"
	"fmt"
	"gobaos/internal/mask"
)

type mqttData struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	CaCert   string `yaml:"ca-cert"`
}

func (s *mqttData) String() string {
	return fmt.Sprintf("mqtt[host='%s',user='%s',password='%s',ca-cert='%s']",
		s.Host, s.User, mask.Apply(s.Password), s.CaCert)
}

func (s *mqttData) valid() error {
	if s.Host == "" {
		return errors.New("empty mqtt host")
	}
	if (s.User == "" && s.Password != "") || (s.User != "" && s.Password == "") {
		return errors.New("empty mqtt user or password")
	}
	return nil
}
