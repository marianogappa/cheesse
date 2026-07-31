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
		isDraw bool
	}{
		{"ASCII 1-0", "1-0", false},
		{"en-dash 1–0", "1–0", false},
		{"ASCII 0-1", "0-1", false},
		{"en-dash 0–1", "0–1", false},
		{"ASCII ½-½", "½-½", true},
		{"en-dash ½–½", "½–½", true},
		{"ASCII 1/2-1/2", "1/2-1/2", true},
		{"en-dash 1/2–1/2", "1/2–1/2", true},
		{"resigns", "resigns", false},
		{"Resigns", "Resigns", false},
		{"White resigns", "White resigns", false},
		{"Black resigns", "Black resigns", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := "1. e4 e5\n2. " + tc.marker
			steps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(core.NewDefaultGame(), game)
			require.NoError(t, err, "result marker %q should parse", tc.marker)
			require.Len(t, steps, 3)
			if tc.isDraw {
				assert.True(t, steps[2].StepAction.IsDraw, "draw marker should match the draw action")
			} else {
				assert.True(t, steps[2].StepAction.IsResign, "decisive marker should match the resign action")
			}
		})
	}
}

func TestDescriptiveParserResultMarkers(t *testing.T) {
	testCases := []struct {
		name   string
		marker string
		isDraw bool
	}{
		{"ASCII 1-0", "1-0", false},
		{"en-dash 1–0", "1–0", false},
		{"ASCII ½-½", "½-½", true},
		{"Resigns", "Resigns", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := "1. P-K4 P-K4\n2. " + tc.marker
			steps, err := NewNotationParserDescriptive(Characteristics{}).Parse(core.NewDefaultGame(), game)
			require.NoError(t, err, "result marker %q should parse", tc.marker)
			require.Len(t, steps, 3)
			if tc.isDraw {
				assert.True(t, steps[2].StepAction.IsDraw, "draw marker should match the draw action")
			} else {
				assert.True(t, steps[2].StepAction.IsResign, "decisive marker should match the resign action")
			}
		})
	}
}
