package api

import (
	"fmt"
	"strings"
)

type board struct {
	board                   []string
	canWhiteKingsideCastle  bool
	canWhiteQueensideCastle bool
	canBlackKingsideCastle  bool
	canBlackQueensideCastle bool
	halfMoveClock           int
	fullMoveNumber          int
	enPassantTargetSquare   string // in Algebraic notation, or empty string
	turn                    string // "Black" or "White"
}

var (
	errBoardInvalidEnPassantTargetSquare   = fmt.Errorf("enPassantTargetSquare must be either empty string or valid algebraic notation square e.g. d6")
	errBoardTurnMustBeBlackOrWhite         = fmt.Errorf("turn must be either Black or White")
	errBoardDuplicateKing                  = fmt.Errorf("board has two kings of the same color")
	errBoardKingMissing                    = fmt.Errorf("board is missing one of the kings")
	errBoardDimensionsWrong                = fmt.Errorf("board dimensions are wrong; should be 8x8")
	errBoardImpossibleBlackCastle          = fmt.Errorf("impossible for black to castle since king has moved")
	errBoardImpossibleBlackQueensideCastle = fmt.Errorf("impossible for black to queenside castle since rook has moved")
	errBoardImpossibleBlackKingsideCastle  = fmt.Errorf("impossible for black to kingside castle since rook has moved")
	errBoardImpossibleWhiteCastle          = fmt.Errorf("impossible for white to castle since king has moved")
	errBoardImpossibleWhiteQueensideCastle = fmt.Errorf("impossible for white to queenside castle since rook has moved")
	errBoardImpossibleWhiteKingsideCastle  = fmt.Errorf("impossible for white to kingside castle since rook has moved")
	errBoardImpossibleEnPassant            = fmt.Errorf("impossible en passant target square, since there's no pawn of the right color next to it")
	errBoardPawnInImpossibleRank           = fmt.Errorf("impossible rank for pawn")
	errBoardBlackHasMoreThan16Pieces       = fmt.Errorf("black has more than 16 pieces")
	errBoardWhiteHasMoreThan16Pieces       = fmt.Errorf("white has more than 16 pieces")
	// TODO check if King is in checkmate that couldn't have been reached
	// TODO don't allow more than 8 pawns of any color
)

func newGameFromBoard(b board) (game, error) {
	g := game{
		canWhiteCastle:          b.canWhiteKingsideCastle && b.canWhiteQueensideCastle,
		canWhiteKingsideCastle:  b.canWhiteKingsideCastle,
		canWhiteQueensideCastle: b.canWhiteQueensideCastle,
		canBlackCastle:          b.canBlackKingsideCastle && b.canBlackQueensideCastle,
		canBlackKingsideCastle:  b.canBlackKingsideCastle,
		canBlackQueensideCastle: b.canBlackQueensideCastle,
		fullMoveNumber:          b.fullMoveNumber,
		halfMoveClock:           b.halfMoveClock,
		moveNumber:              (b.fullMoveNumber - 1) * 2,
		pieces:                  []map[xy]piece{{}, {}},
		kings:                   []piece{{}, {}},
	}

	// Move number
	if b.turn != "Black" && b.turn != "White" {
		return game{}, errBoardTurnMustBeBlackOrWhite
	}
	if b.turn == "Black" {
		g.moveNumber++
	}

	// Pieces and kings
	pieceTypeMap := map[rune]pieceType{
		'♕': pieceQueen,
		'♔': pieceKing,
		'♖': pieceRook,
		'♗': pieceBishop,
		'♘': pieceKnight,
		'♙': piecePawn,
		'♛': pieceQueen,
		'♚': pieceKing,
		'♜': pieceRook,
		'♝': pieceBishop,
		'♞': pieceKnight,
		'♟': piecePawn,
	}
	lenY := 0
	for y := range b.board {
		lenX := 0
		for _, p := range b.board[y] {
			switch p {
			case '♛', '♚', '♜', '♝', '♞', '♟':
				g.pieces[colorBlack][xy{lenX, lenY}] = piece{pieceType: pieceTypeMap[p], owner: colorBlack, xy: xy{lenX, lenY}}
				if p == '♚' && g.kings[colorBlack].pieceType != pieceNone {
					return game{}, errBoardDuplicateKing
				}
				if p == '♚' {
					g.kings[colorBlack] = piece{pieceType: pieceTypeMap[p], owner: colorBlack, xy: xy{lenX, lenY}}
				}
			case '♕', '♔', '♖', '♗', '♘', '♙':
				g.pieces[colorWhite][xy{lenX, lenY}] = piece{pieceType: pieceTypeMap[p], owner: colorWhite, xy: xy{lenX, lenY}}
				if p == '♔' && g.kings[colorWhite].pieceType != pieceNone {
					return game{}, errBoardDuplicateKing
				}
				if p == '♔' {
					g.kings[colorWhite] = piece{pieceType: pieceTypeMap[p], owner: colorWhite, xy: xy{lenX, lenY}}
				}
			default:
			}
			if (p == '♟' || p == '♙') && (lenY == 0 || lenY == 7) {
				return game{}, errBoardPawnInImpossibleRank
			}
			lenX++
		}
		if lenX != 8 {
			return game{}, errBoardDimensionsWrong
		}
		lenY++
	}
	if lenY != 8 {
		return game{}, errBoardDimensionsWrong
	}
	if g.kings[colorBlack].pieceType == pieceNone || g.kings[colorWhite].pieceType == pieceNone {
		return game{}, errBoardKingMissing
	}
	if len(g.pieces[colorBlack]) > 16 {
		return game{}, errBoardBlackHasMoreThan16Pieces
	}
	if len(g.pieces[colorWhite]) > 16 {
		return game{}, errBoardWhiteHasMoreThan16Pieces
	}

	// Castling validation
	if g.canBlackCastle && !g.kings[colorBlack].xy.eq(xy{4, 0}) {
		return game{}, errBoardImpossibleBlackCastle
	}
	if g.canBlackQueensideCastle && g.pieces[colorBlack][xy{0, 0}].pieceType != pieceRook {
		return game{}, errBoardImpossibleBlackQueensideCastle
	}
	if g.canBlackKingsideCastle && g.pieces[colorBlack][xy{7, 0}].pieceType != pieceRook {
		return game{}, errBoardImpossibleBlackKingsideCastle
	}
	if g.canWhiteCastle && !g.kings[colorWhite].xy.eq(xy{4, 7}) {
		return game{}, errBoardImpossibleWhiteCastle
	}
	if g.canWhiteQueensideCastle && g.pieces[colorWhite][xy{0, 7}].pieceType != pieceRook {
		return game{}, errBoardImpossibleWhiteQueensideCastle
	}
	if g.canWhiteKingsideCastle && g.pieces[colorWhite][xy{7, 7}].pieceType != pieceRook {
		return game{}, errBoardImpossibleWhiteKingsideCastle
	}

	// En passant
	switch {
	case b.enPassantTargetSquare == "":
		g.isLastMoveEnPassant = false
	case len(b.enPassantTargetSquare) != 2,
		(b.enPassantTargetSquare[1] != '6' && b.enPassantTargetSquare[1] != '3'),
		(b.enPassantTargetSquare[0] < 'a' && b.enPassantTargetSquare[0] > 'h'):
		return game{}, errBoardInvalidEnPassantTargetSquare
	default:
		g.isLastMoveEnPassant = true
		g.enPassantTargetSquare = xy{x: int(b.enPassantTargetSquare[0] - 'a'), y: int('8' - b.enPassantTargetSquare[1])}
	}
	if g.isLastMoveEnPassant && g.turn() == colorBlack && g.pieces[colorWhite][g.enPassantTargetSquare.add(xy{0, -1})].pieceType != piecePawn {
		return game{}, errBoardImpossibleEnPassant
	}
	if g.isLastMoveEnPassant && g.turn() == colorWhite && g.pieces[colorBlack][g.enPassantTargetSquare.add(xy{0, 1})].pieceType != piecePawn {
		return game{}, errBoardImpossibleEnPassant
	}

	return g.calculateCriticalFlags(), nil
}

func (g game) toBoard() board {
	pieceTypeMap := []map[pieceType]rune{
		{
			pieceQueen:  '♛',
			pieceKing:   '♚',
			pieceRook:   '♜',
			pieceBishop: '♝',
			pieceKnight: '♞',
			piecePawn:   '♟',
		},
		{
			pieceQueen:  '♕',
			pieceKing:   '♔',
			pieceRook:   '♖',
			pieceBishop: '♗',
			pieceKnight: '♘',
			piecePawn:   '♙',
		},
	}
	b := make([]string, 8)
	for y := 0; y < 8; y++ {
		var sb strings.Builder
		for x := 0; x < 8; x++ {
			bp, bpExists := g.pieces[colorBlack][xy{x, y}]
			wp, wpExists := g.pieces[colorWhite][xy{x, y}]
			switch {
			case bpExists:
				sb.WriteRune(pieceTypeMap[colorBlack][bp.pieceType])
			case wpExists:
				sb.WriteRune(pieceTypeMap[colorWhite][wp.pieceType])
			default:
				sb.WriteByte(' ')
			}
		}
		b[y] = sb.String()
	}
	result := board{
		board:                   b,
		canWhiteKingsideCastle:  g.canWhiteKingsideCastle,
		canWhiteQueensideCastle: g.canWhiteQueensideCastle,
		canBlackKingsideCastle:  g.canBlackKingsideCastle,
		canBlackQueensideCastle: g.canBlackQueensideCastle,
		halfMoveClock:           g.halfMoveClock,
		fullMoveNumber:          g.fullMoveNumber,
		enPassantTargetSquare:   "",
		turn:                    "Black",
	}
	if g.turn() == colorWhite {
		result.turn = "White"
	}
	if g.isLastMoveEnPassant {
		result.enPassantTargetSquare = fmt.Sprintf("%v%v", "abcdefgh"[g.enPassantTargetSquare.x], 8-g.enPassantTargetSquare.y)
	}
	return result
}
