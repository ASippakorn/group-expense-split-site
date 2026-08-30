package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPercentageParsingAndFormatting(t *testing.T) {
	basisPoints, err := parsePercentageBasisPoints("10.05")
	require.NoError(t, err)
	require.Equal(t, int64(1005), basisPoints)
	require.Equal(t, "10.05", formatPercentage(basisPoints))
	require.Equal(t, "50", formatPercentage(5000))
}

func TestPercentageParsingRejectsMoreThanTwoDecimalPlaces(t *testing.T) {
	_, err := parsePercentageBasisPoints("33.333")
	require.Error(t, err)
}
