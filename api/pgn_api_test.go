package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNotation_PGNMetadataExposed(t *testing.T) {
	pgn := `[Event "API Test"]
[White "Alice"]
[Black "Bob"]

1. e4 e5 2. Nf3 Nc6 1-0`
	_, result, err := New().ParseNotation(InputGame{}, pgn)
	require.NoError(t, err)
	require.True(t, result.ParseWasSuccessful, "parse failed: %v", result.Error)
	assert.Equal(t, "PGN", result.NotationName)
	assert.Equal(t, "API Test", result.Metadata["Event"])
	assert.Equal(t, "Alice", result.Metadata["White"])
	assert.Equal(t, "Bob", result.Metadata["Black"])
}

func TestParseNotation_PGNWithVariationsAndNAGs(t *testing.T) {
	pgn := `[Event "Annotated"]

1. e4 $1 {King's pawn} e5 (1... c5 {Sicilian} 2. Nf3) 2. Nf3 Nc6 *`
	_, result, err := New().ParseNotation(InputGame{}, pgn)
	require.NoError(t, err)
	require.True(t, result.ParseWasSuccessful, "parse failed: %v", result.Error)
	assert.Equal(t, "PGN", result.NotationName)
	assert.Equal(t, 5, result.ValidActionCount) // 4 moves + result marker
}

func TestConvertNotation_ToPGN(t *testing.T) {
	_, result, err := New().ConvertNotation(InputGame{}, "1. e4 e5 2. Nf3 Nc6", "PGN")
	require.NoError(t, err)
	require.True(t, result.ParseWasSuccessful, "parse failed: %v", result.Error)
	actionStrings := make([]string, 0, len(result.Steps))
	for _, s := range result.Steps {
		actionStrings = append(actionStrings, s.ActionString)
	}
	assert.Equal(t, []string{"e4", "e5", "Nf3", "Nc6"}, actionStrings)
}

func TestParseNotation_PGNPartialParse(t *testing.T) {
	pgn := `1. e4 e5 2. Qxf7 Nc6`
	_, result, err := New().ParseNotation(InputGame{}, pgn)
	require.NoError(t, err)
	assert.False(t, result.ParseWasSuccessful)
	assert.Equal(t, 2, result.ValidActionCount, "the two valid moves before the failure should be counted")
	assert.True(t, strings.Contains(result.Error, "Qxf7"), "error should name the failing move: %v", result.Error)
}
