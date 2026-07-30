# cheesse
Simple package, server, CLI tool and WebAssembly binary for all things chess.

Please note that this library is NOT YET ready for mainstream use. Its API is not final, two of its API methods are not fully implemented, and it hasn't yet been battle-tested against a massive corpus of games (only about 300).

## API

```go
DefaultGame() OutputGame
ParseGame(game InputGame) (OutputGame, error)
DoAction(game InputGame, action InputAction) (OutputGame, OutputAction, error)

// Auto-detects the notation: Algebraic (incl. figurine and PGN), Coordinate, Descriptive, ICCF, Smith
ParseNotation(game InputGame, notationString string) (OutputGame, OutputParseResult, error)

// Auto-detects the source notation and re-renders every move in the target notation:
// one of {Algebraic|Figurine|Descriptive|Coordinate|ICCF|Smith}
ConvertNotation(game InputGame, notationString string, targetNotation string) (OutputGame, OutputParseResult, error)
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

## WebAssembly

Build:

```bash
$ GOOS=js GOARCH=wasm go build -tags tinygo -o cheesse.wasm .
```

The binary exposes the full API as synchronous JS globals. Every function takes a
`Uint8Array` containing a JSON request and returns a `Uint8Array` containing a JSON
response (same shapes as the HTTP endpoints; `{"error": "..."}` on failure):

```js
const enc = new TextEncoder(), dec = new TextDecoder();
const call = (fn, obj) => JSON.parse(dec.decode(fn(enc.encode(JSON.stringify(obj)))));

JSON.parse(dec.decode(cheesseDefaultGame()));
call(cheesseParseGame,       {game: {fenString: "..."}});
call(cheesseDoAction,        {game: {}, action: {fromSquare: "e2", toSquare: "e4"}});
call(cheesseParseNotation,   {game: {}, notationString: "1. e4 e5"});
call(cheesseConvertNotation, {game: {}, notationString: "1. e4 e5", targetNotation: "ICCF"});
call(cheesseAIMove,          {game: {}, mode: "random"}); // random|easy|medium|hard
```

[Auto-play example](https://marianogappa.github.io/cheesse-examples/)

## Why is it called "cheesse"?

That's roughly how kiwi people pronounce chess.
