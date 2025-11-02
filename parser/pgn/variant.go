package pgn

import (
	"fmt"
	"strings"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
)

// VariantPGN implements PGN (Portable Game Notation) parsing variant.
type VariantPGN struct{}

// NewVariantPGN creates a new PGN variant.
func NewVariantPGN() *VariantPGN {
	return &VariantPGN{}
}

// Initialize implements ParserVariant.Initialize for PGN.
func (p *VariantPGN) Initialize(initialGame core.Game, s string) (*parser.ParsingGame, error) {
	tagPairs, remaining, err := extractTagPairs(s)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tag pairs: %w", err)
	}

	alternatives := []parser.GameAlternative{
		{
			InitialGame: initialGame,
			GameSteps:   []core.GameStep{},
		},
	}

	return &parser.ParsingGame{
		Alternatives: alternatives,
		Metadata:     tagPairs,
		Remaining:    remaining,
		Pos:          0,
	}, nil
}

// PopHalfMove implements ParserVariant.PopHalfMove for PGN.
// It handles all token peeking and position management internally.
// Only the token extraction logic is PGN-specific; game state management is handled generically.
func (p *VariantPGN) PopHalfMove(pg *parser.ParsingGame) (*parser.Token, bool, error) {
	// First, skip any move numbers
	for {
		tokenValue, tokenType, newPos, err := peekToken(pg.Remaining, pg.Pos)
		if err != nil {
			return nil, false, err
		}

		// If we've reached the end, break
		if tokenValue == "" && newPos >= len(pg.Remaining) {
			return nil, false, nil
		}

		if tokenType == tokenTypeMoveNumber {
			// Skip move numbers, but check if position actually advanced
			if newPos <= pg.Pos {
				// Position didn't advance, break to avoid infinite loop
				break
			}
			pg.Pos = newPos
			continue
		}
		break
	}

	// Now peek at the next token - should be a half move or result
	tokenValue, tokenType, newPos, err := peekToken(pg.Remaining, pg.Pos)
	if err != nil {
		return nil, false, err
	}

	// If no token or not a half move/result, we're done
	if tokenValue == "" || (tokenType != tokenTypeHalfMove && tokenType != tokenTypeResult) {
		return nil, false, nil
	}

	// Consume the token
	pg.Pos = newPos

	// Convert internal token type to generic token type
	var genericType parser.TokenType
	if tokenType == tokenTypeHalfMove {
		genericType = parser.TokenTypeHalfMove
		// Strip move annotations from the end of the token
		// Valid annotations: ?, !, ??, !!, ?!, !?
		tokenValue = stripMoveAnnotations(tokenValue)
	} else if tokenType == tokenTypeResult {
		genericType = parser.TokenTypeResult
	}

	token := &parser.Token{
		Value: tokenValue,
		Type:  genericType,
	}

	// Now process any comments or annotations that follow the half move
	for {
		nextTokenValue, nextType, nextPos, err := peekToken(pg.Remaining, pg.Pos)
		if err != nil {
			return nil, false, err
		}

		// If we've reached the end (empty token), we're done
		if nextTokenValue == "" && nextPos >= len(pg.Remaining) {
			return token, false, nil
		}

		// If it's a comment or annotation, consume it
		if nextType == tokenTypeAnnotation || nextType == tokenTypeComment || nextType == tokenTypeCurlyComment {
			// Consume the comment/annotation (for now, just skip it)
			// In a later version, they will be added to the action
			// Check if position actually advanced to avoid infinite loop
			if nextPos <= pg.Pos {
				return token, false, nil
			}
			pg.Pos = nextPos
			continue
		}

		// If it's a move number, skip it
		if nextType == tokenTypeMoveNumber {
			// Check if position actually advanced to avoid infinite loop
			if nextPos <= pg.Pos {
				return token, false, nil
			}
			pg.Pos = nextPos
			continue
		}

		// We've reached something else - peek at what's next
		// If the next token is another half move or result, return true (there's more to process)
		if nextType == tokenTypeHalfMove || nextType == tokenTypeResult {
			return token, true, nil
		}

		// Otherwise, we're done (no more half moves)
		return token, false, nil
	}
}

// Finalize implements ParserVariant.Finalize for PGN.
func (p *VariantPGN) Finalize(pg *parser.ParsingGame) error {
	// Does nothing for now
	return nil
}

// ActionToStringVariants implements ParserVariant.ActionToStringVariants for PGN.
func (p *VariantPGN) ActionToStringVariants(a core.Action, g core.Game) []string {
	if a.IsKingsideCastle {
		return []string{"O-O", "O-O+", "O-O#"}
	}
	if a.IsQueensideCastle {
		return []string{"O-O-O", "O-O-O+", "O-O-O#"}
	}
	capture := ""
	if a.IsCapture {
		capture = "x"
	}

	_g := g.Clone()
	__g := _g.DoAction(a)
	check, checkMate := "", ""
	if __g.IsCheck {
		check = "+"
	}
	if __g.IsCheckmate {
		check = ""
		checkMate = "#"
	}

	if a.IsPromotion {
		return []string{a.ToXY.ToAlgebraic() + "=" + a.PromotionPieceType.ToAlgebraic() + check + checkMate}
	}

	variants := []string{}

	pieceOptions := []string{a.FromPiece.PieceType.ToAlgebraic()}
	if a.FromPiece.PieceType == core.PiecePawn {
		pieceOptions = append(pieceOptions, "P")
	}
	pieceXYOptions := []string{"", a.FromPiece.XY.ToAlgebraic()[0:1], a.FromPiece.XY.ToAlgebraic()[1:2], a.FromPiece.XY.ToAlgebraic()}

	for _, pieceOption := range pieceOptions {
		for _, pieceXYOption := range pieceXYOptions {
			variants = append(variants,
				fmt.Sprintf(
					"%s%s%s%s%s%s",
					pieceOption,
					pieceXYOption,
					capture,
					a.ToXY.ToAlgebraic(),
					check,
					checkMate,
				),
			)
		}
	}

	return variants
}

// extractTagPairs extracts all tag pairs from the PGN string and returns them along with
// the remaining string (with tag pairs removed).
func extractTagPairs(pgn string) (map[string]string, string, error) {
	tagPairs := make(map[string]string)
	var result strings.Builder
	result.Grow(len(pgn))

	state := stateNormal
	pos := 0

	// Track tag name and value
	var tagName strings.Builder
	var tagValue strings.Builder
	inTagValue := false

	for pos < len(pgn) {
		char := pgn[pos]

		switch state {
		case stateNormal:
			if char == '[' {
				state = stateTagPair
				tagName.Reset()
				tagValue.Reset()
				inTagValue = false
				pos++
			} else {
				result.WriteByte(char)
				pos++
			}
		case stateTagPair:
			if char == ']' {
				// End of tag pair
				if tagName.Len() > 0 {
					// Remove quotes from tag value if present
					value := tagValue.String()
					if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
						value = value[1 : len(value)-1]
					}
					tagPairs[tagName.String()] = value
				}
				state = stateNormal
				pos++
			} else if char == '"' {
				inTagValue = !inTagValue
				pos++
			} else if !inTagValue && char != ' ' {
				tagName.WriteByte(char)
				pos++
			} else if inTagValue {
				tagValue.WriteByte(char)
				pos++
			} else {
				pos++
			}
		}
	}

	return tagPairs, strings.TrimSpace(result.String()), nil
}

// peekToken peeks at the next token in the PGN string without consuming it.
// Returns the token string, token type, new position, and error.
func peekToken(pgn string, pos int) (string, tokenType, int, error) {
	// Skip whitespace
	startPos := pos
	for startPos < len(pgn) && (pgn[startPos] == ' ' || pgn[startPos] == '\t' || pgn[startPos] == '\n' || pgn[startPos] == '\r') {
		startPos++
	}

	if startPos >= len(pgn) {
		return "", tokenTypeMoveNumber, startPos, nil // No more tokens
	}

	char := pgn[startPos]
	pos = startPos

	// Check for move numbers: 1., 2., etc. or 1... (for black moves)
	if char >= '0' && char <= '9' {
		for pos < len(pgn) && pgn[pos] >= '0' && pgn[pos] <= '9' {
			pos++
		}
		if pos < len(pgn) && pgn[pos] == '.' {
			pos++
			// Check for ellipsis notation: 1... (three dots for black moves)
			if pos < len(pgn) && pgn[pos] == '.' {
				pos++
				if pos < len(pgn) && pgn[pos] == '.' {
					pos++
				}
			}
			// Skip any whitespace after the dot(s)
			for pos < len(pgn) && (pgn[pos] == ' ' || pgn[pos] == '\t') {
				pos++
			}
			return "", tokenTypeMoveNumber, pos, nil
		}
		// If no dot, it might be part of a half move (like "1e4" - unusual but possible)
		// Reset pos and continue
		pos = startPos
	}

	// Check for annotations: ($1), ($2), etc.
	if char == '(' && pos+1 < len(pgn) && pgn[pos+1] == '$' {
		pos += 2 // Skip "($"
		for pos < len(pgn) && pgn[pos] >= '0' && pgn[pos] <= '9' {
			pos++
		}
		if pos < len(pgn) && pgn[pos] == ')' {
			pos++
			return pgn[startPos:pos], tokenTypeAnnotation, pos, nil
		}
	}

	// Check for curly comments: {comment}
	if char == '{' {
		pos++
		for pos < len(pgn) && pgn[pos] != '}' {
			pos++
		}
		if pos < len(pgn) {
			pos++ // Include the closing brace
		}
		return pgn[startPos:pos], tokenTypeCurlyComment, pos, nil
	}

	// Check for semicolon comments: ; comment to end of line
	if char == ';' {
		pos++
		for pos < len(pgn) && pgn[pos] != '\n' {
			pos++
		}
		if pos < len(pgn) {
			pos++ // Include the newline
		}
		return pgn[startPos:pos], tokenTypeComment, pos, nil
	}

	// Check for result markers: 1-0, 0-1, 1/2-1/2, *
	if char == '*' {
		return "*", tokenTypeResult, pos + 1, nil
	}
	if pos+2 < len(pgn) {
		threeChar := pgn[pos : pos+3]
		if threeChar == "1-0" || threeChar == "0-1" {
			return threeChar, tokenTypeResult, pos + 3, nil
		}
		if pos+6 < len(pgn) {
			sevenChar := pgn[pos : pos+7]
			if sevenChar == "1/2-1/2" {
				return sevenChar, tokenTypeResult, pos + 7, nil
			}
		}
	}

	// It's a half move - collect until whitespace or end
	moveEnd := pos
	for moveEnd < len(pgn) && pgn[moveEnd] != ' ' && pgn[moveEnd] != '\t' && pgn[moveEnd] != '\n' && pgn[moveEnd] != '\r' {
		moveEnd++
	}

	return pgn[startPos:moveEnd], tokenTypeHalfMove, moveEnd, nil
}

// stripMoveAnnotations removes PGN move annotations from the end of a move string.
// Valid annotations: ?, !, ??, !!, ?!, !?
func stripMoveAnnotations(move string) string {
	if len(move) == 0 {
		return move
	}

	// Check for two-character annotations first
	if len(move) >= 2 {
		lastTwo := move[len(move)-2:]
		if lastTwo == "??" || lastTwo == "!!" || lastTwo == "?!" || lastTwo == "!?" {
			return move[:len(move)-2]
		}
	}

	// Check for single-character annotations
	if len(move) >= 1 {
		lastChar := move[len(move)-1]
		if lastChar == '?' || lastChar == '!' {
			return move[:len(move)-1]
		}
	}

	return move
}

type parserState int

const (
	stateNormal parserState = iota
	stateTagPair
	stateCurlyComment
	stateSemicolonComment
	stateAnnotation
)

type tokenType int

const (
	tokenTypeHalfMove tokenType = iota
	tokenTypeAnnotation
	tokenTypeComment
	tokenTypeCurlyComment
	tokenTypeMoveNumber
	tokenTypeResult
)
