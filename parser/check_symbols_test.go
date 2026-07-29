package parser

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlgebraicParserCheckSymbols(t *testing.T) {
	testCases := []struct {
		name string
		fen  string
		s    string
	}{
		{
			name: "dblch symbol accepted for double check",
			fen:  "4k3/8/8/3N4/8/8/8/4R2K w - - 0 1",
			s:    "1. Nf6dblch",
		},
		{
			name: "dbl.ch symbol accepted for double check",
			fen:  "4k3/8/8/3N4/8/8/8/4R2K w - - 0 1",
			s:    "1. Nf6dbl.ch",
		},
		{
			name: "++ symbol accepted for double check",
			fen:  "4k3/8/8/3N4/8/8/8/4R2K w - - 0 1",
			s:    "1. Nf6++",
		},
		{
			name: "disch symbol accepted for discovered check",
			fen:  "4k3/8/8/8/8/8/3N4/4R2K w - - 0 1",
			s:    "1. Nf3disch",
		},
		{
			name: "dis.ch symbol accepted for discovered check",
			fen:  "4k3/8/8/8/8/8/3N4/4R2K w - - 0 1",
			s:    "1. Nf3dis.ch",
		},
		{
			name: "+ symbol accepted for any check",
			fen:  "4k3/8/8/3N4/8/8/8/4R2K w - - 0 1",
			s:    "1. Nf6+",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			steps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.s)
			require.NoError(t, err)
			require.Len(t, steps, 1)
			assert.True(t, steps[0].StepGame.IsCheck)
		})
	}
}

func TestAlgebraicParserCheckSymbolsGameState(t *testing.T) {
	t.Run("double check sets IsDoubleCheck on game", func(t *testing.T) {
		g, err := core.NewGameFromFEN("4k3/8/8/3N4/8/8/8/4R2K w - - 0 1")
		require.NoError(t, err)
		steps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, "1. Nf6++")
		require.NoError(t, err)
		assert.True(t, steps[0].StepGame.IsDoubleCheck)
		assert.True(t, steps[0].StepGame.IsDiscoverCheck)
	})

	t.Run("discovered check sets IsDiscoverCheck on game", func(t *testing.T) {
		g, err := core.NewGameFromFEN("4k3/8/8/8/8/8/3N4/4R2K w - - 0 1")
		require.NoError(t, err)
		steps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, "1. Nf3+")
		require.NoError(t, err)
		assert.True(t, steps[0].StepGame.IsDiscoverCheck)
		assert.False(t, steps[0].StepGame.IsDoubleCheck)
	})
}
