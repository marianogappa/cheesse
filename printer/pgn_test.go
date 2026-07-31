package printer

import (
	"strings"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/marianogappa/cheesse/parser/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parsePGNForPrinting(t *testing.T, s string) *parser.ParsedGame {
	t.Helper()
	initialGame, err := core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.NoError(t, err)
	parsed, err := parser.NewGenericNotationParser(pgn.NewVariantPGN()).Parse(initialGame, s)
	require.NoError(t, err)
	return parsed
}

func TestPGNPrinterTagSection(t *testing.T) {
	parsed := parsePGNForPrinting(t, `[Event "Casual"]
[White "Alice"]
[Black "Bob"]
[ECO "C50"]

1. e4 e5 1/2-1/2`)

	p := PGNPrinter{Metadata: parsed.Metadata}
	lines, err := p.PrintGame(parsed.GameSteps, SANCharacteristics())
	require.NoError(t, err)
	doc := strings.Join(lines, "\n")

	assert.Contains(t, doc, `[Event "Casual"]`)
	assert.Contains(t, doc, `[White "Alice"]`)
	assert.Contains(t, doc, `[Black "Bob"]`)
	assert.Contains(t, doc, `[ECO "C50"]`, "non-STR metadata should pass through")
	assert.Contains(t, doc, `[Site "?"]`, "missing STR tags get placeholders")
	assert.Contains(t, doc, `[Date "????.??.??"]`)
	assert.Contains(t, doc, "1. e4 e5 1/2-1/2")
}

func TestPGNPrinterMovetext(t *testing.T) {
	t.Run("move numbers and result", func(t *testing.T) {
		parsed := parsePGNForPrinting(t, "1. e4 e5 2. Nf3 Nc6 1-0")
		lines, err := PGNPrinter{}.PrintGame(parsed.GameSteps, SANCharacteristics())
		require.NoError(t, err)
		doc := strings.Join(lines, "\n")
		assert.Contains(t, doc, "1. e4 e5 2. Nf3 Nc6 1-0")
	})

	t.Run("comments are preserved", func(t *testing.T) {
		parsed := parsePGNForPrinting(t, "1. e4 {King's pawn} e5 2. Nf3")
		lines, err := PGNPrinter{}.PrintGame(parsed.GameSteps, SANCharacteristics())
		require.NoError(t, err)
		doc := strings.Join(lines, "\n")
		assert.Contains(t, doc, "1. e4 {King's pawn} 1... e5 2. Nf3")
	})

	t.Run("no explicit result yields *", func(t *testing.T) {
		parsed := parsePGNForPrinting(t, "1. e4 e5")
		lines, err := PGNPrinter{}.PrintGame(parsed.GameSteps, SANCharacteristics())
		require.NoError(t, err)
		doc := strings.Join(lines, "\n")
		assert.Contains(t, doc, "1. e4 e5 *")
		assert.Contains(t, doc, `[Result "*"]`)
	})

	t.Run("long games wrap at 80 columns", func(t *testing.T) {
		parsed := parsePGNForPrinting(t, "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7")
		lines, err := PGNPrinter{}.PrintGame(parsed.GameSteps, SANCharacteristics())
		require.NoError(t, err)
		for _, line := range lines {
			assert.LessOrEqual(t, len(line), 80, "line exceeds 80 columns: %q", line)
		}
	})
}

func TestPGNPrinterRoundTrip(t *testing.T) {
	original := `[Event "RT"]
[White "A"]
[Black "B"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 1-0`
	parsed := parsePGNForPrinting(t, original)
	lines, err := PGNPrinter{Metadata: parsed.Metadata}.PrintGame(parsed.GameSteps, SANCharacteristics())
	require.NoError(t, err)

	reparsed := parsePGNForPrinting(t, strings.Join(lines, "\n"))
	require.Equal(t, len(parsed.GameSteps), len(reparsed.GameSteps))
	for i := range parsed.GameSteps {
		assert.Equal(t, parsed.GameSteps[i].StepAction, reparsed.GameSteps[i].StepAction, "step %d differs", i)
	}
	assert.Equal(t, "RT", reparsed.Metadata["Event"])
}
