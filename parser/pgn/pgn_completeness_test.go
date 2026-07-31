package pgn

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parsePGN(t *testing.T, pgn string) (*parser.ParsedGame, error) {
	t.Helper()
	initialGame, err := core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.NoError(t, err)
	return parser.NewGenericNotationParser(NewVariantPGN()).Parse(initialGame, pgn)
}

func TestPGNRAVVariations(t *testing.T) {
	t.Run("skips a simple variation", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 e5 (1... c5 2. Nf3) 2. Nf3 Nc6")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 4)
		assert.Equal(t, "e4", parsed.GameSteps[0].StepString)
		assert.Equal(t, "e5", parsed.GameSteps[1].StepString)
		assert.Equal(t, "Nf3", parsed.GameSteps[2].StepString)
		assert.Equal(t, "Nc6", parsed.GameSteps[3].StepString)
	})

	t.Run("skips nested variations", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 e5 (1... c5 (1... e6 2. d4) 2. Nf3) 2. Nf3")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 3)
		assert.Equal(t, "Nf3", parsed.GameSteps[2].StepString)
	})

	t.Run("skips variation containing comments with parens", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 e5 (1... c5 {sicilian (sharp)} 2. Nf3) 2. Nf3")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 3)
	})

	t.Run("variation before first move of a game continuation", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 (1. d4 d5) 1... e5 2. Nf3")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 3)
		assert.Equal(t, "e5", parsed.GameSteps[1].StepString)
	})
}

func TestPGNBareNAGs(t *testing.T) {
	t.Run("bare NAG after move", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 $1 e5 $14 2. Nf3")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 3)
		assert.Equal(t, "e4", parsed.GameSteps[0].StepString)
		assert.Equal(t, "e5", parsed.GameSteps[1].StepString)
	})

	t.Run("parenthesized NAG still works", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 ($1) e5")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 2)
	})
}

func TestPGNCommentsAttached(t *testing.T) {
	t.Run("curly comment attaches to preceding move", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 {best by test} e5 2. Nf3 {develops}")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 3)
		assert.Equal(t, "best by test", parsed.GameSteps[0].StepComment)
		assert.Equal(t, "", parsed.GameSteps[1].StepComment)
		assert.Equal(t, "develops", parsed.GameSteps[2].StepComment)
	})

	t.Run("semicolon comment attaches to preceding move", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 e5 ; the open game\n2. Nf3")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 3)
		assert.Equal(t, "the open game", parsed.GameSteps[1].StepComment)
	})

	t.Run("multiple comments concatenate", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 {one} {two} e5")
		require.NoError(t, err)
		require.Len(t, parsed.GameSteps, 2)
		assert.Equal(t, "one two", parsed.GameSteps[0].StepComment)
	})
}

func TestPGNPartialParse(t *testing.T) {
	t.Run("invalid move mid-game returns valid prefix and error", func(t *testing.T) {
		parsed, err := parsePGN(t, "1. e4 e5 2. Qxf7 Nc6")
		require.Error(t, err)
		require.NotNil(t, parsed)
		assert.Len(t, parsed.GameSteps, 2, "the two valid moves before the failure should be returned")
	})
}

func TestPGNRealWorldAnnotatedGame(t *testing.T) {
	pgn := `[Event "Annotated"]
[White "A"]
[Black "B"]

1. e4 $1 {King's pawn} e5 (1... c5 {Sicilian} 2. Nf3 (2. c3 $6) 2... d6) 2. Nf3 {attacks e5} Nc6 3. Bb5 $14 a6 1/2-1/2`
	parsed, err := parsePGN(t, pgn)
	require.NoError(t, err)
	require.Len(t, parsed.GameSteps, 7) // 6 moves + result marker
	assert.Equal(t, "King's pawn", parsed.GameSteps[0].StepComment)
	assert.Equal(t, "attacks e5", parsed.GameSteps[2].StepComment)
	assert.Equal(t, "1/2-1/2", parsed.GameSteps[6].StepString)
	assert.Equal(t, "Annotated", parsed.Metadata["Event"])
}
