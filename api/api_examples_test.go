package api

import (
	"fmt"
	"strings"
)

func ExampleAPI_DefaultGame() {
	var (
		a    = New()
		game = a.DefaultGame()
	)
	fmt.Println(game.Board.Board[0])
	fmt.Println(game.Board.Board[1])
	fmt.Println(game.Board.Board[6])
	fmt.Println(game.Board.Board[7])

	// Output:
	// ♜♞♝♛♚♝♞♜
	// ♟♟♟♟♟♟♟♟
	// ♙♙♙♙♙♙♙♙
	// ♖♘♗♕♔♗♘♖
}

func ExampleAPI_ParseGame_default_game() {
	var (
		a       = New()
		game, _ = a.ParseGame(InputGame{}) // Default game inferred
	)
	fmt.Println(game.Board.Board[0])
	fmt.Println(game.Board.Board[1])
	fmt.Println(game.Board.Board[6])
	fmt.Println(game.Board.Board[7])

	// Output:
	// ♜♞♝♛♚♝♞♜
	// ♟♟♟♟♟♟♟♟
	// ♙♙♙♙♙♙♙♙
	// ♖♘♗♕♔♗♘♖
}

func ExampleAPI_ParseGame_fen_string() {
	var (
		a       = New()
		game, _ = a.ParseGame(InputGame{FENString: "4k3/4p3/P2p4/7p/2bP4/p7/2P5/K2B4 w - - 0 1"})
	)
	for _, line := range game.Board.Board {
		// Sorry about the TrimRight. Example matches strings but trims spaces on the right.
		fmt.Println(strings.TrimRight(line, " "))
	}

	// Output:
	//     ♚
	//     ♟
	// ♙  ♟
	//        ♟
	//   ♝♙
	// ♟
	//   ♙
	// ♔  ♗
}

func ExampleAPI_ParseGame_board() {
	var (
		a       = New()
		game, _ = a.ParseGame(InputGame{Board: Board{
			Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟    ♟",
				"       ♟",
				"       ♟",
				"       ♟",
				"       ♟",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			},
			CanWhiteKingsideCastle:  true,
			CanWhiteQueensideCastle: true,
			CanBlackKingsideCastle:  true,
			CanBlackQueensideCastle: true,
			HalfMoveClock:           0,
			FullMoveNumber:          1,
			EnPassantTargetSquare:   "",
			Turn:                    "White",
		}})
	)
	for _, line := range game.Board.Board {
		fmt.Println(line)
	}

	// Output:
	// ♜♞♝♛♚♝♞♜
	// ♟♟♟    ♟
	//        ♟
	//        ♟
	//        ♟
	//        ♟
	// ♙♙♙♙♙♙♙♙
	// ♖♘♗♕♔♗♘♖
}
