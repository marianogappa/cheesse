package parser

import (
	"strings"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotationParserFigurine(t *testing.T) {
	testCases := []struct {
		name        string
		fen         string
		s           string
		expectedFEN string
	}{
		{
			name: "French Defense (ostinato six-notation suite)",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. e4 e6
				2. d4 d5
				3. ♘c3 ♝b4
				4. ♗b5+ ♝d7
				5. ♗xd7+ ♛xd7
				6. ♘ge2 dxe4
				7. 0-0`,
			expectedFEN: "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7",
		},
		{
			name:        "white knight move with white glyph",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:           "1. ♘f3",
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
		},
		{
			name:        "black knight move with black glyph",
			fen:         "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
			s:           "1. ♞f6",
			expectedFEN: "rnbqkb1r/pppppppp/5n2/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 2 2",
		},
		{
			name:        "glyphs are color-agnostic (black glyph for white piece)",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:           "1. ♞f3",
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
		},
		{
			name:        "queen capture with glyph",
			fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			s:           "1. exd5 ♛xd5",
			expectedFEN: "rnb1kbnr/ppp1pppp/8/3q4/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 3",
		},
		{
			name:        "promotion to queen glyph",
			fen:         "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. f8=♕",
			expectedFEN: "5Q2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "promotion to knight glyph with check",
			fen:         "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. f8=♘+",
			expectedFEN: "5N2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "rook move with disambiguation and glyph",
			fen:         "3k4/8/8/8/8/8/K7/R6R w - - 0 1",
			s:           "1. ♖hd1+",
			expectedFEN: "3k4/8/8/8/8/8/K7/R2R4 b - - 1 1",
		},
		{
			name:        "king move with glyph",
			fen:         "4k3/8/8/8/8/8/8/4K3 w - - 0 1",
			s:           "1. ♔e2",
			expectedFEN: "4k3/8/8/8/8/8/4K3/8 b - - 1 1",
		},
		{
			name:        "mixed figurine and letter moves in one game",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. ♘f3 Nc6
				2. e4 e5`,
			expectedFEN: "r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq e6 0 3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.s)
			require.NoError(t, err)
			require.NotEmpty(t, gameSteps)
			lastGame := gameSteps[len(gameSteps)-1].StepGame
			assert.Equal(t, tc.expectedFEN, lastGame.ToFEN(), "FEN mismatch")
		})
	}
}

func TestFigurinePrinter(t *testing.T) {
	testCases := []struct {
		name          string
		algGame       string
		expectedLines []string
	}{
		{
			name: "French Defense prints with color-aware glyphs",
			algGame: `1. e4 e6
				2. d4 d5
				3. Nc3 Bb4
				4. Bb5+ Bd7
				5. Bxd7+ Qxd7
				6. Ne2 dxe4
				7. 0-0`,
			expectedLines: []string{
				"1. e4 e6",
				"2. d4 d5",
				"3. ♘c3 ♝b4",
				"4. ♗b5+ ♝d7",
				"5. ♗xd7+ ♛xd7",
				// N.B. ♘e2 not ♘ge2: the c3 knight is pinned by the b4 bishop, so
				// only the g1 knight can reach e2 and no disambiguation is needed.
				"6. ♘e2 dxe4",
				"7. O-O",
			},
		},
		{
			name: "Scholar's mate with queen glyphs",
			algGame: `1. e4 e5
				2. Bc4 Nc6
				3. Qh5 Nf6
				4. Qxf7#`,
			expectedLines: []string{
				"1. e4 e5",
				"2. ♗c4 ♞c6",
				"3. ♕h5 ♞f6",
				"4. ♕xf7#",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := core.NewDefaultGame()
			algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.algGame)
			require.NoError(t, err)

			printed, err := printer.AlgebraicPrinter{}.PrintGame(algSteps, printer.FigurineCharacteristics())
			require.NoError(t, err)
			assert.Equal(t, tc.expectedLines, printed)
		})
	}
}

func TestFigurine_RoundTrip(t *testing.T) {
	testCases := []struct {
		name    string
		algGame string
	}{
		{
			name: "French Defense opening",
			algGame: `1. e4 e6
				2. d4 d5
				3. Nc3 Bb4
				4. Bb5+ Bd7
				5. Bxd7+ Qxd7
				6. Ne2 dxe4
				7. 0-0`,
		},
		{
			name: "Italian Game with castling",
			algGame: `1. e4 e5
				2. Nf3 Nc6
				3. Bc4 Bc5
				4. 0-0 Nf6
				5. d3 0-0`,
		},
		{
			name: "Scholar's mate",
			algGame: `1. e4 e5
				2. Bc4 Nc6
				3. Qh5 Nf6
				4. Qxf7#`,
		},
		{
			name: "Game with promotion",
			algGame: `1. d4 c5
				2. d5 e6
				3. dxe6 d5
				4. exf7+ Kd7
				5. fxg8=N`,
		},
		{
			name: "Ruy Lopez with knight disambiguation",
			algGame: `1. e4 e5
				2. Nf3 Nc6
				3. Bb5 a6
				4. Ba4 Nf6
				5. O-O Be7
				6. Re1 b5
				7. Bb3 O-O
				8. c3 d5`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := core.NewDefaultGame()

			algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.algGame)
			require.NoError(t, err, "failed to parse algebraic game")

			figurinePrinted, err := printer.AlgebraicPrinter{}.PrintGame(algSteps, printer.FigurineCharacteristics())
			require.NoError(t, err, "failed to print as figurine")

			figurineStr := strings.Join(figurinePrinted, "\n")

			figurineSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, figurineStr)
			require.NoError(t, err, "failed to re-parse figurine notation: %s", figurineStr)

			require.Len(t, figurineSteps, len(algSteps), "move count mismatch after round-trip")

			for i, orig := range algSteps {
				final := figurineSteps[i]
				assert.Equal(t, orig.StepAction.FromPiece.XY, final.StepAction.FromPiece.XY,
					"Move %d: from square mismatch", i+1)
				assert.Equal(t, orig.StepAction.ToXY, final.StepAction.ToXY,
					"Move %d: to square mismatch", i+1)
				assert.Equal(t, orig.StepGame.ToFEN(), final.StepGame.ToFEN(),
					"Move %d: resulting FEN mismatch", i+1)
			}
		})
	}
}

func TestFigurine_CrossNotationEquivalence(t *testing.T) {
	algGame := `1. e4 e6
		2. d4 d5
		3. Nc3 Bb4
		4. Bb5+ Bd7
		5. Bxd7+ Qxd7
		6. Ne2 dxe4
		7. 0-0`

	figurineGame := `1. e4 e6
		2. d4 d5
		3. ♘c3 ♝b4
		4. ♗b5+ ♝d7
		5. ♗xd7+ ♛xd7
		6. ♘ge2 dxe4
		7. 0-0`

	g := core.NewDefaultGame()

	algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, algGame)
	require.NoError(t, err)

	figurineSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, figurineGame)
	require.NoError(t, err)

	require.Len(t, figurineSteps, len(algSteps))

	for i, alg := range algSteps {
		fig := figurineSteps[i]
		assert.Equal(t, alg.StepAction.FromPiece.XY, fig.StepAction.FromPiece.XY,
			"Move %d: from square mismatch", i+1)
		assert.Equal(t, alg.StepAction.ToXY, fig.StepAction.ToXY,
			"Move %d: to square mismatch", i+1)
		assert.Equal(t, alg.StepGame.ToFEN(), fig.StepGame.ToFEN(),
			"Move %d: FEN mismatch", i+1)
	}

	expectedFEN := "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7"
	assert.Equal(t, expectedFEN, figurineSteps[len(figurineSteps)-1].StepGame.ToFEN())
}
