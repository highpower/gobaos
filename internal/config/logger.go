package config

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap/zapcore"
)

type loggerSettings struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func (s *loggerSettings) String() string {
	return fmt.Sprintf("logger[level='%s',output='%s']", s.Level, s.Output)
}

func (s *loggerSettings) valid() error {
	if _, err := zapcore.ParseLevel(s.Level); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	switch {
	case s.Output == "":
		return fmt.Errorf("empty log output")
	case strings.EqualFold(s.Output, "stderr") || strings.EqualFold(s.Output, "stdout"):
		return nil
	}
	stat, err := os.Stat(s.Output)
	switch {
	case err != nil:
		return fmt.Errorf("invalid log output: %w", err)
	case !stat.Mode().IsRegular():
		return fmt.Errorf("invalid log output: %s is not regular file", s.Output)
	}
	return nil
}
