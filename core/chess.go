package core

// updateBoardLayout updates a game's layout-only (i.e. pieces and kings) after a given action, so that the resulting
// layout can be checked for checks, checkmates, etc.  This method is meant to be used as a dry-run to decide if the
// action can actually be executed. Should only be called by piece.buildAction and game.doAction.
//
// Note that this method assumes things like:
// - The destination xy is within bounds.
// - The action does not uncover an opponent's check, checkmate, etc.
// - The destination xy doesn't put fromPiece on top of a friendly piece, or does not jump another piece when invalid.
func (g Game) updateBoardLayout(a Action) Game {
	clonedGame := g.Clone()
	opponent := opponent(a.FromPiece.Owner)

	// Special case for resignation action, because it doesn't require board changes
	if a.IsResign {
		return clonedGame
	}

	// Update fromPiece's properties to reflect action
	fromPiece := g.Pieces[a.FromPiece.Owner][a.FromPiece.XY]
	fromPiece.XY = a.ToXY
	if a.IsPromotion {
		fromPiece.PieceType = a.PromotionPieceType
	}

	// Remove pieces at {from, to} locations
	delete(clonedGame.Pieces[a.FromPiece.Owner], a.FromPiece.XY)
	delete(clonedGame.Pieces[opponent], fromPiece.XY)
	// Place fromPiece at "to" location
	clonedGame.Pieces[a.FromPiece.Owner][fromPiece.XY] = fromPiece

	// Extra movements and deletions in the case of en passant capture and castling
	switch {
	case a.IsCapture && g.IsLastMoveEnPassant && g.EnPassantTargetSquare.eq(a.ToXY) && a.FromPiece.Owner == ColorBlack:
		delete(clonedGame.Pieces[opponent], a.ToXY.add(XY{0, -1}))
	case a.IsCapture && g.IsLastMoveEnPassant && g.EnPassantTargetSquare.eq(a.ToXY) && a.FromPiece.Owner == ColorWhite:
		delete(clonedGame.Pieces[opponent], a.ToXY.add(XY{0, 1}))
	case a.IsQueensideCastle && a.FromPiece.Owner == ColorBlack:
		newRook := clonedGame.Pieces[a.FromPiece.Owner][XY{0, 0}]
		newRook.XY = XY{3, 0}
		clonedGame.Pieces[a.FromPiece.Owner][XY{3, 0}] = newRook
		delete(clonedGame.Pieces[a.FromPiece.Owner], XY{0, 0})
	case a.IsQueensideCastle && a.FromPiece.Owner == ColorWhite:
		newRook := clonedGame.Pieces[a.FromPiece.Owner][XY{0, 7}]
		newRook.XY = XY{3, 7}
		clonedGame.Pieces[a.FromPiece.Owner][XY{3, 7}] = newRook
		delete(clonedGame.Pieces[a.FromPiece.Owner], XY{0, 7})
	case a.IsKingsideCastle && a.FromPiece.Owner == ColorBlack:
		newRook := clonedGame.Pieces[a.FromPiece.Owner][XY{7, 0}]
		newRook.XY = XY{5, 0}
		clonedGame.Pieces[a.FromPiece.Owner][XY{5, 0}] = newRook
		delete(clonedGame.Pieces[a.FromPiece.Owner], XY{7, 0})
	case a.IsKingsideCastle && a.FromPiece.Owner == ColorWhite:
		newRook := clonedGame.Pieces[a.FromPiece.Owner][XY{7, 7}]
		newRook.XY = XY{5, 7}
		clonedGame.Pieces[a.FromPiece.Owner][XY{5, 7}] = newRook
		delete(clonedGame.Pieces[a.FromPiece.Owner], XY{7, 7})
	}

	// If it's a king, game.kings needs to be updated
	if fromPiece.PieceType == PieceKing {
		clonedGame.Kings[a.FromPiece.Owner] = fromPiece
	}

	return clonedGame
}

func (g Game) calculateAllActions() []Action {
	if g.IsGameOver {
		return []Action{}
	}
	actions := []Action{}
	// TODO these can be checked in parallel
	for _, piece := range g.Pieces[g.Turn()] {
		actions = append(actions, piece.calculateAllActions(g)...)
	}
	actions = append(actions, Action{FromPiece: Piece{Owner: g.Turn()}, IsResign: true})
	return actions
}

func (g Game) Turn() color {
	if g.MoveNumber%2 == 0 {
		return ColorWhite
	}
	return ColorBlack
}

func (p Piece) calculateAllActions(g Game) []Action {
	if g.IsGameOver || g.Turn() != p.Owner || p.PieceType == PieceNone || g.Pieces[p.Owner][p.XY] != p {
		return []Action{}
	}

	// Naive deltas (special case for Pawn)
	deltas := movementDeltasByPieceType[p.PieceType]
	switch {
	case p.PieceType == PiecePawn && p.Owner == ColorBlack:
		deltas = []XY{{0, 1}, {0, 2}, {-1, 1}, {1, 1}} // Includes en passant and captures
	case p.PieceType == PiecePawn && p.Owner == ColorWhite:
		deltas = []XY{{0, -1}, {0, -2}, {-1, -1}, {1, -1}} // Includes en passant and captures
	}

	// TODO These can be checked in parallel
	actions := []Action{}
	for _, delta := range deltas {
		toXY := p.XY.add(delta)

		// This for-loop covers the case of Bishop, Rook and Queen being able to move multiple times on a given delta
		for {
			// If this action is a promotion, then 4 possible actions should be created, one for each promotion piece
			promotionPieces := []PieceType{PieceNone}
			if p.PieceType == PiecePawn && ((p.Owner == ColorBlack && p.XY.add(delta).Y == 7) || (p.Owner == ColorWhite && p.XY.add(delta).Y == 0)) {
				promotionPieces = []PieceType{PieceQueen, PieceBishop, PieceKnight, PieceRook}
			}

			var (
				a   Action
				err error
			)
			// This for-loop covers the case of up to 4 actions created due to promotion pieces
			for _, promotionPiece := range promotionPieces {
				a, err = p.buildAction(toXY, g, promotionPiece)
				if err != nil {
					break
				}
				actions = append(actions, a)
			}

			toXY = toXY.add(delta)
			if (p.PieceType != PieceQueen && p.PieceType != PieceBishop && p.PieceType != PieceRook) ||
				!isInBounds(toXY) ||
				a.IsCapture ||
				err == errFriendlyPieceInDestination ||
				err == errPieceInBetween {
				break
			}
		}
	}
	return actions
}

// buildAction tries to create an action given a piece, a game, a destination xy and optionally a promotion piece type.
// It will do an exhaustive check of validity, including en passant, en passant capture, castling, promotions, etc.
//
// The only thing that is mostly assumed valid is that the piece's destination is reachable. This is because the
// piece's deltas and the eventuality of "jumping" over pieces are both calculated by piece.calculateAllActions. For
// this reason, this method should only be called internally by piece.calculateAllActions.
func (p Piece) buildAction(toXY XY, g Game, promotionPieceType PieceType) (Action, error) {
	if !p.isInBounds(toXY) {
		return Action{}, errNotInBounds
	}

	// There's a friendly piece in the destination position
	if _, ok := g.Pieces[p.Owner][toXY]; ok {
		return Action{}, errFriendlyPieceInDestination
	}

	// There's a friendly/opponent piece between piece.xy and toXY
	for _, xyBetween := range p.xysTowards(toXY) {
		if g.Pieces[p.Owner][xyBetween].PieceType != PieceNone || g.Pieces[opponent(p.Owner)][xyBetween].PieceType != PieceNone {
			return Action{}, errPieceInBetween
		}
	}

	a := Action{FromPiece: p, ToXY: toXY}
	opponentPieceAtToXY, hasOpponentPieceAtToXY := g.Pieces[opponent(p.Owner)][toXY]

	// Edge case: Pawn is the only piece that cannot capture while moving forwards
	// If any piece (i.e. owner by any player) is either in the destination or in the middle, the action is invalid
	if p.PieceType == PiecePawn && toXY.X == p.XY.X {
		// Bail if any constraint is not met
		switch {
		case !g.isEmptyAt(toXY),
			abs(toXY.Y-p.XY.Y) == 2 && p.Owner == ColorBlack && !g.isEmptyAt(XY{X: toXY.X, Y: toXY.Y - 1}),
			abs(toXY.Y-p.XY.Y) == 2 && p.Owner == ColorWhite && !g.isEmptyAt(XY{X: toXY.X, Y: toXY.Y + 1}),
			abs(toXY.Y-p.XY.Y) == 2 && p.Owner == ColorBlack && p.XY.Y != 1,
			abs(toXY.Y-p.XY.Y) == 2 && p.Owner == ColorWhite && p.XY.Y != 6:
			return Action{}, errPieceBlockingPawn
		}
	}

	// Edge case: Pawn can only move diagonally if there's an opponent piece in that position, or if it's en passant
	if p.PieceType == PiecePawn && toXY.X != p.XY.X {
		// Bail if any constraint is not met
		switch {
		case abs(toXY.X-p.XY.X) != 1,
			abs(toXY.Y-p.XY.Y) != 1,
			!hasOpponentPieceAtToXY && (!g.IsLastMoveEnPassant || !g.EnPassantTargetSquare.eq(toXY)):
			return Action{}, errPawnCantCapture
		}
	}

	// Set the isEnPassant flag
	if p.PieceType == PiecePawn && ((p.Owner == ColorBlack && p.XY.Y == 1 && toXY.Y == 3) || (p.Owner == ColorWhite && p.XY.Y == 6 && toXY.Y == 4)) {
		a.IsEnPassant = true
	}

	// Castling context
	if p.PieceType == PieceKing && abs(toXY.X-p.XY.X) > 1 { // It's a castle attempt
		// Set castle type context
		var castleType castleType = castleTypeQueenside
		if toXY.X == 6 {
			castleType = castleTypeKingside
		}

		// Bail if any constraint is not met
		switch {
		case p.Owner == ColorBlack && !g.CanBlackCastle,
			p.Owner == ColorWhite && !g.CanWhiteCastle,
			p.Owner == ColorBlack && castleType == castleTypeQueenside && !g.CanBlackQueensideCastle,
			p.Owner == ColorBlack && castleType == castleTypeKingside && !g.CanBlackKingsideCastle,
			p.Owner == ColorWhite && castleType == castleTypeQueenside && !g.CanWhiteQueensideCastle,
			p.Owner == ColorWhite && castleType == castleTypeKingside && !g.CanWhiteKingsideCastle,
			p.Owner == ColorBlack && (p.XY.Y != 0 || toXY.Y != 0 || p.XY.X != 4 || (toXY.X != 6 && toXY.X != 2)),
			p.Owner == ColorWhite && (p.XY.Y != 7 || toXY.Y != 7 || p.XY.X != 4 || (toXY.X != 6 && toXY.X != 2)),
			!g.isEmptyAtAllOf(emptyXYsForCastlingByColorAndCastleType[p.Owner][castleType]),
			g.isAnyXYThreatened(unthreatenedXYsForCastlingByColorAndCastleType[p.Owner][castleType], p.Owner):
			return Action{}, errCantCastle
		}

		// Set castling context
		a.IsCastle = true
		switch {
		case p.Owner == ColorBlack && p.XY.Y == 0 && toXY.X == 2:
			a.IsQueensideCastle = true
		case p.Owner == ColorWhite && p.XY.Y == 7 && toXY.X == 2:
			a.IsQueensideCastle = true
		case p.Owner == ColorBlack && p.XY.Y == 0 && toXY.X == 6:
			a.IsKingsideCastle = true
		case p.Owner == ColorWhite && p.XY.Y == 7 && toXY.X == 6:
			a.IsKingsideCastle = true
		}
	}

	// Set capture context
	if hasOpponentPieceAtToXY {
		a.IsCapture = true
		a.CapturedPiece = opponentPieceAtToXY
	}
	// Set capture context in the case of an en passant capture
	if p.PieceType == PiecePawn && g.IsLastMoveEnPassant && g.EnPassantTargetSquare.eq(toXY) {
		a.IsCapture = true
		a.IsEnPassantCapture = true
		switch p.Owner {
		case ColorBlack:
			a.CapturedPiece = g.Pieces[opponent(p.Owner)][XY{X: toXY.X, Y: toXY.Y - 1}]
		case ColorWhite:
			a.CapturedPiece = g.Pieces[opponent(p.Owner)][XY{X: toXY.X, Y: toXY.Y + 1}]
		}
	}

	// Set promotion context
	if p.PieceType == PiecePawn && ((p.Owner == ColorBlack && toXY.Y == 7) || (p.Owner == ColorWhite && toXY.Y == 0)) {
		if promotionPieceType == PiecePawn || promotionPieceType == PieceKing || promotionPieceType == PieceNone {
			return Action{}, errCantPromote
		}
		a.IsPromotion = true
		a.PromotionPieceType = promotionPieceType
	}

	newGame := g.updateBoardLayout(a)

	// check if moving puts the owner's King in check
	if len(newGame.Kings[p.Owner].threatenedBy(newGame)) > 0 { // N.B. this is an expensive operation!
		return Action{}, errActionLeavesKingThreatened
	}

	return a, nil
}

func (g Game) isEmptyAtAllOf(xys []XY) bool {
	for _, xy := range xys {
		if !g.isEmptyAt(xy) {
			return false
		}
	}
	return true
}

func (g Game) isEmptyAt(xy XY) bool {
	_, blackPieceExists := g.Pieces[ColorBlack][xy]
	_, whitePieceExists := g.Pieces[ColorWhite][xy]
	return !blackPieceExists && !whitePieceExists
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

	// Castling context update
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

	newGame.MoveNumber = g.MoveNumber + 1
	if lastTurn == ColorBlack {
		newGame.FullMoveNumber = g.FullMoveNumber + 1
	}
	newGame.IsLastMoveEnPassant = a.IsEnPassant
	if a.IsEnPassant && lastTurn == ColorBlack {
		newGame.EnPassantTargetSquare = XY{X: a.ToXY.X, Y: a.ToXY.Y - 1}
	}
	if a.IsEnPassant && lastTurn == ColorWhite {
		newGame.EnPassantTargetSquare = XY{X: a.ToXY.X, Y: a.ToXY.Y + 1}
	}

	newGame.HalfMoveClock = g.HalfMoveClock + 1
	if a.IsCapture || a.FromPiece.PieceType == PiecePawn {
		newGame.HalfMoveClock = 0
	}

	return newGame.calculateCriticalFlags()
}

func (g Game) calculateCriticalFlags() Game {
	turn := g.Turn()

	g.IsCheck = false
	g.IsCheckmate = false
	g.IsStalemate = false
	g.IsDraw = false
	g.IsGameOver = false
	g.GameOverWinner = -1
	g.InCheckBy = []Piece{}

	g.InCheckBy = g.Kings[turn].threatenedBy(g) // This is expensive!
	if len(g.InCheckBy) > 0 {
		g.IsCheck = true
	}

	g.Actions = g.calculateAllActions() // This is incredibly expensive!
	if len(g.Actions) == 1 && g.Actions[0].IsResign {
		g.IsCheckmate = g.IsCheck
		g.IsStalemate = !g.IsCheck
	}

	if g.HalfMoveClock == 100 {
		g.IsDraw = true // N.B. this forces draw rather than allowing a claim for draw. Is that ok?
	}
	if g.IsCheckmate || g.IsStalemate || g.IsDraw {
		g.IsGameOver = true
	}
	if g.IsCheckmate {
		g.GameOverWinner = opponent(turn)
	}

	return g
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

func (g Game) isAnyXYThreatened(xys []XY, owner color) bool {
	for _, xy := range xys {
		if len(g.xyThreatenedBy(xy, owner, false /* checkAllThreats */)) > 0 {
			return true
		}
	}
	return false
}

func (g Game) xyThreatenedBy(sq XY, owner color, checkAllThreats bool) []Piece {
	pieces := []Piece{}
	opponent := opponent(owner)

	// Knights
	for _, delta := range movementDeltasByPieceType[PieceKnight] {
		xy := sq.add(delta)
		if !isInBounds(xy) {
			continue
		}
		if p, ok := g.Pieces[opponent][xy]; !ok || p.PieceType != PieceKnight {
			continue
		}
		pieces = append(pieces, g.Pieces[opponent][xy])
		if !checkAllThreats {
			return pieces
		}
	}

	// Vertical and horizontal (Rook & Queen)
	for _, delta := range movementDeltasByPieceType[PieceRook] {
		for sq := sq.add(delta); isInBounds(sq); sq = sq.add(delta) {
			// There's a friendly piece in this delta; it will block further pieces
			if _, ok := g.Pieces[owner][sq]; ok {
				break
			}
			// There's an opponent piece in this delta but it's not Rook/Queen; it will block further pieces
			if p, ok := g.Pieces[opponent][sq]; ok && p.PieceType != PieceQueen && p.PieceType != PieceRook {
				break
			}
			if p, ok := g.Pieces[opponent][sq]; ok && p.PieceType == PieceQueen || p.PieceType == PieceRook {
				pieces = append(pieces, p)
				if !checkAllThreats {
					return pieces
				}
				break
			}
		}
	}

	// Diagonal (Bishop & Queen)
	for _, delta := range movementDeltasByPieceType[PieceBishop] {
		for sq := sq.add(delta); isInBounds(sq); sq = sq.add(delta) {
			// There's a friendly piece in this delta; it will block further pieces
			if _, ok := g.Pieces[owner][sq]; ok {
				break
			}
			// There's an opponent piece in this delta but it's not Bishop/Queen; it will block further pieces
			if p, ok := g.Pieces[opponent][sq]; ok && p.PieceType != PieceQueen && p.PieceType != PieceBishop {
				break
			}
			if p, ok := g.Pieces[opponent][sq]; ok && p.PieceType == PieceQueen || p.PieceType == PieceBishop {
				pieces = append(pieces, p)
				if !checkAllThreats {
					return pieces
				}
				break
			}
		}
	}

	// King
	if abs(sq.X-g.Kings[opponent].XY.X) <= 1 && abs(sq.Y-g.Kings[opponent].XY.Y) <= 1 {
		pieces = append(pieces, g.Kings[opponent])
		if !checkAllThreats {
			return pieces
		}
	}

	// Pawns
	pawnXYs := []XY{sq.add(XY{-1, 1}), sq.add(XY{1, 1})}
	if owner == ColorWhite {
		pawnXYs = []XY{sq.add(XY{-1, -1}), sq.add(XY{1, -1})}
	}
	if piece, ok := g.Pieces[opponent][pawnXYs[0]]; ok && piece.PieceType == PiecePawn {
		pieces = append(pieces, piece)
		if !checkAllThreats {
			return pieces
		}
	}
	if piece, ok := g.Pieces[opponent][pawnXYs[1]]; ok && piece.PieceType == PiecePawn {
		pieces = append(pieces, piece)
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
