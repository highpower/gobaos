package mqttbridge

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
}

type nopLogger struct{}

func (n nopLogger) Infof(_ string, _ ...any) {
}

func (n nopLogger) Errorf(_ string, _ ...any) {
}

func (n nopLogger) Debugf(_ string, _ ...any) {
}
