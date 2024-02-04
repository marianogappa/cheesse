package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO test way more cases
func TestThreats(t *testing.T) {
	ts := []struct {
		name              string
		board             Board
		xy                XY
		owner             color
		threateningPieces []Piece
	}{
		{
			name: "white king is not threatened",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"  ♔     ",
					"  ♙     ",
					"        ",
					"♙♙  ♙♙♙♙",
					"♖♘♗♕ ♗♘♖",
				},
				Turn: "White",
			},
			owner:             ColorWhite,
			xy:                XY{2, 3},
			threateningPieces: []Piece{},
		},
		{
			name: "white king is threatened by black queen",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟ ♟♟♟♟",
					"        ",
					"        ",
					"  ♙♔    ",
					"        ",
					"♙♙  ♙♙♙♙",
					"♖♘♗♕ ♗♘♖",
				},
				Turn: "White",
			},
			owner:             ColorWhite,
			xy:                XY{3, 4},
			threateningPieces: []Piece{{PieceType: PieceQueen, Owner: ColorBlack, XY: XY{3, 0}}},
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromBoard(tc.board)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.threateningPieces, g.xyThreatenedBy(tc.xy, tc.owner, true))
		})
	}
}

func TestCheckmate(t *testing.T) {
	ts := []struct {
		name   string
		board  Board
		winner color
	}{
		{
			name: "black does checkmate",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞ ",
					"♟♟♟♟♟♟♟ ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙ ",
					"♖♘♗♕♔  ♜",
				},
				Turn: "White",
			},
			winner: ColorBlack,
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromBoard(tc.board)
			require.NoError(t, err)
			assert.True(t, g.IsCheckmate)
			assert.True(t, g.IsGameOver)
			assert.Equal(t, tc.winner, g.GameOverWinner)
		})
	}
}
