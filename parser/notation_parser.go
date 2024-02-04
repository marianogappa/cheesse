package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

var (
	errInvalidRegexp       = errors.New("invalid regexp; this should not happen")
	errDidntMatchAnyRegexp = errors.New("invalid state at step TODO: didn't match any regexp")
)

type actionPattern struct {
	fromPieceType      core.PieceType
	fromX              *int
	fromY              *int
	toX                *int
	toY                *int
	isCapture          *bool
	isResign           *bool
	isPromotion        *bool
	isEnPassant        *bool
	isEnPassantCapture *bool
	isCastle           *bool
	isKingsideCastle   *bool
	isQueensideCastle  *bool
	promotionPieceType core.PieceType
	capturedPieceType  core.PieceType
	capturedPieceX     *int
	capturedPieceY     *int
	isCheck            *bool
	isCheckmate        *bool
}

func (p actionPattern) Clone() actionPattern {
	return actionPattern{
		fromPieceType:      p.fromPieceType,
		fromX:              cloneInt(p.fromX),
		fromY:              cloneInt(p.fromY),
		toX:                cloneInt(p.toX),
		toY:                cloneInt(p.toY),
		isCapture:          cloneBool(p.isCapture),
		isResign:           cloneBool(p.isResign),
		isPromotion:        cloneBool(p.isPromotion),
		isEnPassant:        cloneBool(p.isEnPassant),
		isEnPassantCapture: cloneBool(p.isEnPassantCapture),
		isCastle:           cloneBool(p.isCastle),
		isKingsideCastle:   cloneBool(p.isKingsideCastle),
		isQueensideCastle:  cloneBool(p.isQueensideCastle),
		promotionPieceType: p.promotionPieceType,
		capturedPieceType:  p.capturedPieceType,
		capturedPieceX:     cloneInt(p.capturedPieceX),
		capturedPieceY:     cloneInt(p.capturedPieceY),
		isCheck:            cloneBool(p.isCheck),
		isCheckmate:        cloneBool(p.isCheckmate),
	}
}

func cloneInt(i *int) *int {
	if i == nil {
		return nil
	}
	cloned := *i
	return &cloned
}

func cloneBool(b *bool) *bool {
	if b == nil {
		return nil
	}
	cloned := *b
	return &cloned
}

func (p actionPattern) String() string {
	var sb strings.Builder
	sb.WriteString("Action Pattern:\n")
	if p.fromPieceType != core.PieceNone {
		sb.WriteString(fmt.Sprintf("{a.fromPiece.pieceType}:%v\n", p.fromPieceType))
	}
	if p.fromX != nil {
		sb.WriteString(fmt.Sprintf("{a.fromPiece.xy.x}:%v\n", *p.fromX))
	}
	if p.fromY != nil {
		sb.WriteString(fmt.Sprintf("{a.fromPiece.xy.y}:%v\n", *p.fromY))
	}
	if p.toX != nil {
		sb.WriteString(fmt.Sprintf("{a.toXY.x}:%v\n", *p.toX))
	}
	if p.toY != nil {
		sb.WriteString(fmt.Sprintf("{a.toXY.y}:%v\n", *p.toY))
	}
	if p.isCapture != nil {
		sb.WriteString(fmt.Sprintf("{a.isCapture}:%v\n", *p.isCapture))
	}
	if p.isResign != nil {
		sb.WriteString(fmt.Sprintf("{a.isResign}:%v\n", *p.isResign))
	}
	if p.isPromotion != nil {
		sb.WriteString(fmt.Sprintf("{a.isPromotion}:%v\n", *p.isPromotion))
	}
	if p.isEnPassant != nil {
		sb.WriteString(fmt.Sprintf("{a.isEnPassant}:%v\n", *p.isEnPassant))
	}
	if p.isEnPassantCapture != nil {
		sb.WriteString(fmt.Sprintf("{a.isEnPassantCapture}:%v\n", *p.isEnPassantCapture))
	}
	if p.isCastle != nil {
		sb.WriteString(fmt.Sprintf("{a.isCastle}:%v\n", *p.isCastle))
	}
	if p.isKingsideCastle != nil {
		sb.WriteString(fmt.Sprintf("{a.isKingsideCastle}:%v\n", *p.isKingsideCastle))
	}
	if p.isQueensideCastle != nil {
		sb.WriteString(fmt.Sprintf("{a.isQueensideCastle}:%v\n", *p.isQueensideCastle))
	}
	if p.promotionPieceType != core.PieceNone {
		sb.WriteString(fmt.Sprintf("{a.promotionPiece.pieceType}:%v\n", p.promotionPieceType))
	}
	if p.capturedPieceType != core.PieceNone {
		sb.WriteString(fmt.Sprintf("{a.capturedPiece.pieceType}:%v\n", p.capturedPieceType))
	}
	if p.capturedPieceX != nil {
		sb.WriteString(fmt.Sprintf("{a.capturedPiece.xy.x}:%v\n", *p.capturedPieceX))
	}
	if p.capturedPieceY != nil {
		sb.WriteString(fmt.Sprintf("{a.capturedPiece.xy.y}:%v\n", *p.capturedPieceY))
	}
	return sb.String()
}

func (p *actionPattern) isMatch(a core.Action) bool {
	if !pieceTypeMatcher(p.fromPieceType)(a.FromPiece.PieceType) ||
		!intMatcher(p.fromX)(a.FromPiece.XY.X) ||
		!intMatcher(p.fromY)(a.FromPiece.XY.Y) ||
		!intMatcher(p.toX)(a.ToXY.X) ||
		!intMatcher(p.toY)(a.ToXY.Y) ||
		!boolMatcher(p.isCapture)(a.IsCapture) ||
		!boolMatcher(p.isResign)(a.IsResign) ||
		!boolMatcher(p.isPromotion)(a.IsPromotion) ||
		!boolMatcher(p.isEnPassant)(a.IsEnPassant) ||
		!boolMatcher(p.isEnPassantCapture)(a.IsEnPassantCapture) ||
		!boolMatcher(p.isCastle)(a.IsCastle) ||
		!boolMatcher(p.isKingsideCastle)(a.IsKingsideCastle) ||
		!boolMatcher(p.isQueensideCastle)(a.IsQueensideCastle) ||
		!pieceTypeMatcher(p.promotionPieceType)(a.PromotionPieceType) ||
		!pieceTypeMatcher(p.capturedPieceType)(a.CapturedPiece.PieceType) ||
		!intMatcher(p.capturedPieceX)(a.CapturedPiece.XY.X) ||
		!intMatcher(p.capturedPieceY)(a.CapturedPiece.XY.Y) {
		return false
	}
	return true
}

type Characteristics struct {
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
}

func boolMatcher(v *bool) func(interface{}) bool {
	return func(w interface{}) bool { return v == nil || *v == w.(bool) }
}
func intMatcher(v *int) func(interface{}) bool {
	return func(w interface{}) bool { return v == nil || *v == w.(int) }
}
func pieceTypeMatcher(v core.PieceType) func(interface{}) bool {
	return func(w interface{}) bool { return v == core.PieceNone || v == w.(core.PieceType) }
}

type gameAlternative struct {
	initialGame core.Game
	gameSteps   []core.GameStep
}

func (a gameAlternative) clone() gameAlternative {
	clonedGameSteps := make([]core.GameStep, len(a.gameSteps))
	for i := range a.gameSteps {
		clonedGameSteps[i] = a.gameSteps[i].Clone()
	}
	return gameAlternative{
		initialGame: a.initialGame.Clone(),
		gameSteps:   clonedGameSteps,
	}
}

func (a gameAlternative) currentGame() core.Game {
	if len(a.gameSteps) > 0 {
		return a.gameSteps[len(a.gameSteps)-1].StepGame
	}
	return a.initialGame
}

type gameStepParser struct {
	alternatives        []gameAlternative
	isSuccess           bool
	parsedGame          gameAlternative
	possibleNextActions []core.Action
}

func newGameStepParser(initialGame core.Game) *gameStepParser {
	return &gameStepParser{
		alternatives: []gameAlternative{{initialGame: initialGame}},
	}
}

func (p *gameStepParser) next(ap actionPattern, actionString string) bool {
	newAlternatives := []gameAlternative{}
	for _, alternative := range p.alternatives {
		for _, action := range alternative.currentGame().Actions {
			if ap.isMatch(action) {
				newGame := alternative.currentGame().DoAction(action)
				if ap.isCheck != nil && newGame.IsCheck != *ap.isCheck {
					continue
				}
				if ap.isCheckmate != nil && newGame.IsCheckmate != *ap.isCheckmate {
					continue
				}
				newAlternative := alternative.clone()
				newAlternative.gameSteps = append(newAlternative.gameSteps, core.GameStep{StepString: actionString, StepAction: action, StepGame: newGame})
				newAlternatives = append(newAlternatives, newAlternative)
			}
		}
	}
	// fmt.Printf("Alternatives for %v: %v\n", actionString, len(newAlternatives))
	// for _, a := range newAlternatives {
	// 	fmt.Println(a.currentGame())
	// 	fmt.Println()
	// }

	if len(newAlternatives) == 0 {
		actionSet := map[core.Action]struct{}{}
		for _, alternative := range p.alternatives {
			for _, action := range alternative.currentGame().Actions {
				if _, ok := actionSet[action]; ok {
					continue
				}
				p.possibleNextActions = append(p.possibleNextActions, action)
				actionSet[action] = struct{}{}
			}
		}
		p.isSuccess = false
		p.parsedGame = p.alternatives[0]
		return false
	}
	p.isSuccess = true
	p.alternatives = newAlternatives
	p.parsedGame = p.alternatives[0]
	p.possibleNextActions = []core.Action{}
	return true
}

type tokenMatch struct {
	match string
	ap    *actionPattern
	ch    Characteristics
}

type NotationParser struct {
	s          string
	stepParser *gameStepParser

	transitions           map[string]map[string]func([]string, core.Game) []tokenMatch
	evolveCharacteristics func(ch Characteristics, sc Characteristics) (Characteristics, error)
	characteristics       Characteristics
}

func newNotationParser(
	transitions map[string]map[string]func([]string, core.Game) []tokenMatch,
	evolveCharacteristics func(ch Characteristics, sc Characteristics) (Characteristics, error),
	initialCharacteristics Characteristics) *NotationParser {
	return &NotationParser{
		transitions:           transitions,
		evolveCharacteristics: evolveCharacteristics,
		characteristics:       initialCharacteristics,
	}
}

func (p *NotationParser) Parse(initialGame core.Game, s string) ([]core.GameStep, error) {
	p.stepParser = newGameStepParser(initialGame)
	p.s = s

	// Compile all regexes
	rxs := map[string]*regexp.Regexp{}
	for step := range p.transitions {
		for srx := range p.transitions[step] {
			rxs[srx] = regexp.MustCompile(fmt.Sprintf("^%v", srx))
		}
	}

	stepOrder := []string{"full_move_start", "move", "half_move_separator", "move", "full_move_separator"}
	stepI := 0
	i := 0
	for i < len(p.s) {
		// Calculate all tokens that match
		var tokenMatches []tokenMatch
		for rx, fs := range p.transitions[stepOrder[stepI]] {
			matches := rxs[rx].FindStringSubmatch(p.s[i:])
			if matches != nil {
				// fmt.Println("At", p.s[i:], "matched", rx, "with", matches[0], "for", stepOrder[stepI])
				tokenMatches = append(tokenMatches, fs(matches, p.stepParser.parsedGame.currentGame())...)
			}
		}

		// Bail if no token matches
		if len(tokenMatches) == 0 {
			err := fmt.Errorf("at index %v [%v] didn't match any token", i, p.s[i:])
			return p.stepParser.parsedGame.gameSteps, err
		}

		tokenMatch := tokenMatches[0]

		// Move steps will advance the game
		if stepOrder[stepI] == "move" {
			var ok bool
			for _, tm := range tokenMatches {
				if ok = p.stepParser.next(*tm.ap, tm.match); ok {
					tokenMatch = tm
					break // Many regexes may match the token, but only one should match any actions
				}
			}
			if !ok {
				err := fmt.Errorf("at %v matched token %v but no valid action found for it; options were: %v", i, tokenMatch.match, p.stepParser.possibleNextActions)
				return p.stepParser.parsedGame.gameSteps, err
			}
		}

		// Calculate characteristics of the notation as we go through the string.
		// Two things here:
		// 1. For a generic parser, any variation is valid, e.g. 0-0 castling, but if then we find an O-O that's an error.
		// 2. A custom parser can already set that castling has to be e.g. O-O, so that if we find 0-0 that's an error.
		newCharacteristics, err := p.evolveCharacteristics(p.characteristics, tokenMatch.ch)
		if err != nil {
			return p.stepParser.parsedGame.gameSteps, err
		}
		p.characteristics = newCharacteristics

		// Advance the parser to the next token
		i += len(tokenMatch.match)

		// Cycle the step
		stepI++
		if stepI >= len(stepOrder) {
			stepI = 0
		}
	}
	if len(p.stepParser.parsedGame.gameSteps) == 0 {
		return p.stepParser.parsedGame.gameSteps, fmt.Errorf("found 0 game steps")
	}
	return p.stepParser.parsedGame.gameSteps, nil
}
