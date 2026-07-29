package parser

import (
	"strings"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotationParserCoordinate(t *testing.T) {
	testCases := []struct {
		name        string
		fen         string
		s           string
		expectedFEN string
	}{
		{
			name: "French Defense (ostinato six-notation suite)",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. e2-e4 e7-e6
				2. d2-d4 d7-d5
				3. b1-c3 f8-b4
				4. f1-b5+ c8-d7
				5. b5xd7+ d8xd7
				6. g1-e2 d5xe4
				7. 0-0`,
			expectedFEN: "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7",
		},
		{
			name:        "simple opening moves with dashes",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:           "1. e2-e4 e7-e5",
			expectedFEN: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name:        "capture with x delimiter",
			fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			s:           "1. e4xd5",
			expectedFEN: "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		},
		{
			name:        "capture with colon delimiter",
			fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			s:           "1. e4:d5",
			expectedFEN: "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		},
		{
			name:        "capture with dash delimiter also works",
			fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2",
			s:           "1. e4-d5",
			expectedFEN: "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
		},
		{
			name:        "kingside castling white with zeroes",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1",
			s:           "1. 0-0",
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 b kq - 1 1",
		},
		{
			name:        "queenside castling white with Os",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3KBNR w KQkq - 0 1",
			s:           "1. O-O-O",
			expectedFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/2KR1BNR b kq - 1 1",
		},
		{
			name:        "promotion with = symbol",
			fen:         "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. f7-f8=Q",
			expectedFEN: "5Q2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "promotion bare letter",
			fen:         "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. f7-f8Q",
			expectedFEN: "5Q2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "promotion with parentheses",
			fen:         "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. f7-f8(Q)",
			expectedFEN: "5Q2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "promotion with slash",
			fen:         "8/5P1k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. f7-f8/N",
			expectedFEN: "5N2/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "capture-promotion",
			fen:         "3r4/4P2k/8/8/8/8/8/K7 w - - 0 1",
			s:           "1. e7xd8=Q",
			expectedFEN: "3Q4/7k/8/8/8/8/8/K7 b - - 0 1",
		},
		{
			name:        "en passant capture",
			fen:         "rnbqkbnr/pppp1ppp/8/4pP2/8/8/PPPPP1PP/RNBQKBNR w KQkq e6 0 1",
			s:           "1. f5xe6",
			expectedFEN: "rnbqkbnr/pppp1ppp/4P3/8/8/8/PPPPP1PP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:        "check suffix with +",
			fen:         "4k3/8/8/8/8/8/8/3Q2K1 w - - 0 1",
			s:           "1. d1-d8+",
			expectedFEN: "3Qk3/8/8/8/8/8/8/6K1 b - - 1 1",
		},
		{
			name:        "check suffix with ch",
			fen:         "4k3/8/8/8/8/8/8/3Q2K1 w - - 0 1",
			s:           "1. d1-d8ch",
			expectedFEN: "3Qk3/8/8/8/8/8/8/6K1 b - - 1 1",
		},
		{
			name:        "no delimiter (compact coordinate)",
			fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:           "1. e2e4 e7e5",
			expectedFEN: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name:        "++ suffix accepted on checkmate move",
			fen:         "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 3",
			s:           "1. h5xf7++",
			expectedFEN: "r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 3",
		},
		{
			name:        "++ suffix accepted on double check move",
			fen:         "4k3/8/8/3N4/8/8/8/4R2K w - - 0 1",
			s:           "1. d5-f6++",
			expectedFEN: "4k3/8/5N2/8/8/8/8/4R2K b - - 1 1",
		},
		{
			name:        "mate suffix accepted on checkmate move",
			fen:         "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 3",
			s:           "1. h5xf7mate",
			expectedFEN: "r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserCoordinate(Characteristics{}).Parse(g, tc.s)
			require.NoError(t, err)
			require.NotEmpty(t, gameSteps)
			lastGame := gameSteps[len(gameSteps)-1].StepGame
			assert.Equal(t, tc.expectedFEN, lastGame.ToFEN(), "FEN mismatch")
		})
	}
}

func TestNotationParserCoordinate_CheckSuffixValidation(t *testing.T) {
	t.Run("wrong check suffix on non-check move fails", func(t *testing.T) {
		g, err := core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.NoError(t, err)
		_, err = NewNotationParserCoordinate(Characteristics{}).Parse(g, "1. e2-e4+")
		require.Error(t, err, "e2-e4 is not check so the + suffix should fail to match")
	})
}

func TestCoordinateParser_RoundTrip(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := core.NewDefaultGame()

			algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, tc.algGame)
			require.NoError(t, err, "failed to parse algebraic game")

			coordPrinted, err := printer.CoordinatePrinter{}.PrintGame(algSteps, printer.GameCharacteristics{})
			require.NoError(t, err, "failed to print as coordinate")

			coordStr := strings.Join(coordPrinted, "\n")

			coordSteps, err := NewNotationParserCoordinate(Characteristics{}).Parse(g, coordStr)
			require.NoError(t, err, "failed to parse coordinate notation: %s", coordStr)

			require.Len(t, coordSteps, len(algSteps), "move count mismatch after round-trip")

			for i, orig := range algSteps {
				final := coordSteps[i]
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

func TestCoordinateParser_CrossNotationEquivalence(t *testing.T) {
	algGame := `1. e4 e6
		2. d4 d5
		3. Nc3 Bb4
		4. Bb5+ Bd7
		5. Bxd7+ Qxd7
		6. Ne2 dxe4
		7. 0-0`

	coordGame := `1. e2-e4 e7-e6
		2. d2-d4 d7-d5
		3. b1-c3 f8-b4
		4. f1-b5+ c8-d7
		5. b5xd7+ d8xd7
		6. g1-e2 d5xe4
		7. 0-0`

	g := core.NewDefaultGame()

	algSteps, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, algGame)
	require.NoError(t, err)

	coordSteps, err := NewNotationParserCoordinate(Characteristics{}).Parse(g, coordGame)
	require.NoError(t, err)

	require.Len(t, coordSteps, len(algSteps))

	for i, alg := range algSteps {
		coord := coordSteps[i]
		assert.Equal(t, alg.StepAction.FromPiece.XY, coord.StepAction.FromPiece.XY,
			"Move %d: from square mismatch", i+1)
		assert.Equal(t, alg.StepAction.ToXY, coord.StepAction.ToXY,
			"Move %d: to square mismatch", i+1)
		assert.Equal(t, alg.StepGame.ToFEN(), coord.StepGame.ToFEN(),
			"Move %d: FEN mismatch", i+1)
	}

	expectedFEN := "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7"
	assert.Equal(t, expectedFEN, coordSteps[len(coordSteps)-1].StepGame.ToFEN())
}
