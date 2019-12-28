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
	flagServe       = flag.Int("serve", 0, "Start a server on the specified port.")
	flagDefaultGame = flag.Bool("defaultGame", false, "Default API call. Returns a default game.")
	flagParseGame   = flag.String("parseGame", "", "ParseGame API call. Requires a JSON string with arguments. Please review spec.")
	flagDoAction    = flag.String("doAction", "", "DoAction API call. Requires a JSON string with arguments. Please review spec.")
	a               = api.New()
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

func handleServerDefaultGame(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(a.DefaultGame())
}

func handleCliDefaultGame() {
	byts, _ := json.Marshal(a.DefaultGame())
	fmt.Println(string(byts))
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

func handleServerDoAction(w http.ResponseWriter, r *http.Request) {
	type args struct {
		Game   api.InputGame   `json:"game"`
		Action api.InputAction `json:"action"`
	}
	var input args
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Fprintln(w, formatError(err))
		return
	}
	defer r.Body.Close()
	outputGame, outputAction, err := a.DoAction(input.Game, input.Action)
	if err != nil {
		fmt.Fprintln(w, formatError(err))
		return
	}
	type out struct {
		Game   api.OutputGame   `json:"game"`
		Action api.OutputAction `json:"action"`
	}
	json.NewEncoder(w).Encode(out{outputGame, outputAction})
}

func handleCliDoAction(flagDoAction *string) {
	type args struct {
		Game   api.InputGame   `json:"game"`
		Action api.InputAction `json:"action"`
	}
	var input args
	if err := json.Unmarshal([]byte(*flagDoAction), &input); err != nil {
		mustCliFatal(err)
	}
	outputGame, outputAction, err := a.DoAction(input.Game, input.Action)
	if err != nil {
		mustCliFatal(err)
	}
	type out struct {
		Game   api.OutputGame   `json:"game"`
		Action api.OutputAction `json:"action"`
	}
	byts, _ := json.Marshal(out{outputGame, outputAction})
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
