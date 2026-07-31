package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ported from ostinato's CastlingTest.scala
func TestOstinatoCastling(t *testing.T) {
	ts := []struct {
		name       string
		fen        string
		hasCastle  bool
		castleSide string
	}{
		{"black king can kingside castle", "4k2r/8/8/8/8/8/8/4K3 b k - 0 1", true, "kingside"},
		{"black can't castle: white's turn", "4k2r/8/8/8/8/8/8/4K3 w k - 0 1", false, ""},
		{"black can't castle: king not on initial square", "3k3r/8/8/8/8/8/8/4K3 b k - 0 1", false, ""},
		{"black can't castle: rook not on initial square", "4k1r1/8/8/8/8/8/8/4K3 b - - 0 1", false, ""},
		{"black can't castle: king is threatened", "4k2r/8/8/8/8/8/4R3/4K3 b k - 0 1", false, ""},
		{"black can't castle: pass-through square f8 threatened", "4k2r/8/8/8/8/8/5R2/4K3 b k - 0 1", false, ""},
		{"black can't long castle: pass-through square d8 threatened", "r3k3/8/8/8/8/8/3R4/4K3 b q - 0 1", false, ""},
		{"black can long castle", "r3k3/8/8/8/8/8/5R2/4K3 b q - 0 1", true, "queenside"},
		{"blocked long castle: knight on b8", "rn2k3/8/8/8/8/8/8/4K3 b q - 0 1", false, ""},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			found := false
			for _, a := range g.Actions {
				if a.IsCastle {
					found = true
					if tc.castleSide == "kingside" {
						assert.True(t, a.IsKingsideCastle)
					} else {
						assert.True(t, a.IsQueensideCastle)
					}
				}
			}
			assert.Equal(t, tc.hasCastle, found)
		})
	}

	t.Run("castling rights revoked after black queenside castle", func(t *testing.T) {
		g, err := NewGameFromFEN("r3k3/8/8/8/8/8/8/4K3 b q - 0 1")
		require.NoError(t, err)
		var castle Action
		for _, a := range g.Actions {
			if a.IsQueensideCastle {
				castle = a
			}
		}
		newGame := g.DoAction(castle)
		assert.False(t, newGame.CanBlackQueensideCastle)
		assert.False(t, newGame.CanBlackKingsideCastle)
		assert.Equal(t, "2kr4/8/8/8/8/8/8/4K3 w - - 1 2", newGame.ToFEN())
	})

	t.Run("castling rights revoked after white queenside castle", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/8/8/8/R3K3 w Q - 0 1")
		require.NoError(t, err)
		var castle Action
		for _, a := range g.Actions {
			if a.IsQueensideCastle {
				castle = a
			}
		}
		newGame := g.DoAction(castle)
		assert.False(t, newGame.CanWhiteQueensideCastle)
		assert.False(t, newGame.CanWhiteKingsideCastle)
	})
}

// Ported from ostinato's EnPassantTest.scala
func TestOstinatoEnPassant(t *testing.T) {
	t.Run("en passant capture available for white pawn", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/3p4/3Pp3/8/8/8/4K3 w - e6 0 1")
		require.NoError(t, err)
		found := false
		for _, a := range g.Actions {
			if a.IsEnPassantCapture {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("en passant illegal when it exposes own king", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/3PpK1r/8/8/8/8 w - e6 0 1")
		require.NoError(t, err)
		for _, a := range g.Actions {
			assert.False(t, a.IsEnPassantCapture, "en passant should be illegal: would expose king to rook")
		}
	})
}

// Ported from ostinato's PawnsTest.scala
func TestOstinatoPawns(t *testing.T) {
	ts := []struct {
		name        string
		fen         string
		pieceXY     XY
		actionCount int
	}{
		{"1 action for black pawn (single advance)", "4k3/8/2p5/8/8/8/8/4K3 b - - 0 1", XY{2, 2}, 1},
		{"2 actions for black pawn (starting rank)", "4k3/2p5/8/8/8/8/8/4K3 b - - 0 1", XY{2, 1}, 2},
		{"4 actions for black pawn (captures both sides)", "4k3/2p5/1R1R4/8/8/8/8/4K3 b - - 0 1", XY{2, 1}, 4},
		{"3 actions for non-starting pawn (advance + 2 captures)", "4k3/8/2p5/1R1R4/8/8/8/4K3 b - - 0 1", XY{2, 2}, 3},
		{"0 actions for blocked pawn", "4k3/2p5/2R5/8/8/8/8/4K3 b - - 0 1", XY{2, 1}, 0},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			count := 0
			for _, a := range g.Actions {
				if !a.IsResign && !a.IsDraw && a.FromPiece.XY == tc.pieceXY {
					count++
				}
			}
			assert.Equal(t, tc.actionCount, count)
		})
	}

	t.Run("4 promotion actions for white pawn on 7th rank", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/2P5/8/8/8/8/8/4K3 w - - 0 1")
		require.NoError(t, err)
		promoCount := 0
		for _, a := range g.Actions {
			if a.IsPromotion && a.FromPiece.XY == (XY{2, 1}) {
				promoCount++
			}
		}
		assert.Equal(t, 4, promoCount)
	})

	t.Run("4 capture-promotion actions for white pawn", func(t *testing.T) {
		g, err := NewGameFromFEN("3n4/2P5/8/8/8/8/8/4K2k w - - 0 1")
		require.NoError(t, err)
		capturePromoCount := 0
		for _, a := range g.Actions {
			if a.IsPromotion && a.IsCapture {
				capturePromoCount++
			}
		}
		assert.Equal(t, 4, capturePromoCount)
	})

	t.Run("promote to each piece type", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/2P5/8/8/8/8/8/4K3 w - - 0 40")
		require.NoError(t, err)
		promotedTo := map[PieceType]bool{}
		for _, a := range g.Actions {
			if a.IsPromotion {
				newGame := g.DoAction(a)
				piece := newGame.PieceAt(a.ToXY)
				promotedTo[piece.PieceType] = true
			}
		}
		assert.True(t, promotedTo[PieceQueen])
		assert.True(t, promotedTo[PieceRook])
		assert.True(t, promotedTo[PieceBishop])
		assert.True(t, promotedTo[PieceKnight])
	})
}

// Ported from ostinato's GameEndingTest.scala
func TestOstinatoGameEnding(t *testing.T) {
	t.Run("checkmate: game over but not draw", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/8/8/5qq1/7K w - - 0 1")
		require.NoError(t, err)
		assert.True(t, g.IsGameOver)
		assert.True(t, g.IsCheckmate)
		assert.False(t, g.IsDraw)
	})

	t.Run("stalemate: game over with no winner", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/8/6q1/8/7K w - - 0 1")
		require.NoError(t, err)
		assert.True(t, g.IsStalemate)
		assert.True(t, g.IsGameOver)
		assert.False(t, g.IsCheckmate)
	})

	t.Run("not draw nor game over", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/6q1/8/8/7K w - - 0 1")
		require.NoError(t, err)
		assert.False(t, g.IsDraw)
		assert.False(t, g.IsGameOver)
	})

	t.Run("not game over even if threatened", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/8/8/6q1/7K w - - 0 1")
		require.NoError(t, err)
		assert.False(t, g.IsGameOver)
		assert.True(t, g.IsCheck)
	})

	t.Run("black has seven checkmate actions available", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/5q2/6q1/8/7K b - - 0 1")
		require.NoError(t, err)
		mateCount := 0
		for _, a := range g.Actions {
			if a.IsResign || a.IsDraw {
				continue
			}
			newGame := g.DoAction(a)
			if newGame.IsCheckmate {
				mateCount++
			}
		}
		assert.Equal(t, 7, mateCount)
	})

	t.Run("lower queen to h3 is check but not checkmate", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/5q2/6q1/8/7K b - - 0 1")
		require.NoError(t, err)
		var action Action
		for _, a := range g.Actions {
			if a.FromPiece.XY == (XY{6, 5}) && a.ToXY == (XY{7, 5}) {
				action = a
			}
		}
		require.NotEqual(t, Action{}, action)
		newGame := g.DoAction(action)
		assert.True(t, newGame.IsCheck)
		assert.False(t, newGame.IsCheckmate)
	})

	t.Run("upper queen to h4 is check and checkmate", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/5q2/6q1/8/7K b - - 0 1")
		require.NoError(t, err)
		var action Action
		for _, a := range g.Actions {
			if a.FromPiece.XY == (XY{5, 4}) && a.ToXY == (XY{7, 4}) {
				action = a
			}
		}
		require.NotEqual(t, Action{}, action)
		newGame := g.DoAction(action)
		assert.True(t, newGame.IsCheck)
		assert.True(t, newGame.IsCheckmate)
	})

	t.Run("upper queen to a4 is neither check nor checkmate", func(t *testing.T) {
		g, err := NewGameFromFEN("4k3/8/8/8/5q2/6q1/8/7K b - - 0 1")
		require.NoError(t, err)
		var action Action
		for _, a := range g.Actions {
			if a.FromPiece.XY == (XY{5, 4}) && a.ToXY == (XY{0, 4}) {
				action = a
			}
		}
		require.NotEqual(t, Action{}, action)
		newGame := g.DoAction(action)
		assert.False(t, newGame.IsCheck)
		assert.False(t, newGame.IsCheckmate)
	})

	t.Run("draw and resign are always available", func(t *testing.T) {
		g := NewDefaultGame()
		hasDraw, hasResign := false, false
		for _, a := range g.Actions {
			if a.IsDraw {
				hasDraw = true
			}
			if a.IsResign {
				hasResign = true
			}
		}
		assert.True(t, hasDraw)
		assert.True(t, hasResign)

		firstMove := g.Actions[0]
		g2 := g.DoAction(firstMove)
		hasDraw, hasResign = false, false
		for _, a := range g2.Actions {
			if a.IsDraw {
				hasDraw = true
			}
			if a.IsResign {
				hasResign = true
			}
		}
		assert.True(t, hasDraw)
		assert.True(t, hasResign)
	})
}

// Ported from ostinato's ChessBoardTest.scala
func TestOstinatoChessBoard(t *testing.T) {
	t.Run("halfMoveClock increments on non-capture non-pawn move", func(t *testing.T) {
		g := NewDefaultGame()
		assert.Equal(t, 0, g.HalfMoveClock)
		var knightAction Action
		for _, a := range g.Actions {
			if a.FromPiece.PieceType == PieceKnight {
				knightAction = a
				break
			}
		}
		newGame := g.DoAction(knightAction)
		assert.Equal(t, 1, newGame.HalfMoveClock)
	})

	t.Run("halfMoveClock resets on pawn move", func(t *testing.T) {
		g := NewDefaultGame()
		var pawnAction Action
		for _, a := range g.Actions {
			if a.FromPiece.PieceType == PiecePawn {
				pawnAction = a
				break
			}
		}
		newGame := g.DoAction(pawnAction)
		assert.Equal(t, 0, newGame.HalfMoveClock)
	})

	t.Run("fullMoveNumber increments after black move", func(t *testing.T) {
		g := NewDefaultGame()
		assert.Equal(t, color(ColorWhite), g.Turn())
		assert.Equal(t, 1, g.FullMoveNumber)

		whiteMove := g.Actions[0]
		g2 := g.DoAction(whiteMove)
		assert.Equal(t, color(ColorBlack), g2.Turn())
		assert.Equal(t, 1, g2.FullMoveNumber)

		blackMove := g2.Actions[0]
		g3 := g2.DoAction(blackMove)
		assert.Equal(t, color(ColorWhite), g3.Turn())
		assert.Equal(t, 2, g3.FullMoveNumber)
	})

	t.Run("turn flips after white move", func(t *testing.T) {
		g := NewDefaultGame()
		g2 := g.DoAction(g.Actions[0])
		assert.Equal(t, color(ColorBlack), g2.Turn())
	})

	t.Run("turn flips after black move", func(t *testing.T) {
		g, err := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")
		require.NoError(t, err)
		g2 := g.DoAction(g.Actions[0])
		assert.Equal(t, color(ColorWhite), g2.Turn())
	})

	t.Run("play 200 moves without panic", func(t *testing.T) {
		g := NewDefaultGame()
		for i := 0; i < 200; i++ {
			if g.IsGameOver {
				break
			}
			var move Action
			for _, a := range g.Actions {
				if !a.IsResign && !a.IsDraw {
					move = a
					break
				}
			}
			if move == (Action{}) {
				break
			}
			g = g.DoAction(move)
		}
	})
}

// Ported from ostinato's ChessGameTest.scala
func TestOstinatoChessGame(t *testing.T) {
	countPieceActions := func(g Game, xy XY) int {
		count := 0
		for _, a := range g.Actions {
			if !a.IsResign && !a.IsDraw && a.FromPiece.XY == xy {
				count++
			}
		}
		return count
	}

	ts := []struct {
		name  string
		fen   string
		xy    XY
		count int
	}{
		{"rook has 14 actions (open board)", "4k3/8/8/3r4/8/8/8/4K3 b - - 0 1", XY{3, 3}, 14},
		{"rook has 0 actions (surrounded by own pieces)", "4k3/8/3n4/2nrn3/3n4/8/8/4K3 b - - 0 1", XY{3, 3}, 0},
		{"rook has 1 action (capture white rook)", "4k3/8/3n4/2nrn3/3R4/8/8/4K3 b - - 0 1", XY{3, 3}, 1},
		{"bishop has 13 actions (open board)", "4k3/8/8/3b4/8/8/8/4K3 b - - 0 1", XY{3, 3}, 13},
		{"bishop has 0 actions (all 8 neighbors are own pieces)", "4k3/8/2nnn3/2nbn3/2nnn3/8/8/4K3 b - - 0 1", XY{3, 3}, 0},
		{"bishop has 3 actions (one diagonal open)", "4k3/8/2nnn3/2nbn3/4n3/8/8/4K3 b - - 0 1", XY{3, 3}, 3},
		{"bishop has 1 action (capture white bishop)", "4k3/8/2nnn3/2nbn3/2B1n3/8/8/4K3 b - - 0 1", XY{3, 3}, 1},
		{"knight has 8 actions (open board)", "4k3/8/8/3n4/8/8/8/4K3 b - - 0 1", XY{3, 3}, 8},
		{"knight has 0 actions (all landing squares blocked by own)", "n3k3/2b5/1b6/8/8/8/8/4K3 b - - 0 1", XY{0, 0}, 0},
		{"knight has 1 action (capture white bishop)", "n3k3/2b5/1B6/8/8/8/8/4K3 b - - 0 1", XY{0, 0}, 1},
		{"queen has 27 actions (open board)", "4k3/8/8/3q4/8/8/8/4K3 b - - 0 1", XY{3, 3}, 27},
		{"queen has 7 actions (corner, blocked by own)", "4k1bq/6b1/8/8/8/8/8/4K3 b - - 0 1", XY{7, 0}, 7},
		{"queen has 1 action (capture white bishop)", "4k1bq/6bB/8/8/8/8/8/4K3 b - - 0 1", XY{7, 0}, 1},
		// Queen on h8 hemmed in by own bishops on g7/g8 and white king on h7: the queen
		// cannot take the king (an illegal "move"), so she has 0 legal actions. The
		// FEN uses a separate white king position to keep the position legal.
		{"queen has 0 actions (hemmed in by own pieces)", "4k1bq/6bn/8/8/8/8/8/1K6 b - - 0 1", XY{7, 0}, 0},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			assert.Equal(t, tc.count, countPieceActions(g, tc.xy))
		})
	}

	t.Run("king has 8 actions on open board", func(t *testing.T) {
		g, err := NewGameFromFEN("8/8/8/3K4/8/8/8/4k3 w - - 0 1")
		require.NoError(t, err)
		boardActions := 0
		for _, a := range g.Actions {
			if !a.IsResign && !a.IsDraw {
				boardActions++
			}
		}
		assert.Equal(t, 8, boardActions)
	})

	t.Run("king can capture bishop", func(t *testing.T) {
		g, err := NewGameFromFEN("8/5Kb1/5p2/8/8/8/8/4k3 w - - 0 1")
		require.NoError(t, err)
		hasCapture := false
		for _, a := range g.Actions {
			if a.IsCapture && a.CapturedPiece.PieceType == PieceBishop {
				hasCapture = true
			}
		}
		assert.True(t, hasCapture)
	})

	t.Run("default game string representation", func(t *testing.T) {
		g := NewDefaultGame()
		expected := "♜♞♝♛♚♝♞♜\n♟♟♟♟♟♟♟♟\n........\n........\n........\n........\n♙♙♙♙♙♙♙♙\n♖♘♗♕♔♗♘♖\n"
		assert.Equal(t, expected, g.String())
	})
}
