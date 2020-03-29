// +build !tinygo

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/marianogappa/cheesse/api"
)

var a = api.New()

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

func handleServerParseNotation(w http.ResponseWriter, r *http.Request) {
	type args struct {
		Game           api.InputGame `json:"game"`
		NotationString string        `json:"notationString"`
	}
	var input args
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Fprintln(w, formatError(err))
		return
	}
	defer r.Body.Close()
	outputGame, outputGameSteps, err := a.ParseNotation(input.Game, input.NotationString)
	if err != nil {
		fmt.Fprintln(w, formatError(err))
		return
	}
	type out struct {
		Game            api.OutputGame       `json:"game"`
		OutputGameSteps []api.OutputGameStep `json:"outputGameSteps"`
	}
	json.NewEncoder(w).Encode(out{outputGame, outputGameSteps})
}

func handleCliParseNotation(flagParseNotation *string) {
	type args struct {
		Game           api.InputGame `json:"game"`
		NotationString string        `json:"notationString"`
	}
	var input args
	if err := json.Unmarshal([]byte(*flagParseNotation), &input); err != nil {
		mustCliFatal(err)
	}
	outputGame, outputGameSteps, err := a.ParseNotation(input.Game, input.NotationString)
	if err != nil {
		mustCliFatal(err)
	}
	type out struct {
		Game            api.OutputGame       `json:"game"`
		OutputGameSteps []api.OutputGameStep `json:"outputGameSteps"`
	}
	byts, _ := json.Marshal(out{outputGame, outputGameSteps})
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
