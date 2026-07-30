package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The same French Defense game in all six supported notations (ostinato's
// six-notation suite), all ending at the same final position.
const sixNotationsFinalFEN = "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7"

var sixNotationGames = []struct {
	notationName string
	game         string
}{
	{"Algebraic Notation", `1. e4 e6
2. d4 d5
3. Nc3 Bb4
4. Bb5+ Bd7
5. Bxd7+ Qxd7
6. Ne2 dxe4
7. 0-0`},
	{"Algebraic Notation", `1. e4 e6
2. d4 d5
3. ♘c3 ♝b4
4. ♗b5+ ♝d7
5. ♗xd7+ ♛xd7
6. ♘ge2 dxe4
7. 0-0`}, // Figurine is a variant of Algebraic
	{"ICCF Notation", `1. 5254 5756
2. 4244 4745
3. 2133 6824
4. 6125 3847
5. 2547 4847
6. 7152 4554
7. 5171`},
	{"Smith Notation", `1. e2e4  e7e6
2. d2d4  d7d5
3. b1c3  f8b4
4. f1b5  c8d7
5. b5d7b d8d7b
6. g1e2  d5e4p
7. e1g1c`},
	{"Coordinate Notation", `1. e2-e4 e7-e6
2. d2-d4 d7-d5
3. b1-c3 f8-b4
4. f1-b5+ c8-d7
5. b5xd7+ d8xd7
6. g1-e2 d5xe4
7. 0-0`},
	{"Descriptive Notation", `1. P-K4 P-K3
2. P-Q4 P-Q4
3. N-QB3 B-N5
4. B-N5ch B-Q2
5. BxBch QxB
6. KN-K2 PxP
7. 0-0`},
	{"PGN", `[Event "Ostinato Testing"]
[Site "Buenos Aires, Argentina"]
[Date "2015.??.??"]
[Round "1"]
[Result "1/2-1/2"]
[White "Fake Player 1"]
[Black "Fake Player 2"]

1. e4 e6 2. d4 d5 3. Nc3 Bb4 4. Bb5+ Bd7 5. Bxd7+ Qxd7 6. Nge2
dxe4 7. 0-0`},
}

func TestParseNotation_AutoDetectsSixNotations(t *testing.T) {
	for _, tc := range sixNotationGames {
		t.Run(tc.notationName+" / "+tc.game[:20], func(t *testing.T) {
			_, result, err := New().ParseNotation(InputGame{}, tc.game)
			require.NoError(t, err)
			assert.True(t, result.ParseWasSuccessful, "parse should succeed; error: %v", result.Error)
			assert.Equal(t, tc.notationName, result.NotationName)
			require.NotEmpty(t, result.Steps)
			// PGN games may end with a result-marker step whose game is unchanged;
			// find the last step with the expected final position.
			lastStep := result.Steps[len(result.Steps)-1]
			assert.Equal(t, sixNotationsFinalFEN, lastStep.Game.FENString)
		})
	}
}

func TestParseNotation_PartialParse(t *testing.T) {
	t.Run("invalid move mid-game returns valid prefix", func(t *testing.T) {
		// Qh7 on move 3 is illegal (queen can't reach h7)
		game := "1. e4 e5\n2. Bc4 Nc6\n3. Qh7 Nf6"
		_, result, err := New().ParseNotation(InputGame{}, game)
		require.NoError(t, err)
		assert.False(t, result.ParseWasSuccessful)
		assert.Equal(t, "Algebraic Notation", result.NotationName)
		assert.Equal(t, 4, result.ValidActionCount, "e4, e5, Bc4, Nc6 are valid")
		assert.Len(t, result.Steps, 4)
		assert.NotEmpty(t, result.Error)
	})

	t.Run("garbage input returns zero valid actions", func(t *testing.T) {
		_, result, err := New().ParseNotation(InputGame{}, "hello world this is not chess")
		require.NoError(t, err)
		assert.False(t, result.ParseWasSuccessful)
		assert.Equal(t, 0, result.ValidActionCount)
		assert.Empty(t, result.Steps)
		assert.NotEmpty(t, result.Error)
	})

	t.Run("truncated ICCF game returns valid prefix", func(t *testing.T) {
		// Third move 9999 is invalid ICCF
		game := "1. 5254 5755\n2. 9999"
		_, result, err := New().ParseNotation(InputGame{}, game)
		require.NoError(t, err)
		assert.False(t, result.ParseWasSuccessful)
		assert.Equal(t, 2, result.ValidActionCount)
		assert.Len(t, result.Steps, 2)
	})

	t.Run("descriptive game with invalid tail returns valid prefix and notation name", func(t *testing.T) {
		game := "1. P-K4 P-K3\n2. P-Q4 P-Q9"
		_, result, err := New().ParseNotation(InputGame{}, game)
		require.NoError(t, err)
		assert.False(t, result.ParseWasSuccessful)
		assert.Equal(t, "Descriptive Notation", result.NotationName)
		assert.Equal(t, 3, result.ValidActionCount, "P-K4, P-K3, P-Q4 are valid")
	})
}

func TestParseNotation_CustomInitialPosition(t *testing.T) {
	// Chernev Game 1 starting position
	inputGame := InputGame{FENString: "8/8/8/8/8/1k5P/8/2K5 w - - 0 1"}
	game := "1. P-R4 K-B5\n2. P-R5 K-Q4"
	_, result, err := New().ParseNotation(inputGame, game)
	require.NoError(t, err)
	assert.True(t, result.ParseWasSuccessful, "error: %v", result.Error)
	assert.Equal(t, "Descriptive Notation", result.NotationName)
	assert.Equal(t, 4, result.ValidActionCount)
}

func TestParseNotation_InvalidInputGame(t *testing.T) {
	_, _, err := New().ParseNotation(InputGame{FENString: "not a fen"}, "1. e4")
	require.Error(t, err)
}

func TestParseNotation_ActionStrings(t *testing.T) {
	// The steps should carry the original action strings for frontend display
	_, result, err := New().ParseNotation(InputGame{}, "1. e4 e5\n2. Nf3 Nc6")
	require.NoError(t, err)
	require.Len(t, result.Steps, 4)
	assert.Equal(t, "e4", result.Steps[0].ActionString)
	assert.Equal(t, "e5", result.Steps[1].ActionString)
	assert.Equal(t, "Nf3", result.Steps[2].ActionString)
	assert.Equal(t, "Nc6", result.Steps[3].ActionString)
}
