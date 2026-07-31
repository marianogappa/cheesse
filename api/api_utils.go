package api

import (
	"strconv"
	"strings"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/marianogappa/cheesse/parser/pgn"
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
	if len(g.PositionHistory) > 0 {
		history := make([]uint64, len(g.PositionHistory))
		for i, s := range g.PositionHistory {
			if history[i], err = strconv.ParseUint(s, 16, 64); err != nil {
				return core.Game{}, errInvalidPositionHistory
			}
		}
		parsedGame = parsedGame.WithPositionHistory(history)
	}
	return parsedGame, nil
}

func (a API) parseAction(ia InputAction, g core.Game) (core.Action, error) {
	// Terminal actions don't carry squares; match them directly.
	if ia.IsResign || ia.IsDraw {
		for _, action := range g.Actions {
			if action.IsResign == ia.IsResign && action.IsDraw == ia.IsDraw {
				return action, nil
			}
		}
		return core.Action{}, errInvalidActionForGivenGame
	}

	// An action supplied as a move string in any notation: auto-detect and parse
	// it as a one-move match. Some parsers (e.g. ICCF) require a move-number
	// prefix, so retry with one if the bare string doesn't parse.
	if ia.ActionString != "" {
		// An ambiguous SAN string (e.g. "Nd2" with two knights reaching d2) must be
		// rejected rather than silently resolved to an arbitrary action.
		if matches, err := parser.NewGenericNotationParser(pgn.NewVariantPGN()).MatchHalfMove(ia.ActionString, g); err == nil && len(matches) > 1 {
			return core.Action{}, errAmbiguousActionString
		}
		for _, s := range []string{ia.ActionString, "1. " + ia.ActionString} {
			gameSteps, result := parseNotationAutoDetect(g, s)
			if result.ParseWasSuccessful && len(gameSteps) == 1 {
				return gameSteps[0].StepAction, nil
			}
		}
		return core.Action{}, errInvalidActionForGivenGame
	}

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
	// Board UIs often supply only from/to squares, so promotions default to Queen.
	if promotionPieceType == core.PieceNone {
		promotionPieceType = core.PieceQueen
	}

	for _, action := range g.Actions {
		if action.IsResign || action.IsDraw {
			continue // Terminal actions carry no squares; matched above.
		}
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
