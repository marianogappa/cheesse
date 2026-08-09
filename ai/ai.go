package ai

import (
	"math/rand"
	"sort"
	"time"

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

// BasicAIAction selects the best action for the side to move using negamax with
// alpha-beta pruning, MVV-LVA move ordering and quiescence search at the given
// depth. Depth 0 evaluates each candidate move with quiescence (captures and
// promotions are resolved); depth 1 additionally looks one ply ahead; etc.
// Returns the action, resulting game, and true; or (Action{}, game, false) when
// the game is already over.
func BasicAIAction(g core.Game, depth int) (core.Action, core.Game, bool) {
	moves := orderedMoves(g)
	if len(moves) == 0 {
		return core.Action{}, g, false
	}
	s := &searcher{}
	bestIdx, _ := searchRoot(s, g, moves, depth)
	chosen := moves[bestIdx]
	return chosen, g.DoAction(chosen), true
}

// TimedAIAction selects the best action for the side to move using iterative
// deepening within the given time budget: it searches at increasing depths and
// returns the best move from the deepest fully completed iteration. Partial
// iterations are discarded, so the result is deterministic per depth reached.
// Returns the action, resulting game, and true; or (Action{}, game, false) when
// the game is already over.
func TimedAIAction(g core.Game, budget time.Duration) (core.Action, core.Game, bool) {
	moves := orderedMoves(g)
	if len(moves) == 0 {
		return core.Action{}, g, false
	}
	deadline := time.Now().Add(budget)

	// Depth 0 runs without a deadline so a move is always available.
	s := &searcher{}
	bestIdx, bestScore := searchRoot(s, g, moves, 0)
	bestAction := moves[bestIdx]

	for depth := 1; depth <= maxSearchDepth; depth++ {
		// Searching the previous iteration's best move first maximizes pruning.
		moveToFront(moves, bestAction)
		sd := &searcher{deadline: deadline, hasDeadline: true}
		idx, score := searchRoot(sd, g, moves, depth)
		if sd.stopped {
			break
		}
		bestAction, bestScore = moves[idx], score
		// A forced mate was found; deeper iterations cannot find a shorter one.
		if bestScore >= mateScore-2*maxSearchDepth {
			break
		}
		if time.Now().After(deadline) {
			break
		}
	}

	return bestAction, g.DoAction(bestAction), true
}

const (
	infScore       = 2_000_000
	mateScore      = 1_000_000
	maxSearchDepth = 30
	maxQuiescePly  = 60
	timeCheckMask  = 1023 // check the clock every 1024 nodes
)

type searcher struct {
	deadline    time.Time
	hasDeadline bool
	nodes       int
	stopped     bool
}

func (s *searcher) outOfTime() bool {
	if s.stopped {
		return true
	}
	s.nodes++
	if s.hasDeadline && s.nodes&timeCheckMask == 0 && time.Now().After(s.deadline) {
		s.stopped = true
	}
	return s.stopped
}

// searchRoot scores every root move with negamax and returns the index of the
// best one and its score. When the searcher runs out of time the result is
// meaningless and must be discarded (check s.stopped).
func searchRoot(s *searcher, g core.Game, moves []core.Action, depth int) (int, int) {
	bestIdx, bestScore := 0, -infScore
	alpha, beta := -infScore, infScore
	for i, a := range moves {
		child := g.DoAction(a)
		score := -s.negamax(child, 1, depth, -beta, -alpha)
		if s.stopped {
			return bestIdx, bestScore
		}
		if score > bestScore || (score == bestScore && tieBreakPrefer(a, moves[bestIdx])) {
			bestScore, bestIdx = score, i
		}
		if score > alpha {
			alpha = score
		}
	}
	return bestIdx, bestScore
}

// negamax returns the score of the position from the perspective of the side to
// move. Mates are scored as mateScore - ply so shorter mates score higher, and
// claimable draws score 0 so a winning engine avoids repetition/50-move draws.
func (s *searcher) negamax(g core.Game, ply, depth int, alpha, beta int) int {
	if s.outOfTime() {
		return 0
	}
	if g.IsCheckmate {
		return -(mateScore - ply)
	}
	if g.IsStalemate || g.IsDraw || g.CanClaimDraw {
		return 0
	}
	if depth <= 0 {
		return s.quiesce(g, ply, alpha, beta)
	}

	best := -infScore
	for _, a := range orderedMoves(g) {
		child := g.DoAction(a)
		score := -s.negamax(child, ply+1, depth-1, -beta, -alpha)
		if s.stopped {
			return 0
		}
		if score > best {
			best = score
		}
		if best > alpha {
			alpha = best
		}
		if alpha >= beta {
			break
		}
	}
	return best
}

// quiesce extends the search through captures and promotions until the position
// is quiet, avoiding the horizon effect of statically evaluating mid-exchange.
func (s *searcher) quiesce(g core.Game, ply int, alpha, beta int) int {
	if s.outOfTime() {
		return 0
	}
	if g.IsCheckmate {
		return -(mateScore - ply)
	}
	if g.IsStalemate || g.IsDraw || g.CanClaimDraw {
		return 0
	}

	standPat := evaluate(g)
	if ply >= maxQuiescePly || standPat >= beta {
		return standPat
	}
	if standPat > alpha {
		alpha = standPat
	}

	best := standPat
	for _, a := range orderedMoves(g) {
		if !a.IsCapture && !a.IsPromotion {
			continue
		}
		child := g.DoAction(a)
		score := -s.quiesce(child, ply+1, -beta, -alpha)
		if s.stopped {
			return 0
		}
		if score > best {
			best = score
		}
		if best > alpha {
			alpha = best
		}
		if alpha >= beta {
			break
		}
	}
	return best
}

func nonResignActions(g core.Game) []core.Action {
	out := make([]core.Action, 0, len(g.Actions))
	for _, a := range g.Actions {
		if !a.IsResign && !a.IsDraw {
			out = append(out, a)
		}
	}
	return out
}

// orderedMoves returns the legal board moves sorted for alpha-beta efficiency:
// captures first by MVV-LVA (most valuable victim, least valuable attacker),
// promotions next, quiet moves last.
func orderedMoves(g core.Game) []core.Action {
	moves := nonResignActions(g)
	sort.SliceStable(moves, func(i, j int) bool {
		return moveOrderScore(moves[i]) > moveOrderScore(moves[j])
	})
	return moves
}

func moveOrderScore(a core.Action) int {
	score := 0
	if a.IsCapture {
		score = 1_000_000 + 10*pieceValue(a.CapturedPiece.PieceType) - pieceValue(a.FromPiece.PieceType)
	}
	if a.IsPromotion {
		score += 500_000 + pieceValue(a.PromotionPieceType)
	}
	return score
}

func moveToFront(moves []core.Action, a core.Action) {
	for i := range moves {
		if moves[i] == a {
			copy(moves[1:i+1], moves[:i])
			moves[0] = a
			return
		}
	}
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
