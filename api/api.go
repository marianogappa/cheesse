package api

import "errors"

type API struct{}

func New() API { return API{} }

var (
	errInvalidInputGame                    = errors.New("invalid input game: please supply a valid fenString or a board")
	errAlgebraicSquareInvalidOrOutOfBounds = errors.New("invalid algebraic square: empty or out of bounds")
	errInvalidPieceTypeName                = errors.New("invalid piece type name: please use one of {Queen|King|Bishop|Knight|Rook|Pawn} or empty string")
	errInvalidActionForGivenGame           = errors.New("the specified action is invalid for the specified game")
)

func (a API) DefaultGame() OutputGame {
	var defaultGame, _ = newGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
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

func (a API) ParseNotation(game InputGame, notationString string) (OutputGame, []OutputGameStep, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, []OutputGameStep{}, err
	}

	// TODO at the moment there only exists an algebraic notation parser
	gameSteps, err := newNotationParserAlgebraic(characteristics{}).parse(parsedGame, notationString)
	return mapGameToOutputGame(parsedGame), mapGameStepsToOutputGameSteps(gameSteps), err
}
