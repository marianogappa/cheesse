package core

import (
	"fmt"
	"math/bits"
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
	bb                      [2][7]uint64 // one bitboard per color per piece type (PieceNone index unused)
	occ                     [2]uint64    // occupancy per color
	squares                 [64]uint8    // pieceType<<1|color per square, 0 if empty
	kingSq                  [2]int8
	IsCheck                 bool
	IsDoubleCheck           bool
	IsDiscoverCheck         bool
	IsCheckmate             bool
	IsStalemate             bool
	IsDraw                  bool
	CanClaimDraw            bool
	IsGameOver              bool
	GameOverWinner          color
	InCheckBy               []Piece
	Actions                 []Action
	// positionHistory holds the Zobrist hashes of positions reached since the last
	// irreversible move (pawn move, capture or castling-right change), most recent
	// last. Used for threefold/fivefold repetition detection.
	positionHistory []uint64
}

// PositionHistory returns the Zobrist hashes of the positions reached since the
// last irreversible move, most recent last. It can be persisted and restored via
// WithPositionHistory to detect repetitions across stateless API calls.
func (g Game) PositionHistory() []uint64 {
	history := make([]uint64, len(g.positionHistory))
	copy(history, g.positionHistory)
	return history
}

// WithPositionHistory returns a copy of the game whose position history is the given
// hashes (as returned by PositionHistory of a previous game), and with draw flags
// recalculated accordingly. The current position is appended if it isn't already the
// last entry, so both "history up to and including this position" and "history of
// prior positions" are accepted.
func (g Game) WithPositionHistory(history []uint64) Game {
	currentHash := g.positionHash()
	g.positionHistory = make([]uint64, 0, len(history)+1)
	g.positionHistory = append(g.positionHistory, history...)
	if len(g.positionHistory) == 0 || g.positionHistory[len(g.positionHistory)-1] != currentHash {
		g.positionHistory = append(g.positionHistory, currentHash)
	}
	return g.calculateCriticalFlags()
}

// Color is the exported name for the color of a player or piece (e.g. core.ColorWhite).
type Color = color

// Pieces returns all pieces of the given color.
func (g Game) Pieces(c color) []Piece {
	pieces := make([]Piece, 0, bits.OnesCount64(g.occ[c]))
	for occ := g.occ[c]; occ != 0; occ &= occ - 1 {
		pieces = append(pieces, g.pieceAtSq(bits.TrailingZeros64(occ)))
	}
	return pieces
}

// PieceAt returns the piece at the given square, or the zero Piece (PieceNone) if empty.
func (g Game) PieceAt(xy XY) Piece {
	return g.pieceAtSq(sqOf(xy))
}

// King returns the king of the given color.
func (g Game) King(c color) Piece {
	return g.pieceAtSq(int(g.kingSq[c]))
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
	// The board layout lives in value-type arrays, so the struct copy is already a
	// deep copy; only the slices need explicit copying. IsLastMoveEnPassant is
	// deliberately dropped, matching the historical Clone behavior.
	clonedGame := g
	clonedGame.IsLastMoveEnPassant = false
	clonedGame.InCheckBy = make([]Piece, len(g.InCheckBy))
	copy(clonedGame.InCheckBy, g.InCheckBy)
	clonedGame.Actions = make([]Action, len(g.Actions))
	copy(clonedGame.Actions, g.Actions)
	clonedGame.positionHistory = make([]uint64, len(g.positionHistory))
	copy(clonedGame.positionHistory, g.positionHistory)
	return clonedGame
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
	IsDraw             bool
	IsPromotion        bool
	IsEnPassantCapture bool
	IsCastle           bool
	IsKingsideCastle   bool
	IsQueensideCastle  bool
	PromotionPieceType PieceType
	CapturedPiece      Piece
}

func (a Action) ICCF() string {
	return fmt.Sprintf("%v%v%v", a.FromPiece.XY.ToICCF(), a.ToXY.ToICCF(), a.PromotionPieceType.ToICCF())
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
	case a.IsDraw:
		return fmt.Sprintf("%s draws", a.FromPiece.Owner)
	case a.IsPromotion:
		return fmt.Sprintf("%s's Pawn at %v promotes to %v", a.FromPiece.Owner, a.FromPiece.XY.ToAlgebraic(), a.PromotionPieceType)
	case a.IsKingsideCastle:
		return fmt.Sprintf("%s kingside castles", a.FromPiece.Owner)
	case a.IsQueensideCastle:
		return fmt.Sprintf("%s queenside castles", a.FromPiece.Owner)
	}
	return fmt.Sprintf("%s's %s at %v moves to %v", a.FromPiece.Owner, a.FromPiece.PieceType, a.FromPiece.XY.ToAlgebraic(), a.ToXY.ToAlgebraic())
}

func (a Action) DebugString() string {
	return fmt.Sprintf("%v at (%v, %v) to (%v, %v), isCapture: %v , isResign: %v , isPromotion: %v , isEnPassantCapture: %v , isCastle: %v , isKingsideCastle: %v , isQueensideCastle: %v, promotionPieceType: %v, capturedPiece: %v at (%v, %v)",
		a.FromPiece.PieceType,
		a.FromPiece.XY.X,
		a.FromPiece.XY.Y,
		a.ToXY.X,
		a.ToXY.Y,
		a.IsCapture,
		a.IsResign,
		a.IsPromotion,
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
	case PieceRook:
		return "2"
	case PieceBishop:
		return "3"
	case PieceKnight:
		return "4"
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

func (t PieceType) ToColorFigurine(c color) string {
	if c == ColorBlack {
		return t.ToFigurine()
	}
	switch t {
	case PieceQueen:
		return "♕"
	case PieceKing:
		return "♔"
	case PieceBishop:
		return "♗"
	case PieceKnight:
		return "♘"
	case PieceRook:
		return "♖"
	case PiecePawn:
		return "♙"
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

func (c XY) ToICCF() string {
	return fmt.Sprintf("%v%v", c.X+1, 8-c.Y)
}

func (c XY) ToDescriptive(turn color, useKt ...bool) string {
	fileNames := []string{"QR", "QN", "QB", "Q", "K", "KB", "KN", "KR"}
	if len(useKt) > 0 && useKt[0] {
		fileNames = []string{"QR", "QKt", "QB", "Q", "K", "KB", "KKt", "KR"}
	}
	y := c.Y + 1
	if turn == ColorWhite {
		y = 8 - c.Y
	}
	return fmt.Sprintf("%v%v", fileNames[c.X], y)
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

type GameStep struct {
	StepString      string
	StepComment     string
	StepAction      Action
	StepGame        Game
	StepPreMoveGame Game
}

func (s GameStep) Clone() GameStep {
	return GameStep{
		StepString:      s.StepString,
		StepComment:     s.StepComment,
		StepAction:      s.StepAction,
		StepGame:        s.StepGame.Clone(),
		StepPreMoveGame: s.StepPreMoveGame.Clone(),
	}
}
