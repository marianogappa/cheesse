package printer

import (
	"reflect"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/stretchr/testify/require"
)

func TestDescriptivePrinter_PrintAction(t *testing.T) {
	testCases := []struct {
		name                string
		gameStep            core.GameStep
		gameCharacteristics GameCharacteristics
		expectedResult      string
	}{
		{
			name: "Test case 1",
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
			expectedResult:      "P-K4",
		},
		// Add more test cases here...
	}

	printer := DescriptivePrinter{}

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

func TestDescriptivePrinter_PrintGame(t *testing.T) {
	testCases := []struct {
		name                string
		fen                 string
		gameDescriptive     string
		gameCharacteristics GameCharacteristics
		expectedResult      []string
	}{
		{
			name: "Test case 1",
			fen:  "8/7k/8/8/8/3P4/8/6K1 w - - 0 1",
			gameDescriptive: `1. K-B2 K-Kt3
				2. K-K3 K-B4
				3. K-Q4! K-K3
				4. K-B5 K-Q2
				5. K-Q5 K-K2
				6. K-B6 K-K3
				7. P-Q4 K-K2
				8. P-Q5 K-Q1
				9. K-Q6 Resigns`,
			gameCharacteristics: GameCharacteristics{},
			expectedResult: []string{
				"1. K-KB2 K-KN3",
				"2. K-K3 K-KB4",
				"3. K-Q4 K-K3", // TODO: enable storing annotations on game steps
				"4. K-QB5 K-Q2",
				"5. K-Q5 K-K2",
				"6. K-QB6 K-K3",
				"7. P-Q4 K-K2",
				"8. P-Q5 K-Q1",
				"9. K-Q6 resigns",
			},
		},
		// Add more test cases here...
	}

	printer := DescriptivePrinter{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := parser.NewNotationParserDescriptive(parser.Characteristics{}).Parse(g, tc.gameDescriptive)
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
