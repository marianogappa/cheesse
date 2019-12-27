package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO test way more cases
func TestThreats(t *testing.T) {
	ts := []struct {
		name              string
		board             board
		xy                xy
		owner             color
		threateningPieces []piece
	}{
		{
			name: "white king is not threatened",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"  ♔     ",
					"  ♙     ",
					"        ",
					"♙♙  ♙♙♙♙",
					"♖♘♗♕ ♗♘♖",
				},
				turn: 'w',
			},
			owner:             colorWhite,
			xy:                xy{2, 3},
			threateningPieces: []piece{},
		},
		{
			name: "white king is threatened by black queen",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟ ♟♟♟♟",
					"        ",
					"        ",
					"  ♙♔    ",
					"        ",
					"♙♙  ♙♙♙♙",
					"♖♘♗♕ ♗♘♖",
				},
				turn: 'w',
			},
			owner:             colorWhite,
			xy:                xy{3, 4},
			threateningPieces: []piece{{pieceType: pieceQueen, owner: colorBlack, xy: xy{3, 0}}},
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := newGameFromBoard(tc.board)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.threateningPieces, g.xyThreatenedBy(tc.xy, tc.owner, true))
		})
	}
}
