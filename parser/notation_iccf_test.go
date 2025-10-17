package parser

import (
	"fmt"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotationParserICCF(t *testing.T) {
	testCases := []struct {
		fen                   string
		s                     string
		expectedErr           error
		expectedMatchedTokens []string
		expectedAlgebraic     []string
		expectedFEN           string
	}{
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. 5254 5755
2. 7163 2836
3. 6125 1716
4. 2514 7866
5. 5171 6857
6. 6151 2725
7. 1423 4746
8. 3233 5878`,
			expectedErr: nil,
			expectedMatchedTokens: []string{
				"5254",
				"5755",
				"7163",
				"2836",
				"6125",
				"1716",
				"2514",
				"7866",
				"5171",
				"6857",
				"6151",
				"2725",
				"1423",
				"4746",
				"3233",
				"5878",
			},
			expectedAlgebraic: []string{
				"1. e4 e5",
				"2. Nf3 Nc6",
				"3. Bb5 a6",
				"4. Ba4 Nf6",
				"5. 0-0 Be7",
				"6. Re1 b5",
				"7. Bb3 d6",
				"8. c3 0-0",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test notation parser ICCF %v", i), func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserICCF(Characteristics{}).Parse(g, tc.s)
			require.Equal(t, tc.expectedErr, err)
			if tc.expectedErr != nil {
				return
			}
			if tc.expectedMatchedTokens != nil {
				require.Len(t, tc.expectedMatchedTokens, len(gameSteps))
				for i, gameStep := range gameSteps {
					assert.Equal(t, tc.expectedMatchedTokens[i], gameStep.StepString)
				}
			}
			actualPrinted, err := printer.AlgebraicPrinter{}.PrintGame(gameSteps, printer.GameCharacteristics{})
			require.Nil(t, err)
			if tc.expectedAlgebraic != nil {
				assert.Equal(t, tc.expectedAlgebraic, actualPrinted)
			}
			if tc.expectedFEN != "" {
				assert.Equal(t, tc.expectedFEN, gameSteps[len(gameSteps)-1].StepGame.ToFEN())
			}
		})
	}
}

func TestNotationParserICCFRoundTrip(t *testing.T) {
	// Test round-trip: ICCF → Algebraic → ICCF
	originalICCF := `1. 5254 5755
2. 7163 2836
3. 6125 1716
4. 2514 7866
5. 5171 6857
6. 6151 2725
7. 1423 4746
8. 3233 5878`

	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	t.Run("ICCF to Algebraic to ICCF round-trip", func(t *testing.T) {
		// Step 1: Parse ICCF notation
		g, err := core.NewGameFromFEN(fen)
		require.NoError(t, err)

		gameSteps, err := NewNotationParserICCF(Characteristics{}).Parse(g, originalICCF)
		require.NoError(t, err)

		// Step 2: Convert to Algebraic notation
		algebraicPrinted, err := printer.AlgebraicPrinter{}.PrintGame(gameSteps, printer.GameCharacteristics{})
		require.NoError(t, err)

		// Step 3: Parse Algebraic notation back to game steps
		algebraicString := ""
		for _, move := range algebraicPrinted {
			algebraicString += move + "\n"
		}

		gameStepsFromAlgebraic, err := NewNotationParserAlgebraic(Characteristics{}).Parse(g, algebraicString)
		require.NoError(t, err)

		// Step 4: Convert back to ICCF notation
		iccfPrinted, err := printer.ICCFPrinter{}.PrintGame(gameStepsFromAlgebraic, printer.GameCharacteristics{})
		require.NoError(t, err)

		// Step 5: Parse the ICCF notation again
		iccfString := ""
		for _, move := range iccfPrinted {
			iccfString += move + "\n"
		}

		finalGameSteps, err := NewNotationParserICCF(Characteristics{}).Parse(g, iccfString)
		require.NoError(t, err)

		// Verify that the final ICCF moves match the original ICCF moves
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
		t.Logf("Algebraic: %v", algebraicPrinted)
		t.Logf("Final ICCF: %v", iccfPrinted)
	})
}
