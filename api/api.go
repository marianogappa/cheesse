package api

import (
	"fmt"
)

// // Cli
// // Server
// // Library

type API struct{}

func New() API { return API{} }

// func (a API) DoAction(game InputGame, action InputAction, options DoActionInputOptions) (OutputGame, OutputAction, error) {

// }

// func (a API) DefaultGame(options DefaultGameInputOptions) (OutputGame, error) {

// }
var (
	errInvalidInputGame = fmt.Errorf("invalid input game: please supply a valid fenString or a board")
)

func (a API) ParseGame(g InputGame) (OutputGame, error) {
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
		return OutputGame{}, err
	}

	return mapGameToOutputGame(parsedGame), nil
}
