package printer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

// PGNPrinter renders a game as a PGN document: a tag-pair section (the Seven Tag
// Roster plus any extra metadata) followed by the movetext with move numbers,
// comments, the result marker, and lines wrapped at 80 columns.
type PGNPrinter struct {
	// Metadata holds tag pairs to render in the tag section (e.g. from a parsed
	// PGN's headers). Seven Tag Roster keys missing from it get "?" placeholders.
	Metadata map[string]string
}

var pgnSevenTagRoster = []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"}

// PrintGame renders the full PGN document as lines: the tag section, a blank
// line, then the movetext wrapped at 80 columns.
func (p PGNPrinter) PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error) {
	result := p.resultMarker(gameSteps)

	lines := []string{}
	for _, tag := range pgnSevenTagRoster {
		value, ok := p.Metadata[tag]
		if !ok {
			switch tag {
			case "Date":
				value = "????.??.??"
			case "Result":
				value = result
			default:
				value = "?"
			}
		}
		lines = append(lines, fmt.Sprintf("[%s %q]", tag, value))
	}
	// Passthrough of any non-STR metadata, in deterministic order.
	extraKeys := []string{}
	for k := range p.Metadata {
		isSTR := false
		for _, tag := range pgnSevenTagRoster {
			if k == tag {
				isSTR = true
				break
			}
		}
		if !isSTR {
			extraKeys = append(extraKeys, k)
		}
	}
	sort.Strings(extraKeys)
	for _, k := range extraKeys {
		lines = append(lines, fmt.Sprintf("[%s %q]", k, p.Metadata[k]))
	}
	lines = append(lines, "")

	movetext, err := p.movetext(gameSteps, result)
	if err != nil {
		return nil, err
	}
	lines = append(lines, wrapText(movetext, 80)...)
	return lines, nil
}

// PrintAction renders a single action in SAN (the movetext language of PGN).
func (p PGNPrinter) PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error) {
	return AlgebraicPrinter{}.PrintAction(gameStep, gameCharacteristics)
}

// resultMarker derives the game's result marker from the final step: an explicit
// result-marker step wins; otherwise the final game state decides.
func (p PGNPrinter) resultMarker(gameSteps []core.GameStep) string {
	if len(gameSteps) == 0 {
		return "*"
	}
	last := gameSteps[len(gameSteps)-1]
	switch last.StepString {
	case "1-0", "0-1", "1/2-1/2", "*":
		return last.StepString
	}
	finalGame := last.StepGame
	switch {
	case finalGame.IsCheckmate && finalGame.GameOverWinner.String() == "White",
		last.StepAction.IsResign && last.StepAction.FromPiece.Owner.String() == "Black":
		return "1-0"
	case finalGame.IsCheckmate && finalGame.GameOverWinner.String() == "Black",
		last.StepAction.IsResign && last.StepAction.FromPiece.Owner.String() == "White":
		return "0-1"
	case finalGame.IsDraw, finalGame.IsStalemate, last.StepAction.IsDraw:
		return "1/2-1/2"
	}
	return "*"
}

func (p PGNPrinter) movetext(gameSteps []core.GameStep, result string) (string, error) {
	sanCharacteristics := SANCharacteristics()
	var sb strings.Builder
	needsMoveNumber := true
	for _, gameStep := range gameSteps {
		// Result markers and terminal actions are rendered via the result marker at the end.
		if gameStep.StepAction == (core.Action{}) || gameStep.StepAction.IsResign || gameStep.StepAction.IsDraw {
			continue
		}
		san, err := AlgebraicPrinter{}.PrintAction(gameStep, sanCharacteristics)
		if err != nil {
			return "", err
		}
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		preMoveGame := gameStep.StepPreMoveGame
		if preMoveGame.Turn() == core.ColorWhite {
			fmt.Fprintf(&sb, "%d. ", preMoveGame.FullMoveNumber)
		} else if needsMoveNumber {
			fmt.Fprintf(&sb, "%d... ", preMoveGame.FullMoveNumber)
		}
		sb.WriteString(san)
		needsMoveNumber = false
		if gameStep.StepComment != "" {
			fmt.Fprintf(&sb, " {%s}", gameStep.StepComment)
			needsMoveNumber = true // Conventionally re-state the move number after a comment
		}
	}
	if sb.Len() > 0 {
		sb.WriteByte(' ')
	}
	sb.WriteString(result)
	return sb.String(), nil
}

// wrapText greedily wraps the text at the given column, breaking on spaces.
func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	lines := []string{}
	current := words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) > width {
			lines = append(lines, current)
			current = word
			continue
		}
		current += " " + word
	}
	return append(lines, current)
}
