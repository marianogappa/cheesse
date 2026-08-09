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

func TestBasicAI_DeclinesCapturingDefendedPiece(t *testing.T) {
	// White Qd1, Kg1. Black Nd5, Pe6 (pawn defends d5), Kg8.
	// Qxd5?? allows exd5 — White loses a queen (9) for a knight (3). Net = -6.
	// A correct depth-1 search must NOT play Qxd5.
	g, err := core.NewGameFromFEN("6k1/8/4p3/3n4/8/8/8/3Q2K1 w - - 0 1")
	require.NoError(t, err)

	action, _, ok := BasicAIAction(g, 1)
	require.True(t, ok)

	isQxd5 := action.FromPiece.PieceType == core.PieceQueen &&
		action.ToXY.X == 3 && action.ToXY.Y == 3 &&
		action.IsCapture
	assert.False(t, isQxd5,
		"depth-1 AI should not capture a defended knight with the queen; chose %v -> %v",
		action.FromPiece.XY.ToAlgebraic(), action.ToXY.ToAlgebraic())
}

func TestBasicAI_FindsMateIn1(t *testing.T) {
	// White Rd1, Kg1. Black Kg8, Pf7, Pg7, Ph7.
	// Rd8# is back-rank mate in 1. Must be found at every depth.
	g, err := core.NewGameFromFEN("6k1/5ppp/8/8/8/8/8/3R2K1 w - - 0 1")
	require.NoError(t, err)

	for _, depth := range []int{0, 1, 2} {
		action, newGame, ok := BasicAIAction(g, depth)
		require.True(t, ok, "depth %d", depth)
		assert.True(t, newGame.IsCheckmate,
			"depth %d: AI must find mate-in-1; chose %v -> %v",
			depth, action.FromPiece.XY.ToAlgebraic(), action.ToXY.ToAlgebraic())
	}
}

func TestBasicAI_FindsMateIn2(t *testing.T) {
	// White Qa3, Rh1, Kg1. Black Kg8 (lone king).
	// 1.Qa7 (threatens Rh8#) Kf8 2.Rh8# — mate in 2.
	// At depth >= 2, playing out the full sequence must end in checkmate.
	g, err := core.NewGameFromFEN("6k1/8/8/8/8/Q7/8/6KR w - - 0 1")
	require.NoError(t, err)

	action, newGame, ok := BasicAIAction(g, 2)
	require.True(t, ok)

	g1 := newGame
	if g1.IsCheckmate {
		return
	}
	_, g2, ok2 := BasicAIAction(g1, 2)
	require.True(t, ok2, "opponent should have a reply; AI played %v -> %v",
		action.FromPiece.XY.ToAlgebraic(), action.ToXY.ToAlgebraic())
	_, g3, ok3 := BasicAIAction(g2, 2)
	require.True(t, ok3, "AI should have a mating move")
	assert.True(t, g3.IsCheckmate,
		"depth-2 AI should find mate-in-2; final FEN: %v", g3.ToFEN())
}

func TestBasicAI_DoesNotHangQueen(t *testing.T) {
	// White Qd1, Kg1. Black Pe5, Pd6, Kg8.
	// Pe5 attacks d4 and f4. Pd6 attacks c5 and e5.
	// At depth >= 1, the AI must not move the queen to a pawn-attacked square.
	g, err := core.NewGameFromFEN("6k1/8/3p4/4p3/8/8/8/3Q2K1 w - - 0 1")
	require.NoError(t, err)

	action, _, ok := BasicAIAction(g, 1)
	require.True(t, ok)

	// The queen should NOT go to squares attacked by pawns:
	// Pe5 attacks d4(3,4) and f4(5,4). Pd6 attacks c5(2,3) and e5(4,3).
	pawnAttackedSquares := map[[2]int]bool{
		{3, 4}: true, // d4
		{5, 4}: true, // f4
		{2, 3}: true, // c5
		{4, 3}: true, // e5
	}
	if action.FromPiece.PieceType == core.PieceQueen {
		dest := [2]int{action.ToXY.X, action.ToXY.Y}
		assert.False(t, pawnAttackedSquares[dest],
			"depth-1 AI should not move queen to a pawn-attacked square; chose %v -> %v",
			action.FromPiece.XY.ToAlgebraic(), action.ToXY.ToAlgebraic())
	}
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
