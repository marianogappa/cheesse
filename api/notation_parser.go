package api

import "fmt"

import "strings"

type actionPattern struct {
	fromPieceType      pieceType
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
	promotionPieceType pieceType
	capturedPieceType  pieceType
	capturedPieceX     *int
	capturedPieceY     *int
	isCheck            *bool
	isCheckmate        *bool
}

func (p actionPattern) String() string {
	var sb strings.Builder
	sb.WriteString("Action Pattern:\n")
	if p.fromPieceType != pieceNone {
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
	if p.promotionPieceType != pieceNone {
		sb.WriteString(fmt.Sprintf("{a.promotionPiece.pieceType}:%v\n", p.promotionPieceType))
	}
	if p.capturedPieceType != pieceNone {
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

func (p *actionPattern) isMatch(a action) bool {
	if !pieceTypeMatcher(p.fromPieceType)(a.fromPiece.pieceType) ||
		!intMatcher(p.fromX)(a.fromPiece.xy.x) ||
		!intMatcher(p.fromY)(a.fromPiece.xy.y) ||
		!intMatcher(p.toX)(a.toXY.x) ||
		!intMatcher(p.toY)(a.toXY.y) ||
		!boolMatcher(p.isCapture)(a.isCapture) ||
		!boolMatcher(p.isResign)(a.isResign) ||
		!boolMatcher(p.isPromotion)(a.isPromotion) ||
		!boolMatcher(p.isEnPassant)(a.isEnPassant) ||
		!boolMatcher(p.isEnPassantCapture)(a.isEnPassantCapture) ||
		!boolMatcher(p.isCastle)(a.isCastle) ||
		!boolMatcher(p.isKingsideCastle)(a.isKingsideCastle) ||
		!boolMatcher(p.isQueensideCastle)(a.isQueensideCastle) ||
		!pieceTypeMatcher(p.promotionPieceType)(a.promotionPieceType) ||
		!pieceTypeMatcher(p.capturedPieceType)(a.capturedPiece.pieceType) ||
		!intMatcher(p.capturedPieceX)(a.capturedPiece.xy.x) ||
		!intMatcher(p.capturedPieceY)(a.capturedPiece.xy.y) {
		return false
	}
	return true
}

func boolMatcher(v *bool) func(interface{}) bool {
	return func(w interface{}) bool { return v == nil || *v == w.(bool) }
}
func intMatcher(v *int) func(interface{}) bool {
	return func(w interface{}) bool { return v == nil || *v == w.(int) }
}
func pieceTypeMatcher(v pieceType) func(interface{}) bool {
	return func(w interface{}) bool { return v == pieceNone || v == w.(pieceType) }
}

type gameStep struct {
	s string
	a action
	g game
}

func (s gameStep) clone() gameStep {
	return gameStep{
		s: s.s,
		a: s.a,
		g: s.g.clone(),
	}
}

type gameAlternative struct {
	initialGame game
	gameSteps   []gameStep
}

func (a gameAlternative) clone() gameAlternative {
	clonedGameSteps := make([]gameStep, len(a.gameSteps))
	for i := range a.gameSteps {
		clonedGameSteps[i] = a.gameSteps[i].clone()
	}
	return gameAlternative{
		initialGame: a.initialGame.clone(),
		gameSteps:   clonedGameSteps,
	}
}

func (a gameAlternative) currentGame() game {
	if len(a.gameSteps) > 0 {
		return a.gameSteps[len(a.gameSteps)-1].g
	}
	return a.initialGame
}

type gameStepParser struct {
	alternatives        []gameAlternative
	isSuccess           bool
	parsedGame          gameAlternative
	possibleNextActions []action
}

func newGameStepParser(initialGame game) *gameStepParser {
	return &gameStepParser{
		alternatives: []gameAlternative{gameAlternative{initialGame: initialGame}},
	}
}

func (p *gameStepParser) next(ap actionPattern, actionString string) bool {
	newAlternatives := []gameAlternative{}
	for _, alternative := range p.alternatives {
		for _, action := range alternative.currentGame().actions {
			if ap.isMatch(action) {
				newGame := alternative.currentGame().doAction(action)
				if ap.isCheck != nil && newGame.isCheck != *ap.isCheck {
					continue
				}
				if ap.isCheckmate != nil && newGame.isCheckmate != *ap.isCheckmate {
					continue
				}
				newAlternative := alternative.clone()
				newAlternative.gameSteps = append(newAlternative.gameSteps, gameStep{s: actionString, a: action, g: newGame})
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
		actionSet := map[action]struct{}{}
		for _, alternative := range p.alternatives {
			for _, action := range alternative.currentGame().actions {
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
	p.possibleNextActions = []action{}
	return true
}
