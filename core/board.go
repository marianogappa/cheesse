package core

import (
	"errors"
	"fmt"
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
	errBoardInvalidEnPassantTargetSquare   = errors.New("enPassantTargetSquare must be either empty string or valid algebraic notation square e.g. d6")
	errBoardTurnMustBeBlackOrWhite         = errors.New("turn must be either Black or White")
	errBoardDuplicateKing                  = errors.New("board has two kings of the same color")
	errBoardKingMissing                    = errors.New("board is missing one of the kings")
	errBoardDimensionsWrong                = errors.New("board dimensions are wrong; should be 8x8")
	errBoardImpossibleBlackCastle          = errors.New("impossible for black to castle since king has moved")
	errBoardImpossibleBlackQueensideCastle = errors.New("impossible for black to queenside castle since rook has moved")
	errBoardImpossibleBlackKingsideCastle  = errors.New("impossible for black to kingside castle since rook has moved")
	errBoardImpossibleWhiteCastle          = errors.New("impossible for white to castle since king has moved")
	errBoardImpossibleWhiteQueensideCastle = errors.New("impossible for white to queenside castle since rook has moved")
	errBoardImpossibleWhiteKingsideCastle  = errors.New("impossible for white to kingside castle since rook has moved")
	errBoardImpossibleEnPassant            = errors.New("impossible en passant target square, since there's no pawn of the right color next to it")
	errBoardPawnInImpossibleRank           = errors.New("impossible rank for pawn")
	errBoardBlackHasMoreThan16Pieces       = errors.New("black has more than 16 pieces")
	errBoardWhiteHasMoreThan16Pieces       = errors.New("white has more than 16 pieces")
	// TODO check if King is in checkmate that couldn't have been reached
	// TODO don't allow more than 8 pawns of any color
	// TODO check if both are in check
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
		Pieces:                  []map[XY]Piece{{}, {}},
		Kings:                   []Piece{{}, {}},
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
				g.Pieces[ColorBlack][XY{lenX, lenY}] = Piece{PieceType: pieceTypeMap[p], Owner: ColorBlack, XY: XY{lenX, lenY}}
				if p == '♚' && g.Kings[ColorBlack].PieceType != PieceNone {
					return Game{}, errBoardDuplicateKing
				}
				if p == '♚' {
					g.Kings[ColorBlack] = Piece{PieceType: pieceTypeMap[p], Owner: ColorBlack, XY: XY{lenX, lenY}}
				}
			case '♕', '♔', '♖', '♗', '♘', '♙':
				g.Pieces[ColorWhite][XY{lenX, lenY}] = Piece{PieceType: pieceTypeMap[p], Owner: ColorWhite, XY: XY{lenX, lenY}}
				if p == '♔' && g.Kings[ColorWhite].PieceType != PieceNone {
					return Game{}, errBoardDuplicateKing
				}
				if p == '♔' {
					g.Kings[ColorWhite] = Piece{PieceType: pieceTypeMap[p], Owner: ColorWhite, XY: XY{lenX, lenY}}
				}
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
	if g.Kings[ColorBlack].PieceType == PieceNone || g.Kings[ColorWhite].PieceType == PieceNone {
		return Game{}, errBoardKingMissing
	}
	if len(g.Pieces[ColorBlack]) > 16 {
		return Game{}, errBoardBlackHasMoreThan16Pieces
	}
	if len(g.Pieces[ColorWhite]) > 16 {
		return Game{}, errBoardWhiteHasMoreThan16Pieces
	}

	// Castling validation
	if g.CanBlackCastle && !g.Kings[ColorBlack].XY.eq(XY{4, 0}) {
		return Game{}, errBoardImpossibleBlackCastle
	}
	if g.CanBlackQueensideCastle && g.Pieces[ColorBlack][XY{0, 0}].PieceType != PieceRook {
		return Game{}, errBoardImpossibleBlackQueensideCastle
	}
	if g.CanBlackKingsideCastle && g.Pieces[ColorBlack][XY{7, 0}].PieceType != PieceRook {
		return Game{}, errBoardImpossibleBlackKingsideCastle
	}
	if g.CanWhiteCastle && !g.Kings[ColorWhite].XY.eq(XY{4, 7}) {
		return Game{}, errBoardImpossibleWhiteCastle
	}
	if g.CanWhiteQueensideCastle && g.Pieces[ColorWhite][XY{0, 7}].PieceType != PieceRook {
		return Game{}, errBoardImpossibleWhiteQueensideCastle
	}
	if g.CanWhiteKingsideCastle && g.Pieces[ColorWhite][XY{7, 7}].PieceType != PieceRook {
		return Game{}, errBoardImpossibleWhiteKingsideCastle
	}

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
	if g.IsLastMoveEnPassant && g.Turn() == ColorBlack && g.Pieces[ColorWhite][g.EnPassantTargetSquare.add(XY{0, -1})].PieceType != PiecePawn {
		return Game{}, errBoardImpossibleEnPassant
	}
	if g.IsLastMoveEnPassant && g.Turn() == ColorWhite && g.Pieces[ColorBlack][g.EnPassantTargetSquare.add(XY{0, 1})].PieceType != PiecePawn {
		return Game{}, errBoardImpossibleEnPassant
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
			bp, bpExists := g.Pieces[ColorBlack][XY{x, y}]
			wp, wpExists := g.Pieces[ColorWhite][XY{x, y}]
			switch {
			case bpExists:
				sb.WriteRune(pieceTypeMap[ColorBlack][bp.PieceType])
			case wpExists:
				sb.WriteRune(pieceTypeMap[ColorWhite][wp.PieceType])
			default:
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
