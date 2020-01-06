package api

import (
	"fmt"
	"strings"
)

type game struct {
	canWhiteCastle          bool
	canWhiteKingsideCastle  bool
	canWhiteQueensideCastle bool
	canBlackCastle          bool
	canBlackKingsideCastle  bool
	canBlackQueensideCastle bool
	halfMoveClock           int
	fullMoveNumber          int
	isLastMoveEnPassant     bool
	enPassantTargetSquare   xy
	moveNumber              int
	pieces                  []map[xy]piece
	kings                   []piece
	isCheck                 bool
	isCheckmate             bool
	isStalemate             bool
	isDraw                  bool
	isGameOver              bool
	gameOverWinner          color
	inCheckBy               []piece
	actions                 []action
}

func (g game) String() string {
	var sb strings.Builder
	for _, s := range g.toBoard().board {
		sb.WriteString(strings.Replace(s, " ", ".", -1))
		sb.WriteString("\n")
	}
	return sb.String()
}

var defaultGame, _ = newGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

func (g game) clone() game {
	clonedPieces := make([]map[xy]piece, len(g.pieces))
	clonedKings := make([]piece, len(g.kings))
	for color, ownerPieces := range g.pieces {
		clonedOwnerPieces := make(map[xy]piece, len(ownerPieces))
		for _, piece := range ownerPieces {
			clonedPiece := piece
			clonedOwnerPieces[piece.xy] = clonedPiece
			if clonedPiece.pieceType == pieceKing {
				clonedKings[color] = clonedPiece
			}
		}
		clonedPieces[color] = clonedOwnerPieces
	}
	clonedInCheckBy := make([]piece, len(g.inCheckBy))
	for i, piece := range g.inCheckBy {
		clonedInCheckBy[i] = piece
	}
	return game{
		canWhiteCastle:          g.canWhiteCastle,
		canWhiteKingsideCastle:  g.canWhiteKingsideCastle,
		canWhiteQueensideCastle: g.canWhiteQueensideCastle,
		canBlackCastle:          g.canBlackCastle,
		canBlackKingsideCastle:  g.canBlackKingsideCastle,
		canBlackQueensideCastle: g.canBlackQueensideCastle,
		halfMoveClock:           g.halfMoveClock,
		fullMoveNumber:          g.fullMoveNumber,
		enPassantTargetSquare:   g.enPassantTargetSquare,
		moveNumber:              g.moveNumber,
		pieces:                  clonedPieces,
		kings:                   clonedKings,
		isCheck:                 g.isCheck,
		isCheckmate:             g.isCheckmate,
		isStalemate:             g.isStalemate,
		isDraw:                  g.isDraw,
		isGameOver:              g.isGameOver,
		gameOverWinner:          g.gameOverWinner,
		inCheckBy:               clonedInCheckBy,
	}
}

type castleType int

const (
	castleTypeQueenside = iota
	castleTypeKingside
)

type action struct {
	fromPiece          piece
	toXY               xy
	isCapture          bool
	isResign           bool
	isPromotion        bool
	isEnPassant        bool
	isEnPassantCapture bool
	isCastle           bool
	isKingsideCastle   bool
	isQueensideCastle  bool
	promotionPieceType pieceType
	capturedPiece      piece
}

func (a action) String() string {
	switch {
	case a.isEnPassantCapture:
		return fmt.Sprintf("%s's Pawn at %v captures %v's Pawn at %v which was doing en passant", a.fromPiece.owner, a.fromPiece.xy.toAlgebraic(), a.capturedPiece.owner, a.capturedPiece.xy.toAlgebraic())
	case a.isCapture:
		return fmt.Sprintf("%s's %s at %v captures %s's %s at %v", a.fromPiece.owner, a.fromPiece.pieceType, a.fromPiece.xy.toAlgebraic(), a.capturedPiece.owner, a.capturedPiece.pieceType, a.capturedPiece.xy.toAlgebraic())
	case a.isResign:
		return fmt.Sprintf("%s resigns", a.fromPiece.owner)
	case a.isPromotion:
		return fmt.Sprintf("%s's Pawn at %v promotes to %v", a.fromPiece.owner, a.fromPiece.xy.toAlgebraic(), a.promotionPieceType)
	case a.isEnPassant:
		return fmt.Sprintf("%s's Pawn at %v does en passant", a.fromPiece.owner, a.fromPiece.xy.toAlgebraic())
	case a.isKingsideCastle:
		return fmt.Sprintf("%s kingside castles", a.fromPiece.owner)
	case a.isQueensideCastle:
		return fmt.Sprintf("%s queenside castles", a.fromPiece.owner)
	}
	return fmt.Sprintf("%s's %s at %v moves to %v", a.fromPiece.owner, a.fromPiece.pieceType, a.fromPiece.xy.toAlgebraic(), a.toXY.toAlgebraic())
}

func (a action) DebugString() string {
	return fmt.Sprintf("%v at (%v, %v) to (%v, %v), isCapture: %v , isResign: %v , isPromotion: %v , isEnPassant: %v , isEnPassantCapture: %v , isCastle: %v , isKingsideCastle: %v , isQueensideCastle: %v, promotionPieceType: %v, capturedPiece: %v at (%v, %v)",
		a.fromPiece.pieceType,
		a.fromPiece.xy.x,
		a.fromPiece.xy.y,
		a.toXY.x,
		a.toXY.y,
		a.isCapture,
		a.isResign,
		a.isPromotion,
		a.isEnPassant,
		a.isEnPassantCapture,
		a.isCastle,
		a.isKingsideCastle,
		a.isQueensideCastle,
		a.promotionPieceType,
		a.capturedPiece.pieceType,
		a.capturedPiece.xy.x,
		a.capturedPiece.xy.y,
	)
}

type pieceType int

const (
	pieceNone = iota
	pieceQueen
	pieceKing
	pieceBishop
	pieceKnight
	pieceRook
	piecePawn
)

func (t pieceType) String() string {
	switch t {
	case pieceQueen:
		return "Queen"
	case pieceKing:
		return "King"
	case pieceBishop:
		return "Bishop"
	case pieceKnight:
		return "Knight"
	case pieceRook:
		return "Rook"
	case piecePawn:
		return "Pawn"
	}
	return ""
}

type color int

const (
	colorBlack = iota
	colorWhite
)

func (t color) String() string {
	switch t {
	case colorBlack:
		return "Black"
	case colorWhite:
		return "White"
	}
	return "Unknown"
}

type xy struct {
	x, y int
}

func (c xy) add(c2 xy) xy {
	return xy{c.x + c2.x, c.y + c2.y}
}

func (c xy) eq(c2 xy) bool {
	return c.x == c2.x && c.y == c2.y
}

func (c xy) toAlgebraic() string {
	return fmt.Sprintf("%v%v", string("abcdefgh"[c.x]), 8-c.y)
}

func (c xy) deltaTowards(c2 xy) xy {
	var delta xy
	switch {
	case c2.x < c.x:
		delta.x = -1
	case c2.x > c.x:
		delta.x = 1
	}
	switch {
	case c2.y < c.y:
		delta.y = -1
	case c2.y > c.y:
		delta.y = 1
	}
	return delta
}

func (p piece) xysTowards(sq xy) []xy {
	if p.pieceType != pieceQueen && p.pieceType != pieceBishop && p.pieceType != pieceRook {
		return []xy{}
	}
	var (
		rookLike   = p.xy.x == sq.x || p.xy.y == sq.y
		absDist    = xy{abs(p.xy.x - sq.x), abs(p.xy.y - sq.y)}
		bishopLike = (absDist.x - absDist.y) == 0
	)
	if (p.pieceType == pieceRook && !rookLike) || (p.pieceType == pieceBishop && !bishopLike) || (p.pieceType == pieceQueen && !rookLike && !bishopLike) {
		return []xy{}
	}
	var (
		delta = p.xy.deltaTowards(sq)
		cur   = p.xy.add(delta)
		xys   = []xy{}
	)
	for cur != sq {
		xys = append(xys, cur)
		cur = cur.add(delta)
	}
	return xys
}

type piece struct {
	pieceType pieceType
	owner     color
	xy        xy
}

func (p piece) String() string {
	return fmt.Sprintf("%v's %v at %v", p.owner, p.pieceType, p.xy.toAlgebraic())
}

var movementDeltasByPieceType = map[pieceType][]xy{
	pieceQueen:  {{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}},
	pieceKing:   {{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}, {-2, 0}, {2, 0}}, // Castling
	pieceBishop: {{-1, -1}, {1, 1}, {-1, 1}, {1, -1}},
	pieceKnight: {{-1, -2}, {1, -2}, {-1, 2}, {1, 2}, {-2, -1}, {-2, 1}, {2, -1}, {2, 1}},
	pieceRook:   {{-1, 0}, {1, 0}, {0, -1}, {0, 1}},
	// N.B. Pawn will be dealt with separately, because it's dependant on color
}

var (
	errNotInBounds                = fmt.Errorf("piece is not in board bounds")
	errFriendlyPieceInDestination = fmt.Errorf("there is a friendly piece in destination")
	errPieceInBetween             = fmt.Errorf("there is a piece between source and destination")
	errPieceBlockingPawn          = fmt.Errorf("there is a piece blocking pawn from moving forwards or en passant")
	errPawnCantCapture            = fmt.Errorf("pawn can't capture because there is no opponent piece (including en passant)")
	errCantCastle                 = fmt.Errorf("king can't castle, because pieces moved, pieces in middle or squares threatened")
	errCantPromote                = fmt.Errorf("pawn can't promote because wrong position or invalid promotion piece type")
	errActionLeavesKingThreatened = fmt.Errorf("action leaves king in a check")
)

var (
	emptyXYsForCastlingByColorAndCastleType = map[color]map[castleType][]xy{
		colorBlack: {
			castleTypeQueenside: {{1, 0}, {2, 0}, {3, 0}},
			castleTypeKingside:  {{5, 0}, {6, 0}},
		},
		colorWhite: {
			castleTypeQueenside: {{1, 7}, {2, 7}, {3, 7}},
			castleTypeKingside:  {{5, 7}, {6, 7}},
		},
	}
	unthreatenedXYsForCastlingByColorAndCastleType = map[color]map[castleType][]xy{
		colorBlack: {
			castleTypeQueenside: {{1, 0}, {2, 0}, {3, 0}, {4, 0}},
			castleTypeKingside:  {{4, 0}, {5, 0}, {6, 0}},
		},
		colorWhite: {
			castleTypeQueenside: {{1, 7}, {2, 7}, {3, 7}, {4, 7}},
			castleTypeKingside:  {{4, 7}, {5, 7}, {6, 7}},
		},
	}
)
