// +build tinygo

package main

import (
	"syscall/js"

	"github.com/marianogappa/cheesse/api"
)

var a = api.New()

func DefaultGame(this js.Value, p []js.Value) interface{} {
	return js.ValueOf(convertOutputGame(a.DefaultGame()))
}

func ParseGame(this js.Value, p []js.Value) interface{} {
	og, err := a.ParseGame(convertToInputGame(p[0]))
	return js.ValueOf(map[string]interface{}{
		"outputGame": convertOutputGame(og),
		"error":      convertError(err),
	})
}

func DoAction(this js.Value, p []js.Value) interface{} {
	og, oa, err := a.DoAction(convertToInputGame(p[0]), convertToInputAction(p[1]))
	return js.ValueOf(map[string]interface{}{
		"outputGame":   convertOutputGame(og),
		"outputAction": convertOutputAction(oa),
		"error":        convertError(err),
	})
}

func ParseNotation(this js.Value, p []js.Value) interface{} {
	og, ogs, err := a.ParseNotation(convertToInputGame(p[0]), p[1].String())
	return js.ValueOf(map[string]interface{}{
		"outputGame":      convertOutputGame(og),
		"outputGameSteps": convertOutputGameSteps(ogs),
		"error":           convertError(err),
	})
}

func main() {
	js.Global().Set("DefaultGame", js.FuncOf(DefaultGame))
	js.Global().Set("ParseGame", js.FuncOf(ParseGame))
	js.Global().Set("DoAction", js.FuncOf(DoAction))
	js.Global().Set("ParseNotation", js.FuncOf(ParseNotation))
	select {}
}

func convertToInputAction(v js.Value) api.InputAction {
	return api.InputAction{
		FromSquare:         jsString(v.Get("fromSquare")),
		ToSquare:           jsString(v.Get("toSquare")),
		PromotionPieceType: jsString(v.Get("promotionPieceType")),
	}
}

func convertToInputGame(v js.Value) api.InputGame {
	board := v.Get("board")
	var (
		innerBoard []string
		outerBoard = api.Board{}
	)
	if board != js.Null() && board != js.Undefined() {
		innerBoard = make([]string, board.Length())
		for i := range innerBoard {
			innerBoard[i] = jsString(board.Index(i))
		}
		outerBoard = api.Board{
			Board:                   innerBoard,
			CanWhiteKingsideCastle:  jsBool(board.Get("canWhiteKingsideCastle")),
			CanWhiteQueensideCastle: jsBool(board.Get("canWhiteQueensideCastle")),
			CanBlackKingsideCastle:  jsBool(board.Get("canBlackKingsideCastle")),
			CanBlackQueensideCastle: jsBool(board.Get("canBlackQueensideCastle")),
			HalfMoveClock:           jsInt(board.Get("halfMoveClock")),
			FullMoveNumber:          jsInt(board.Get("fullMoveNumber")),
			EnPassantTargetSquare:   jsString(board.Get("enPassantTargetSquare")),
			Turn:                    jsString(board.Get("turn")),
		}
	}
	return api.InputGame{
		DefaultGame: jsBool(v.Get("defaultGame")),
		FENString:   jsString(v.Get("fenString")),
		Board:       outerBoard,
	}
}

func jsBool(j js.Value) bool {
	if j == js.Undefined() || j == js.Null() {
		return false
	}
	return j.Bool()
}

func jsInt(j js.Value) int {
	if j == js.Undefined() || j == js.Null() {
		return 0
	}
	return j.Int()
}

func jsString(j js.Value) string {
	if j == js.Undefined() || j == js.Null() {
		return ""
	}
	return j.String()
}

func convertError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func convertOutputGameSteps(ogs []api.OutputGameStep) []interface{} {
	is := make([]interface{}, len(ogs))
	for i := range ogs {
		is[i] = convertOutputGameStep(ogs[i])
	}
	return is
}

func convertOutputGameStep(ogs api.OutputGameStep) map[string]interface{} {
	return map[string]interface{}{
		"game":         convertOutputGame(ogs.Game),
		"action":       convertOutputAction(ogs.Action),
		"actionString": ogs.ActionString,
	}
}

func convertOutputGame(og api.OutputGame) map[string]interface{} {
	return map[string]interface{}{
		"fenString":               og.FENString,
		"board":                   convertBoard(og.Board),
		"actions":                 convertOutputActions(og.Actions),
		"canWhiteCastle":          og.CanWhiteCastle,
		"canWhiteKingsideCastle":  og.CanWhiteKingsideCastle,
		"canWhiteQueensideCastle": og.CanWhiteQueensideCastle,
		"canBlackCastle":          og.CanBlackCastle,
		"canBlackKingsideCastle":  og.CanBlackKingsideCastle,
		"canBlackQueensideCastle": og.CanBlackQueensideCastle,
		"halfMoveClock":           og.HalfMoveClock,
		"fullMoveNumber":          og.FullMoveNumber,
		"isLastMoveEnPassant":     og.IsLastMoveEnPassant,
		"enPassantTargetSquare":   og.EnPassantTargetSquare,
		"moveNumber":              og.MoveNumber,
		"blackPieces":             convertMapStringToString(og.BlackPieces),
		"whitePieces":             convertMapStringToString(og.WhitePieces),
		"blackKing":               og.BlackKing,
		"whiteKing":               og.WhiteKing,
		"isCheck":                 og.IsCheck,
		"isCheckmate":             og.IsCheckmate,
		"isStalemate":             og.IsStalemate,
		"isDraw":                  og.IsDraw,
		"isGameOver":              og.IsGameOver,
		"gameOverWinner":          og.GameOverWinner,
		"inCheckBy":               convertStringArr(og.InCheckBy),
	}
}

func convertStringArr(ss []string) []interface{} {
	is := make([]interface{}, len(ss))
	for i := 0; i < len(ss); i++ {
		is[i] = ss[i]
	}
	return is
}

func convertMapStringToString(mss map[string]string) map[string]interface{} {
	m := make(map[string]interface{}, len(mss))
	for k, v := range mss {
		m[k] = v
	}
	return m
}

func convertOutputActions(as []api.OutputAction) []interface{} {
	is := make([]interface{}, len(as))
	for i := range as {
		is[i] = convertOutputAction(as[i])
	}
	return is
}

func convertOutputAction(a api.OutputAction) map[string]interface{} {
	return map[string]interface{}{
		"fromPieceOwner":     a.FromPieceOwner,
		"fromPieceType":      a.FromPieceType,
		"fromPieceSquare":    a.FromPieceSquare,
		"fromSquare":         a.FromPieceSquare, // TODO alias so that it can be used as input action
		"toSquare":           a.ToSquare,
		"isCapture":          a.IsCapture,
		"isResign":           a.IsResign,
		"isPromotion":        a.IsPromotion,
		"isEnPassant":        a.IsEnPassant,
		"isEnPassantCapture": a.IsEnPassantCapture,
		"isCastle":           a.IsCastle,
		"isKingsideCastle":   a.IsKingsideCastle,
		"isQueensideCastle":  a.IsQueensideCastle,
		"promotionPieceType": a.PromotionPieceType,
		"capturedPieceType":  a.CapturedPieceType,
	}
}

func convertBoard(b api.Board) map[string]interface{} {
	return map[string]interface{}{
		"board":                   convertStringArr(b.Board),
		"canWhiteKingsideCastle":  b.CanWhiteKingsideCastle,
		"canWhiteQueensideCastle": b.CanWhiteQueensideCastle,
		"canBlackKingsideCastle":  b.CanBlackKingsideCastle,
		"canBlackQueensideCastle": b.CanBlackQueensideCastle,
		"halfMoveClock":           b.HalfMoveClock,
		"fullMoveNumber":          b.FullMoveNumber,
		"enPassantTargetSquare":   b.EnPassantTargetSquare,
		"turn":                    b.Turn,
	}
}
