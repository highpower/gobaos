package config

import (
	"fmt"
	"thermo-watchdog/internal/pattern"
	"time"
)

type idMapData struct {
	ReadOnly  int `yaml:"read-only"`
	Updatable int `yaml:"updatable"`
}

type datapointsData struct {
	Pattern          string        `yaml:"pattern"`
	UpdatablePattern string        `yaml:"updatable-pattern"`
	Every            time.Duration `yaml:"every"`
	IDMap            []idMapData   `yaml:"id-map"`
}

func (s *datapointsData) valid() error {
	if _, err := pattern.Parse(s.Pattern); err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}
	if _, err := pattern.Parse(s.UpdatablePattern); err != nil {
		return fmt.Errorf("invalid updatable pattern: %w", err)
	}
	if s.Every == 0 {
		return fmt.Errorf("positive value for every is required")
	}
	all := make(map[int]bool)
	for _, m := range s.IDMap {
		if m.ReadOnly == 0 || m.Updatable == 0 {
			return fmt.Errorf("invalid id mapping: %d, %d", m.ReadOnly, m.Updatable)
		}
		if _, ok := all[m.ReadOnly]; ok {
			return fmt.Errorf("duplicate read-only id: %d", m.ReadOnly)
		}
		if _, ok := all[m.Updatable]; ok {
			return fmt.Errorf("duplicate updatable id: %d", m.Updatable)
		}
		all[m.ReadOnly], all[m.Updatable] = true, true
	}
	return nil
}
