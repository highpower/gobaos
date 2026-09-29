package config

type baosData struct {
	Device             string `yaml:"device"`
	IndicationHandlers int    `yaml:"indication-handlers"`
}

func (d *baosData) valid() error {
	return nil
}
