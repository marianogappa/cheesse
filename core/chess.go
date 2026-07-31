package core

import "math/bits"

// shallowCloneForMove creates a clone optimized for move-legality checking: the
// board layout (bitboards) is copied by the struct copy, and the derived fields
// (flags, InCheckBy, Actions) are dropped, since they are not needed for the
// "does this move leave the king in check?" test.
func (g Game) shallowCloneForMove() Game {
	clonedGame := g
	clonedGame.IsCheck = false
	clonedGame.IsDoubleCheck = false
	clonedGame.IsDiscoverCheck = false
	clonedGame.IsCheckmate = false
	clonedGame.IsStalemate = false
	clonedGame.IsDraw = false
	clonedGame.CanClaimDraw = false
	clonedGame.IsGameOver = false
	clonedGame.GameOverWinner = 0
	clonedGame.InCheckBy = nil
	clonedGame.Actions = nil
	clonedGame.positionHistory = nil
	return clonedGame
}

// updateBoardLayout updates a game's layout-only (i.e. bitboards) after a given action, so that the resulting
// layout can be checked for checks, checkmates, etc.  This method is meant to be used as a dry-run to decide if the
// action can actually be executed. Should only be called by piece.calculateAllActions and game.DoAction.
//
// Note that this method assumes the action is fully built and valid (bounds, no friendly piece at destination, etc.).
func (g Game) updateBoardLayout(a Action) Game {
	clonedGame := g.shallowCloneForMove()

	// Special case for resignation action, because it doesn't require board changes
	if a.IsResign {
		return clonedGame
	}

	owner := a.FromPiece.Owner
	toPieceType := a.FromPiece.PieceType
	if a.IsPromotion {
		toPieceType = a.PromotionPieceType
	}

	// Remove pieces at {from, to} locations, and place fromPiece at "to" location
	clonedGame.clearSq(owner, a.FromPiece.PieceType, sqOf(a.FromPiece.XY))
	if a.IsCapture && !a.IsEnPassantCapture {
		clonedGame.clearSq(a.CapturedPiece.Owner, a.CapturedPiece.PieceType, sqOf(a.CapturedPiece.XY))
	}
	clonedGame.setSq(owner, toPieceType, sqOf(a.ToXY))

	// Extra movements and deletions in the case of en passant capture and castling
	homeRank := 0
	if owner == ColorWhite {
		homeRank = 7
	}
	switch {
	case a.IsEnPassantCapture:
		clonedGame.clearSq(a.CapturedPiece.Owner, PiecePawn, sqOf(a.CapturedPiece.XY))
	case a.IsQueensideCastle:
		clonedGame.clearSq(owner, PieceRook, sqOf(XY{0, homeRank}))
		clonedGame.setSq(owner, PieceRook, sqOf(XY{3, homeRank}))
	case a.IsKingsideCastle:
		clonedGame.clearSq(owner, PieceRook, sqOf(XY{7, homeRank}))
		clonedGame.setSq(owner, PieceRook, sqOf(XY{5, homeRank}))
	}

	return clonedGame
}

func (g Game) calculateAllActions() []Action {
	if g.IsGameOver {
		return []Action{}
	}
	turn := g.Turn()
	actions := make([]Action, 0, 64)
	for occ := g.occ[turn]; occ != 0; occ &= occ - 1 {
		actions = g.pieceAtSq(bits.TrailingZeros64(occ)).appendActions(actions, g)
	}
	actions = append(actions, Action{FromPiece: Piece{Owner: turn}, IsResign: true})
	return actions
}

func (g Game) Turn() color {
	if g.MoveNumber%2 == 0 {
		return ColorWhite
	}
	return ColorBlack
}

func (p Piece) calculateAllActions(g Game) []Action {
	return p.appendActions(make([]Action, 0, 8), g)
}

var promotionPieceTypes = [4]PieceType{PieceQueen, PieceBishop, PieceKnight, PieceRook}

// appendActions appends all legal actions for the piece to the given slice. Pseudo-legal
// destinations are computed from the precomputed attack bitboards; each is then filtered
// by the "does this move leave the king in check?" test.
func (p Piece) appendActions(actions []Action, g Game) []Action {
	if g.IsGameOver || g.Turn() != p.Owner || p.PieceType == PieceNone || g.pieceAtSq(sqOf(p.XY)) != p {
		return actions
	}
	sq := sqOf(p.XY)
	own := g.occ[p.Owner]
	occ := g.occAll()

	var targets uint64
	switch p.PieceType {
	case PieceQueen:
		targets = queenAttacks(sq, occ) &^ own
	case PieceKing:
		targets = kingAttacks[sq] &^ own
	case PieceBishop:
		targets = bishopAttacks(sq, occ) &^ own
	case PieceKnight:
		targets = knightAttacks[sq] &^ own
	case PieceRook:
		targets = rookAttacks(sq, occ) &^ own
	case PiecePawn:
		targets = p.pawnTargets(g, sq, occ)
	}

	for t := targets; t != 0; t &= t - 1 {
		actions = p.appendActionsTo(actions, g, xyOfSq(bits.TrailingZeros64(t)))
	}

	if p.PieceType == PieceKing {
		actions = p.appendCastleActions(actions, g)
	}
	return actions
}

// pawnTargets computes the pseudo-legal destination squares for a pawn: forward pushes
// (single, and double from the home rank) onto empty squares, plus diagonal captures
// onto opponent pieces or the en passant target square.
func (p Piece) pawnTargets(g Game, sq int, occ uint64) uint64 {
	var targets uint64
	if p.Owner == ColorBlack && p.XY.Y < 7 {
		fwd := sq + 8
		if occ&sqBit(fwd) == 0 {
			targets |= sqBit(fwd)
			if p.XY.Y == 1 && occ&sqBit(fwd+8) == 0 {
				targets |= sqBit(fwd + 8)
			}
		}
	}
	if p.Owner == ColorWhite && p.XY.Y > 0 {
		fwd := sq - 8
		if occ&sqBit(fwd) == 0 {
			targets |= sqBit(fwd)
			if p.XY.Y == 6 && occ&sqBit(fwd-8) == 0 {
				targets |= sqBit(fwd - 8)
			}
		}
	}
	captureTargets := g.occ[opponent(p.Owner)]
	if g.IsLastMoveEnPassant {
		captureTargets |= sqBit(sqOf(g.EnPassantTargetSquare))
	}
	return targets | pawnCaptureAttacks[p.Owner][sq]&captureTargets
}

// appendActionsTo builds the action(s) for a pseudo-legal destination square (4 actions
// in the case of a promotion), filtering out those that leave the own king in check.
func (p Piece) appendActionsTo(actions []Action, g Game, toXY XY) []Action {
	a := Action{FromPiece: p, ToXY: toXY}

	// Set capture context
	if capturedPiece := g.PieceAt(toXY); capturedPiece.PieceType != PieceNone {
		a.IsCapture = true
		a.CapturedPiece = capturedPiece
	}
	// Set capture context in the case of an en passant capture
	if p.PieceType == PiecePawn && toXY.X != p.XY.X && !a.IsCapture {
		a.IsCapture = true
		a.IsEnPassantCapture = true
		switch p.Owner {
		case ColorBlack:
			a.CapturedPiece = g.PieceAt(XY{X: toXY.X, Y: toXY.Y - 1})
		case ColorWhite:
			a.CapturedPiece = g.PieceAt(XY{X: toXY.X, Y: toXY.Y + 1})
		}
	}

	// check if moving puts the owner's King in check (the promoted piece type cannot
	// affect this, so the check is done once for all 4 promotion actions)
	newGame := g.updateBoardLayout(a)
	if newGame.attackersOf(int(newGame.kingSq[p.Owner]), p.Owner) != 0 {
		return actions
	}

	// Set promotion context
	if p.PieceType == PiecePawn && (toXY.Y == 0 || toXY.Y == 7) {
		a.IsPromotion = true
		for _, promotionPieceType := range promotionPieceTypes {
			a.PromotionPieceType = promotionPieceType
			actions = append(actions, a)
		}
		return actions
	}

	return append(actions, a)
}

func (p Piece) appendCastleActions(actions []Action, g Game) []Action {
	var canQueenside, canKingside bool
	homeRank := 0
	switch p.Owner {
	case ColorBlack:
		canQueenside, canKingside = g.CanBlackCastle && g.CanBlackQueensideCastle, g.CanBlackCastle && g.CanBlackKingsideCastle
	case ColorWhite:
		canQueenside, canKingside = g.CanWhiteCastle && g.CanWhiteQueensideCastle, g.CanWhiteCastle && g.CanWhiteKingsideCastle
		homeRank = 7
	}
	if (!canQueenside && !canKingside) || !p.XY.eq(XY{4, homeRank}) {
		return actions
	}

	occ := g.occAll()
	castles := []struct {
		allowed    bool
		castleType castleType
		toX        int
	}{
		{canQueenside, castleTypeQueenside, 2},
		{canKingside, castleTypeKingside, 6},
	}
	for _, c := range castles {
		if !c.allowed ||
			occ&castleEmptyMasks[p.Owner][c.castleType] != 0 ||
			g.anySqThreatened(castleUnthreatenedSqs[p.Owner][c.castleType], p.Owner) {
			continue
		}
		a := Action{
			FromPiece:         p,
			ToXY:              XY{c.toX, homeRank},
			IsCastle:          true,
			IsQueensideCastle: c.castleType == castleTypeQueenside,
			IsKingsideCastle:  c.castleType == castleTypeKingside,
		}
		actions = append(actions, a)
	}
	return actions
}

// doAction executes the given action on the given game.
// It assumes that the game is in a state where this action can be executed.
// It assumes that the action is fully-correctly created and it's valid.
// It fully updates the game context.
// This is an expensive method (due to having to check for check and checkmate), so use only if needed.
func (g Game) DoAction(a Action) Game {
	newGame := g.updateBoardLayout(a)
	lastTurn := g.Turn()

	// Special case for resignation action
	if a.IsResign {
		newGame.IsGameOver = true
		newGame.GameOverWinner = opponent(lastTurn)

		// TODO is it necessary to update other things?
		return newGame
	}

	// Castling context update: moving player's king or rook
	switch {
	case lastTurn == ColorBlack && (a.IsCastle || a.FromPiece.PieceType == PieceKing):
		newGame.CanBlackCastle = false
		newGame.CanBlackQueensideCastle = false
		newGame.CanBlackKingsideCastle = false
	case lastTurn == ColorWhite && (a.IsCastle || a.FromPiece.PieceType == PieceKing):
		newGame.CanWhiteCastle = false
		newGame.CanWhiteQueensideCastle = false
		newGame.CanWhiteKingsideCastle = false
	case lastTurn == ColorBlack && a.FromPiece.PieceType == PieceRook && a.FromPiece.XY == XY{X: 0, Y: 0}:
		newGame.CanBlackQueensideCastle = false
		newGame.CanBlackCastle = newGame.CanBlackKingsideCastle
	case lastTurn == ColorBlack && a.FromPiece.PieceType == PieceRook && a.FromPiece.XY == XY{X: 7, Y: 0}:
		newGame.CanBlackKingsideCastle = false
		newGame.CanBlackCastle = newGame.CanBlackQueensideCastle
	case lastTurn == ColorWhite && a.FromPiece.PieceType == PieceRook && a.FromPiece.XY == XY{X: 0, Y: 7}:
		newGame.CanWhiteQueensideCastle = false
		newGame.CanWhiteCastle = newGame.CanWhiteKingsideCastle
	case lastTurn == ColorWhite && a.FromPiece.PieceType == PieceRook && a.FromPiece.XY == XY{X: 7, Y: 7}:
		newGame.CanWhiteKingsideCastle = false
		newGame.CanWhiteCastle = newGame.CanWhiteQueensideCastle
	}

	// Castling context update: opponent's rook captured on its home square
	if a.IsCapture {
		switch a.ToXY {
		case XY{X: 0, Y: 0}:
			newGame.CanBlackQueensideCastle = false
			newGame.CanBlackCastle = newGame.CanBlackKingsideCastle
		case XY{X: 7, Y: 0}:
			newGame.CanBlackKingsideCastle = false
			newGame.CanBlackCastle = newGame.CanBlackQueensideCastle
		case XY{X: 0, Y: 7}:
			newGame.CanWhiteQueensideCastle = false
			newGame.CanWhiteCastle = newGame.CanWhiteKingsideCastle
		case XY{X: 7, Y: 7}:
			newGame.CanWhiteKingsideCastle = false
			newGame.CanWhiteCastle = newGame.CanWhiteQueensideCastle
		}
	}

	newGame.MoveNumber = g.MoveNumber + 1
	if lastTurn == ColorBlack {
		newGame.FullMoveNumber = g.FullMoveNumber + 1
	}
	isDoubleAdvance := a.FromPiece.PieceType == PiecePawn && abs(a.ToXY.Y-a.FromPiece.XY.Y) == 2
	newGame.IsLastMoveEnPassant = isDoubleAdvance
	if isDoubleAdvance && lastTurn == ColorBlack {
		newGame.EnPassantTargetSquare = XY{X: a.ToXY.X, Y: a.ToXY.Y - 1}
	}
	if isDoubleAdvance && lastTurn == ColorWhite {
		newGame.EnPassantTargetSquare = XY{X: a.ToXY.X, Y: a.ToXY.Y + 1}
	}

	newGame.HalfMoveClock = g.HalfMoveClock + 1
	if a.IsCapture || a.FromPiece.PieceType == PiecePawn {
		newGame.HalfMoveClock = 0
	}

	// Maintain position history for repetition detection. Captures and pawn moves are
	// irreversible, so earlier positions can never repeat and the history restarts.
	if newGame.HalfMoveClock > 0 {
		newGame.positionHistory = append(newGame.positionHistory, g.positionHistory...)
	}
	newGame.positionHistory = append(newGame.positionHistory, newGame.positionHash())

	newGame = newGame.calculateCriticalFlags()

	if newGame.IsCheck {
		newGame.IsDoubleCheck = len(newGame.InCheckBy) >= 2

		hasRevealedChecker := false
		for _, checker := range newGame.InCheckBy {
			if checker.XY != a.ToXY {
				hasRevealedChecker = true
				break
			}
		}
		newGame.IsDiscoverCheck = hasRevealedChecker
	}

	return newGame
}

func (g Game) calculateCriticalFlags() Game {
	turn := g.Turn()

	g.IsCheck = false
	g.IsDoubleCheck = false
	g.IsDiscoverCheck = false
	g.IsCheckmate = false
	g.IsStalemate = false
	g.IsDraw = false
	g.CanClaimDraw = false
	g.IsGameOver = false
	g.GameOverWinner = -1
	g.InCheckBy = []Piece{}

	g.InCheckBy = g.King(turn).threatenedBy(g)
	if len(g.InCheckBy) > 0 {
		g.IsCheck = true
	}

	g.Actions = g.calculateAllActions() // This is incredibly expensive!
	if len(g.Actions) == 1 && g.Actions[0].IsResign {
		g.IsCheckmate = g.IsCheck
		g.IsStalemate = !g.IsCheck
	}

	// Draw rules, per FIDE: the 75-move rule, fivefold repetition and insufficient
	// material (dead position) end the game automatically; the 50-move rule and
	// threefold repetition make a draw claimable by the player to move.
	repetitions := g.repetitionCount()
	switch {
	case g.HalfMoveClock >= 150, repetitions >= 5, g.isInsufficientMaterial():
		g.IsDraw = true
	case g.HalfMoveClock >= 100, repetitions >= 3:
		g.CanClaimDraw = true
	}

	if g.IsCheckmate || g.IsStalemate || g.IsDraw {
		g.IsGameOver = true
	}
	if g.IsCheckmate {
		g.GameOverWinner = opponent(turn)
	}

	return g
}

// repetitionCount returns how many times the current position has occurred, based on
// the position history since the last irreversible move.
func (g Game) repetitionCount() int {
	if len(g.positionHistory) == 0 {
		return 0
	}
	current := g.positionHistory[len(g.positionHistory)-1]
	count := 0
	for _, h := range g.positionHistory {
		if h == current {
			count++
		}
	}
	return count
}

// isInsufficientMaterial reports whether neither side can possibly checkmate: K vs K,
// KB vs K, KN vs K, and KB vs KB with both bishops on the same color complex.
func (g Game) isInsufficientMaterial() bool {
	// Any pawn, rook or queen is (potentially) sufficient material.
	if g.bb[ColorBlack][PiecePawn]|g.bb[ColorWhite][PiecePawn]|
		g.bb[ColorBlack][PieceRook]|g.bb[ColorWhite][PieceRook]|
		g.bb[ColorBlack][PieceQueen]|g.bb[ColorWhite][PieceQueen] != 0 {
		return false
	}
	bishops := g.bb[ColorBlack][PieceBishop] | g.bb[ColorWhite][PieceBishop]
	knights := g.bb[ColorBlack][PieceKnight] | g.bb[ColorWhite][PieceKnight]
	minorCount := bits.OnesCount64(bishops | knights)
	if minorCount <= 1 {
		return true // K vs K, KB vs K, KN vs K
	}
	// KB vs KB with same-color bishops: only bishops remain, one per side, on the
	// same color complex.
	const lightSquares = 0x55AA55AA55AA55AA
	if knights == 0 &&
		bits.OnesCount64(g.bb[ColorBlack][PieceBishop]) == 1 &&
		bits.OnesCount64(g.bb[ColorWhite][PieceBishop]) == 1 &&
		(bishops&lightSquares == bishops || bishops&lightSquares == 0) {
		return true
	}
	return false
}

func opponent(c color) color {
	if c == ColorBlack {
		return ColorWhite
	}
	return ColorBlack
}

func isInBounds(xy XY) bool {
	return xy.X >= 0 && xy.Y >= 0 && xy.X <= 7 && xy.Y <= 7
}

func (p Piece) isInBounds(xy XY) bool {
	if !isInBounds(xy) {
		return false
	}
	if p.PieceType == PiecePawn && ((p.Owner == ColorWhite && p.XY.Y == 7) || (p.Owner == ColorBlack && p.XY.Y == 0)) {
		return false
	}
	return true
}

func (g Game) anySqThreatened(sqs [3]int8, owner color) bool {
	for _, sq := range sqs {
		if g.attackersOf(int(sq), owner) != 0 {
			return true
		}
	}
	return false
}

// attackersOf returns the bitboard of opponent pieces attacking the given square.
func (g Game) attackersOf(sq int, owner color) uint64 {
	opp := opponent(owner)
	occ := g.occAll()
	return knightAttacks[sq]&g.bb[opp][PieceKnight] |
		rookAttacks(sq, occ)&(g.bb[opp][PieceRook]|g.bb[opp][PieceQueen]) |
		bishopAttacks(sq, occ)&(g.bb[opp][PieceBishop]|g.bb[opp][PieceQueen]) |
		kingAttacks[sq]&g.bb[opp][PieceKing] |
		pawnCaptureAttacks[owner][sq]&g.bb[opp][PiecePawn]
}

func (g Game) xyThreatenedBy(sq XY, owner color, checkAllThreats bool) []Piece {
	pieces := []Piece{}
	for attackers := g.attackersOf(sqOf(sq), owner); attackers != 0; attackers &= attackers - 1 {
		pieces = append(pieces, g.pieceAtSq(bits.TrailingZeros64(attackers)))
		if !checkAllThreats {
			return pieces
		}
	}
	return pieces
}

func (p Piece) threatenedBy(g Game) []Piece {
	return g.xyThreatenedBy(p.XY, p.Owner, true /* checkAllThreats */)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
