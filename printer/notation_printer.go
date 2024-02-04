package printer

import (
	"fmt"

	"github.com/marianogappa/cheesse/core"
)

type NotationPrinter interface {
	PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error)
	PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error)
}

type GameCharacteristics struct {
	isFigurine                     bool
	isCheck                        bool
	isCheckmate                    bool
	usesCheckSymbol                *string
	usesCheckmateSymbol            *string
	fullMoveNumber                 *int
	usesFullMoveDot                *bool
	usesNewlineAsFullMoveSeparator *bool
	usesThreatenSymbol             *string
	usesAnnotationSymbol           *string
	usesCaptureSymbol              *string
	usesEndGameSymbol              *string
	usesPromotionSymbol            *string
	usesCastlingSymbol             *string
	usesEnPassantSymbol            *string
	descriptiveUseKt               *bool
}

func pstr(s string) *string {
	return &s
}

func pbool(b bool) *bool {
	return &b
}

func applyDefaultGameCharacteristics(gameCharacteristics GameCharacteristics) GameCharacteristics {
	if gameCharacteristics.usesCheckSymbol == nil {
		gameCharacteristics.usesCheckSymbol = pstr("+")
	}
	if gameCharacteristics.usesCheckmateSymbol == nil {
		gameCharacteristics.usesCheckmateSymbol = pstr("++")
	}
	if gameCharacteristics.usesFullMoveDot == nil {
		gameCharacteristics.usesFullMoveDot = pbool(true)
	}
	if gameCharacteristics.usesNewlineAsFullMoveSeparator == nil {
		gameCharacteristics.usesNewlineAsFullMoveSeparator = pbool(true)
	}
	if gameCharacteristics.usesThreatenSymbol == nil {
		gameCharacteristics.usesThreatenSymbol = pstr("")
	}
	if gameCharacteristics.usesAnnotationSymbol == nil {
		gameCharacteristics.usesAnnotationSymbol = pstr("")
	}
	if gameCharacteristics.usesCaptureSymbol == nil {
		gameCharacteristics.usesCaptureSymbol = pstr("x")
	}
	if gameCharacteristics.usesEndGameSymbol == nil {
		gameCharacteristics.usesEndGameSymbol = pstr("resigns") // TODO
	}
	if gameCharacteristics.usesPromotionSymbol == nil {
		gameCharacteristics.usesPromotionSymbol = pstr("=")
	}
	if gameCharacteristics.usesCastlingSymbol == nil {
		gameCharacteristics.usesCastlingSymbol = pstr("0-0")
	}
	if gameCharacteristics.usesEnPassantSymbol == nil {
		gameCharacteristics.usesEnPassantSymbol = pstr("e.p.")
	}
	return gameCharacteristics
}

func genericGamePrinter(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics, p NotationPrinter) ([]string, error) {
	rawStepStrings := []string{}
	for _, gameStep := range gameSteps {
		stepString, err := p.PrintAction(gameStep, gameCharacteristics)
		if err != nil {
			return nil, err
		}
		rawStepStrings = append(rawStepStrings, stepString)
	}

	stepStrings := []string{}
	i := 0
	fullMoveNumber := 1
	for i < len(rawStepStrings) {
		secondMovePart := ""
		if i+1 < len(rawStepStrings) {
			secondMovePart = fmt.Sprintf(" %v", rawStepStrings[i+1])
		}
		stepStrings = append(stepStrings, fmt.Sprintf("%v. %v%v", fullMoveNumber, rawStepStrings[i], secondMovePart))
		fullMoveNumber++
		i += 2
	}

	return stepStrings, nil
}
