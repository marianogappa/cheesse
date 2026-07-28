package printer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlgebraicPrinter_PrintAction(t *testing.T) {
	testCases := []struct {
		name                string
		gameStep            core.GameStep
		gameCharacteristics GameCharacteristics
		expectedResult      string
	}{
		{
			name: "Test case 1",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame(),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 6}},
					ToXY:      core.XY{X: 4, Y: 4},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "e4",
		},
		// Add more test cases here...
	}

	printer := AlgebraicPrinter{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := printer.PrintAction(tc.gameStep, tc.gameCharacteristics)

			if err != nil {
				t.Errorf("PrintAction returned an error: %v", err)
			}

			if result != tc.expectedResult {
				t.Errorf("PrintAction returned incorrect result. Expected: %v, got: %v", tc.expectedResult, result)
			}
		})
	}
}

func TestAlgebraicPrinter_PrintGame(t *testing.T) {
	testCases := []struct {
		name                string
		gameAlgebraic       string
		gameCharacteristics GameCharacteristics
		expectedResult      []string
	}{
		{
			name: "Test case 1",
			gameAlgebraic: `1. e4 e6
				2. d4 d5
				3. Nc3 Bb4
				4. Bb5+ Bd7
				5. Bxd7+ Qxd7
				6. Ne2 dxe4
				7. 0-0`,
			gameCharacteristics: GameCharacteristics{},
			expectedResult: []string{
				"1. e4 e6",
				"2. d4 d5",
				"3. Nc3 Bb4",
				"4. Bb5+ Bd7",
				"5. Bxd7+ Qxd7",
				"6. Ne2 dxe4",
				"7. 0-0",
			},
		},
		// Add more test cases here...
	}

	printer := AlgebraicPrinter{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gameSteps, err := parser.NewNotationParserAlgebraic(parser.Characteristics{}).Parse(core.NewDefaultGame(), tc.gameAlgebraic)
			require.NoError(t, err)
			result, err := printer.PrintGame(gameSteps, tc.gameCharacteristics)

			if err != nil {
				t.Errorf("PrintGame returned an error: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("PrintGame returned incorrect result. Expected: %v, got: %v", tc.expectedResult, result)
			}
		})
	}
}

func buildGameStep(t *testing.T, fen string, fromSq, toSq string) core.GameStep {
	t.Helper()
	g, err := core.NewGameFromFEN(fen)
	require.NoError(t, err)

	from := algebraicToXY(t, fromSq)
	to := algebraicToXY(t, toSq)

	for _, action := range g.Actions {
		if action.FromPiece.XY == from && action.ToXY == to {
			newGame := g.DoAction(action)
			return core.GameStep{
				StepAction:      action,
				StepGame:        newGame,
				StepPreMoveGame: g,
			}
		}
	}
	t.Fatalf("no legal action from %s to %s in FEN %s", fromSq, toSq, fen)
	return core.GameStep{}
}

func buildGameStepWithPromotion(t *testing.T, fen string, fromSq, toSq string, promoPiece core.PieceType) core.GameStep {
	t.Helper()
	g, err := core.NewGameFromFEN(fen)
	require.NoError(t, err)

	from := algebraicToXY(t, fromSq)
	to := algebraicToXY(t, toSq)

	for _, action := range g.Actions {
		if action.FromPiece.XY == from && action.ToXY == to && action.PromotionPieceType == promoPiece {
			newGame := g.DoAction(action)
			return core.GameStep{
				StepAction:      action,
				StepGame:        newGame,
				StepPreMoveGame: g,
			}
		}
	}
	t.Fatalf("no legal action from %s to %s promoting to %v in FEN %s", fromSq, toSq, promoPiece, fen)
	return core.GameStep{}
}

func algebraicToXY(t *testing.T, sq string) core.XY {
	t.Helper()
	require.Len(t, sq, 2)
	return core.XY{X: int(sq[0] - 'a'), Y: int('8' - sq[1])}
}

func TestAlgebraicPrinter_Disambiguation(t *testing.T) {
	p := AlgebraicPrinter{}
	gc := GameCharacteristics{}

	t.Run("two knights different files", func(t *testing.T) {
		// Knights at d2 and f2 both can reach e4
		gs := buildGameStep(t, "4k3/8/8/8/8/8/3N1N2/4K3 w - - 0 1", "d2", "e4")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "Nde4", result)
	})

	t.Run("disambiguation not needed when only one piece can reach the square", func(t *testing.T) {
		gs := buildGameStep(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "g1", "f3")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "Nf3", result)
	})

	t.Run("two rooks same rank different files", func(t *testing.T) {
		// King at a2, rooks at a1 and h1 — both can reach d1
		gs := buildGameStep(t, "3k4/8/8/8/8/8/K7/R6R w - - 0 1", "h1", "d1")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "Rhd1+", result)
	})

	t.Run("two rooks same file need rank disambiguation", func(t *testing.T) {
		// Rooks at a1 and a8, king at h1
		gs := buildGameStep(t, "R2k4/8/8/8/8/8/8/R6K w - - 0 1", "a1", "a4")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "R1a4+", result)
	})

	t.Run("no disambiguation for pawns", func(t *testing.T) {
		gs := buildGameStep(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2", "e4")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "e4", result)
	})

	t.Run("no disambiguation for kings", func(t *testing.T) {
		gs := buildGameStep(t, "4k3/8/8/8/8/8/8/4K3 w - - 0 1", "e1", "e2")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "Ke2", result)
	})
}

func TestAlgebraicPrinter_EnPassant(t *testing.T) {
	p := AlgebraicPrinter{}
	gc := GameCharacteristics{}

	t.Run("en passant prints destination square not captured pawn square", func(t *testing.T) {
		gs := buildGameStep(t, "rnbqkbnr/pppp1ppp/8/4pP2/8/8/PPPPP1PP/RNBQKBNR w KQkq e6 0 1", "f5", "e6")
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "fxe6 e.p.", result)
	})
}

func TestAlgebraicPrinter_CheckmateSymbol(t *testing.T) {
	p := AlgebraicPrinter{}

	t.Run("default checkmate symbol is #", func(t *testing.T) {
		// Scholar's mate final move: Qxf7#
		gs := buildGameStep(t, "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 3", "h5", "f7")
		result, err := p.PrintAction(gs, GameCharacteristics{})
		require.NoError(t, err)
		assert.Equal(t, "Qxf7#", result)
	})
}

func TestAlgebraicPrinter_Promotion(t *testing.T) {
	p := AlgebraicPrinter{}
	gc := GameCharacteristics{}

	t.Run("promotion to queen", func(t *testing.T) {
		gs := buildGameStepWithPromotion(t, "8/5P1k/8/8/8/8/8/K7 w - - 0 1", "f7", "f8", core.PieceQueen)
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "f8=Q", result)
	})

	t.Run("promotion to knight with check", func(t *testing.T) {
		gs := buildGameStepWithPromotion(t, "8/5P1k/8/8/8/8/8/K7 w - - 0 1", "f7", "f8", core.PieceKnight)
		result, err := p.PrintAction(gs, gc)
		require.NoError(t, err)
		assert.Equal(t, "f8=N+", result)
	})
}

func TestAlgebraicPrinter_RoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		game string
	}{
		{
			name: "French Defense opening",
			game: `1. e4 e6
				2. d4 d5
				3. Nc3 Bb4
				4. Bb5+ Bd7
				5. Bxd7+ Qxd7
				6. Ne2 dxe4
				7. 0-0`,
		},
		{
			name: "Scholar's mate",
			game: `1. e4 e5
				2. Bc4 Nc6
				3. Qh5 Nf6
				4. Qxf7#`,
		},
		{
			name: "Italian Game with castling",
			game: `1. e4 e5
				2. Nf3 Nc6
				3. Bc4 Bc5
				4. 0-0 Nf6
				5. d3 0-0`,
		},
		{
			name: "Game with en passant",
			game: `1. e4 Nf6
				2. e5 d5
				3. exd6`,
		},
		{
			name: "Game with promotion",
			game: `1. d4 c5
				2. d5 e6
				3. dxe6 d5
				4. exf7+ Kd7
				5. fxg8=N`,
		},
		{
			name: "Ruy Lopez with knight disambiguation",
			game: `1. e4 e5
				2. Nf3 Nc6
				3. Bb5 a6
				4. Ba4 Nf6
				5. O-O Be7
				6. Re1 b5
				7. Bb3 O-O
				8. c3 d5`,
		},
		{
			name: "Queen's Gambit Accepted",
			game: `1. d4 d5
				2. c4 dxc4
				3. Nf3 Nf6
				4. e3 e6
				5. Bxc4 c5
				6. O-O a6`,
		},
	}

	p := AlgebraicPrinter{}
	roundTripGC := SANCharacteristics()
	roundTripGC.usesEnPassantSymbol = pstr("")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := core.NewDefaultGame()
			gameSteps, err := parser.NewNotationParserAlgebraic(parser.Characteristics{}).Parse(g, tc.game)
			require.NoError(t, err, "failed to parse original game")

			printed, err := p.PrintGame(gameSteps, roundTripGC)
			require.NoError(t, err, "failed to print game")

			printedStr := strings.Join(printed, "\n")
			reParsed, err := parser.NewNotationParserAlgebraic(parser.Characteristics{}).Parse(g, printedStr)
			require.NoError(t, err, "failed to re-parse printed game: %s", printedStr)

			require.Len(t, reParsed, len(gameSteps), "move count mismatch after round-trip")

			for i, orig := range gameSteps {
				final := reParsed[i]
				assert.Equal(t, orig.StepAction.FromPiece.XY, final.StepAction.FromPiece.XY,
					"Move %d: from square mismatch", i+1)
				assert.Equal(t, orig.StepAction.ToXY, final.StepAction.ToXY,
					"Move %d: to square mismatch", i+1)
				assert.Equal(t, orig.StepAction.FromPiece.PieceType, final.StepAction.FromPiece.PieceType,
					"Move %d: piece type mismatch", i+1)
				assert.Equal(t, orig.StepGame.ToFEN(), final.StepGame.ToFEN(),
					"Move %d: resulting FEN mismatch", i+1)
			}
		})
	}
}
