# cheesse

Chess engine as a Go package, HTTP server, CLI and WebAssembly binary: parse, convert and play chess in any notation.

[![Release](https://img.shields.io/github/v/release/marianogappa/cheesse)](https://github.com/marianogappa/cheesse/releases)
[![CI](https://github.com/marianogappa/cheesse/actions/workflows/ci.yml/badge.svg)](https://github.com/marianogappa/cheesse/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/marianogappa/cheesse)](go.mod)
[![Go Reference](https://pkg.go.dev/badge/github.com/marianogappa/cheesse/api.svg)](https://pkg.go.dev/github.com/marianogappa/cheesse/api)

## Live demos

[Notation Converter](https://marianogappa.github.io/cheesse/convert.html) · [Match Parser](https://marianogappa.github.io/cheesse/parser.html) · [Play vs AI](https://marianogappa.github.io/cheesse/play.html) · [AI Auto-play](https://marianogappa.github.io/cheesse/autoplay.html)

All four run the engine as WebAssembly in your browser and are redeployed on every merge to master.

## Features

- **Full rules engine** on bitboards: legal move generation (perft-validated), castling, en passant, promotions, check/checkmate/stalemate, double and discovered check detection.
- **Complete draw handling**: threefold repetition, insufficient material, the 50-move rule, and claimable/agreed draws.
- **Six notations, parsed and printed**: Algebraic (SAN), Figurine, Descriptive, Coordinate, ICCF and Smith, plus PGN with tag pairs, comments, NAGs and `(...)` variations.
- **Notation auto-detection with partial-parse tolerance**: invalid input still yields the valid prefix, the most likely notation, and why parsing stopped.
- **Conversion** of a game from any notation to any notation, optionally from a custom starting position (FEN or board strings).
- **AI move selection**: `random`, or minimax with alpha-beta pruning at `easy`/`medium`/`hard`.
- **Battle-tested**: 5,553 real PGN games in CI, plus ostinato's engine and notation corpora (Chernev endings, six-notation suite).

## API

Six functions, identical across the Go package, HTTP server, CLI and WASM binary:

```go
DefaultGame() OutputGame
ParseGame(game InputGame) (OutputGame, error)
DoAction(game InputGame, action InputAction) (OutputGame, OutputAction, error)

// Auto-detects the notation: Algebraic (incl. Figurine and PGN), Coordinate, Descriptive, ICCF, Smith
ParseNotation(game InputGame, notationString string) (OutputGame, OutputParseResult, error)

// Auto-detects the source notation and re-renders every move in the target notation:
// one of {Algebraic|Figurine|Descriptive|Coordinate|ICCF|Smith}
ConvertNotation(game InputGame, notationString string, targetNotation string) (OutputGame, OutputParseResult, error)

// mode is one of {random|easy|medium|hard}; ok is false when the game is over
AIMove(game InputGame, mode string) (game OutputGame, action OutputAction, ok bool, err error)
```

## Go package

```go
package main

import (
	"fmt"

	"github.com/marianogappa/cheesse/api"
)

func main() {
	a := api.New()

	// Convert a Descriptive-notation game (auto-detected) to ICCF
	_, result, _ := a.ConvertNotation(api.InputGame{}, "1. P-K4 P-K3\n2. P-Q4 P-Q4", "ICCF")
	fmt.Println(result.NotationName) // Descriptive Notation
	for _, step := range result.Steps {
		fmt.Print(step.ActionString, " ") // 5254 5756 4244 4745
	}

	// Make a move and inspect the resulting position
	game, action, _ := a.DoAction(api.InputGame{}, api.InputAction{FromSquare: "e2", ToSquare: "e4"})
	fmt.Println(action.ActionString) // e4
	fmt.Println(game.FENString)      // rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1
}
```

## Server

```bash
cheesse -serve 8080
```

One endpoint per API function (`/defaultGame`, `/parseGame`, `/doAction`, `/parseNotation`, `/convertNotation`, `/aiMove`); requests and responses are JSON:

```bash
curl -s localhost:8080/aiMove -d '{"game":{},"mode":"hard"}' | jq .action.actionString
curl -s localhost:8080/defaultGame | jq .game.board.board
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

## CLI

Same JSON in, JSON out, one flag per API function:

```bash
cheesse -parseNotation '{"game":{},"notationString":"1. e4 e5 2. Nf3"}' | jq .parseResult.notationName
cheesse -convertNotation '{"game":{},"notationString":"1. P-K4 P-K3","targetNotation":"Smith"}' | jq '.parseResult.steps[].actionString'
```

## WebAssembly

Prebuilt `cheesse.wasm`, `cheesse.js` and `wasm_exec.js` are attached to every [release](https://github.com/marianogappa/cheesse/releases) (with `checksums.txt`), or build from source:

```bash
GOOS=js GOARCH=wasm go build -tags tinygo -o cheesse.wasm .
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

## Performance

<!-- bench-start -->
<details>
<summary>Benchmark results (auto-updated on every push)</summary>

Updated by CI on the next merge to master.

</details>
<!-- bench-end -->

## Why is it called "cheesse"?

That's roughly how kiwi people pronounce chess.
