// +build !tinygo

package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	flagServe         = flag.Int("serve", 0, "Start a server on the specified port.")
	flagDefaultGame   = flag.Bool("defaultGame", false, "Default API call. Returns a default game.")
	flagParseGame     = flag.String("parseGame", "", "ParseGame API call. Requires a JSON string with arguments. Please review spec.")
	flagDoAction      = flag.String("doAction", "", "DoAction API call. Requires a JSON string with arguments. Please review spec.")
	flagParseNotation = flag.String("parseNotation", "", "ParseNotation API call. Requires a JSON string with arguments. Please review spec.")
)

func main() {
	flag.Parse()

	http.HandleFunc("/parseGame", handleServerParseGame)
	http.HandleFunc("/defaultGame", handleServerDefaultGame)
	http.HandleFunc("/doAction", handleServerDoAction)
	http.HandleFunc("/parseNotation", handleServerParseNotation)

	switch {
	case *flagServe != 0:
		http.ListenAndServe(fmt.Sprintf(":%v", *flagServe), nil)
	case *flagDefaultGame:
		handleCliDefaultGame()
	case *flagParseGame != "":
		handleCliParseGame(flagParseGame)
	case *flagDoAction != "":
		handleCliDoAction(flagDoAction)
	case *flagParseNotation != "":
		handleCliParseNotation(flagParseNotation)
	}
}
