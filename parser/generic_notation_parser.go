package parser

import (
	"fmt"

	"github.com/marianogappa/cheesse/core"
)

// TokenType represents the type of a parsed token.
type TokenType int

const (
	TokenTypeHalfMove TokenType = iota
	TokenTypeResult
	TokenTypeComment
	TokenTypeAnnotation
	TokenTypeMoveNumber
)

// Token represents a parsed token from the notation string.
type Token struct {
	Value string
	Type  TokenType
	// Comment holds the text of any comment(s) attached to this token (e.g. PGN
	// {...} or ; comments following a move), without delimiters.
	Comment string
}

// ParserVariant defines the interface for notation-specific parsing logic.
// Each notation type (PGN, Algebraic, Descriptive, etc.) implements this interface.
type ParserVariant interface {
	// Initialize extracts notation-specific metadata (like tag pairs) and sets up the initial parsing state.
	Initialize(initialGame core.Game, s string) (*ParsingGame, error)

	// PopHalfMove pops the next half move token and any following comments or annotations.
	// It handles all token peeking and position management internally.
	// Returns:
	//   - token: the half move or result token (nil if no more half moves)
	//   - hasMore: true if another half move is available after this one, false otherwise
	//   - error: any error that occurred during token extraction
	PopHalfMove(pg *ParsingGame) (*Token, bool, error)

	// Finalize completes the parsing process (does nothing for now, but reserved for future use).
	Finalize(pg *ParsingGame) error

	// ActionToStringVariants converts an action to all possible string representations in this notation.
	// For example, in PGN, an action might have multiple representations like "e4", "Pe4", "P4e4", etc.
	ActionToStringVariants(a core.Action, g core.Game) []string
}

// ParsingGame represents a game in the process of being parsed.
// It contains the game alternatives and other parsing state.
type ParsingGame struct {
	Alternatives []GameAlternative
	Metadata     map[string]string
	Remaining    string // Remaining notation string after metadata removed
	Pos          int    // Current position in remaining string
}

// ParsedGame represents the final parsed result.
type ParsedGame struct {
	GameSteps []core.GameStep
	Metadata  map[string]string
}

// Build constructs the final ParsedGame from the ParsingGame.
// If multiple branches remain, it picks the first one.
func (pg *ParsingGame) Build() *ParsedGame {
	var gameSteps []core.GameStep
	if len(pg.Alternatives) > 0 {
		gameSteps = pg.Alternatives[0].GameSteps
	}
	return &ParsedGame{
		GameSteps: gameSteps,
		Metadata:  pg.Metadata,
	}
}

// GenericNotationParser is a generic parser that works with any notation variant.
type GenericNotationParser struct {
	variant ParserVariant
}

// NewGenericNotationParser creates a new generic notation parser with the specified variant.
func NewGenericNotationParser(variant ParserVariant) *GenericNotationParser {
	return &GenericNotationParser{
		variant: variant,
	}
}

// Parse performs the complete 3-step parsing process:
// 1. Initialize - extracts metadata and sets up parsing state
// 2. Loop through ParseHalfMove - processes all half moves
// 3. Finalize - completes parsing
func (p *GenericNotationParser) Parse(initialGame core.Game, s string) (*ParsedGame, error) {
	// Step 1: Initialize
	parsingGame, err := p.variant.Initialize(initialGame, s)
	if err != nil {
		return nil, err
	}

	// Step 2: Loop through parseHalfMoves
	// The variant pops tokens, and we handle game state management.
	// On a mid-game failure, parsing stops but the valid prefix of steps is
	// returned along with the error, matching the other parsers' behavior.
	for {
		token, hasMore, err := p.variant.PopHalfMove(parsingGame)
		if err != nil {
			return parsingGame.Build(), err
		}
		if token == nil {
			break
		}

		// Process the token based on its type
		if err := p.ProcessToken(parsingGame, token); err != nil {
			return parsingGame.Build(), err
		}

		if !hasMore {
			break
		}
	}

	// Step 3: Finalize
	if err := p.variant.Finalize(parsingGame); err != nil {
		return nil, err
	}

	return parsingGame.Build(), nil
}

// MatchHalfMove attempts to match a half move string against all possible actions in the current game.
// It uses the variant's ActionToStringVariants method to convert actions to notation strings.
// Returns a slice of actions that match the half move notation.
func (p *GenericNotationParser) MatchHalfMove(halfMove string, g core.Game) ([]core.Action, error) {
	// Get all possible actions from the current game state
	allActions := g.Actions
	if len(allActions) == 0 {
		// Actions is not populated - this can happen for the initial game
		return nil, fmt.Errorf("game Actions field is not populated - ensure the game has Actions calculated (e.g., by calling a method that populates it)")
	}

	// Build a map from action string to slice of action pointers
	// Multiple actions can yield the same notation string, so we need a map from string to []*Action
	// Using pointers to avoid duplicating Action structs in memory
	actionStringMap := make(map[string][]*core.Action)

	// For each action, get all its notation variant strings and add them to the map
	for i := range allActions {
		action := &allActions[i]
		variants := p.variant.ActionToStringVariants(*action, g)

		for _, variant := range variants {
			actionStringMap[variant] = append(actionStringMap[variant], action)
		}
	}

	// Try to find a match for the half move
	matchingActions, found := actionStringMap[halfMove]
	if !found {
		return nil, fmt.Errorf("no action matches half move %q", halfMove)
	}

	// Deduplicate actions (same action might appear via multiple variant strings)
	// Convert pointers to values and deduplicate
	seen := make(map[core.Action]bool)
	result := []core.Action{}
	for _, actionPtr := range matchingActions {
		action := *actionPtr
		if !seen[action] {
			seen[action] = true
			result = append(result, action)
		}
	}

	return result, nil
}

// ProcessToken processes a token and updates the parsing game state.
// This handles all the generic logic: matching actions, branching, and applying moves.
func (p *GenericNotationParser) ProcessToken(pg *ParsingGame, token *Token) error {
	switch token.Type {
	case TokenTypeHalfMove:
		// Process the half move - this is the generic branching logic
		newAlternatives := []GameAlternative{}

		// For each current alternative, try to match the half move
		for _, alt := range pg.Alternatives {
			currentGame := alt.CurrentGame()
			matches, err := p.MatchHalfMove(token.Value, currentGame)
			if err != nil {
				// No match for this alternative - branch dies (not an error, there may be other valid branches)
				continue
			}

			// Create a new alternative for each matching action
			for _, action := range matches {
				newGame := currentGame.DoAction(action)
				newAlternative := alt.Clone()
				newAlternative.GameSteps = append(newAlternative.GameSteps, core.GameStep{
					StepString:      token.Value,
					StepComment:     token.Comment,
					StepAction:      action,
					StepGame:        newGame,
					StepPreMoveGame: currentGame,
				})

				newAlternatives = append(newAlternatives, newAlternative)
			}
		}

		// If all branches died, return an error
		if len(newAlternatives) == 0 {
			return fmt.Errorf("could not match half move %q against any valid action", token.Value)
		}

		pg.Alternatives = newAlternatives

	case TokenTypeResult:
		// Result tokens (1-0, 0-1, 1/2-1/2, *) should be added as half moves
		newAlternatives := []GameAlternative{}
		for _, alt := range pg.Alternatives {
			newAlternative := alt.Clone()
			newAlternative.GameSteps = append(newAlternative.GameSteps, core.GameStep{
				StepString:      token.Value,
				StepComment:     token.Comment,
				StepAction:      core.Action{},     // Empty action for result markers
				StepGame:        alt.CurrentGame(), // Game state doesn't change
				StepPreMoveGame: alt.CurrentGame(),
			})
			newAlternatives = append(newAlternatives, newAlternative)
		}
		pg.Alternatives = newAlternatives

	default:
		// Comments, annotations, move numbers are handled by the variant during token popping
		return fmt.Errorf("unexpected token type: %v", token.Type)
	}

	return nil
}
