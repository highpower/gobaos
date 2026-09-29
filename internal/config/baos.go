package config

type baosData struct {
	IndicationHandlers int `yaml:"indication-handlers"`
}

func (d *baosData) valid() error {
	return nil
}
