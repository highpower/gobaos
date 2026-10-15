package config

import (
	"fmt"
	"strconv"
	"strings"

	"gobaos/internal/pattern"
)

type datapointSettings struct {
	IDs     []string `yaml:"ids"`
	Pattern string   `yaml:"pattern"`
}

func (s *datapointSettings) valid() error {
	for _, id := range s.IDs {
		if _, err := parseIDs(id); err != nil {
			return fmt.Errorf("invalid ids: %w", err)
		}
	}
	if _, err := pattern.Parse(s.Pattern); err != nil {
		return fmt.Errorf("invalid read-only pattern: %w", err)
	}
	return nil
}

func parseIDs(s string) ([]int, error) {
	result := make([]int, 0)
	parts := strings.SplitSeq(s, ",")
	for p := range parts {
		p = strings.TrimSpace(p)
		if strings.ContainsRune(p, '-') {
			ids, err := parseIDRange(p)
			if err != nil {
				return nil, fmt.Errorf("invalid ids: %w", err)
			}
			result = append(result, ids...)
			continue
		}
		value, err := strconv.ParseInt(p, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid ids: %w", err)
		}
		result = append(result, int(value))
	}
	return result, nil
}

func parseIDRange(s string) ([]int, error) {
	result := make([]int, 0)
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id range: '%s'", s)
	}
	left, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid id range: %w", err)
	}
	right, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid id range: %w", err)
	}
	if left > right {
		return nil, fmt.Errorf("invalid id range: %d > %d", left, right)
	}
	for i := left; i <= right; i++ {
		result = append(result, i)
	}
	return result, nil
}
