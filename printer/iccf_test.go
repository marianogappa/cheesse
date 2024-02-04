package printer

import (
	"reflect"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
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
			name: "Test case 1",
			gameStep: core.GameStep{
				StepGame: core.NewDefaultGame(),
				StepAction: core.Action{
					FromPiece: core.Piece{PieceType: core.PiecePawn, Owner: core.ColorWhite, XY: core.XY{X: 4, Y: 6}},
					ToXY:      core.XY{X: 4, Y: 4},
				},
			},
			gameCharacteristics: GameCharacteristics{},
			expectedResult:      "5254",
		},
		// Add more test cases here...
	}

	printer := ICCFPrinter{}

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

func TestICCFPrinter_PrintGame(t *testing.T) {
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
				"1. 5254 5756",
				"2. 4244 4745",
				"3. 2133 6824",
				"4. 6125 3847",
				"5. 2547 4847",
				"6. 7152 4554",
				"7. 5171",
			},
		},
		// Add more test cases here...
	}

	printer := ICCFPrinter{}

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
