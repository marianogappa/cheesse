package api

import "strings"

func (a API) parseGame(g InputGame) (game, error) {
	var (
		parsedGame game
		err        error
	)
	switch {
	case g.FENString != "":
		parsedGame, err = newGameFromFEN(g.FENString)
	case len(g.Board.Board) > 0:
		parsedGame, err = newGameFromBoard(mapBoardToInternalBoard(g.Board))
	default:
		err = errInvalidInputGame
	}
	if err != nil {
		return game{}, err
	}
	return parsedGame, nil
}

func (a API) parseAction(ia InputAction, g game) (action, error) {
	// TODO eventually accept other forms of action input
	fromXY, err := a.algebraicToXY(strings.ToLower(ia.FromSquare))
	if err != nil {
		return action{}, err
	}
	toXY, err := a.algebraicToXY(strings.ToLower(ia.ToSquare))
	if err != nil {
		return action{}, err
	}
	promotePieceType, err := a.stringToPieceType(ia.PromotePieceType)
	if err != nil {
		return action{}, err
	}

	for _, action := range g.actions {
		if action.fromPiece.xy != fromXY || action.toXY != toXY || (action.isPromotion && action.promotionPieceType != promotePieceType) {
			continue
		}
		return action, nil
	}

	return action{}, errInvalidActionForGivenGame
}

func (a API) algebraicToXY(sq string) (xy, error) {
	if len(sq) != 2 || sq[0] < 'a' || sq[0] > 'h' || sq[1] < '1' || sq[1] > '8' {
		return xy{}, errAlgebraicSquareInvalidOrOutOfBounds
	}
	return xy{int(sq[0] - 'a'), int('8' - sq[1])}, nil
}

func (a API) stringToPieceType(s string) (pieceType, error) {
	m := map[string]pieceType{
		"Queen":  pieceQueen,
		"King":   pieceKing,
		"Bishop": pieceBishop,
		"Knight": pieceKnight,
		"Rook":   pieceRook,
		"Pawn":   piecePawn,
		"":       pieceNone,
	}
	pt, ok := m[s]
	if !ok {
		return pieceNone, errInvalidPieceTypeName
	}
	return pt, nil
}
