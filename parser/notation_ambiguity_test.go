package parser

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotationParserAlgebraicAmbiguity(t *testing.T) {
	testCases := []struct {
		name        string
		fen         string
		s           string
		expectError bool
	}{
		{
			name:        "two knights reaching the same square require disambiguation",
			fen:         "4k3/8/8/8/8/8/3N1N2/4K3 w - - 0 1",
			s:           "1. Ne4",
			expectError: true,
		},
		{
			name:        "file disambiguation resolves two knights",
			fen:         "4k3/8/8/8/8/8/3N1N2/4K3 w - - 0 1",
			s:           "1. Nde4",
			expectError: false,
		},
		{
			name:        "two rooks on same rank require disambiguation",
			fen:         "3k4/8/8/8/8/8/K7/R6R w - - 0 1",
			s:           "1. Rd1",
			expectError: true,
		},
		{
			name:        "file disambiguation resolves two rooks on same rank",
			fen:         "3k4/8/8/8/8/8/K7/R6R w - - 0 1",
			s:           "1. Rhd1",
			expectError: false,
		},
		{
			name:        "two rooks on same file require rank disambiguation",
			fen:         "R2k4/8/8/8/8/8/8/R6K w - - 0 1",
			s:           "1. Ra4",
			expectError: true,
		},
		{
			name:        "rank disambiguation resolves two rooks on same file",
			fen:         "R2k4/8/8/8/8/8/8/R6K w - - 0 1",
			s:           "1. R1a4",
			expectError: false,
		},
		{
			name:        "unambiguous knight move parses fine",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:           "1. Nf3",
			expectError: false,
		},
		{
			name:        "two knights capturing the same piece require disambiguation",
			fen:         "4k3/8/8/8/4p3/8/3N1N2/4K3 w - - 0 1",
			s:           "1. Nxe4",
			expectError: true,
		},
		{
			name:        "file disambiguation resolves two knights capturing the same piece",
			fen:         "4k3/8/8/8/4p3/8/3N1N2/4K3 w - - 0 1",
			s:           "1. Ndxe4",
			expectError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			_, err = NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.s)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "ambiguous")
				assert.Contains(t, err.Error(), "please disambiguate")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNotationParserDescriptiveAmbiguity(t *testing.T) {
	testCases := []struct {
		name        string
		fen         string
		s           string
		expectError bool
	}{
		{
			name: "two rooks reaching the same square require disambiguation",
			// Rooks on QR1 (a1) and KR1 (h1); R-Q1 (d1) is reachable by both.
			fen:         "3k4/8/8/8/8/8/K7/R6R w - - 0 1",
			s:           "1. R-Q1",
			expectError: true,
		},
		{
			name: "unambiguous rook move parses fine",
			// Only the KR1 (h1) rook can reach KR5 (h5).
			fen:         "3k4/8/8/8/8/8/K7/R6R w - - 0 1",
			s:           "1. R-KR5",
			expectError: false,
		},
		{
			name:        "unambiguous pawn move parses fine",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:           "1. P-K4",
			expectError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			_, err = NewNotationParserDescriptive(Characteristics{}).Parse(g, tc.s)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "ambiguous")
				assert.Contains(t, err.Error(), "please disambiguate")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
