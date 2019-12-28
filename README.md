# cheesse
Simple package, server and CLI tool for all things Chess.

## API

```go
DefaultGame() OutputGame
ParseGame(game InputGame) (OutputGame, error)
DoAction(game InputGame, action InputAction) (OutputGame, OutputAction, error)
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
