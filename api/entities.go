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
	opponentInCheckBy       []piece
	actions                 []action
}

func (g game) String() string {
	var sb strings.Builder
	for _, s := range g.toBoard().board {
		sb.WriteString(strings.Replace(s, " ", ".", -1))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	return sb.String()
}

var defaultGame = game{
	canWhiteCastle:          true,
	canWhiteKingsideCastle:  true,
	canWhiteQueensideCastle: true,
	canBlackCastle:          true,
	canBlackKingsideCastle:  true,
	canBlackQueensideCastle: true,
	halfMoveClock:           0,
	fullMoveNumber:          0,
	isLastMoveEnPassant:     false,
	enPassantTargetSquare:   xy{},
	moveNumber:              0,
	pieces: []map[xy]piece{
		{{0, 0}: {pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}},
		{{1, 0}: {pieceType: pieceKnight, owner: colorBlack, xy: xy{1, 0}}},
		{{2, 0}: {pieceType: pieceBishop, owner: colorBlack, xy: xy{2, 0}}},
		{{3, 0}: {pieceType: pieceQueen, owner: colorBlack, xy: xy{3, 0}}},
		{{4, 0}: {pieceType: pieceKing, owner: colorBlack, xy: xy{4, 0}}},
		{{5, 0}: {pieceType: pieceBishop, owner: colorBlack, xy: xy{5, 0}}},
		{{6, 0}: {pieceType: pieceKnight, owner: colorBlack, xy: xy{6, 0}}},
		{{7, 0}: {pieceType: piecePawn, owner: colorBlack, xy: xy{7, 0}}},
		{{0, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{0, 1}}},
		{{1, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{1, 1}}},
		{{2, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{2, 1}}},
		{{3, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{3, 1}}},
		{{4, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{4, 1}}},
		{{5, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{5, 1}}},
		{{6, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{6, 1}}},
		{{7, 1}: {pieceType: piecePawn, owner: colorBlack, xy: xy{7, 1}}},
		{{0, 7}: {pieceType: pieceRook, owner: colorWhite, xy: xy{0, 7}}},
		{{1, 7}: {pieceType: pieceKnight, owner: colorWhite, xy: xy{1, 7}}},
		{{2, 7}: {pieceType: pieceBishop, owner: colorWhite, xy: xy{2, 7}}},
		{{3, 7}: {pieceType: pieceQueen, owner: colorWhite, xy: xy{3, 7}}},
		{{4, 7}: {pieceType: pieceKing, owner: colorWhite, xy: xy{4, 7}}},
		{{5, 7}: {pieceType: pieceBishop, owner: colorWhite, xy: xy{5, 7}}},
		{{6, 7}: {pieceType: pieceKnight, owner: colorWhite, xy: xy{6, 7}}},
		{{7, 7}: {pieceType: piecePawn, owner: colorWhite, xy: xy{7, 7}}},
		{{0, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{0, 6}}},
		{{1, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{1, 6}}},
		{{2, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{2, 6}}},
		{{3, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{3, 6}}},
		{{4, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{4, 6}}},
		{{5, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{5, 6}}},
		{{6, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{6, 6}}},
		{{7, 6}: {pieceType: piecePawn, owner: colorWhite, xy: xy{7, 6}}},
	},
	kings: []piece{
		{pieceType: pieceKing, owner: colorBlack, xy: xy{4, 0}},
		{pieceType: pieceKing, owner: colorWhite, xy: xy{4, 7}},
	},
	isCheck:           false,
	isCheckmate:       false,
	isStalemate:       false,
	isDraw:            false,
	isGameOver:        false,
	gameOverWinner:    -1,
	opponentInCheckBy: []piece{},
}

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
	clonedOpponentInCheckBy := make([]piece, len(g.opponentInCheckBy))
	for i, piece := range g.opponentInCheckBy {
		clonedOpponentInCheckBy[i] = piece
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
		opponentInCheckBy:       clonedOpponentInCheckBy,
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
	case a.isCapture:
		return fmt.Sprintf("%s's %s at (%v, %v) captures %s's %s at (%v, %v)", a.fromPiece.owner, a.fromPiece.pieceType, a.fromPiece.xy.x, a.fromPiece.xy.y, a.capturedPiece.owner, a.capturedPiece.pieceType, a.capturedPiece.xy.x, a.capturedPiece.xy.y)
	case a.isResign:
		return fmt.Sprintf("%s resigns", a.fromPiece.owner)
	case a.isPromotion:
		return fmt.Sprintf("%s's Pawn at (%v, %v) promotes to %v", a.fromPiece.owner, a.fromPiece.xy.x, a.fromPiece.xy.y, a.promotionPieceType)
	case a.isEnPassant:
		return fmt.Sprintf("%s Pawn at (%v,%v) does en passant", a.fromPiece.owner, a.fromPiece.xy.x, a.fromPiece.xy.y)
	case a.isEnPassantCapture:
		return fmt.Sprintf("%s's Pawn at (%v, %v) captures %v's Pawn at (%v, %v) who was doing en passant", a.fromPiece.owner, a.fromPiece.xy.x, a.fromPiece.xy.y, a.capturedPiece.owner, a.capturedPiece.xy.x, a.capturedPiece.xy.y)
	case a.isKingsideCastle:
		return fmt.Sprintf("%s kingside castles", a.fromPiece.owner)
	case a.isQueensideCastle:
		return fmt.Sprintf("%s queenside castles", a.fromPiece.owner)
	}
	return fmt.Sprintf("%s's %s at (%v, %v) moves to (%v, %v)", a.fromPiece.owner, a.fromPiece.pieceType, a.fromPiece.xy.x, a.fromPiece.xy.y, a.toXY.x, a.toXY.y)
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
	return fmt.Sprintf("%v's %v at (%v, %v)", p.owner, p.pieceType, p.xy.x, p.xy.y)
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
