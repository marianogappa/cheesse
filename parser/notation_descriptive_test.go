package parser

import (
	"fmt"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotationParserDescriptive(t *testing.T) {
	testCases := []struct {
		fen                   string
		s                     string
		expectedErr           error
		expectedMatchedTokens []string
		expectedAlgebraic     []string
		expectedFEN           string
	}{
		{
			fen: "8/7k/8/8/8/3P4/8/6K1 w - - 0 1",
			s: `1. K-B2 K-Kt3
				2. K-K3 K-B4
				3. K-Q4! K-K3
				4. K-B5 K-Q2
				5. K-Q5 K-K2
				6. K-B6 K-K3
				7. P-Q4 K-K2
				8. P-Q5 K-Q1
				9. K-Q6 Resigns`,
			expectedErr: nil,
			expectedMatchedTokens: []string{
				"K-B2",
				"K-Kt3",
				"K-K3",
				"K-B4",
				"K-Q4!", // TODO: enable storing annotations on game steps
				"K-K3",
				"K-B5",
				"K-Q2",
				"K-Q5",
				"K-K2",
				"K-B6",
				"K-K3",
				"P-Q4",
				"K-K2",
				"P-Q5",
				"K-Q1",
				"K-Q6",
				"Resigns",
			},
			expectedAlgebraic: []string{
				"1. Kf2 Kg6",
				"2. Ke3 Kf5",
				"3. Kd4 Ke6", // TODO: enable storing annotations on game steps
				"4. Kc5 Kd7",
				"5. Kd5 Ke7",
				"6. Kc6 Ke6",
				"7. d4 Ke7",
				"8. d5 Kd8",
				"9. Kd6 resigns", // TODO: handle resigns vs 1-0
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test notation parser algebraic %v", i), func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, tc.s)
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
			fmt.Println(gameSteps)
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
