package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func parseLevel(s string) (zapcore.Level, error) {
	result := zapcore.InfoLevel
	if s != "" {
		level, err := zapcore.ParseLevel(s)
		if err != nil {
			return zapcore.InvalidLevel, err
		}
		result = level
	}
	return result, nil
}

func newLogger(l, out string) (*zap.SugaredLogger, error) {
	level, err := parseLevel(l)
	if err != nil {
		return nil, err
	}
	atomicLevel := zap.NewAtomicLevelAt(level)
	cf := zap.Config{
		Level:             atomicLevel,
		Encoding:          "json",
		OutputPaths:       []string{out},
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
