package core

import (
	"errors"
	"fmt"
	"math/bits"
	"strings"
)

type Board struct {
	Board                   []string
	CanWhiteKingsideCastle  bool
	CanWhiteQueensideCastle bool
	CanBlackKingsideCastle  bool
	CanBlackQueensideCastle bool
	HalfMoveClock           int
	FullMoveNumber          int
	EnPassantTargetSquare   string // in Algebraic notation, or empty string
	Turn                    string // "Black" or "White"
}

var (
	errBoardInvalidEnPassantTargetSquare = errors.New("enPassantTargetSquare must be either empty string or valid algebraic notation square e.g. d6")
	errBoardTurnMustBeBlackOrWhite       = errors.New("turn must be either Black or White")
	errBoardDuplicateKing                = errors.New("board has two kings of the same color")
	errBoardKingMissing                  = errors.New("board is missing one of the kings")
	errBoardDimensionsWrong              = errors.New("board dimensions are wrong; should be 8x8")
	errBoardPawnInImpossibleRank         = errors.New("impossible rank for pawn")
	errBoardBlackHasMoreThan16Pieces     = errors.New("black has more than 16 pieces")
	errBoardWhiteHasMoreThan16Pieces     = errors.New("white has more than 16 pieces")
	errBoardSideNotToMoveInCheck         = errors.New("side not to move is in check")
	// TODO check if King is in checkmate that couldn't have been reached
	// TODO don't allow more than 8 pawns of any color
)

func NewDefaultGame() Game {
	g, _ := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return g
}

func NewGameFromBoard(b Board) (Game, error) {
	g := Game{
		CanWhiteCastle:          b.CanWhiteKingsideCastle && b.CanWhiteQueensideCastle,
		CanWhiteKingsideCastle:  b.CanWhiteKingsideCastle,
		CanWhiteQueensideCastle: b.CanWhiteQueensideCastle,
		CanBlackCastle:          b.CanBlackKingsideCastle && b.CanBlackQueensideCastle,
		CanBlackKingsideCastle:  b.CanBlackKingsideCastle,
		CanBlackQueensideCastle: b.CanBlackQueensideCastle,
		FullMoveNumber:          b.FullMoveNumber,
		HalfMoveClock:           b.HalfMoveClock,
		MoveNumber:              (b.FullMoveNumber - 1) * 2,
		kingSq:                  [2]int8{-1, -1},
	}

	// Move number
	if b.Turn != "Black" && b.Turn != "White" {
		return Game{}, errBoardTurnMustBeBlackOrWhite
	}
	if b.Turn == "Black" {
		g.MoveNumber++
	}

	// Pieces and kings
	pieceTypeMap := map[rune]PieceType{
		'♕': PieceQueen,
		'♔': PieceKing,
		'♖': PieceRook,
		'♗': PieceBishop,
		'♘': PieceKnight,
		'♙': PiecePawn,
		'♛': PieceQueen,
		'♚': PieceKing,
		'♜': PieceRook,
		'♝': PieceBishop,
		'♞': PieceKnight,
		'♟': PiecePawn,
	}
	lenY := 0
	for y := range b.Board {
		lenX := 0
		for _, p := range b.Board[y] {
			switch p {
			case '♛', '♚', '♜', '♝', '♞', '♟':
				if p == '♚' && g.kingSq[ColorBlack] >= 0 {
					return Game{}, errBoardDuplicateKing
				}
				g.setSq(ColorBlack, pieceTypeMap[p], sqOf(XY{lenX, lenY}))
			case '♕', '♔', '♖', '♗', '♘', '♙':
				if p == '♔' && g.kingSq[ColorWhite] >= 0 {
					return Game{}, errBoardDuplicateKing
				}
				g.setSq(ColorWhite, pieceTypeMap[p], sqOf(XY{lenX, lenY}))
			default:
			}
			if (p == '♟' || p == '♙') && (lenY == 0 || lenY == 7) {
				return Game{}, errBoardPawnInImpossibleRank
			}
			lenX++
		}
		if lenX != 8 {
			return Game{}, errBoardDimensionsWrong
		}
		lenY++
	}
	if lenY != 8 {
		return Game{}, errBoardDimensionsWrong
	}
	if g.kingSq[ColorBlack] < 0 || g.kingSq[ColorWhite] < 0 {
		return Game{}, errBoardKingMissing
	}
	if bits.OnesCount64(g.occ[ColorBlack]) > 16 {
		return Game{}, errBoardBlackHasMoreThan16Pieces
	}
	if bits.OnesCount64(g.occ[ColorWhite]) > 16 {
		return Game{}, errBoardWhiteHasMoreThan16Pieces
	}

	// Castling auto-correction: rights inconsistent with king/rook placement are
	// silently narrowed rather than rejected (a board editor's "all rights" with a
	// moved king just means no castling).
	g.CanWhiteKingsideCastle, g.CanWhiteQueensideCastle, g.CanBlackKingsideCastle, g.CanBlackQueensideCastle = narrowCastlingRights(
		g,
		g.CanWhiteKingsideCastle, g.CanWhiteQueensideCastle, g.CanBlackKingsideCastle, g.CanBlackQueensideCastle,
	)
	g.CanWhiteCastle = g.CanWhiteKingsideCastle || g.CanWhiteQueensideCastle
	g.CanBlackCastle = g.CanBlackKingsideCastle || g.CanBlackQueensideCastle

	// En passant
	switch {
	case b.EnPassantTargetSquare == "":
		g.IsLastMoveEnPassant = false
	case len(b.EnPassantTargetSquare) != 2,
		(b.EnPassantTargetSquare[1] != '6' && b.EnPassantTargetSquare[1] != '3'),
		(b.EnPassantTargetSquare[0] < 'a' && b.EnPassantTargetSquare[0] > 'h'):
		return Game{}, errBoardInvalidEnPassantTargetSquare
	default:
		g.IsLastMoveEnPassant = true
		g.EnPassantTargetSquare = XY{X: int(b.EnPassantTargetSquare[0] - 'a'), Y: int('8' - b.EnPassantTargetSquare[1])}
	}
	// En passant auto-correction: an impossible e.p. target is silently dropped.
	if g.IsLastMoveEnPassant && g.Turn() == ColorBlack && !g.hasPieceAt(ColorWhite, PiecePawn, g.EnPassantTargetSquare.add(XY{0, -1})) {
		g.IsLastMoveEnPassant = false
		g.EnPassantTargetSquare = XY{}
	}
	if g.IsLastMoveEnPassant && g.Turn() == ColorWhite && !g.hasPieceAt(ColorBlack, PiecePawn, g.EnPassantTargetSquare.add(XY{0, 1})) {
		g.IsLastMoveEnPassant = false
		g.EnPassantTargetSquare = XY{}
	}

	// The side not to move must not be in check: such a position is unreachable
	// and move generation semantics break down.
	sideNotToMove := opponent(g.Turn())
	if g.attackersOf(int(g.kingSq[sideNotToMove]), sideNotToMove) != 0 {
		return Game{}, errBoardSideNotToMoveInCheck
	}

	return g.calculateCriticalFlags(), nil
}

func (g Game) ToBoard() Board {
	pieceTypeMap := []map[PieceType]rune{
		{
			PieceQueen:  '♛',
			PieceKing:   '♚',
			PieceRook:   '♜',
			PieceBishop: '♝',
			PieceKnight: '♞',
			PiecePawn:   '♟',
		},
		{
			PieceQueen:  '♕',
			PieceKing:   '♔',
			PieceRook:   '♖',
			PieceBishop: '♗',
			PieceKnight: '♘',
			PiecePawn:   '♙',
		},
	}
	b := make([]string, 8)
	for y := 0; y < 8; y++ {
		var sb strings.Builder
		for x := 0; x < 8; x++ {
			if p := g.PieceAt(XY{x, y}); p.PieceType != PieceNone {
				sb.WriteRune(pieceTypeMap[p.Owner][p.PieceType])
			} else {
				sb.WriteByte(' ')
			}
		}
		b[y] = sb.String()
	}
	result := Board{
		Board:                   b,
		CanWhiteKingsideCastle:  g.CanWhiteKingsideCastle,
		CanWhiteQueensideCastle: g.CanWhiteQueensideCastle,
		CanBlackKingsideCastle:  g.CanBlackKingsideCastle,
		CanBlackQueensideCastle: g.CanBlackQueensideCastle,
		HalfMoveClock:           g.HalfMoveClock,
		FullMoveNumber:          g.FullMoveNumber,
		EnPassantTargetSquare:   "",
		Turn:                    "Black",
	}
	if g.Turn() == ColorWhite {
		result.Turn = "White"
	}
	if g.IsLastMoveEnPassant {
		result.EnPassantTargetSquare = fmt.Sprintf("%v%v", "abcdefgh"[g.EnPassantTargetSquare.X], 8-g.EnPassantTargetSquare.Y)
	}
	return result
}
