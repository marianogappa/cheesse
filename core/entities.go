package core

import (
	"errors"
	"fmt"
	"strings"
)

type Game struct {
	CanWhiteCastle          bool
	CanWhiteKingsideCastle  bool
	CanWhiteQueensideCastle bool
	CanBlackCastle          bool
	CanBlackKingsideCastle  bool
	CanBlackQueensideCastle bool
	HalfMoveClock           int
	FullMoveNumber          int
	IsLastMoveEnPassant     bool
	EnPassantTargetSquare   XY
	MoveNumber              int
	Pieces                  []map[XY]Piece
	Kings                   []Piece
	IsCheck                 bool
	IsCheckmate             bool
	IsStalemate             bool
	IsDraw                  bool
	IsGameOver              bool
	GameOverWinner          color
	InCheckBy               []Piece
	Actions                 []Action
}

func (g Game) String() string {
	var sb strings.Builder
	for _, s := range g.ToBoard().Board {
		sb.WriteString(strings.Replace(s, " ", ".", -1))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g Game) Clone() Game {
	clonedPieces := make([]map[XY]Piece, len(g.Pieces))
	clonedKings := make([]Piece, len(g.Kings))
	for color, ownerPieces := range g.Pieces {
		clonedOwnerPieces := make(map[XY]Piece, len(ownerPieces))
		for _, piece := range ownerPieces {
			clonedPiece := piece
			clonedOwnerPieces[piece.XY] = clonedPiece
			if clonedPiece.PieceType == PieceKing {
				clonedKings[color] = clonedPiece
			}
		}
		clonedPieces[color] = clonedOwnerPieces
	}
	clonedInCheckBy := make([]Piece, len(g.InCheckBy))
	for i, piece := range g.InCheckBy {
		clonedInCheckBy[i] = piece
	}
	return Game{
		CanWhiteCastle:          g.CanWhiteCastle,
		CanWhiteKingsideCastle:  g.CanWhiteKingsideCastle,
		CanWhiteQueensideCastle: g.CanWhiteQueensideCastle,
		CanBlackCastle:          g.CanBlackCastle,
		CanBlackKingsideCastle:  g.CanBlackKingsideCastle,
		CanBlackQueensideCastle: g.CanBlackQueensideCastle,
		HalfMoveClock:           g.HalfMoveClock,
		FullMoveNumber:          g.FullMoveNumber,
		EnPassantTargetSquare:   g.EnPassantTargetSquare,
		MoveNumber:              g.MoveNumber,
		Pieces:                  clonedPieces,
		Kings:                   clonedKings,
		IsCheck:                 g.IsCheck,
		IsCheckmate:             g.IsCheckmate,
		IsStalemate:             g.IsStalemate,
		IsDraw:                  g.IsDraw,
		IsGameOver:              g.IsGameOver,
		GameOverWinner:          g.GameOverWinner,
		InCheckBy:               clonedInCheckBy,
	}
}

type castleType int

const (
	castleTypeQueenside = iota
	castleTypeKingside
)

type Action struct {
	FromPiece          Piece
	ToXY               XY
	IsCapture          bool
	IsResign           bool
	IsPromotion        bool
	IsEnPassant        bool
	IsEnPassantCapture bool
	IsCastle           bool
	IsKingsideCastle   bool
	IsQueensideCastle  bool
	PromotionPieceType PieceType
	CapturedPiece      Piece
}

func (a Action) String() string {
	switch {
	case a.IsEnPassantCapture:
		return fmt.Sprintf("%s's Pawn at %v captures %v's Pawn at %v which was doing en passant", a.FromPiece.Owner, a.FromPiece.XY.ToAlgebraic(), a.CapturedPiece.Owner, a.CapturedPiece.XY.ToAlgebraic())
	case a.IsCapture && a.IsPromotion:
		return fmt.Sprintf("%s's %s at %v captures %s's %s at %v while promoting to %v", a.FromPiece.Owner, a.FromPiece.PieceType, a.FromPiece.XY.ToAlgebraic(), a.CapturedPiece.Owner, a.CapturedPiece.PieceType, a.CapturedPiece.XY.ToAlgebraic(), a.PromotionPieceType)
	case a.IsCapture:
		return fmt.Sprintf("%s's %s at %v captures %s's %s at %v", a.FromPiece.Owner, a.FromPiece.PieceType, a.FromPiece.XY.ToAlgebraic(), a.CapturedPiece.Owner, a.CapturedPiece.PieceType, a.CapturedPiece.XY.ToAlgebraic())
	case a.IsResign:
		return fmt.Sprintf("%s resigns", a.FromPiece.Owner)
	case a.IsPromotion:
		return fmt.Sprintf("%s's Pawn at %v promotes to %v", a.FromPiece.Owner, a.FromPiece.XY.ToAlgebraic(), a.PromotionPieceType)
	case a.IsEnPassant:
		return fmt.Sprintf("%s's Pawn at %v does en passant", a.FromPiece.Owner, a.FromPiece.XY.ToAlgebraic())
	case a.IsKingsideCastle:
		return fmt.Sprintf("%s kingside castles", a.FromPiece.Owner)
	case a.IsQueensideCastle:
		return fmt.Sprintf("%s queenside castles", a.FromPiece.Owner)
	}
	return fmt.Sprintf("%s's %s at %v moves to %v", a.FromPiece.Owner, a.FromPiece.PieceType, a.FromPiece.XY.ToAlgebraic(), a.ToXY.ToAlgebraic())
}

func (a Action) DebugString() string {
	return fmt.Sprintf("%v at (%v, %v) to (%v, %v), isCapture: %v , isResign: %v , isPromotion: %v , isEnPassant: %v , isEnPassantCapture: %v , isCastle: %v , isKingsideCastle: %v , isQueensideCastle: %v, promotionPieceType: %v, capturedPiece: %v at (%v, %v)",
		a.FromPiece.PieceType,
		a.FromPiece.XY.X,
		a.FromPiece.XY.Y,
		a.ToXY.X,
		a.ToXY.Y,
		a.IsCapture,
		a.IsResign,
		a.IsPromotion,
		a.IsEnPassant,
		a.IsEnPassantCapture,
		a.IsCastle,
		a.IsKingsideCastle,
		a.IsQueensideCastle,
		a.PromotionPieceType,
		a.CapturedPiece.PieceType,
		a.CapturedPiece.XY.X,
		a.CapturedPiece.XY.Y,
	)
}

type PieceType int

const (
	PieceNone = iota
	PieceQueen
	PieceKing
	PieceBishop
	PieceKnight
	PieceRook
	PiecePawn
)

func (t PieceType) String() string {
	switch t {
	case PieceQueen:
		return "Queen"
	case PieceKing:
		return "King"
	case PieceBishop:
		return "Bishop"
	case PieceKnight:
		return "Knight"
	case PieceRook:
		return "Rook"
	case PiecePawn:
		return "Pawn"
	}
	return ""
}

func (t PieceType) ToICCF() string {
	switch t {
	case PieceQueen:
		return "1"
	case PieceBishop:
		return "3"
	case PieceKnight:
		return "4"
	case PieceRook:
		return "2"
	}
	return ""
}

func (t PieceType) ToSmith() string {
	switch t {
	case PieceQueen:
		return "q"
	case PieceKing:
		return "k"
	case PieceBishop:
		return "b"
	case PieceKnight:
		return "n"
	case PieceRook:
		return "r"
	case PiecePawn:
		return "p"
	}
	return ""
}

func (t PieceType) ToDescriptive(useKt bool) string {
	switch t {
	case PieceQueen:
		return "Q"
	case PieceKing:
		return "K"
	case PieceBishop:
		return "B"
	case PieceKnight:
		if useKt {
			return "Kt"
		}
		return "N"
	case PieceRook:
		return "R"
	case PiecePawn:
		return "P"
	}
	return ""
}

func (t PieceType) ToAlgebraic() string {
	if t == PiecePawn {
		return ""
	}
	return strings.ToUpper(t.ToSmith())
}

func (t PieceType) ToFigurine() string {
	switch t {
	case PieceQueen:
		return "♛"
	case PieceKing:
		return "♚"
	case PieceBishop:
		return "♝"
	case PieceKnight:
		return "♞"
	case PieceRook:
		return "♜"
	case PiecePawn:
		return "♟"
	}
	return ""
}

type color int

const (
	ColorBlack = iota
	ColorWhite
)

func (c color) String() string {
	switch c {
	case ColorBlack:
		return "Black"
	case ColorWhite:
		return "White"
	}
	return "Unknown"
}

func (c color) Opponent() color {
	if c == ColorBlack {
		return ColorWhite
	}
	return ColorBlack
}

type XY struct {
	X, Y int
}

func (c XY) add(c2 XY) XY {
	return XY{c.X + c2.X, c.Y + c2.Y}
}

func (c XY) eq(c2 XY) bool {
	return c.X == c2.X && c.Y == c2.Y
}

func (c XY) ToAlgebraic() string {
	return fmt.Sprintf("%v%v", string("abcdefgh"[c.X]), 8-c.Y)
}

func (c XY) ToDescriptive(turn color) string {
	fileNames := []string{"QR", "QN", "QB", "Q", "K", "KB", "KN", "KR"} // TODO: this doesn't respect Kt characteristic
	y := c.Y + 1
	if turn == ColorWhite {
		y = 8 - c.Y
	}
	return fmt.Sprintf("%v%v", fileNames[c.X], y)
}

func (c XY) ToICCF() string {
	return fmt.Sprintf("%v%v", c.X+1, 8-c.Y)
}

func (c XY) deltaTowards(c2 XY) XY {
	var delta XY
	switch {
	case c2.X < c.X:
		delta.X = -1
	case c2.X > c.X:
		delta.X = 1
	}
	switch {
	case c2.Y < c.Y:
		delta.Y = -1
	case c2.Y > c.Y:
		delta.Y = 1
	}
	return delta
}

func (p Piece) xysTowards(sq XY) []XY {
	if p.PieceType != PieceQueen && p.PieceType != PieceBishop && p.PieceType != PieceRook {
		return []XY{}
	}
	var (
		rookLike   = p.XY.X == sq.X || p.XY.Y == sq.Y
		absDist    = XY{abs(p.XY.X - sq.X), abs(p.XY.Y - sq.Y)}
		bishopLike = (absDist.X - absDist.Y) == 0
	)
	if (p.PieceType == PieceRook && !rookLike) || (p.PieceType == PieceBishop && !bishopLike) || (p.PieceType == PieceQueen && !rookLike && !bishopLike) {
		return []XY{}
	}
	var (
		delta = p.XY.deltaTowards(sq)
		cur   = p.XY.add(delta)
		xys   = []XY{}
	)
	for cur != sq {
		xys = append(xys, cur)
		cur = cur.add(delta)
	}
	return xys
}

type Piece struct {
	PieceType PieceType
	Owner     color
	XY        XY
}

func (p Piece) String() string {
	return fmt.Sprintf("%v's %v at %v", p.Owner, p.PieceType, p.XY.ToAlgebraic())
}

var movementDeltasByPieceType = map[PieceType][]XY{
	PieceQueen:  {{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}},
	PieceKing:   {{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}, {-2, 0}, {2, 0}}, // Castling
	PieceBishop: {{-1, -1}, {1, 1}, {-1, 1}, {1, -1}},
	PieceKnight: {{-1, -2}, {1, -2}, {-1, 2}, {1, 2}, {-2, -1}, {-2, 1}, {2, -1}, {2, 1}},
	PieceRook:   {{-1, 0}, {1, 0}, {0, -1}, {0, 1}},
	// N.B. Pawn will be dealt with separately, because it's dependant on color
}

var (
	errNotInBounds                = errors.New("piece is not in board bounds")
	errFriendlyPieceInDestination = errors.New("there is a friendly piece in destination")
	errPieceInBetween             = errors.New("there is a piece between source and destination")
	errPieceBlockingPawn          = errors.New("there is a piece blocking pawn from moving forwards or en passant")
	errPawnCantCapture            = errors.New("pawn can't capture because there is no opponent piece (including en passant)")
	errCantCastle                 = errors.New("king can't castle, because pieces moved, pieces in middle or squares threatened")
	errCantPromote                = errors.New("pawn can't promote because wrong position or invalid promotion piece type")
	errActionLeavesKingThreatened = errors.New("action leaves king in a check")
)

var (
	emptyXYsForCastlingByColorAndCastleType = map[color]map[castleType][]XY{
		ColorBlack: {
			castleTypeQueenside: {{1, 0}, {2, 0}, {3, 0}},
			castleTypeKingside:  {{5, 0}, {6, 0}},
		},
		ColorWhite: {
			castleTypeQueenside: {{1, 7}, {2, 7}, {3, 7}},
			castleTypeKingside:  {{5, 7}, {6, 7}},
		},
	}
	unthreatenedXYsForCastlingByColorAndCastleType = map[color]map[castleType][]XY{
		ColorBlack: {
			castleTypeQueenside: {{1, 0}, {2, 0}, {3, 0}, {4, 0}},
			castleTypeKingside:  {{4, 0}, {5, 0}, {6, 0}},
		},
		ColorWhite: {
			castleTypeQueenside: {{1, 7}, {2, 7}, {3, 7}, {4, 7}},
			castleTypeKingside:  {{4, 7}, {5, 7}, {6, 7}},
		},
	}
)

type GameStep struct {
	StepString string
	StepAction Action
	StepGame   Game
}

func (s GameStep) Clone() GameStep {
	return GameStep{
		StepString: s.StepString,
		StepAction: s.StepAction,
		StepGame:   s.StepGame.Clone(),
	}
}
