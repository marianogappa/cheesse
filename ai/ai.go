package ai

import (
	"math"
	"math/rand"

	"github.com/marianogappa/cheesse/core"
)

// RandomAction picks a uniformly random legal non-resign action for the side to
// move. Returns the action and the resulting game, or (Action{}, game, false)
// when the game is already over.
func RandomAction(g core.Game, rng *rand.Rand) (core.Action, core.Game, bool) {
	nonResign := nonResignActions(g)
	if len(nonResign) == 0 {
		return core.Action{}, g, false
	}
	action := nonResign[rng.Intn(len(nonResign))]
	return action, g.DoAction(action), true
}

// BasicAIAction selects the best action for the side to move using minimax with
// alpha-beta pruning at the given depth. Depth 0 evaluates the position after
// each candidate move; depth 1 looks one full move (two plies) ahead; etc.
// Returns the action, resulting game, and true; or (Action{}, game, false) when
// the game is already over.
func BasicAIAction(g core.Game, depth int) (core.Action, core.Game, bool) {
	nonResign := nonResignActions(g)
	if len(nonResign) == 0 {
		return core.Action{}, g, false
	}

	player := int(g.Turn())
	bestScore := int64(math.MinInt64)
	bestIdx := 0

	for i, action := range nonResign {
		newGame := g.DoAction(action)
		score := alphabeta(newGame, action, player, depth, math.MinInt64, math.MaxInt64)
		if score > bestScore || (score == bestScore && tieBreakPrefer(action, nonResign[bestIdx])) {
			bestScore = score
			bestIdx = i
		}
	}

	chosen := nonResign[bestIdx]
	return chosen, g.DoAction(chosen), true
}

func nonResignActions(g core.Game) []core.Action {
	out := make([]core.Action, 0, len(g.Actions))
	for _, a := range g.Actions {
		if !a.IsResign {
			out = append(out, a)
		}
	}
	return out
}

func alphabeta(g core.Game, lastAction core.Action, player int, depth int, alpha, beta int64) int64 {
	if depth == 0 {
		return evaluate(g, lastAction, player)
	}

	nonResign := nonResignActions(g)
	if len(nonResign) == 0 {
		return evaluate(g, lastAction, player)
	}

	if int(g.Turn()) != player {
		// Maximizing (our turn next is opponent's perspective → we maximize our score)
		v := int64(math.MinInt64)
		for _, action := range nonResign {
			newGame := g.DoAction(action)
			score := alphabeta(newGame, action, player, depth-1, alpha, beta)
			if score > v {
				v = score
			}
			if v > alpha {
				alpha = v
			}
			if beta <= alpha {
				break
			}
		}
		return v
	}

	// Minimizing (opponent's turn)
	v := int64(math.MaxInt64)
	for _, action := range nonResign {
		newGame := g.DoAction(action)
		score := alphabeta(newGame, action, player, depth-1, alpha, beta)
		if score < v {
			v = score
		}
		if v < beta {
			beta = v
		}
		if beta <= alpha {
			break
		}
	}
	return v
}

func sign(owner, player int) int {
	if owner == player {
		return 1
	}
	return -1
}

func materialValue(pt core.PieceType) int {
	switch pt {
	case core.PieceQueen:
		return 9
	case core.PieceRook:
		return 5
	case core.PieceBishop, core.PieceKnight:
		return 3
	case core.PiecePawn:
		return 1
	}
	return 0
}

func evaluate(g core.Game, lastAction core.Action, player int) int64 {
	actionTurn := int(lastAction.FromPiece.Owner)
	opponent := 1 - player

	// Checkmate is decisive
	if g.IsCheckmate {
		return math.MaxInt64 / 2 * int64(sign(actionTurn, player))
	}
	if g.IsStalemate || g.IsDraw {
		return 0
	}

	// Material
	playerMaterial, opponentMaterial := 0, 0
	for _, piece := range g.Pieces(core.Color(player)) {
		playerMaterial += materialValue(piece.PieceType)
	}
	for _, piece := range g.Pieces(core.Color(opponent)) {
		opponentMaterial += materialValue(piece.PieceType)
	}
	totalMaterial := int64(playerMaterial-opponentMaterial) * 1_000_000_000

	// Check bonus
	checkBonus := int64(0)
	if g.IsCheck {
		checkBonus = 100_000_000 * int64(sign(actionTurn, player))
	}

	// Center control: pieces in the center (files c-f, ranks 4-5)
	centerPlayer, centerOpponent := 0, 0
	for _, piece := range g.Pieces(core.Color(player)) {
		if piece.XY.X >= 2 && piece.XY.X <= 5 && piece.XY.Y >= 3 && piece.XY.Y <= 4 {
			centerPlayer++
		}
	}
	for _, piece := range g.Pieces(core.Color(opponent)) {
		if piece.XY.X >= 2 && piece.XY.X <= 5 && piece.XY.Y >= 3 && piece.XY.Y <= 4 {
			centerOpponent++
		}
	}
	centerBonus := int64(centerPlayer-centerOpponent) * 100

	// Castling incentive
	castleBonus := int64(0)
	if lastAction.IsCastle {
		castleBonus = 100 * int64(sign(actionTurn, player))
	}

	return totalMaterial + checkBonus + centerBonus + castleBonus
}

// tieBreakPrefer returns true if a should be preferred over b when scores are equal.
func tieBreakPrefer(a, b core.Action) bool {
	return actionPriority(a) > actionPriority(b)
}

func actionPriority(a core.Action) int {
	switch {
	case a.IsPromotion && a.IsCapture:
		return 6
	case a.IsPromotion:
		return 5
	case a.IsCapture:
		return 4
	case a.IsEnPassantCapture:
		return 3
	case a.IsCastle:
		return 2
	default:
		return 1
	}
}
