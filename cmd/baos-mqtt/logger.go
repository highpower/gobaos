package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerSettingsProvider interface {
	Level() string
	Output() string
}

func newLogger(sp loggerSettingsProvider) (*zap.SugaredLogger, error) {
	level, err := zapcore.ParseLevel(sp.Level())
	if err != nil {
		return nil, err
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	cf := zap.Config{
		Level:             atomicLevel,
		Encoding:          "json",
		OutputPaths:       []string{sp.Output()},
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeTime:  zapcore.TimeEncoderOfLayout(time.RFC3339),
			MessageKey:  "message",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}
	log, err := cf.Build()
	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}
