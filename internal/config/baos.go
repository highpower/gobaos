package config

type baosSettings struct {
	Device             string `yaml:"device"`
	IndicationHandlers int    `yaml:"indication-handlers"`
}

func (s *baosSettings) valid() error {
	return nil
}
