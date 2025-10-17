# cheesse
Simple package, server, CLI tool and WebAssembly binary for all things chess.

Please note that this library is NOT YET ready for mainstream use. Its API is not final, two of its API methods are not fully implemented, and it hasn't yet been battle-tested against a massive corpus of games (only about 300).

## API

```go
DefaultGame() OutputGame
ParseGame(game InputGame) (OutputGame, error)
DoAction(game InputGame, action InputAction) (OutputGame, OutputAction, error)

// Currently only supporting Algebraic Notation; others coming soon
ParseNotation(game InputGame, notationString string) (OutputGame, []OutputGameStep, error)

// Coming soon
ConvertNotation(game InputGame, notationString string, toNotation string) (OutputGame, []OutputGameStep, error)
```

## Server example

```bash
$ ./cheesse -serve 8080
```

```bash
$ curl localhost:8080/defaultGame | jq .game.board.board
```

```json
[
  "♜♞♝♛♚♝♞♜",
  "♟♟♟♟♟♟♟♟",
  "        ",
  "        ",
  "        ",
  "        ",
  "♙♙♙♙♙♙♙♙",
  "♖♘♗♕♔♗♘♖"
]
```

## CLI example

```bash
$ ./cheesse -defaultGame | jq .game.board.board
```

```json
[
  "♜♞♝♛♚♝♞♜",
  "♟♟♟♟♟♟♟♟",
  "        ",
  "        ",
  "        ",
  "        ",
  "♙♙♙♙♙♙♙♙",
  "♖♘♗♕♔♗♘♖"
]
```
## Package import example

```go
package main

import (
  "fmt"

  "github.com/marianogappa/cheesse/api"
)

func main() {
	for _, s := range a.DefaultGame().Board.Board {
		fmt.Println(s)
	}
}
```

```
♜♞♝♛♚♝♞♜
♟♟♟♟♟♟♟♟
        
        
        
        
♙♙♙♙♙♙♙♙
♖♘♗♕♔♗♘♖
```

## WebAssembly example (using TinyGo compiler)

[Auto-play](https://marianogappa.github.io/cheesse-examples/)

## Why is it called "cheesse"?

That's roughly how kiwi people pronounce chess.
