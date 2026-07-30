//go:build tinygo
// +build tinygo

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/marianogappa/cheesse/api"
)

// The WASM bindings expose the full cheesse API as synchronous JS globals.
//
// Every function takes a single Uint8Array argument containing a JSON request and
// returns a Uint8Array containing a JSON response. Field casing follows the JSON
// tags of the api package entities (camelCase), so the JS shapes are exactly the
// same as the HTTP server's.
//
// Response envelope: {"error": "..."} on failure, otherwise the same shape as the
// corresponding HTTP endpoint's response.
func main() {
	js.Global().Set("cheesseDefaultGame", js.FuncOf(jsDefaultGame))
	js.Global().Set("cheesseParseGame", js.FuncOf(jsParseGame))
	js.Global().Set("cheesseDoAction", js.FuncOf(jsDoAction))
	js.Global().Set("cheesseParseNotation", js.FuncOf(jsParseNotation))
	js.Global().Set("cheesseConvertNotation", js.FuncOf(jsConvertNotation))
	js.Global().Set("cheesseAIMove", js.FuncOf(jsAIMove))
	select {}
}

var a = api.New()

func jsDefaultGame(this js.Value, p []js.Value) interface{} {
	type out struct {
		Game api.OutputGame `json:"game"`
	}
	return toJS(out{a.DefaultGame()}, nil)
}

func jsParseGame(this js.Value, p []js.Value) interface{} {
	type args struct {
		Game api.InputGame `json:"game"`
	}
	var input args
	if err := fromJS(p[0], &input); err != nil {
		return toJS(nil, err)
	}
	outputGame, err := a.ParseGame(input.Game)
	if err != nil {
		return toJS(nil, err)
	}
	type out struct {
		Game api.OutputGame `json:"game"`
	}
	return toJS(out{outputGame}, nil)
}

func jsDoAction(this js.Value, p []js.Value) interface{} {
	type args struct {
		Game   api.InputGame   `json:"game"`
		Action api.InputAction `json:"action"`
	}
	var input args
	if err := fromJS(p[0], &input); err != nil {
		return toJS(nil, err)
	}
	outputGame, outputAction, err := a.DoAction(input.Game, input.Action)
	if err != nil {
		return toJS(nil, err)
	}
	type out struct {
		Game   api.OutputGame   `json:"game"`
		Action api.OutputAction `json:"action"`
	}
	return toJS(out{outputGame, outputAction}, nil)
}

func jsParseNotation(this js.Value, p []js.Value) interface{} {
	type args struct {
		Game           api.InputGame `json:"game"`
		NotationString string        `json:"notationString"`
	}
	var input args
	if err := fromJS(p[0], &input); err != nil {
		return toJS(nil, err)
	}
	outputGame, parseResult, err := a.ParseNotation(input.Game, input.NotationString)
	if err != nil {
		return toJS(nil, err)
	}
	type out struct {
		Game        api.OutputGame        `json:"game"`
		ParseResult api.OutputParseResult `json:"parseResult"`
	}
	return toJS(out{outputGame, parseResult}, nil)
}

func jsConvertNotation(this js.Value, p []js.Value) interface{} {
	type args struct {
		Game           api.InputGame `json:"game"`
		NotationString string        `json:"notationString"`
		TargetNotation string        `json:"targetNotation"`
	}
	var input args
	if err := fromJS(p[0], &input); err != nil {
		return toJS(nil, err)
	}
	outputGame, parseResult, err := a.ConvertNotation(input.Game, input.NotationString, input.TargetNotation)
	if err != nil {
		return toJS(nil, err)
	}
	type out struct {
		Game        api.OutputGame        `json:"game"`
		ParseResult api.OutputParseResult `json:"parseResult"`
	}
	return toJS(out{outputGame, parseResult}, nil)
}

func jsAIMove(this js.Value, p []js.Value) interface{} {
	type args struct {
		Game api.InputGame `json:"game"`
		Mode string        `json:"mode"`
	}
	var input args
	if err := fromJS(p[0], &input); err != nil {
		return toJS(nil, err)
	}
	outputGame, outputAction, moveAvailable, err := a.AIMove(input.Game, input.Mode)
	if err != nil {
		return toJS(nil, err)
	}
	type out struct {
		Game          api.OutputGame   `json:"game"`
		Action        api.OutputAction `json:"action"`
		MoveAvailable bool             `json:"moveAvailable"`
	}
	return toJS(out{outputGame, outputAction, moveAvailable}, nil)
}

// fromJS reads a Uint8Array JS value containing JSON into dst.
func fromJS(v js.Value, dst interface{}) error {
	jsonBytes := make([]byte, v.Length())
	js.CopyBytesToGo(jsonBytes, v)
	return json.Unmarshal(jsonBytes, dst)
}

type errOut struct {
	Error string `json:"error"`
}

// toJS marshals a response (or an error envelope) to JSON and returns it as a
// Uint8Array JS value.
func toJS(response interface{}, err error) js.Value {
	var bs []byte
	if err != nil {
		bs, _ = json.Marshal(errOut{err.Error()})
	} else {
		bs, err = json.Marshal(response)
		if err != nil {
			bs, _ = json.Marshal(errOut{err.Error()})
		}
	}
	buffer := js.Global().Get("Uint8Array").New(len(bs))
	js.CopyBytesToJS(buffer, bs)
	return buffer
}
