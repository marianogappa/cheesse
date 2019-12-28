package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	flagServe       = flag.Int("serve", 0, "Start a server on the specified port.")
	flagDefaultGame = flag.Bool("defaultGame", false, "Default API call. Returns a default game.")
	flagParseGame   = flag.String("parseGame", "", "ParseGame API call. Requires a JSON string with arguments. Please review spec.")
	flagDoAction    = flag.String("doAction", "", "DoAction API call. Requires a JSON string with arguments. Please review spec.")
)

// api.InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},

func main() {
	flag.Parse()

	http.HandleFunc("/parseGame", handleServerParseGame)
	http.HandleFunc("/defaultGame", handleServerDefaultGame)
	http.HandleFunc("/doAction", handleServerDoAction)

	switch {
	case *flagServe != 0:
		http.ListenAndServe(fmt.Sprintf(":%v", *flagServe), nil)
	case *flagDefaultGame:
		handleCliDefaultGame()
	case *flagParseGame != "":
		handleCliParseGame(flagParseGame)
	case *flagDoAction != "":
		handleCliDoAction(flagDoAction)
	}
}
