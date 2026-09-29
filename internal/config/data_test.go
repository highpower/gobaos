package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	_, err := Parse("testdata/config.yaml")
	require.NoError(t, err)

}
