package parser

import (
	"strings"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotationParserSmith(t *testing.T) {
	testCases := []struct {
		name        string
		fen         string
		s           string
		expectedFEN string
		expectedErr bool
	}{
		{
			name: "French Defense (ostinato six-notation suite)",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. e2e4  e7e6
				2. d2d4  d7d5
				3. b1c3  f8b4
				4. f1b5  c8d7
				5. b5d7b d8d7b
				6. g1e2  d5e4p
				7. e1g1c`,
			expectedFEN: "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7",
		},
		{
			name: "simple opening moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:    "1. e2e4 e7e5",
			expectedFEN: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name: "capture with pawn",
			fen:  "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			s:    "1. e4d5p",
			expectedFEN: "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		},
		{
			name: "capture with knight",
			fen:  "r1bqkbnr/pppppppp/2n5/4P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
			s:    "1. c6e5p",
			expectedFEN: "r1bqkbnr/pppppppp/8/4n3/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 3",
		},
		{
			name: "kingside castling white",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1",
			s:    "1. e1g1c",
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 b kq - 1 1",
		},
		{
			name: "queenside castling white",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3KBNR w KQkq - 0 1",
			s:    "1. e1c1C",
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/2KR1BNR b kq - 1 1",
		},
		{
			name: "kingside castling black",
			fen:  "rnbqk2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			s:    "1. e8g8c",
			expectedFEN: "rnbq1rk1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 1 2",
		},
		{
			name: "queenside castling black",
			fen:  "r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			s:    "1. e8c8C",
			expectedFEN: "2kr1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 1 2",
		},
		{
			name: "promotion to queen",
			fen:  "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:    "1. f7f8Q",
			expectedFEN: "5Q2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name: "promotion to knight",
			fen:  "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:    "1. f7f8N",
			expectedFEN: "5N2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name: "capture-promotion",
			fen:  "3r4/4P2k/8/8/8/8/8/K7 w - - 0 1",
			s:    "1. e7d8rQ",
			expectedFEN: "3Q4/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name: "en passant capture",
			fen:  "rnbqkbnr/pppp1ppp/8/4pP2/8/8/PPPPP1PP/RNBQKBNR w KQkq e6 0 1",
			s:    "1. f5e6E",
			expectedFEN: "rnbqkbnr/pppp1ppp/4P3/8/8/8/PPPPP1PP/RNBQKBNR b KQkq - 0 1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserSmith(Characteristics{}).Parse(g, tc.s)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, gameSteps)
			lastGame := gameSteps[len(gameSteps)-1].StepGame
			assert.Equal(t, tc.expectedFEN, lastGame.ToFEN(), "FEN mismatch")
		})
	}
}

func TestSmithParser_RoundTrip(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := core.NewDefaultGame()

			algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.algGame)
			require.NoError(t, err, "failed to parse algebraic game")

			smithPrinted, err := printer.SmithPrinter{}.PrintGame(algSteps, printer.GameCharacteristics{})
			require.NoError(t, err, "failed to print as Smith")

			smithStr := strings.Join(smithPrinted, "\n")

			smithSteps, err := NewNotationParserSmith(Characteristics{}).Parse(g, smithStr)
			require.NoError(t, err, "failed to parse Smith notation: %s", smithStr)

			require.Len(t, smithSteps, len(algSteps), "move count mismatch after round-trip")

			for i, orig := range algSteps {
				final := smithSteps[i]
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

func TestSmithParser_CrossNotationEquivalence(t *testing.T) {
	algGame := `1. e4 e6
		2. d4 d5
		3. Nc3 Bb4
		4. Bb5+ Bd7
		5. Bxd7+ Qxd7
		6. Ne2 dxe4
		7. 0-0`

	smithGame := `1. e2e4  e7e6
		2. d2d4  d7d5
		3. b1c3  f8b4
		4. f1b5  c8d7
		5. b5d7b d8d7b
		6. g1e2  d5e4p
		7. e1g1c`

	g := core.NewDefaultGame()

	algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, algGame)
	require.NoError(t, err)

	smithSteps, err := NewNotationParserSmith(Characteristics{}).Parse(g, smithGame)
	require.NoError(t, err)

	require.Len(t, smithSteps, len(algSteps))

	for i, alg := range algSteps {
		smith := smithSteps[i]
		assert.Equal(t, alg.StepAction.FromPiece.XY, smith.StepAction.FromPiece.XY,
			"Move %d: from square mismatch", i+1)
		assert.Equal(t, alg.StepAction.ToXY, smith.StepAction.ToXY,
			"Move %d: to square mismatch", i+1)
		assert.Equal(t, alg.StepGame.ToFEN(), smith.StepGame.ToFEN(),
			"Move %d: FEN mismatch", i+1)
	}

	expectedFEN := "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7"
	assert.Equal(t, expectedFEN, smithSteps[len(smithSteps)-1].StepGame.ToFEN())
}
