package api

// InputGame is the input interface to supply a chess game.
//
// There are 3 different ways to supply the chess game:
//
// 1. via `fenString`: supply a FEN Notation string.
//
// 2. via `board`: supply the board, together with required aspects of the game state.
//
// 3. via empty struct: assumes the defaultGame.
//
// If you supply both the `fenString` and the `board`, `board` is ignored silently.
type InputGame struct {
	FENString string `json:"fenString"`
	Board     Board  `json:"board"`
}

// InputAction is the input interface to supply a chess action.
//
// - `fromSquare` and `toSquare` are required, and must be board cells described in
// Algebraic Notation (e.g. `e2`). Note that `a1` is where the White Queen's Rook starts.
//
// - `promotionPieceType` is only required if the action is a promotion.
//
// - `promotionPieceType` must be one of: `{Queen|King|Bishop|Knight|Rook|Pawn}`.
type InputAction struct {
	FromSquare         string `json:"fromSquare"`
	ToSquare           string `json:"toSquare"`
	PromotionPieceType string `json:"promotionPieceType"`
}

// Board is one of the input interfaces to supply a chess game.
//
// The `board` struct member must consist of 8 strings of length 8, containing the
// representation of the chess board using unicode characters. For the empty cells,
// any character may be used, but the output board will use spaces.
//
// The other struct members represent all required game state as described in the
// FEN notation.
//
// - `enPassantTargetSquare` must be a board cell described in Algebraic Notation
// (e.g. `e2`). Note that `a1` is where the White Queen's Rook starts. If there's
// no en passant target square, it must be an empty string.
//
// - `turn` must be one of: `{Black|White}`.
type Board struct {
	Board                   []string `json:"board"`
	CanWhiteKingsideCastle  bool     `json:"canWhiteKingsideCastle"`
	CanWhiteQueensideCastle bool     `json:"canWhiteQueensideCastle"`
	CanBlackKingsideCastle  bool     `json:"canBlackKingsideCastle"`
	CanBlackQueensideCastle bool     `json:"canBlackQueensideCastle"`
	HalfMoveClock           int      `json:"halfMoveClock"`
	FullMoveNumber          int      `json:"fullMoveNumber"`
	EnPassantTargetSquare   string   `json:"enPassantTargetSquare"` // in Algebraic notation, or empty string
	Turn                    string   `json:"turn"`                  // "Black" or "White"
}

// OutputGame is the output interface that describes a chess game.
// All API calls that return a chess game represent it with an OutputGame.
//
// - `fenString` represents the chess game as a FEN Notation string.
//
// - `actions` is the exhaustive list of actions that can follow from this game.
//
// - `enPassantTargetSquare` is a board cell described in Algebraic Notation
// (e.g. `e2`). Note that `a1` is where the White Queen's Rook starts. If there's
// no en passant target square, it's an empty string.
//
// - `blackPieces` and `whitePieces` are maps from cells to piece names. The cells
// are represented in Algebraic Notation (e.g `e2`), and the piece names are one
// of `{Queen|King|Bishop|Knight|Rook|Pawn}`.
//
// - `blackKing` and `whiteKing` are the cells where the Kings are located. The
// cells are represented in Algebraic Notation (e.g `e2`).
//
// - `gameOverWinner` is one of `{Black|White|Unknown}`, and represents the winner
// of the game, when `isGameOver` is true. `Unknown` otherwise.
//
// - `inCheckBy` is a list of cells whose pieces are threatening the player whose
// turn it is to move. `board.turn` dictates who this player is. The cells are
// represented in Algebraic Notation (e.g `e2`). To find out which piece is in a
// cell, inspect `blackPieces` and `whitePieces`.
//
// Because OutputGame is a superset of InputGame, you may supply an OutputGame to
// any API call that expects an InputGame.
type OutputGame struct {
	FENString               string            `json:"fenString"`
	Board                   Board             `json:"board"`
	Actions                 []OutputAction    `json:"actions"`
	CanWhiteCastle          bool              `json:"canWhiteCastle"`
	CanWhiteKingsideCastle  bool              `json:"canWhiteKingsideCastle"`
	CanWhiteQueensideCastle bool              `json:"canWhiteQueensideCastle"`
	CanBlackCastle          bool              `json:"canBlackCastle"`
	CanBlackKingsideCastle  bool              `json:"canBlackKingsideCastle"`
	CanBlackQueensideCastle bool              `json:"canBlackQueensideCastle"`
	HalfMoveClock           int               `json:"halfMoveClock"`
	FullMoveNumber          int               `json:"fullMoveNumber"`
	IsLastMoveEnPassant     bool              `json:"isLastMoveEnPassant"`
	EnPassantTargetSquare   string            `json:"enPassantTargetSquare"`
	MoveNumber              int               `json:"moveNumber"`
	BlackPieces             map[string]string `json:"blackPieces"`
	WhitePieces             map[string]string `json:"whitePieces"`
	BlackKing               string            `json:"blackKing"`
	WhiteKing               string            `json:"whiteKing"`
	IsCheck                 bool              `json:"isCheck"`
	IsCheckmate             bool              `json:"isCheckmate"`
	IsStalemate             bool              `json:"isStalemate"`
	IsDraw                  bool              `json:"isDraw"`
	IsGameOver              bool              `json:"isGameOver"`
	GameOverWinner          string            `json:"gameOverWinner"`
	InCheckBy               []string          `json:"inCheckBy"`
}

// OutputAction is the output interface that describes a chess action.
// All API calls that return a chess action represent it with an OutputAction.
//
// - `fromPieceOwner` is one of `{Black|White}`. The owner of the piece doing
// the action.
//
// - `fromPieceType` is one of `{Queen|King|Bishop|Knight|Rook|Pawn}`.
//
// - `fromPieceSquare` and `toSquare` are the source and destination board
// cells for the action, described in Algebraic Notation (e.g. `e2`). Note
// that `a1` is where the White Queen's Rook starts.
//
// - `promotionPieceType` is one of `{Queen|King|Bishop|Knight|Rook|Pawn}`,
// and represents the piece that a Pawn promotes to, if the action is a
// promotion. If the action is not a promotion, it's an empty string.
//
// - `capturedPieceType` is one of `{Queen|King|Bishop|Knight|Rook|Pawn}`,
// and represents the piece that was captured, if the action is a capture.
// If the action is not a capture, it's an empty string.
type OutputAction struct {
	FromPieceOwner     string `json:"fromPieceOwner"`
	FromPieceType      string `json:"fromPieceType"`
	FromPieceSquare    string `json:"fromPieceSquare"`
	ToSquare           string `json:"toSquare"`
	IsCapture          bool   `json:"isCapture"`
	IsResign           bool   `json:"isResign"`
	IsPromotion        bool   `json:"isPromotion"`
	IsEnPassant        bool   `json:"isEnPassant"`
	IsEnPassantCapture bool   `json:"isEnPassantCapture"`
	IsCastle           bool   `json:"isCastle"`
	IsKingsideCastle   bool   `json:"isKingsideCastle"`
	IsQueensideCastle  bool   `json:"isQueensideCastle"`
	PromotionPieceType string `json:"promotionPieceType"`
	CapturedPieceType  string `json:"capturedPieceType"`
}

// OutputGameStep is the output interface that describes a step in a parsed
// or converted match in a given notation string.
//
// A given notation string is parsed into a list of "action strings", each
// one representing each action in the match. There will be an OutputGameStep
// for each one of these "action strings".
//
// Please refer to the docs for the OutputGame and OutputAction formats.
//
// -  `actionString` is a string representing a chess action as supplied by the
// client, so it could be in any of the supported notations, and is not
// modified by the API. It could be incorrectly sliced, though.
//
// - `action` represents the action that the API inferred from the
// `actionString`.
//
// - `game` represents the chess game AFTER applying the inferred action.
type OutputGameStep struct {
	Game         OutputGame   `json:"game"`
	Action       OutputAction `json:"action"`
	ActionString string       `json:"actionString"`
}

func mapGameToOutputGame(g game) OutputGame {
	var o OutputGame

	o.FENString = g.toFEN()
	o.Board = mapInternalBoardToBoard(g.toBoard())
	o.Actions = make([]OutputAction, len(g.actions))
	o.CanWhiteCastle = g.canWhiteCastle
	o.CanWhiteKingsideCastle = g.canWhiteKingsideCastle
	o.CanWhiteQueensideCastle = g.canWhiteQueensideCastle
	o.CanBlackCastle = g.canBlackCastle
	o.CanBlackKingsideCastle = g.canBlackKingsideCastle
	o.CanBlackQueensideCastle = g.canBlackQueensideCastle
	o.HalfMoveClock = g.halfMoveClock
	o.FullMoveNumber = g.fullMoveNumber
	o.IsLastMoveEnPassant = g.isLastMoveEnPassant
	o.EnPassantTargetSquare = ""
	o.MoveNumber = g.moveNumber
	o.BlackPieces = make(map[string]string, len(g.pieces[colorBlack]))
	o.WhitePieces = make(map[string]string, len(g.pieces[colorWhite]))
	o.BlackKing = g.kings[colorBlack].xy.toAlgebraic()
	o.WhiteKing = g.kings[colorWhite].xy.toAlgebraic()
	o.IsCheck = g.isCheck
	o.IsCheckmate = g.isCheckmate
	o.IsStalemate = g.isStalemate
	o.IsDraw = g.isDraw
	o.IsGameOver = g.isGameOver
	o.GameOverWinner = g.gameOverWinner.String()
	o.InCheckBy = make([]string, len(g.inCheckBy))

	for i := range g.actions {
		o.Actions[i] = mapInternalActionToAction(g.actions[i])
	}

	if g.isLastMoveEnPassant {
		o.EnPassantTargetSquare = g.enPassantTargetSquare.toAlgebraic()
	}

	for sq, p := range g.pieces[colorBlack] {
		o.BlackPieces[sq.toAlgebraic()] = p.pieceType.String()
	}

	for sq, p := range g.pieces[colorWhite] {
		o.WhitePieces[sq.toAlgebraic()] = p.pieceType.String()
	}

	for i := range g.inCheckBy {
		o.InCheckBy[i] = g.inCheckBy[i].xy.toAlgebraic()
	}

	return o
}

func mapInternalBoardToBoard(b board) Board {
	return Board{
		Board:                   b.board,
		CanWhiteKingsideCastle:  b.canWhiteKingsideCastle,
		CanWhiteQueensideCastle: b.canWhiteQueensideCastle,
		CanBlackKingsideCastle:  b.canBlackKingsideCastle,
		CanBlackQueensideCastle: b.canBlackQueensideCastle,
		HalfMoveClock:           b.halfMoveClock,
		FullMoveNumber:          b.fullMoveNumber,
		EnPassantTargetSquare:   b.enPassantTargetSquare,
		Turn:                    b.turn,
	}
}

func mapBoardToInternalBoard(b Board) board {
	return board{
		board:                   b.Board,
		canWhiteKingsideCastle:  b.CanWhiteKingsideCastle,
		canWhiteQueensideCastle: b.CanWhiteQueensideCastle,
		canBlackKingsideCastle:  b.CanBlackKingsideCastle,
		canBlackQueensideCastle: b.CanBlackQueensideCastle,
		halfMoveClock:           b.HalfMoveClock,
		fullMoveNumber:          b.FullMoveNumber,
		enPassantTargetSquare:   b.EnPassantTargetSquare,
		turn:                    b.Turn,
	}
}

func mapInternalActionToAction(a action) OutputAction {
	return OutputAction{
		FromPieceOwner:     a.fromPiece.owner.String(),
		FromPieceType:      a.fromPiece.pieceType.String(),
		FromPieceSquare:    a.fromPiece.xy.toAlgebraic(),
		ToSquare:           a.toXY.toAlgebraic(),
		IsCapture:          a.isCapture,
		IsResign:           a.isResign,
		IsPromotion:        a.isPromotion,
		IsEnPassant:        a.isEnPassant,
		IsEnPassantCapture: a.isEnPassantCapture,
		IsCastle:           a.isCastle,
		IsKingsideCastle:   a.isKingsideCastle,
		IsQueensideCastle:  a.isQueensideCastle,
		PromotionPieceType: a.promotionPieceType.String(),
		CapturedPieceType:  a.capturedPiece.pieceType.String(),
	}
}

func mapGameStepsToOutputGameSteps(gss []gameStep) []OutputGameStep {
	ogs := make([]OutputGameStep, len(gss))
	for i, gs := range gss {
		ogs[i] = OutputGameStep{
			Game:         mapGameToOutputGame(gs.g),
			Action:       mapInternalActionToAction(gs.a),
			ActionString: gs.s,
		}
	}
	return ogs
}
