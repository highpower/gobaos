package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerSettings interface {
	Level() string
	Output() string
}

func newLogger(settings loggerSettings) (*zap.SugaredLogger, error) {
	level, err := zapcore.ParseLevel(settings.Level())
	if err != nil {
		return nil, err
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	cf := zap.Config{
		Level:             atomicLevel,
		Encoding:          "json",
		OutputPaths:       []string{settings.Output()},
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
