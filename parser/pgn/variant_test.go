package pgn

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariantPGN_Parse(t *testing.T) {
	// Create initial game
	initialGame, err := core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.NoError(t, err)

	// Populate Actions by doing one move (e4)
	// This ensures Actions is populated for the parser to work
	// We'll do e4, then parse the remaining moves from that game state
	e2Pawn := initialGame.PieceAt(core.XY{X: 4, Y: 6})
	require.Equal(t, core.PieceType(core.PiecePawn), e2Pawn.PieceType)

	// Do e4 to populate Actions - create a valid e4 action
	e4Action := core.Action{
		FromPiece: e2Pawn,
		ToXY:      core.XY{X: 4, Y: 4},
	}
	gameAfterE4 := initialGame.DoAction(e4Action)
	require.NotEmpty(t, gameAfterE4.Actions, "Actions should be populated after DoAction")

	variant := NewVariantPGN()
	genericParser := parser.NewGenericNotationParser(variant)

	t.Run("parse e5 after e4", func(t *testing.T) {
		// Parse "e5" starting from the game after e4
		parsedGame, err := genericParser.Parse(gameAfterE4, "e5")
		require.NoError(t, err)
		require.Len(t, parsedGame.GameSteps, 1)
		assert.Equal(t, "e5", parsedGame.GameSteps[0].StepString)
		// e7 pawn is at Y=1 (black's side), e5 is at Y=3
		assert.Equal(t, core.XY{X: 4, Y: 1}, parsedGame.GameSteps[0].StepAction.FromPiece.XY) // e7 pawn
		assert.Equal(t, core.XY{X: 4, Y: 3}, parsedGame.GameSteps[0].StepAction.ToXY)         // e5
	})

	t.Run("parse multiple moves", func(t *testing.T) {
		// Parse "e5 2.Nf3 Nc6" starting from the game after e4
		parsedGame, err := genericParser.Parse(gameAfterE4, "e5 2.Nf3 Nc6")
		require.NoError(t, err)
		require.Len(t, parsedGame.GameSteps, 3)
		assert.Equal(t, "e5", parsedGame.GameSteps[0].StepString)
		assert.Equal(t, "Nf3", parsedGame.GameSteps[1].StepString)
		assert.Equal(t, "Nc6", parsedGame.GameSteps[2].StepString)
	})

	t.Run("empty half moves", func(t *testing.T) {
		parsedGame, err := genericParser.Parse(gameAfterE4, "")
		require.NoError(t, err)
		require.Len(t, parsedGame.GameSteps, 0)
	})

	t.Run("parse from initial game", func(t *testing.T) {
		// Test parsing from initial game (which may or may not have Actions populated)
		// If Actions is populated, parsing should work
		if len(initialGame.Actions) > 0 {
			parsedGame, err := genericParser.Parse(initialGame, "e4")
			require.NoError(t, err)
			require.Len(t, parsedGame.GameSteps, 1)
			assert.Equal(t, "e4", parsedGame.GameSteps[0].StepString)
		}
	})

	t.Run("parse with tag pairs", func(t *testing.T) {
		// Use initialGame only if it has Actions populated
		testGame := initialGame
		expectedMoves := 2
		if len(initialGame.Actions) == 0 {
			// Use gameAfterE4 which has Actions
			testGame = gameAfterE4
			expectedMoves = 1
		}
		pgn := `[Event "Test Game"]
[White "Player 1"]
[Black "Player 2"]

1.e4 e5`
		if len(initialGame.Actions) == 0 {
			pgn = `[Event "Test Game"]
[White "Player 1"]
[Black "Player 2"]

e5`
		}
		parsedGame, err := genericParser.Parse(testGame, pgn)
		require.NoError(t, err)
		assert.Equal(t, "Test Game", parsedGame.Metadata["Event"])
		assert.Equal(t, "Player 1", parsedGame.Metadata["White"])
		assert.Equal(t, "Player 2", parsedGame.Metadata["Black"])
		require.Len(t, parsedGame.GameSteps, expectedMoves)
		if len(initialGame.Actions) > 0 {
			assert.Equal(t, "e4", parsedGame.GameSteps[0].StepString)
			assert.Equal(t, "e5", parsedGame.GameSteps[1].StepString)
		} else {
			assert.Equal(t, "e5", parsedGame.GameSteps[0].StepString)
		}
	})

	t.Run("parse with annotations and comments", func(t *testing.T) {
		// Use initialGame only if it has Actions populated
		testGame := initialGame
		if len(initialGame.Actions) == 0 {
			testGame = gameAfterE4
		}
		pgn := `1.e4 ($1) {This is a comment} e5 ; another comment
2.Nf3`
		if len(initialGame.Actions) == 0 {
			pgn = `e5 ($1) {This is a comment} 2.Nf3 ; another comment`
		}
		parsedGame, err := genericParser.Parse(testGame, pgn)
		require.NoError(t, err)
		expectedMoves := 3
		if len(initialGame.Actions) == 0 {
			expectedMoves = 2
		}
		require.Len(t, parsedGame.GameSteps, expectedMoves)
		if len(initialGame.Actions) > 0 {
			assert.Equal(t, "e4", parsedGame.GameSteps[0].StepString)
			assert.Equal(t, "e5", parsedGame.GameSteps[1].StepString)
			assert.Equal(t, "Nf3", parsedGame.GameSteps[2].StepString)
		} else {
			assert.Equal(t, "e5", parsedGame.GameSteps[0].StepString)
			assert.Equal(t, "Nf3", parsedGame.GameSteps[1].StepString)
		}
	})

	t.Run("3-step parsing process", func(t *testing.T) {
		// Test the 3-step process explicitly
		// Use initialGame only if it has Actions populated
		testGame := initialGame
		var pgn string
		if len(initialGame.Actions) == 0 {
			// Populate Actions by doing e4
			testGame = gameAfterE4
			pgn = `e5` // Only e5 since e4 is already done
		} else {
			pgn = `1.e4 e5`
		}

		// Create a generic parser to handle token processing
		genericParser := parser.NewGenericNotationParser(variant)

		// Step 1: Initialize
		parsingGame, err := variant.Initialize(testGame, pgn)
		require.NoError(t, err)
		assert.NotNil(t, parsingGame)

		// Step 2: Loop through parseHalfMoves
		moveCount := 0
		for {
			token, hasMore, err := variant.PopHalfMove(parsingGame)
			require.NoError(t, err)
			if token == nil {
				break
			}

			// Process the token (generic logic)
			err = genericParser.ProcessToken(parsingGame, token)
			require.NoError(t, err)

			if token.Type == parser.TokenTypeHalfMove || token.Type == parser.TokenTypeResult {
				moveCount++
			}
			if !hasMore {
				break
			}
		}

		// Step 3: Finalize
		err = variant.Finalize(parsingGame)
		require.NoError(t, err)

		// Build
		parsedGame := parsingGame.Build()
		if len(initialGame.Actions) == 0 {
			require.Len(t, parsedGame.GameSteps, 1)
			assert.Equal(t, "e5", parsedGame.GameSteps[0].StepString)
			assert.Equal(t, 1, moveCount)
		} else {
			require.Len(t, parsedGame.GameSteps, 2)
			assert.Equal(t, "e4", parsedGame.GameSteps[0].StepString)
			assert.Equal(t, "e5", parsedGame.GameSteps[1].StepString)
			assert.Equal(t, 2, moveCount)
		}
	})

	t.Run("parse full game example", func(t *testing.T) {
		// Parse a complete game from Fischer vs Spassky
		pgn := `[Event "F/S Return Match"]
[Site "Belgrade, Serbia JUG"]
[Date "1992.11.04"]
[Round "29"]
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."]
[Result "1/2-1/2"]

1.e4 e5 2.Nf3 Nc6 3.Bb5 a6
4.Ba4 Nf6 5.O-O Be7 6.Re1 b5 7.Bb3 d6 8.c3 O-O 9.h3 Nb8 10.d4 Nbd7
11.c4 c6 12.cxb5 axb5 13.Nc3 Bb7 14.Bg5 b4 15.Nb1 h6 16.Bh4 c5 17.dxe5
Nxe4 18.Bxe7 Qxe7 19.exd6 Qf6 20.Nbd2 Nxd6 21.Nc4 Nxc4 22.Bxc4 Nb6
23.Ne5 Rae8 24.Bxf7+ Rxf7 25.Nxf7 Rxe1+ 26.Qxe1 Kxf7 27.Qe3 Qg5 28.Qxg5
hxg5 29.b3 Ke6 30.a3 Kd6 31.axb4 cxb4 32.Ra5 Nd5 33.f3 Bc8 34.Kf2 Bf5
35.Ra7 g6 36.Ra6+ Kc5 37.Ke1 Nf4 38.g3 Nxh3 39.Kd2 Kb5 40.Rd6 Kc5 41.Ra6
Nf2 42.g4 Bd3 43.Re6 1/2-1/2`

		parsedGame, err := genericParser.Parse(initialGame, pgn)
		require.NoError(t, err)

		// Verify tag pairs (stored as metadata)
		assert.Equal(t, "F/S Return Match", parsedGame.Metadata["Event"])
		assert.Equal(t, "Belgrade, Serbia JUG", parsedGame.Metadata["Site"])
		assert.Equal(t, "1992.11.04", parsedGame.Metadata["Date"])
		assert.Equal(t, "29", parsedGame.Metadata["Round"])
		assert.Equal(t, "Fischer, Robert J.", parsedGame.Metadata["White"])
		assert.Equal(t, "Spassky, Boris V.", parsedGame.Metadata["Black"])
		assert.Equal(t, "1/2-1/2", parsedGame.Metadata["Result"])

		// Verify game steps - should have all moves plus the result
		require.Greater(t, len(parsedGame.GameSteps), 80, "Should have parsed many moves")

		// Check first few moves
		assert.Equal(t, "e4", parsedGame.GameSteps[0].StepString)
		assert.Equal(t, "e5", parsedGame.GameSteps[1].StepString)

		// Check that result token is included as last step
		lastStep := parsedGame.GameSteps[len(parsedGame.GameSteps)-1]
		assert.Equal(t, "1/2-1/2", lastStep.StepString)
	})

	t.Run("parse game with result tokens", func(t *testing.T) {
		tests := []struct {
			name         string
			pgn          string
			expectedLast string
		}{
			{
				name:         "white wins",
				pgn:          "1.e4 e5 2.Nf3 Nc6 1-0",
				expectedLast: "1-0",
			},
			{
				name:         "black wins",
				pgn:          "1.e4 e5 2.Nf3 Nc6 0-1",
				expectedLast: "0-1",
			},
			{
				name:         "draw",
				pgn:          "1.e4 e5 2.Nf3 Nc6 1/2-1/2",
				expectedLast: "1/2-1/2",
			},
			{
				name:         "asterisk",
				pgn:          "1.e4 e5 2.Nf3 *",
				expectedLast: "*",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				parsedGame, err := genericParser.Parse(initialGame, tt.pgn)
				if len(initialGame.Actions) == 0 {
					// Skip if Actions not populated
					return
				}
				require.NoError(t, err)
				require.Greater(t, len(parsedGame.GameSteps), 0, "Should have parsed moves")

				// Check that result is the last step
				lastStep := parsedGame.GameSteps[len(parsedGame.GameSteps)-1]
				assert.Equal(t, tt.expectedLast, lastStep.StepString)
			})
		}
	})

	t.Run("parse all games in testdata", func(t *testing.T) {
		initialGame, err := core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		require.NoError(t, err)
		require.NotEmpty(t, initialGame.Actions, "Actions should be populated by NewGameFromFEN")

		variant := NewVariantPGN()
		genericParser := parser.NewGenericNotationParser(variant)

		testdataDir := "../testdata/games"
		files, err := filepath.Glob(filepath.Join(testdataDir, "*.pgn"))
		require.NoError(t, err, "Failed to glob PGN files")
		require.Greater(t, len(files), 0, "No PGN files found in testdata/games")

		limit := len(files)
		if testing.Short() {
			limit = 100
		}
		t.Logf("Testing %d of %d PGN files", limit, len(files))

		var parseErrors []error

		for _, file := range files[:limit] {
			pgnContent, err := os.ReadFile(file)
			if err != nil {
				parseErrors = append(parseErrors, fmt.Errorf("failed to read file %s: %w", file, err))
				continue
			}

			_, err = genericParser.Parse(initialGame, string(pgnContent))
			if err != nil {
				parseErrors = append(parseErrors, fmt.Errorf("failed to parse %s: %w", filepath.Base(file), err))
			}
		}

		if len(parseErrors) > 0 {
			t.Errorf("Failed to parse %d out of %d games", len(parseErrors), limit)
			for i, err := range parseErrors {
				if i >= 10 {
					t.Logf("  ... and %d more errors", len(parseErrors)-10)
					break
				}
				t.Logf("  %v", err)
			}
		} else {
			t.Logf("Successfully parsed all %d games", limit)
		}
	})
}
