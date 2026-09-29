package mqttbridge

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
}
