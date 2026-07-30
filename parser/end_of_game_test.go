package parser

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlgebraicParserResultMarkers(t *testing.T) {
	testCases := []struct {
		name   string
		marker string
	}{
		{"ASCII 1-0", "1-0"},
		{"en-dash 1–0", "1–0"},
		{"ASCII 0-1", "0-1"},
		{"en-dash 0–1", "0–1"},
		{"ASCII ½-½", "½-½"},
		{"en-dash ½–½", "½–½"},
		{"ASCII 1/2-1/2", "1/2-1/2"},
		{"en-dash 1/2–1/2", "1/2–1/2"},
		{"resigns", "resigns"},
		{"Resigns", "Resigns"},
		{"White resigns", "White resigns"},
		{"Black resigns", "Black resigns"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := "1. e4 e5\n2. " + tc.marker
			steps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(core.NewDefaultGame(), game)
			require.NoError(t, err, "result marker %q should parse", tc.marker)
			require.Len(t, steps, 3)
			assert.True(t, steps[2].StepAction.IsResign, "last step should be the terminal action")
		})
	}
}

func TestDescriptiveParserResultMarkers(t *testing.T) {
	testCases := []struct {
		name   string
		marker string
	}{
		{"ASCII 1-0", "1-0"},
		{"en-dash 1–0", "1–0"},
		{"ASCII ½-½", "½-½"},
		{"Resigns", "Resigns"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := "1. P-K4 P-K4\n2. " + tc.marker
			steps, err := NewNotationParserDescriptive(Characteristics{}).Parse(core.NewDefaultGame(), game)
			require.NoError(t, err, "result marker %q should parse", tc.marker)
			require.Len(t, steps, 3)
			assert.True(t, steps[2].StepAction.IsResign, "last step should be the terminal action")
		})
	}
}
