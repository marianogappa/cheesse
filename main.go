package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/marianogappa/cheesse/api"
)

var (
	flagServe     = flag.Int("serve", 0, "Start a server on the specified port.")
	flagParseGame = flag.String("parseGame", "", "ParseGame api call. Requires a JSON string with arguments. Please review spec.")
	a             = api.New()
)

// api.InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},

func main() {
	flag.Parse()
	http.HandleFunc("/parseGame", handleServerParseGame)
	switch {
	case *flagServe != 0:
		http.ListenAndServe(fmt.Sprintf(":%v", *flagServe), nil)
	case *flagParseGame != "":
		handleCliParseGame(flagParseGame)
	}
}

func handleServerParseGame(w http.ResponseWriter, r *http.Request) {
	var ig api.InputGame
	if err := json.NewDecoder(r.Body).Decode(&ig); err != nil {
		fmt.Fprintln(w, formatError(err))
		return
	}
	defer r.Body.Close()
	outputGame, err := a.ParseGame(ig)
	if err != nil {
		fmt.Fprintln(w, formatError(err))
		return
	}
	json.NewEncoder(w).Encode(outputGame)
}

func handleCliParseGame(flagParseGame *string) {
	var ig api.InputGame
	if err := json.Unmarshal([]byte(*flagParseGame), &ig); err != nil {
		mustCliFatal(err)
	}
	outputGame, err := a.ParseGame(ig)
	if err != nil {
		mustCliFatal(err)
	}
	byts, _ := json.Marshal(outputGame)
	fmt.Println(string(byts))
}

func mustCliFatal(err error) {
	fmt.Println(formatError(err))
	os.Exit(1)
}

func formatError(err error) string {
	errByts, _ := json.Marshal(err.Error())
	return fmt.Sprintf(`{"error": %v}`, string(errByts))
}
