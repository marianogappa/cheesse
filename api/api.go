package api

import (
	"fmt"
)

type API struct{}

func New() API { return API{} }

var (
	errInvalidInputGame                    = fmt.Errorf("invalid input game: please supply a valid fenString or a board")
	errAlgebraicSquareInvalidOrOutOfBounds = fmt.Errorf("invalid algebraic square: empty or out of bounds")
	errInvalidPieceTypeName                = fmt.Errorf("invalid piece type name: please use one of {Queen|King|Bishop|Knight|Rook|Pawn} or empty string")
	errInvalidActionForGivenGame           = fmt.Errorf("the specified action is invalid for the specified game")
)

func (a API) DefaultGame() OutputGame {
	return mapGameToOutputGame(defaultGame)
}

func (a API) ParseGame(game InputGame) (OutputGame, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, err
	}
	return mapGameToOutputGame(parsedGame), nil
}

func (a API) DoAction(game InputGame, action InputAction) (OutputGame, OutputAction, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, OutputAction{}, err
	}
	parsedAction, err := a.parseAction(action, parsedGame)
	if err != nil {
		return OutputGame{}, OutputAction{}, err
	}
	return mapGameToOutputGame(parsedGame.doAction(parsedAction)), mapInternalActionToAction(parsedAction), nil
}
