package core

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	errFENRegexDoesNotMatch        = errors.New("FEN string does not match FEN regexp")
	errFENRankLargerThan8Squares   = errors.New("FEN string has a rank larger than 8 squares")
	errFENDuplicateKing            = errors.New("FEN string has more than one king of the same color")
	errFENKingMissing              = errors.New("FEN string is lacking one of the kings")
	errFENPawnInImpossibleRank     = errors.New("impossible rank for pawn")
	errFENBlackHasMoreThan16Pieces = errors.New("black has more than 16 pieces")
	errFENWhiteHasMoreThan16Pieces = errors.New("white has more than 16 pieces")
	errFENSideNotToMoveInCheck     = errors.New("side not to move is in check")
	// TODO check if King is in checkmate that couldn't have been reached
	// TODO don't allow more than 8 pawns of any color
)

func NewGameFromFEN(s string) (Game, error) {
	rxFEN := regexp.MustCompile(`^([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8})\/([1-8rnbqkpRNBQKP]{1,8}) ([wb]) ([KQkq]{0,4}|-) ([a-h][36]|-) ([0-9]{1,3}) ([0-9]{1,3})$`)
	matches := rxFEN.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return Game{}, errFENRegexDoesNotMatch
	}

	fullMoveNumber := atoi(matches[0][13]) // The regex cannot pass a non-number here
	halfMoveClock := atoi(matches[0][12])  // The regex cannot pass a non-number here

	// moveNumber calculation
	moveNumber := 0
	turn := matches[0][9]
	switch turn { // The regex cannot pass a different case here
	case "w":
		moveNumber = 2 * (fullMoveNumber - 1)
	case "b":
		moveNumber = 2*(fullMoveNumber-1) + 1
	}

	// En passant calculation
	isLastMoveEnPassant := false
	enPassantTargetSquare := XY{}
	if matches[0][11] != "-" { // The regex cannot pass a different string here
		isLastMoveEnPassant = true
		enPassantTargetSquare = XY{int(matches[0][11][0] - 'a'), int('8' - matches[0][11][1])}
	}

	// Castling calculation
	castlingMap := map[byte]bool{}
	for i := 0; i < len(matches[0][10]); i++ {
		castlingMap[matches[0][10][i]] = true
	}
	canWhiteKingsideCastle := castlingMap['K']
	canWhiteQueensideCastle := castlingMap['Q']
	canWhiteCastle := canWhiteKingsideCastle || canWhiteQueensideCastle
	canBlackKingsideCastle := castlingMap['k']
	canBlackQueensideCastle := castlingMap['q']
	canBlackCastle := canBlackKingsideCastle || canBlackQueensideCastle

	// Pieces and kings calculation
	pieceTypeMap := map[byte]PieceType{'Q': PieceQueen, 'K': PieceKing, 'B': PieceBishop, 'N': PieceKnight, 'R': PieceRook, 'P': PiecePawn}
	pieces := []map[XY]Piece{{}, {}}
	kings := []Piece{{}, {}}
	for y, row := range []string{matches[0][1], matches[0][2], matches[0][3], matches[0][4], matches[0][5], matches[0][6], matches[0][7], matches[0][8]} {
		x := 0
		for i := 0; i < len(row); i++ {
			b := row[i]
			if x >= 8 {
				return Game{}, errFENRankLargerThan8Squares
			}
			switch b {
			case 'Q', 'K', 'B', 'N', 'R', 'P':
				pieces[ColorWhite][XY{x, y}] = Piece{PieceType: pieceTypeMap[b], Owner: ColorWhite, XY: XY{x, y}}
				x++
			case 'q', 'k', 'b', 'n', 'r', 'p':
				pieces[ColorBlack][XY{x, y}] = Piece{PieceType: pieceTypeMap[b-'a'+'A'], Owner: ColorBlack, XY: XY{x, y}}
				x++
			case '1', '2', '3', '4', '5', '6', '7', '8':
				x += int(b - '0')
			}
			switch {
			case b == 'k' && kings[ColorBlack].PieceType == PieceKing, b == 'K' && kings[ColorWhite].PieceType == PieceKing:
				return Game{}, errFENDuplicateKing
			case b == 'K':
				kings[ColorWhite] = pieces[ColorWhite][XY{x - 1, y}]
			case b == 'k':
				kings[ColorBlack] = pieces[ColorBlack][XY{x - 1, y}]
			case (b == 'p' || b == 'P') && (y == 0 || y == 7):
				return Game{}, errFENPawnInImpossibleRank
			}
		}
	}
	if kings[ColorBlack].PieceType == PieceNone || kings[ColorWhite].PieceType == PieceNone {
		return Game{}, errFENKingMissing
	}
	if len(pieces[ColorBlack]) > 16 {
		return Game{}, errFENBlackHasMoreThan16Pieces
	}
	if len(pieces[ColorWhite]) > 16 {
		return Game{}, errFENWhiteHasMoreThan16Pieces
	}

	// En passant auto-correction: an impossible e.p. target (no pawn of the right
	// color next to it) is silently dropped rather than rejected, since board
	// editors construct FENs naively.
	if isLastMoveEnPassant && turn == "b" && pieces[ColorWhite][enPassantTargetSquare.add(XY{0, -1})].PieceType != PiecePawn {
		isLastMoveEnPassant = false
		enPassantTargetSquare = XY{}
	}
	if isLastMoveEnPassant && turn == "w" && pieces[ColorBlack][enPassantTargetSquare.add(XY{0, 1})].PieceType != PiecePawn {
		isLastMoveEnPassant = false
		enPassantTargetSquare = XY{}
	}

	// Castling auto-correction: rights inconsistent with king/rook placement are
	// silently narrowed rather than rejected (a board editor's "KQkq" with a moved
	// king just means no castling).
	canWhiteKingsideCastle, canWhiteQueensideCastle, canBlackKingsideCastle, canBlackQueensideCastle = narrowCastlingRights(
		pieces, kings,
		canWhiteKingsideCastle, canWhiteQueensideCastle, canBlackKingsideCastle, canBlackQueensideCastle,
	)
	canWhiteCastle = canWhiteKingsideCastle || canWhiteQueensideCastle
	canBlackCastle = canBlackKingsideCastle || canBlackQueensideCastle

	game := Game{
		CanWhiteCastle:          canWhiteCastle,
		CanWhiteKingsideCastle:  canWhiteKingsideCastle,
		CanWhiteQueensideCastle: canWhiteQueensideCastle,
		CanBlackCastle:          canBlackCastle,
		CanBlackKingsideCastle:  canBlackKingsideCastle,
		CanBlackQueensideCastle: canBlackQueensideCastle,
		HalfMoveClock:           halfMoveClock,
		FullMoveNumber:          fullMoveNumber,
		IsLastMoveEnPassant:     isLastMoveEnPassant,
		EnPassantTargetSquare:   enPassantTargetSquare,
		MoveNumber:              moveNumber,
		Pieces:                  pieces,
		Kings:                   kings,
	}

	// The side not to move must not be in check: such a position is unreachable
	// and move generation semantics break down (the opponent's king is capturable).
	if len(game.xyThreatenedBy(kings[opponentOfTurn(turn)].XY, opponentOfTurn(turn), false)) > 0 {
		return Game{}, errFENSideNotToMoveInCheck
	}

	return game.calculateCriticalFlags(), nil
}

func opponentOfTurn(turn string) color {
	if turn == "w" {
		return ColorBlack
	}
	return ColorWhite
}

// narrowCastlingRights drops any castling right that is inconsistent with the actual
// king and rook placement.
func narrowCastlingRights(pieces []map[XY]Piece, kings []Piece, wk, wq, bk, bq bool) (bool, bool, bool, bool) {
	if !kings[ColorWhite].XY.eq(XY{4, 7}) {
		wk, wq = false, false
	}
	if !kings[ColorBlack].XY.eq(XY{4, 0}) {
		bk, bq = false, false
	}
	if pieces[ColorWhite][XY{7, 7}].PieceType != PieceRook {
		wk = false
	}
	if pieces[ColorWhite][XY{0, 7}].PieceType != PieceRook {
		wq = false
	}
	if pieces[ColorBlack][XY{7, 0}].PieceType != PieceRook {
		bk = false
	}
	if pieces[ColorBlack][XY{0, 0}].PieceType != PieceRook {
		bq = false
	}
	return wk, wq, bk, bq
}

// This assumes that Atoi can't fail because the regex capture cannot return a non-number
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func (g Game) ToFEN() string {
	var sb strings.Builder
	pieceTypeMap := map[PieceType]byte{PieceQueen: 'Q', PieceKing: 'K', PieceBishop: 'B', PieceKnight: 'N', PieceRook: 'R', PiecePawn: 'P'}
	for y := 0; y < 8; y++ {
		count := 0
		for x := 0; x < 8; x++ {
			bp, bExists := g.Pieces[ColorBlack][XY{x, y}]
			wp, wExists := g.Pieces[ColorWhite][XY{x, y}]
			switch {
			case !bExists && !wExists:
				count++
			case bExists:
				if count > 0 {
					sb.WriteString(fmt.Sprintf("%v", count))
				}
				count = 0
				sb.WriteString(strings.ToLower(string(pieceTypeMap[bp.PieceType])))
			case wExists:
				if count > 0 {
					sb.WriteString(fmt.Sprintf("%v", count))
				}
				count = 0
				sb.WriteByte(pieceTypeMap[wp.PieceType])
			}
		}
		if count > 0 {
			sb.WriteString(fmt.Sprintf("%v", count))
		}
		if y < 7 {
			sb.WriteByte('/')
		}
	}

	turn := "b"
	if g.Turn() == ColorWhite {
		turn = "w"
	}

	var castlingSB strings.Builder
	if g.CanWhiteKingsideCastle {
		castlingSB.WriteByte('K')
	}
	if g.CanWhiteQueensideCastle {
		castlingSB.WriteByte('Q')
	}
	if g.CanBlackKingsideCastle {
		castlingSB.WriteByte('k')
	}
	if g.CanBlackQueensideCastle {
		castlingSB.WriteByte('q')
	}
	castling := castlingSB.String()
	if castling == "" {
		castling = "-"
	}

	enPassant := "-"
	if g.IsLastMoveEnPassant {
		enPassant = fmt.Sprintf("%v%v", string("abcdefgh"[g.EnPassantTargetSquare.X]), 8-g.EnPassantTargetSquare.Y)
	}

	sb.WriteString(fmt.Sprintf(" %v %v %v %v %v", turn, castling, enPassant, g.HalfMoveClock, g.FullMoveNumber))

	return sb.String()
}
