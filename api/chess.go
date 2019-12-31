package api

// updateBoardLayout updates a game's layout-only (i.e. pieces and kings) after a given action, so that the resulting
// layout can be checked for checks, checkmates, etc.  This method is meant to be used as a dry-run to decide if the
// action can actually be executed. Should only be called by piece.buildAction and game.doAction.
//
// Note that this method assumes things like:
// - The destination xy is within bounds.
// - The action does not uncover an opponent's check, checkmate, etc.
// - The destination xy doesn't put fromPiece on top of a friendly piece, or does not jump another piece when invalid.
func (g game) updateBoardLayout(a action) game {
	clonedGame := g.clone()
	opponent := opponent(a.fromPiece.owner)

	// Special case for resignation action, because it doesn't require board changes
	if a.isResign {
		return clonedGame
	}

	// Update fromPiece's properties to reflect action
	fromPiece := g.pieces[a.fromPiece.owner][a.fromPiece.xy]
	fromPiece.xy = a.toXY
	if a.isPromotion {
		fromPiece.pieceType = a.promotionPieceType
	}

	// Remove pieces at {from, to} locations
	delete(clonedGame.pieces[a.fromPiece.owner], a.fromPiece.xy)
	delete(clonedGame.pieces[opponent], a.fromPiece.xy)
	// Place fromPiece at "to" location
	clonedGame.pieces[a.fromPiece.owner][fromPiece.xy] = fromPiece

	// Extra movements and deletions in the case of en passant capture and castling
	switch {
	case a.isCapture && g.isLastMoveEnPassant && g.enPassantTargetSquare.eq(a.toXY) && a.fromPiece.owner == colorBlack:
		delete(clonedGame.pieces[opponent], a.toXY.add(xy{0, -1}))
	case a.isCapture && g.isLastMoveEnPassant && g.enPassantTargetSquare.eq(a.toXY) && a.fromPiece.owner == colorWhite:
		delete(clonedGame.pieces[opponent], a.toXY.add(xy{0, 1}))
	case a.isQueensideCastle && a.fromPiece.owner == colorBlack:
		newRook := clonedGame.pieces[a.fromPiece.owner][xy{0, 0}]
		newRook.xy = xy{3, 0}
		clonedGame.pieces[a.fromPiece.owner][xy{3, 0}] = newRook
		delete(clonedGame.pieces[a.fromPiece.owner], xy{0, 0})
	case a.isQueensideCastle && a.fromPiece.owner == colorWhite:
		newRook := clonedGame.pieces[a.fromPiece.owner][xy{0, 7}]
		newRook.xy = xy{3, 7}
		clonedGame.pieces[a.fromPiece.owner][xy{3, 7}] = newRook
		delete(clonedGame.pieces[a.fromPiece.owner], xy{0, 7})
	case a.isKingsideCastle && a.fromPiece.owner == colorBlack:
		newRook := clonedGame.pieces[a.fromPiece.owner][xy{7, 0}]
		newRook.xy = xy{5, 0}
		clonedGame.pieces[a.fromPiece.owner][xy{5, 0}] = newRook
		delete(clonedGame.pieces[a.fromPiece.owner], xy{7, 0})
	case a.isKingsideCastle && a.fromPiece.owner == colorWhite:
		newRook := clonedGame.pieces[a.fromPiece.owner][xy{7, 7}]
		newRook.xy = xy{5, 7}
		clonedGame.pieces[a.fromPiece.owner][xy{5, 7}] = newRook
		delete(clonedGame.pieces[a.fromPiece.owner], xy{7, 7})
	}

	// If it's a king, game.kings needs to be updated
	if fromPiece.pieceType == pieceKing {
		clonedGame.kings[a.fromPiece.owner] = fromPiece
	}

	return clonedGame
}

func (g game) calculateAllActions() []action {
	if g.isGameOver {
		return []action{}
	}
	actions := []action{}
	// TODO these can be checked in parallel
	for _, piece := range g.pieces[g.turn()] {
		actions = append(actions, piece.calculateAllActions(g)...)
	}
	actions = append(actions, action{fromPiece: piece{owner: g.turn()}, isResign: true})
	return actions
}

func (g game) turn() color {
	if g.moveNumber%2 == 0 {
		return colorWhite
	}
	return colorBlack
}

func (p piece) calculateAllActions(g game) []action {
	if g.isGameOver || g.turn() != p.owner || p.pieceType == pieceNone || g.pieces[p.owner][p.xy] != p {
		return []action{}
	}

	// Naive deltas (special case for Pawn)
	deltas := movementDeltasByPieceType[p.pieceType]
	switch {
	case p.pieceType == piecePawn && p.owner == colorBlack:
		deltas = []xy{{0, 1}, {0, 2}, {-1, 1}, {1, 1}} // Includes en passant and captures
	case p.pieceType == piecePawn && p.owner == colorWhite:
		deltas = []xy{{0, -1}, {0, -2}, {-1, -1}, {1, -1}} // Includes en passant and captures
	}

	// TODO These can be checked in parallel
	actions := []action{}
	for _, delta := range deltas {
		toXY := p.xy.add(delta)

		// This for-loop covers the case of Bishop, Rook and Queen being able to move multiple times on a given delta
		for {
			// If this action is a promotion, then 4 possible actions should be created, one for each promotion piece
			promotionPieces := []pieceType{pieceNone}
			if p.pieceType == piecePawn && ((p.owner == colorBlack && p.xy.add(delta).y == 7) || (p.owner == colorWhite && p.xy.add(delta).y == 0)) {
				promotionPieces = []pieceType{pieceQueen, pieceBishop, pieceKnight, pieceRook}
			}

			var (
				a   action
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
			if (p.pieceType != pieceQueen && p.pieceType != pieceBishop && p.pieceType != pieceRook) ||
				!isInBounds(toXY) ||
				a.isCapture ||
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
func (p piece) buildAction(toXY xy, g game, promotionPieceType pieceType) (action, error) {
	if !p.isInBounds(toXY) {
		return action{}, errNotInBounds
	}

	// There's a friendly piece in the destination position
	if _, ok := g.pieces[p.owner][toXY]; ok {
		return action{}, errFriendlyPieceInDestination
	}

	// There's a friendly/opponent piece between piece.xy and toXY
	for _, xyBetween := range p.xysTowards(toXY) {
		if g.pieces[p.owner][xyBetween].pieceType != pieceNone || g.pieces[opponent(p.owner)][xyBetween].pieceType != pieceNone {
			return action{}, errPieceInBetween
		}
	}

	a := action{fromPiece: p, toXY: toXY}
	opponentPieceAtToXY, hasOpponentPieceAtToXY := g.pieces[opponent(p.owner)][toXY]

	// Edge case: Pawn is the only piece that cannot capture while moving forwards
	// If any piece (i.e. owner by any player) is either in the destination or in the middle, the action is invalid
	if p.pieceType == piecePawn && toXY.x == p.xy.x {
		// Bail if any constraint is not met
		switch {
		case !g.isEmptyAt(toXY),
			abs(toXY.y-p.xy.y) == 2 && p.owner == colorBlack && !g.isEmptyAt(xy{x: toXY.x, y: toXY.y - 1}),
			abs(toXY.y-p.xy.y) == 2 && p.owner == colorWhite && !g.isEmptyAt(xy{x: toXY.x, y: toXY.y + 1}),
			abs(toXY.y-p.xy.y) == 2 && p.owner == colorBlack && p.xy.y != 1,
			abs(toXY.y-p.xy.y) == 2 && p.owner == colorWhite && p.xy.y != 6:
			return action{}, errPieceBlockingPawn
		}
	}

	// Edge case: Pawn can only move diagonally if there's an opponent piece in that position, or if it's en passant
	if p.pieceType == piecePawn && toXY.x != p.xy.x {
		// Bail if any constraint is not met
		switch {
		case abs(toXY.x-p.xy.x) != 1,
			abs(toXY.y-p.xy.y) != 1,
			!hasOpponentPieceAtToXY && (!g.isLastMoveEnPassant || !g.enPassantTargetSquare.eq(toXY)):
			return action{}, errPawnCantCapture
		}
	}

	// Set the isEnPassant flag
	if p.pieceType == piecePawn && ((p.owner == colorBlack && p.xy.y == 1 && toXY.y == 3) || (p.owner == colorWhite && p.xy.y == 6 && toXY.y == 4)) {
		a.isEnPassant = true
	}

	// Castling context
	if p.pieceType == pieceKing && abs(toXY.x-p.xy.x) > 1 { // It's a castle attempt
		// Set castle type context
		var castleType castleType = castleTypeQueenside
		if toXY.x == 6 {
			castleType = castleTypeKingside
		}

		// Bail if any constraint is not met
		switch {
		case p.owner == colorBlack && !g.canBlackCastle,
			p.owner == colorWhite && !g.canWhiteCastle,
			p.owner == colorBlack && castleType == castleTypeQueenside && !g.canBlackQueensideCastle,
			p.owner == colorBlack && castleType == castleTypeKingside && !g.canBlackKingsideCastle,
			p.owner == colorWhite && castleType == castleTypeQueenside && !g.canWhiteQueensideCastle,
			p.owner == colorWhite && castleType == castleTypeKingside && !g.canWhiteKingsideCastle,
			p.owner == colorBlack && (p.xy.y != 0 || toXY.y != 0 || p.xy.x != 4 || (toXY.x != 6 && toXY.x != 2)),
			p.owner == colorWhite && (p.xy.y != 7 || toXY.y != 7 || p.xy.x != 4 || (toXY.x != 6 && toXY.x != 2)),
			!g.isEmptyAtAllOf(emptyXYsForCastlingByColorAndCastleType[p.owner][castleType]),
			g.isAnyXYThreatened(unthreatenedXYsForCastlingByColorAndCastleType[p.owner][castleType], p.owner):
			return action{}, errCantCastle
		}

		// Set castling context
		a.isCastle = true
		switch {
		case p.owner == colorBlack && p.xy.y == 0 && toXY.x == 2:
			a.isQueensideCastle = true
		case p.owner == colorWhite && p.xy.y == 7 && toXY.x == 2:
			a.isQueensideCastle = true
		case p.owner == colorBlack && p.xy.y == 0 && toXY.x == 6:
			a.isKingsideCastle = true
		case p.owner == colorWhite && p.xy.y == 7 && toXY.x == 6:
			a.isKingsideCastle = true
		}
	}

	// Set capture context
	if hasOpponentPieceAtToXY {
		a.isCapture = true
		a.capturedPiece = opponentPieceAtToXY
	}
	// Set capture context in the case of an en passant capture
	if p.pieceType == piecePawn && g.isLastMoveEnPassant && g.enPassantTargetSquare.eq(toXY) {
		a.isCapture = true
		switch p.owner {
		case colorBlack:
			a.capturedPiece = g.pieces[opponent(p.owner)][xy{x: toXY.x, y: toXY.y - 1}]
		case colorWhite:
			a.capturedPiece = g.pieces[opponent(p.owner)][xy{x: toXY.x, y: toXY.y + 1}]
		}
	}

	// Set promotion context
	if p.pieceType == piecePawn && ((p.owner == colorBlack && toXY.y == 7) || (p.owner == colorWhite && toXY.y == 0)) {
		if promotionPieceType == piecePawn || promotionPieceType == pieceKing || promotionPieceType == pieceNone {
			return action{}, errCantPromote
		}
		a.isPromotion = true
		a.promotionPieceType = promotionPieceType
	}

	newGame := g.updateBoardLayout(a)

	// check if moving puts the owner's King in check
	if len(newGame.kings[p.owner].threatenedBy(newGame)) > 0 { // N.B. this is an expensive operation!
		return action{}, errActionLeavesKingThreatened
	}

	return a, nil
}

func (g game) isEmptyAtAllOf(xys []xy) bool {
	for _, xy := range xys {
		if !g.isEmptyAt(xy) {
			return false
		}
	}
	return true
}

func (g game) isEmptyAt(xy xy) bool {
	_, blackPieceExists := g.pieces[colorBlack][xy]
	_, whitePieceExists := g.pieces[colorWhite][xy]
	return !blackPieceExists && !whitePieceExists
}

// doAction executes the given action on the given game.
// It assumes that the game is in a state where this action can be executed.
// It assumes that the action is fully-correctly created and it's valid.
// It fully updates the game context.
// This is an expensive method (due to having to check for check and checkmate), so use only if needed.
func (g game) doAction(a action) game {
	newGame := g.updateBoardLayout(a)
	lastTurn := g.turn()

	// Special case for resignation action
	if a.isResign {
		newGame.isGameOver = true
		newGame.gameOverWinner = opponent(lastTurn)

		// TODO is it necessary to update other things?
		return newGame
	}

	// Castling context update
	switch {
	case lastTurn == colorBlack && (a.isCastle || a.fromPiece.pieceType == pieceKing):
		newGame.canBlackCastle = false
		newGame.canBlackQueensideCastle = false
		newGame.canBlackKingsideCastle = false
	case lastTurn == colorWhite && (a.isCastle || a.fromPiece.pieceType == pieceKing):
		newGame.canWhiteCastle = false
		newGame.canWhiteQueensideCastle = false
		newGame.canWhiteKingsideCastle = false
	case lastTurn == colorBlack && a.fromPiece.pieceType == pieceRook && a.fromPiece.xy == xy{x: 0, y: 0}:
		newGame.canBlackQueensideCastle = false
		newGame.canBlackCastle = newGame.canBlackKingsideCastle
	case lastTurn == colorBlack && a.fromPiece.pieceType == pieceRook && a.fromPiece.xy == xy{x: 7, y: 0}:
		newGame.canBlackKingsideCastle = false
		newGame.canBlackCastle = newGame.canBlackQueensideCastle
	case lastTurn == colorWhite && a.fromPiece.pieceType == pieceRook && a.fromPiece.xy == xy{x: 0, y: 7}:
		newGame.canWhiteQueensideCastle = false
		newGame.canWhiteCastle = newGame.canWhiteKingsideCastle
	case lastTurn == colorWhite && a.fromPiece.pieceType == pieceRook && a.fromPiece.xy == xy{x: 7, y: 7}:
		newGame.canWhiteKingsideCastle = false
		newGame.canWhiteCastle = newGame.canWhiteQueensideCastle
	}

	newGame.moveNumber = g.moveNumber + 1
	if lastTurn == colorBlack {
		newGame.fullMoveNumber = g.fullMoveNumber + 1
	}
	newGame.isLastMoveEnPassant = a.isEnPassant
	if a.isEnPassant && lastTurn == colorBlack {
		newGame.enPassantTargetSquare = xy{x: a.toXY.x, y: a.toXY.y + 1}
	}
	if a.isEnPassant && lastTurn == colorWhite {
		newGame.enPassantTargetSquare = xy{x: a.toXY.x, y: a.toXY.y - 1}
	}

	newGame.halfMoveClock = g.halfMoveClock + 1
	if a.isCapture || a.fromPiece.pieceType == piecePawn {
		newGame.halfMoveClock = 0
	}

	return newGame.calculateCriticalFlags()
}

func (g game) calculateCriticalFlags() game {
	turn := g.turn()

	g.isCheck = false
	g.isCheckmate = false
	g.isStalemate = false
	g.isDraw = false
	g.isGameOver = false
	g.gameOverWinner = -1
	g.inCheckBy = []piece{}

	g.inCheckBy = g.kings[turn].threatenedBy(g) // This is expensive!
	if len(g.inCheckBy) > 0 {
		g.isCheck = true
	}

	g.actions = g.calculateAllActions() // This is incredibly expensive!
	if len(g.actions) == 1 && g.actions[0].isResign {
		g.isCheckmate = g.isCheck
		g.isStalemate = !g.isCheck
	}

	if g.halfMoveClock == 100 {
		g.isDraw = true // N.B. this forces draw rather than allowing a claim for draw. Is that ok?
	}
	if g.isCheckmate || g.isStalemate || g.isDraw {
		g.isGameOver = true
	}
	if g.isCheckmate {
		g.gameOverWinner = opponent(turn)
	}

	return g
}

func opponent(c color) color {
	if c == colorBlack {
		return colorWhite
	}
	return colorBlack
}

func isInBounds(xy xy) bool {
	return xy.x >= 0 && xy.y >= 0 && xy.x <= 7 && xy.y <= 7
}

func (p piece) isInBounds(xy xy) bool {
	if !isInBounds(xy) {
		return false
	}
	if p.pieceType == piecePawn && ((p.owner == colorWhite && p.xy.y == 7) || (p.owner == colorBlack && p.xy.y == 0)) {
		return false
	}
	return true
}

func (g game) isAnyXYThreatened(xys []xy, owner color) bool {
	for _, xy := range xys {
		if len(g.xyThreatenedBy(xy, owner, false /* checkAllThreats */)) > 0 {
			return true
		}
	}
	return false
}

func (g game) xyThreatenedBy(sq xy, owner color, checkAllThreats bool) []piece {
	pieces := []piece{}
	opponent := opponent(owner)

	// Knights
	for _, delta := range movementDeltasByPieceType[pieceKnight] {
		xy := sq.add(delta)
		if !isInBounds(xy) {
			continue
		}
		if p, ok := g.pieces[opponent][xy]; !ok || p.pieceType != pieceKnight {
			continue
		}
		pieces = append(pieces, g.pieces[opponent][xy])
		if !checkAllThreats {
			return pieces
		}
	}

	// Vertical and horizontal (Rook & Queen)
	for _, delta := range movementDeltasByPieceType[pieceRook] {
		for sq := sq.add(delta); isInBounds(sq); sq = sq.add(delta) {
			// There's a friendly piece in this delta; it will block further pieces
			if _, ok := g.pieces[owner][sq]; ok {
				break
			}
			// There's an opponent piece in this delta but it's not Rook/Queen; it will block further pieces
			if p, ok := g.pieces[opponent][sq]; ok && p.pieceType != pieceQueen && p.pieceType != pieceRook {
				break
			}
			if p, ok := g.pieces[opponent][sq]; ok && p.pieceType == pieceQueen || p.pieceType == pieceRook {
				pieces = append(pieces, p)
				if !checkAllThreats {
					return pieces
				}
				break
			}
		}
	}

	// Diagonal (Bishop & Queen)
	for _, delta := range movementDeltasByPieceType[pieceBishop] {
		for sq := sq.add(delta); isInBounds(sq); sq = sq.add(delta) {
			// There's a friendly piece in this delta; it will block further pieces
			if _, ok := g.pieces[owner][sq]; ok {
				break
			}
			// There's an opponent piece in this delta but it's not Bishop/Queen; it will block further pieces
			if p, ok := g.pieces[opponent][sq]; ok && p.pieceType != pieceQueen && p.pieceType != pieceBishop {
				break
			}
			if p, ok := g.pieces[opponent][sq]; ok && p.pieceType == pieceQueen || p.pieceType == pieceBishop {
				pieces = append(pieces, p)
				if !checkAllThreats {
					return pieces
				}
				break
			}
		}
	}

	// King
	if abs(sq.x-g.kings[opponent].xy.x) <= 1 && abs(sq.y-g.kings[opponent].xy.y) <= 1 {
		pieces = append(pieces, g.kings[opponent])
		if !checkAllThreats {
			return pieces
		}
	}

	// Pawns
	pawnXYs := []xy{sq.add(xy{-1, -1}), sq.add(xy{1, -1})}
	if owner == colorWhite {
		pawnXYs = []xy{sq.add(xy{-1, 1}), sq.add(xy{1, 1})}
	}
	if piece, ok := g.pieces[opponent][pawnXYs[0]]; ok && piece.pieceType == piecePawn {
		pieces = append(pieces, piece)
		if !checkAllThreats {
			return pieces
		}
	}
	if piece, ok := g.pieces[opponent][pawnXYs[1]]; ok && piece.pieceType == piecePawn {
		pieces = append(pieces, piece)
		if !checkAllThreats {
			return pieces
		}
	}

	return pieces
}

func (p piece) threatenedBy(g game) []piece {
	return g.xyThreatenedBy(p.xy, p.owner, true /* checkAllThreats */)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
