package ai

import (
	"math/rand"
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomAction_ReturnsLegalAction(t *testing.T) {
	g := core.NewDefaultGame()
	rng := rand.New(rand.NewSource(42))
	action, newGame, ok := RandomAction(g, rng)
	require.True(t, ok)
	assert.NotEqual(t, core.Action{}, action)
	assert.False(t, action.IsResign)
	assert.NotEqual(t, g.ToFEN(), newGame.ToFEN())
}

func TestRandomAction_SeededDeterminism(t *testing.T) {
	g := core.NewDefaultGame()
	rng1 := rand.New(rand.NewSource(123))
	rng2 := rand.New(rand.NewSource(123))
	a1, _, _ := RandomAction(g, rng1)
	a2, _, _ := RandomAction(g, rng2)
	assert.Equal(t, a1, a2)
}

func TestRandomAction_GameOver(t *testing.T) {
	// Scholar's mate final position: Qf7# Black is checkmated.
	g, err := core.NewGameFromFEN("r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4")
	require.NoError(t, err)
	require.True(t, g.IsGameOver, "position should be game over (checkmate)")
	_, _, ok := RandomAction(g, rand.New(rand.NewSource(0)))
	assert.False(t, ok)
}

func TestBasicAI_ChoosesCheckmate(t *testing.T) {
	// Scholar's mate position: White Qh5, Bc4 can play Qxf7#.
	g, err := core.NewGameFromFEN("r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 3")
	require.NoError(t, err)
	action, newGame, ok := BasicAIAction(g, 0)
	require.True(t, ok)
	assert.True(t, newGame.IsCheckmate, "AI should choose checkmate; chose %v -> %v", action.FromPiece.XY.ToAlgebraic(), action.ToXY.ToAlgebraic())
}

func TestBasicAI_PrefersCapture(t *testing.T) {
	// White Ka1, Nb3. Black Kh8, Qd4. Nxd4 captures the queen (9 material gain).
	g, err := core.NewGameFromFEN("7k/8/8/8/3q4/1N6/8/K7 w - - 0 1")
	require.NoError(t, err)
	action, _, ok := BasicAIAction(g, 0)
	require.True(t, ok)
	assert.True(t, action.IsCapture, "AI should prefer capturing the queen; chose %v -> %v", action.FromPiece.XY.ToAlgebraic(), action.ToXY.ToAlgebraic())
}

func TestBasicAI_PromotesToQueen(t *testing.T) {
	// White Kb2, Pf7. Black Kh8. f8=Q promotes.
	g, err := core.NewGameFromFEN("7k/5P2/8/8/8/8/1K6/8 w - - 0 1")
	require.NoError(t, err)
	action, _, ok := BasicAIAction(g, 0)
	require.True(t, ok)
	assert.True(t, action.IsPromotion)
	assert.Equal(t, core.PieceType(core.PieceQueen), action.PromotionPieceType)
}

func TestBasicAI_GameOverReturnsNoMove(t *testing.T) {
	// Scholar's mate final position.
	g, err := core.NewGameFromFEN("r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4")
	require.NoError(t, err)
	_, _, ok := BasicAIAction(g, 0)
	assert.False(t, ok)
}

func TestBasicAI_Depth1(t *testing.T) {
	g := core.NewDefaultGame()
	action, newGame, ok := BasicAIAction(g, 1)
	require.True(t, ok)
	assert.NotEqual(t, core.Action{}, action)
	assert.NotEqual(t, g.ToFEN(), newGame.ToFEN())
}

func TestBasicAI_Depth2(t *testing.T) {
	// Endgame position for speed; depth 2 from the default position is slow.
	g, err := core.NewGameFromFEN("7k/5P2/8/8/3n4/1B6/8/K7 w - - 0 1")
	require.NoError(t, err)
	action, newGame, ok := BasicAIAction(g, 2)
	require.True(t, ok)
	assert.NotEqual(t, core.Action{}, action)
	assert.NotEqual(t, g.ToFEN(), newGame.ToFEN())
}

func TestBasicAI_AutoplayGameCompletes(t *testing.T) {
	// Two AI agents play against each other until the game ends or 200 moves.
	g := core.NewDefaultGame()
	for i := 0; i < 200 && !g.IsGameOver; i++ {
		_, newGame, ok := BasicAIAction(g, 0)
		if !ok {
			break
		}
		g = newGame
	}
	// The game should eventually end (checkmate, stalemate, or draw).
	// No panic means the AI produced valid moves throughout.
}
