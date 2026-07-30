package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// actionStrings extracts the per-step action strings of a conversion result.
func actionStrings(result OutputParseResult) []string {
	out := make([]string, len(result.Steps))
	for i, step := range result.Steps {
		out[i] = step.ActionString
	}
	return out
}

// The six-notation-suite French Defense expressed as the expected per-move action
// strings in each target notation. Cross-converting between any source notation in
// sixNotationGames (parse_notation_test.go) and any of these targets must yield
// exactly these strings.
var expectedConversionsByTarget = map[string][]string{
	"Algebraic": {
		"e4", "e6", "d4", "d5", "Nc3", "Bb4", "Bb5+", "Bd7", "Bxd7+", "Qxd7", "Ne2", "dxe4", "O-O",
	},
	"Figurine": {
		"e4", "e6", "d4", "d5", "♘c3", "♝b4", "♗b5+", "♝d7", "♗xd7+", "♛xd7", "♘e2", "dxe4", "O-O",
	},
	"Coordinate": {
		"e2-e4", "e7-e6", "d2-d4", "d7-d5", "b1-c3", "f8-b4", "f1-b5+", "c8-d7", "b5xd7+", "d8xd7", "g1-e2", "d5xe4", "O-O",
	},
	"ICCF": {
		"5254", "5756", "4244", "4745", "2133", "6824", "6125", "3847", "2547", "4847", "7152", "4554", "5171",
	},
	"Smith": {
		"e2e4", "e7e6", "d2d4", "d7d5", "b1c3", "f8b4", "f1b5", "c8d7", "b5d7b", "d8d7b", "g1e2", "d5e4p", "e1g1c",
	},
}

func TestConvertNotation_CrossNotationMatrix(t *testing.T) {
	// Every source notation from the six-notation suite converted to every target
	// notation must produce the exact expected action strings.
	for _, source := range sixNotationGames {
		for target, expected := range expectedConversionsByTarget {
			t.Run(source.notationName+"->"+target, func(t *testing.T) {
				_, result, err := New().ConvertNotation(InputGame{}, source.game, target)
				require.NoError(t, err)
				require.True(t, result.ParseWasSuccessful, "parse failed: %v", result.Error)
				assert.Equal(t, source.notationName, result.NotationName)

				actual := actionStrings(result)
				// PGN sources have a trailing result-marker step with no action;
				// compare only the action-bearing steps.
				if len(actual) == len(expected)+1 && actual[len(actual)-1] == "" {
					actual = actual[:len(actual)-1]
				}
				assert.Equal(t, expected, actual)
			})
		}
	}
}

func TestConvertNotation_ToDescriptive(t *testing.T) {
	// Descriptive output isn't asserted in the cross matrix because its rendering
	// varies more; assert the exact strings once from an algebraic source.
	algebraic := `1. e4 e6
2. d4 d5
3. Nc3 Bb4
4. Bb5+ Bd7
5. Bxd7+ Qxd7
6. Ne2 dxe4
7. 0-0`
	_, result, err := New().ConvertNotation(InputGame{}, algebraic, "Descriptive")
	require.NoError(t, err)
	require.True(t, result.ParseWasSuccessful, "parse failed: %v", result.Error)
	got := actionStrings(result)
	require.Len(t, got, 13)
	assert.Equal(t, "P-K4", got[0])
	assert.Equal(t, "P-K3", got[1])
	assert.Equal(t, "P-Q4", got[2])
	assert.Equal(t, "N-QB3", got[4])
	assert.Equal(t, "0-0", got[12])
}

func TestConvertNotation_PartialInput(t *testing.T) {
	// Invalid move mid-game: valid prefix still converts.
	game := "1. e4 e5\n2. Bc4 Nc6\n3. Qh7 Nf6"
	_, result, err := New().ConvertNotation(InputGame{}, game, "ICCF")
	require.NoError(t, err)
	assert.False(t, result.ParseWasSuccessful)
	assert.Equal(t, "Algebraic Notation", result.NotationName)
	assert.Equal(t, 4, result.ValidActionCount)
	assert.Equal(t, []string{"5254", "5755", "6134", "2836"}, actionStrings(result))
	assert.NotEmpty(t, result.Error)
}

func TestConvertNotation_CustomInitialPosition(t *testing.T) {
	// Chernev Game 1: descriptive source from custom FEN converted to algebraic.
	inputGame := InputGame{FENString: "8/8/8/8/8/1k5P/8/2K5 w - - 0 1"}
	game := "1. P-R4 K-B5\n2. P-R5 K-Q4\n3. P-R6 K-K3\n4. P-R7 K-B2\n5. P-R8(Q) Resigns"
	_, result, err := New().ConvertNotation(inputGame, game, "Algebraic")
	require.NoError(t, err)
	require.True(t, result.ParseWasSuccessful, "parse failed: %v", result.Error)
	got := actionStrings(result)
	require.Len(t, got, 10)
	assert.Equal(t, []string{"h4", "Kc4", "h5", "Kd5", "h6", "Ke6", "h7", "Kf7", "h8=Q", "resigns"}, got)
}

func TestConvertNotation_UnknownTargetNotation(t *testing.T) {
	_, _, err := New().ConvertNotation(InputGame{}, "1. e4", "Klingon")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown target notation")
}

func TestConvertNotation_TargetNotationCaseInsensitive(t *testing.T) {
	for _, target := range []string{"algebraic", "ALGEBRAIC", "Algebraic", "iccf", "ICCF", "smith", "figurine", "coordinate", "descriptive"} {
		_, result, err := New().ConvertNotation(InputGame{}, "1. e4 e5", target)
		require.NoError(t, err, "target %q should be accepted", target)
		assert.True(t, result.ParseWasSuccessful)
	}
}

func TestConvertNotation_InvalidInputGame(t *testing.T) {
	_, _, err := New().ConvertNotation(InputGame{FENString: "not a fen"}, "1. e4", "Algebraic")
	require.Error(t, err)
}

func TestConvertNotation_GarbageInput(t *testing.T) {
	_, result, err := New().ConvertNotation(InputGame{}, "hello world not chess", "Algebraic")
	require.NoError(t, err)
	assert.False(t, result.ParseWasSuccessful)
	assert.Equal(t, 0, result.ValidActionCount)
	assert.Empty(t, result.Steps)
}

func TestConvertNotation_RoundTripThroughAllNotations(t *testing.T) {
	// Convert algebraic → target, then feed the converted text back (joined as one
	// half-move per line pair) and convert to algebraic; the result must match the
	// canonical algebraic strings. This exercises print→parse consistency per notation.
	canonical := []string{"e4", "e6", "d4", "d5", "Nc3", "Bb4", "Bb5+", "Bd7", "Bxd7+", "Qxd7", "Ne2", "dxe4", "O-O"}
	algebraic := `1. e4 e6
2. d4 d5
3. Nc3 Bb4
4. Bb5+ Bd7
5. Bxd7+ Qxd7
6. Ne2 dxe4
7. 0-0`

	for _, target := range []string{"Algebraic", "ICCF", "Smith", "Coordinate", "Figurine"} {
		t.Run(target, func(t *testing.T) {
			_, converted, err := New().ConvertNotation(InputGame{}, algebraic, target)
			require.NoError(t, err)
			require.True(t, converted.ParseWasSuccessful)

			// Reassemble the converted action strings into a game text with move
			// numbers (numberless ICCF is ambiguous: a bare "5254" reads as a move
			// number).
			var sb strings.Builder
			for i := 0; i < len(converted.Steps); i += 2 {
				fmt.Fprintf(&sb, "%d. %s", i/2+1, converted.Steps[i].ActionString)
				if i+1 < len(converted.Steps) {
					sb.WriteString(" ")
					sb.WriteString(converted.Steps[i+1].ActionString)
				}
				sb.WriteString("\n")
			}

			_, back, err := New().ConvertNotation(InputGame{}, sb.String(), "Algebraic")
			require.NoError(t, err)
			require.True(t, back.ParseWasSuccessful, "re-parse of %s output failed: %v\ninput:\n%s", target, back.Error, sb.String())
			assert.Equal(t, canonical, actionStrings(back))
		})
	}
}
