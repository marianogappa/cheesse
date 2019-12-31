package api

type InputGame struct {
	FENString string `json:"fenString"`
	Board     Board  `json:"board"`
}

type InputAction struct {
	FromSquare       string `json:"fromSquare"`
	ToSquare         string `json:"toSquare"`
	PromotePieceType string `json:"promotePieceType"`
}

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
