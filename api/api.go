package api

import "errors"

// API represents the cheesse API. All cheesse API methods are exported methods of this struct.
type API struct{}

// New constructs an API.
func New() API { return API{} }

var (
	errInvalidInputGame                    = errors.New("invalid input game: please supply a valid fenString or a board")
	errAlgebraicSquareInvalidOrOutOfBounds = errors.New("invalid algebraic square: empty or out of bounds")
	errInvalidPieceTypeName                = errors.New("invalid piece type name: please use one of {Queen|King|Bishop|Knight|Rook|Pawn} or empty string")
	errInvalidActionForGivenGame           = errors.New("the specified action is invalid for the specified game")
)

// DefaultGame returns the initial game of chess, with all pieces on their default positions
// and before any action has taken place.
func (a API) DefaultGame() OutputGame {
	var defaultGame, _ = newGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return mapGameToOutputGame(defaultGame)
}

// ParseGame takes any valid input game and parses it, returning an OutputGame, which contains
// a lot of useful information about it, like possible actions, locations of pieces, game state
// in terms of threats, is the game over, etc.
//
// If the input game is invalid, an error will be returned with a description of the problem.
//
// Please refer to InputGame's and OutputGame's docs for format details.
func (a API) ParseGame(game InputGame) (OutputGame, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, err
	}
	return mapGameToOutputGame(parsedGame), nil
}

// DoAction takes any valid input game and any valid input action, parses them and attempts
// to apply the action on the given game. If parsing any of the entities fails or applying
// the action on the parsed game fails an error will be returned.
//
// If applying the action succeeds, it returns the parsed action and the resulting game
// AFTER applying the action.
//
// Please refer to InputGame's, InputAction's, OutputGame's and OutputAction's docs for
// format details.
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

// ParseNotation takes any valid input game and a string representing a match in some
// notation, parses them and attempts to play the match starting from the supplied
// game. If it fails, it returns an error describing the problem.
//
// If parsing the match succeeds, it returns the parsed initial game and a list of
// steps, one per action in the `notationString`.
//
// An example `notationString` (Scholar's mate):
//
// `1. e4 e5\n2. Bc4 Nc6\n3. Qh5 Nf6??\n4. Qxf7#`
//
// At the moment, only Algebraic Notation is supported, but support for most
// notations is planned at a later release.
//
// Please refer to InputGame's, OutputGame's and OutputGameStep's docs for format
// details.
func (a API) ParseNotation(game InputGame, notationString string) (OutputGame, []OutputGameStep, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, []OutputGameStep{}, err
	}

	// TODO at the moment there only exists an algebraic notation parser
	gameSteps, err := newNotationParserAlgebraic(characteristics{}).parse(parsedGame, notationString)
	return mapGameToOutputGame(parsedGame), mapGameStepsToOutputGameSteps(gameSteps), err
}
