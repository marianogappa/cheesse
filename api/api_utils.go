package api

import (
	"strings"

	"github.com/marianogappa/cheesse/core"
)

func (a API) parseGame(g InputGame) (core.Game, error) {
	var (
		parsedGame core.Game
		err        error
	)
	switch {
	case g.FENString != "":
		parsedGame, err = core.NewGameFromFEN(g.FENString)
	case len(g.Board.Board) > 0:
		parsedGame, err = core.NewGameFromBoard(mapBoardToInternalBoard(g.Board))
	default:
		var defaultGame, _ = core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		parsedGame = defaultGame
	}
	if err != nil {
		return core.Game{}, err
	}
	return parsedGame, nil
}

func (a API) parseAction(ia InputAction, g core.Game) (core.Action, error) {
	// TODO eventually accept other forms of action input
	fromXY, err := a.algebraicToXY(strings.ToLower(ia.FromSquare))
	if err != nil {
		return core.Action{}, err
	}
	toXY, err := a.algebraicToXY(strings.ToLower(ia.ToSquare))
	if err != nil {
		return core.Action{}, err
	}
	promotionPieceType, err := a.stringToPieceType(ia.PromotionPieceType)
	if err != nil {
		return core.Action{}, err
	}

	for _, action := range g.Actions {
		if action.FromPiece.XY != fromXY || action.ToXY != toXY || (action.IsPromotion && action.PromotionPieceType != promotionPieceType) {
			continue
		}
		return action, nil
	}

	return core.Action{}, errInvalidActionForGivenGame
}

func (a API) algebraicToXY(sq string) (core.XY, error) {
	if len(sq) != 2 || sq[0] < 'a' || sq[0] > 'h' || sq[1] < '1' || sq[1] > '8' {
		return core.XY{}, errAlgebraicSquareInvalidOrOutOfBounds
	}
	return core.XY{X: int(sq[0] - 'a'), Y: int('8' - sq[1])}, nil
}

func (a API) stringToPieceType(s string) (core.PieceType, error) {
	m := map[string]core.PieceType{
		"Queen":  core.PieceQueen,
		"King":   core.PieceKing,
		"Bishop": core.PieceBishop,
		"Knight": core.PieceKnight,
		"Rook":   core.PieceRook,
		"Pawn":   core.PiecePawn,
		"":       core.PieceNone,
	}
	pt, ok := m[s]
	if !ok {
		return core.PieceNone, errInvalidPieceTypeName
	}
	return pt, nil
}
