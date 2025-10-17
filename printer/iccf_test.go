package printer

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestICCFPrinter_PrintAction(t *testing.T) {
	testCases := []struct {
		name                string
		gameStep            core.GameStep
		gameCharacteristics GameCharacteristics
		expectedResult      string
	}{
		{
			name: "Pawn move",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 6}},
					ToXY:      core.XY{X: 4, Y: 4},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 6}},
					ToXY:      core.XY{X: 4, Y: 4},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "5254",
		},
		{
			name: "Knight move",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKnight, Owner: core.ColorWhite, XY: core.XY{X: 6, Y: 7}},
					ToXY:      core.XY{X: 5, Y: 5},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKnight, Owner: core.ColorWhite, XY: core.XY{X: 6, Y: 7}},
					ToXY:      core.XY{X: 5, Y: 5},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "7163",
		},
		{
			name: "Bishop move",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceBishop, Owner: core.ColorWhite, XY: core.XY{X: 5, Y: 7}},
					ToXY:      core.XY{X: 0, Y: 2},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceBishop, Owner: core.ColorWhite, XY: core.XY{X: 5, Y: 7}},
					ToXY:      core.XY{X: 0, Y: 2},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "6116",
		},
		{
			name: "Rook move",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceRook, Owner: core.ColorWhite, XY: core.XY{X: 0, Y: 7}},
					ToXY:      core.XY{X: 0, Y: 5},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceRook, Owner: core.ColorWhite, XY: core.XY{X: 0, Y: 7}},
					ToXY:      core.XY{X: 0, Y: 5},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "1113",
		},
		{
			name: "Queen move",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceQueen, Owner: core.ColorWhite, XY: core.XY{X: 3, Y: 7}},
					ToXY:      core.XY{X: 7, Y: 3},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceQueen, Owner: core.ColorWhite, XY: core.XY{X: 3, Y: 7}},
					ToXY:      core.XY{X: 7, Y: 3},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "4185",
		},
		{
			name: "King move",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKing, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 7}},
					ToXY:      core.XY{X: 6, Y: 7},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKing, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 7}},
					ToXY:      core.XY{X: 6, Y: 7},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "5171",
		},
		{
			name: "Pawn promotion to Queen",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceQueen,
				}),
				StepAction: core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceQueen,
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "57581",
		},
		{
			name: "Pawn promotion to Rook",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceRook,
				}),
				StepAction: core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceRook,
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "57582",
		},
		{
			name: "Pawn promotion to Bishop",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceBishop,
				}),
				StepAction: core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceBishop,
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "57583",
		},
		{
			name: "Pawn promotion to Knight",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceKnight,
				}),
				StepAction: core.Action{
					FromPiece:          core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 1}},
					ToXY:               core.XY{X: 4, Y: 0},
					IsPromotion:        true,
					PromotionPieceType: core.PieceKnight,
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "57584",
		},
	}

	printer := ICCFPrinter{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := printer.PrintAction(tc.gameStep, tc.gameCharacteristics)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestICCFPrinter_PrintGame(t *testing.T) {
	testCases := []struct {
		name                string
		fen                 string
		gameICCF            string
		gameCharacteristics GameCharacteristics
		expectedResult      []string
	}{
		{
			name: "Basic opening sequence",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			gameICCF: `1. 5254 5755
2. 7163 2836
3. 6125 1716
4. 2514 7866
5. 5171 6857
6. 6151 2725
7. 1423 4746
8. 3233 5878`,
			gameCharacteristics: GameCharacteristics{},
			expectedResult: []string{
				"1. 5254 5755",
				"2. 7163 2836",
				"3. 6125 1716",
				"4. 2514 7866",
				"5. 5171 6857",
				"6. 6151 2725",
				"7. 1423 4746",
				"8. 3233 5878",
			},
		},
		{
			name: "Simple two-move sequence",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			gameICCF: `1. 5254 5755
2. 7163 2836`,
			gameCharacteristics: GameCharacteristics{},
			expectedResult: []string{
				"1. 5254 5755",
				"2. 7163 2836",
			},
		},
	}

	printer := ICCFPrinter{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := parser.NewNotationParserICCF(parser.Characteristics{}).Parse(g, tc.gameICCF)
			require.NoError(t, err)
			result, err := printer.PrintGame(gameSteps, tc.gameCharacteristics)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestICCFPrinter_RoundTrip(t *testing.T) {
	// Test round-trip: ICCF → Game Steps → ICCF
	originalICCF := `1. 5254 5755
2. 7163 2836
3. 6125 1716
4. 2514 7866
5. 5171 6857
6. 6151 2725
7. 1423 4746
8. 3233 5878`

	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	t.Run("ICCF round-trip consistency", func(t *testing.T) {
		// Step 1: Parse ICCF notation
		g, err := core.NewGameFromFEN(fen)
		require.NoError(t, err)

		gameSteps, err := parser.NewNotationParserICCF(parser.Characteristics{}).Parse(g, originalICCF)
		require.NoError(t, err)

		// Step 2: Print back to ICCF notation
		iccfPrinted, err := ICCFPrinter{}.PrintGame(gameSteps, GameCharacteristics{})
		require.NoError(t, err)

		// Step 3: Parse the printed ICCF notation again
		iccfString := ""
		for _, move := range iccfPrinted {
			iccfString += move + "\n"
		}

		finalGameSteps, err := parser.NewNotationParserICCF(parser.Characteristics{}).Parse(g, iccfString)
		require.NoError(t, err)

		// Verify that the final game steps match the original game steps
		require.Len(t, finalGameSteps, len(gameSteps))

		for i, originalStep := range gameSteps {
			finalStep := finalGameSteps[i]
			assert.Equal(t, originalStep.StepAction.FromPiece.XY, finalStep.StepAction.FromPiece.XY,
				"Move %d: From position should match", i+1)
			assert.Equal(t, originalStep.StepAction.ToXY, finalStep.StepAction.ToXY,
				"Move %d: To position should match", i+1)
			assert.Equal(t, originalStep.StepAction.FromPiece.PieceType, finalStep.StepAction.FromPiece.PieceType,
				"Move %d: Piece type should match", i+1)
			if originalStep.StepAction.IsPromotion {
				assert.Equal(t, originalStep.StepAction.PromotionPieceType, finalStep.StepAction.PromotionPieceType,
					"Move %d: Promotion piece type should match", i+1)
			}
		}

		// Also verify that the string representations match exactly
		assert.Equal(t, iccfPrinted, []string{
			"1. 5254 5755",
			"2. 7163 2836",
			"3. 6125 1716",
			"4. 2514 7866",
			"5. 5171 6857",
			"6. 6151 2725",
			"7. 1423 4746",
			"8. 3233 5878",
		}, "Final ICCF string representation should match original")

		t.Logf("Original ICCF: %s", originalICCF)
		t.Logf("Final ICCF: %v", iccfPrinted)
	})
}

func TestICCFPrinter_EdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		gameStep       core.GameStep
		expectedResult string
	}{
		{
			name: "Castling kingside",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKing, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 7}},
					ToXY:      core.XY{X: 6, Y: 7},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKing, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 7}},
					ToXY:      core.XY{X: 6, Y: 7},
				},
			},
			expectedResult: "5171",
		},
		{
			name: "Castling queenside",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKing, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 7}},
					ToXY:      core.XY{X: 2, Y: 7},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PieceKing, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 7}},
					ToXY:      core.XY{X: 2, Y: 7},
				},
			},
			expectedResult: "5131",
		},
		{
			name: "En passant capture",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame().DoAction(core.Action{
					FromPiece: core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 3}},
					ToXY:      core.XY{X: 5, Y: 2},
				}),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 3}},
					ToXY:      core.XY{X: 5, Y: 2},
				},
			},
			expectedResult: "5566",
		},
	}

	printer := ICCFPrinter{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := printer.PrintAction(tc.gameStep, GameCharacteristics{})

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
