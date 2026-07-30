package api

import (
	"errors"
	"strings"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/marianogappa/cheesse/parser/pgn"
	"github.com/marianogappa/cheesse/printer"
)

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
	var defaultGame, _ = core.NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
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
	newGame := parsedGame.DoAction(parsedAction)
	outputAction := mapInternalActionToAction(parsedAction)
	actionString, err := printer.AlgebraicPrinter{}.PrintAction(core.GameStep{StepAction: parsedAction, StepGame: newGame, StepPreMoveGame: parsedGame}, printer.SANCharacteristics())
	if err != nil {
		return OutputGame{}, OutputAction{}, err
	}
	outputAction.ActionString = actionString
	return mapGameToOutputGame(newGame), outputAction, nil
}

// ParseNotation takes any valid input game and a string representing a match in some
// notation, auto-detects the notation and attempts to play the match starting from
// the supplied game.
//
// The notation is auto-detected across all supported notations: Algebraic/SAN
// (including figurine and PGN), Coordinate, Descriptive, ICCF and Smith. All
// supported notations are attempted, and the attempt that parses the furthest wins.
//
// Partial parses are supported: if the notation string stops being valid at some
// point, the result still contains the valid prefix of steps, the count of valid
// actions, the name of the most likely notation, and a description of the parse
// failure.
//
// An example `notationString` (Scholar's mate):
//
// `1. e4 e5\n2. Bc4 Nc6\n3. Qh5 Nf6??\n4. Qxf7#`
//
// An error is only returned if the input game itself is invalid.
//
// Please refer to InputGame's, OutputGame's, OutputGameStep's and OutputParseResult's
// docs for format details.
func (a API) ParseNotation(game InputGame, notationString string) (OutputGame, OutputParseResult, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, OutputParseResult{}, err
	}
	gameSteps, result := parseNotationAutoDetect(parsedGame, notationString)
	result.Steps = mapGameStepsToOutputGameSteps(gameSteps)
	return mapGameToOutputGame(parsedGame), result, nil
}

// parseNotationAutoDetect tries all supported notation parsers and returns the game
// steps of the attempt that parsed the furthest, along with a parse result describing
// the winning attempt (with Steps unset; callers map the steps as needed).
func parseNotationAutoDetect(parsedGame core.Game, notationString string) ([]core.GameStep, OutputParseResult) {
	notationString = strings.TrimSpace(notationString)
	type notationCandidate struct {
		name  string
		parse func() ([]core.GameStep, error)
	}
	candidates := []notationCandidate{
		{"Algebraic Notation", func() ([]core.GameStep, error) {
			return parser.NewNotationParserAlgebraic(parser.Characteristics{}).Parse(parsedGame, notationString)
		}},
		{"ICCF Notation", func() ([]core.GameStep, error) {
			return parser.NewNotationParserICCF(parser.Characteristics{}).Parse(parsedGame, notationString)
		}},
		{"Smith Notation", func() ([]core.GameStep, error) {
			return parser.NewNotationParserSmith(parser.Characteristics{}).Parse(parsedGame, notationString)
		}},
		{"Coordinate Notation", func() ([]core.GameStep, error) {
			return parser.NewNotationParserCoordinate(parser.Characteristics{}).Parse(parsedGame, notationString)
		}},
		{"Descriptive Notation", func() ([]core.GameStep, error) {
			return parser.NewNotationParserDescriptive(parser.Characteristics{}).Parse(parsedGame, notationString)
		}},
		{"PGN", func() ([]core.GameStep, error) {
			parsed, err := parser.NewGenericNotationParser(pgn.NewVariantPGN()).Parse(parsedGame, notationString)
			if parsed == nil {
				return nil, err
			}
			return parsed.GameSteps, err
		}},
	}

	var (
		bestSteps  []core.GameStep
		bestResult *OutputParseResult
	)
	for _, candidate := range candidates {
		gameSteps, err := candidate.parse()
		validSteps := len(gameSteps)
		result := OutputParseResult{
			NotationName:       candidate.name,
			ParseWasSuccessful: err == nil,
			ValidActionCount:   validSteps,
		}
		if err != nil {
			result.Error = err.Error()
		}

		// A fully-successful parse with at least one step wins immediately.
		if result.ParseWasSuccessful && validSteps > 0 {
			return gameSteps, result
		}
		// Otherwise keep the attempt that parsed the furthest.
		if bestResult == nil || validSteps > bestResult.ValidActionCount {
			bestSteps, bestResult = gameSteps, &result
		}
	}

	return bestSteps, *bestResult
}

// ConvertNotation takes any valid input game and a string representing a match in
// some notation, auto-detects the source notation, and re-renders every move in the
// target notation.
//
// `targetNotation` must be one of: `{Algebraic|Figurine|Descriptive|Coordinate|ICCF|Smith}`
// (case-insensitive).
//
// Partial input still converts the valid prefix: the result reports the detected
// source notation, whether the whole input parsed, how many actions were valid, and
// one OutputGameStep per valid action whose `actionString` is rendered in the target
// notation.
//
// An error is only returned if the input game itself is invalid or the target
// notation is unknown.
//
// Please refer to InputGame's, OutputGame's, OutputGameStep's and OutputParseResult's
// docs for format details.
func (a API) ConvertNotation(game InputGame, notationString string, targetNotation string) (OutputGame, OutputParseResult, error) {
	parsedGame, err := a.parseGame(game)
	if err != nil {
		return OutputGame{}, OutputParseResult{}, err
	}

	targetPrinter, targetCharacteristics, err := notationPrinter(targetNotation)
	if err != nil {
		return OutputGame{}, OutputParseResult{}, err
	}

	gameSteps, result := parseNotationAutoDetect(parsedGame, notationString)

	result.Steps = mapGameStepsToOutputGameSteps(gameSteps)
	for i, gameStep := range gameSteps {
		if gameStep.StepAction == (core.Action{}) {
			// Result markers (e.g. "1-0" in PGN) have no action to re-render.
			continue
		}
		actionString, err := targetPrinter.PrintAction(gameStep, targetCharacteristics)
		if err != nil {
			return OutputGame{}, OutputParseResult{}, err
		}
		result.Steps[i].ActionString = actionString
	}

	return mapGameToOutputGame(parsedGame), result, nil
}

var errUnknownTargetNotation = errors.New("unknown target notation: please use one of {Algebraic|Figurine|Descriptive|Coordinate|ICCF|Smith}")

func notationPrinter(targetNotation string) (printer.NotationPrinter, printer.GameCharacteristics, error) {
	switch strings.ToLower(targetNotation) {
	case "algebraic":
		return printer.AlgebraicPrinter{}, printer.SANCharacteristics(), nil
	case "figurine":
		return printer.AlgebraicPrinter{}, printer.FigurineCharacteristics(), nil
	case "descriptive":
		return printer.DescriptivePrinter{}, printer.GameCharacteristics{}, nil
	case "coordinate":
		return printer.CoordinatePrinter{}, printer.SANCharacteristics(), nil
	case "iccf":
		return printer.ICCFPrinter{}, printer.GameCharacteristics{}, nil
	case "smith":
		return printer.SmithPrinter{}, printer.GameCharacteristics{}, nil
	}
	return nil, printer.GameCharacteristics{}, errUnknownTargetNotation
}
