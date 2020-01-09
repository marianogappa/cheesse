package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type parserAlgebraic struct {
	s      string
	parser *gameStepParser

	characteristics characteristics
}

var (
	errInvalidRegexp             = fmt.Errorf("invalid regexp; this should not happen")
	errDidntMatchAnyRegexp       = fmt.Errorf("invalid state at step TODO: didn't match any regexp")
	errParsedInexistentCheck     = fmt.Errorf("parsed action indicates a check, but there's none")
	errParsedInexistentCheckmate = fmt.Errorf("parsed action indicates a checkmate, but there's none")
)

func stringToPieceType(s string) pieceType {
	return map[string]pieceType{
		"Q": pieceQueen,
		"K": pieceKing,
		"B": pieceBishop,
		"N": pieceKnight,
		"R": pieceRook,
		"":  piecePawn,
	}[s]
}

func fileToPInt(file string) *int {
	if file == "" {
		return nil
	}
	v := int(file[0] - 'a')
	return &v
}

func rankToPInt(rank string) *int {
	if rank == "" {
		return nil
	}
	v := (8 - int(rank[0]-'0'))
	return &v
}

func pBool(b bool) *bool {
	return &b
}

type tokenMatch struct {
	match string
	ap    *actionPattern
	ch    characteristics
}

type characteristics struct {
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

func (c characteristics) withThreatenSymbol(threatenSymbol string) characteristics {
	switch threatenSymbol {
	case "+", "†", "ch", "++", "dblch", "dbl ch", "dbl.ch", "disch", "dis ch", "dis.ch":
		c.isCheck = true
		c.isCheckmate = false
		c.usesCheckSymbol = &threatenSymbol
		c.usesCheckmateSymbol = nil
	case "#", "mate", "‡", "≠", "X", "x", "×":
		c.isCheck = false
		c.isCheckmate = true
		c.usesCheckSymbol = nil
		c.usesCheckmateSymbol = &threatenSymbol
	default:
		c.isCheck = false
		c.isCheckmate = false
		c.usesCheckSymbol = nil
		c.usesCheckmateSymbol = nil
	}
	return c
}

var (
	parserAlgebraicTransitions = map[string]map[string]func([]string) tokenMatch{
		"full_move_start": {
			`[\t\f\r ]*([0-9]+)?(\.)?[\t\f\r ]*`: func(ms []string) tokenMatch {
				var fullMoveNumber *int
				if len(ms[1]) > 0 {
					fmn, _ := strconv.Atoi(ms[1])
					fullMoveNumber = &fmn
				}
				var usesFullMoveDot *bool
				if len(ms[2]) == 1 {
					usesFullMoveDot = pBool(true)
				}
				return tokenMatch{ms[0], nil, characteristics{fullMoveNumber: fullMoveNumber, usesFullMoveDot: usesFullMoveDot}}
			},
		},
		"half_move_separator": {
			`[\t\f\r ]+`: func(ms []string) tokenMatch {
				return tokenMatch{ms[0], nil, characteristics{}}
			},
		},
		"full_move_separator": {
			`([\t\f\r ]*?\n|[\t\f\r ]+)`: func(ms []string) tokenMatch {
				var usesNewlineAsFullMoveSeparator *bool
				if strings.Contains(ms[0], "\n") {
					usesNewlineAsFullMoveSeparator = pBool(true)
				}
				return tokenMatch{ms[0], nil, characteristics{usesNewlineAsFullMoveSeparator: usesNewlineAsFullMoveSeparator}}
			},
		},
		"move": {
			// Move
			`([QKBNR]?)([a-h])?([1-8])?([a-h])([1-8])(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				sFromPieceType, fromSquareFile, fromSquareRank, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
				ap := actionPattern{
					fromPieceType:      stringToPieceType(sFromPieceType),
					fromX:              fileToPInt(fromSquareFile),
					fromY:              rankToPInt(fromSquareRank),
					toX:                fileToPInt(toSquareFile),
					toY:                rankToPInt(toSquareRank),
					isCapture:          pBool(false),
					isPromotion:        pBool(false),
					isCastle:           pBool(false),
					isResign:           pBool(false),
					isEnPassantCapture: pBool(false),
				}
				ch := characteristics{}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// Capture
			`([QKBNR])([a-h])?([1-8])?(x|:)?([a-h])([1-8])(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				sFromPieceType, fromSquareFile, fromSquareRank, _, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8]
				ap := actionPattern{
					fromPieceType:      stringToPieceType(sFromPieceType),
					fromX:              fileToPInt(fromSquareFile),
					fromY:              rankToPInt(fromSquareRank),
					toX:                fileToPInt(toSquareFile),
					toY:                rankToPInt(toSquareRank),
					isCapture:          pBool(true),
					capturedPieceX:     fileToPInt(toSquareFile),
					capturedPieceY:     rankToPInt(toSquareRank),
					isPromotion:        pBool(false),
					isCastle:           pBool(false),
					isResign:           pBool(false),
					isEnPassantCapture: pBool(false),
				}
				ch := characteristics{}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// Capture with colon at the end
			`([QKBNR])([a-h])?([1-8])?([a-h])([1-8]):(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				sFromPieceType, fromSquareFile, fromSquareRank, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
				ap := actionPattern{
					fromPieceType:      stringToPieceType(sFromPieceType),
					fromX:              fileToPInt(fromSquareFile),
					fromY:              rankToPInt(fromSquareRank),
					toX:                fileToPInt(toSquareFile),
					toY:                rankToPInt(toSquareRank),
					isCapture:          pBool(true),
					capturedPieceX:     fileToPInt(toSquareFile),
					capturedPieceY:     rankToPInt(toSquareRank),
					isPromotion:        pBool(false),
					isCastle:           pBool(false),
					isResign:           pBool(false),
					isEnPassantCapture: pBool(false),
				}
				ch := characteristics{}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// Capture with pawn, potentially without rank
			`([a-h])(x|:)?([a-h])([1-8]?)( ?e.p.)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				fromSquareFile, _, toSquareFile, toSquareRank, enPassantCapture, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
				ap := actionPattern{
					fromPieceType:      stringToPieceType(""),
					fromX:              fileToPInt(fromSquareFile),
					toX:                fileToPInt(toSquareFile),
					toY:                rankToPInt(toSquareRank),
					isEnPassantCapture: pBool(strings.HasSuffix(enPassantCapture, "e.p.")),
					isCapture:          pBool(true),
					isPromotion:        pBool(false),
					isCastle:           pBool(false),
					isResign:           pBool(false),
				}
				ch := characteristics{}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// Capture and promotion with pawn, potentially without rank
			`([a-h])(x|:)?([a-h])([1-8]?)([=\(])([QBNR])\)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				fromSquareFile, _, toSquareFile, toSquareRank, promotionSymbol, sPromotionPieceType, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8]
				ap := actionPattern{
					fromPieceType:      stringToPieceType(""),
					fromX:              fileToPInt(fromSquareFile),
					toX:                fileToPInt(toSquareFile),
					toY:                rankToPInt(toSquareRank),
					isEnPassantCapture: pBool(false),
					isCapture:          pBool(true),
					isPromotion:        pBool(true),
					isCastle:           pBool(false),
					isResign:           pBool(false),
					promotionPieceType: stringToPieceType(sPromotionPieceType),
				}
				ch := characteristics{usesPromotionSymbol: &promotionSymbol}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// Promotion
			`([a-h])([1-8])([=\(])([QBNR])\)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				toSquareFile, toSquareRank, promotionSymbol, sPromotionPieceType, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6]
				ap := actionPattern{
					fromPieceType:      stringToPieceType(""),
					toX:                fileToPInt(toSquareFile),
					toY:                rankToPInt(toSquareRank),
					isCapture:          pBool(false),
					isPromotion:        pBool(true),
					isCastle:           pBool(false),
					isResign:           pBool(false),
					isEnPassantCapture: pBool(false),
					promotionPieceType: stringToPieceType(sPromotionPieceType),
				}
				ch := characteristics{usesPromotionSymbol: &promotionSymbol}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// Castling
			`(0-0|0-0-0|O-O|O-O-O)(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				castlingSymbol, threatenSymbol, _ := ms[1], ms[2], ms[3]
				ap := actionPattern{
					isCastle:           pBool(true),
					isQueensideCastle:  pBool(castlingSymbol == "0-0-0" || castlingSymbol == "O-O-O"),
					isKingsideCastle:   pBool(castlingSymbol == "0-0" || castlingSymbol == "O-O"),
					isCapture:          pBool(false),
					isPromotion:        pBool(false),
					isResign:           pBool(false),
					isEnPassantCapture: pBool(false),
				}
				cs := string(castlingSymbol[0])
				ch := characteristics{usesCastlingSymbol: &cs}.withThreatenSymbol(threatenSymbol)
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},

			// End of game
			`(1–0|0–1|½–½|resigns|White resigns|Black resigns)`: func(ms []string) tokenMatch {
				var usesEndGameSymbol string
				switch ms[1] {
				case "1-0", "0-1", "½–½":
					usesEndGameSymbol = "numbers"
				case "resigns":
					usesEndGameSymbol = "resigns"
				case "White resigns", "Black resigns":
					usesEndGameSymbol = "color resigns"
				}
				ap := actionPattern{
					isResign:           pBool(strings.HasSuffix(usesEndGameSymbol, "resigns")),
					isPromotion:        pBool(false),
					isCastle:           pBool(false),
					isCapture:          pBool(false),
					isEnPassantCapture: pBool(false),
				}
				ch := characteristics{usesEndGameSymbol: &usesEndGameSymbol}
				if ch.isCheck {
					ap.isCheck = pBool(true)
				}
				if ch.isCheckmate {
					ap.isCheckmate = pBool(true)
				}
				return tokenMatch{ms[0], &ap, ch}
			},
		},
	}
)

func newParserAlgebraic(initialGame game, s string) *parserAlgebraic {
	return &parserAlgebraic{s: s, parser: newGameStepParser(initialGame)}
}

func (p *parserAlgebraic) parse() ([]gameStep, error) {
	// Compile all regexes
	rxs := map[string]*regexp.Regexp{}
	for step := range parserAlgebraicTransitions {
		for srx := range parserAlgebraicTransitions[step] {
			rxs[srx] = regexp.MustCompile(fmt.Sprintf("^%v", srx))
		}
	}

	stepOrder := []string{"full_move_start", "move", "half_move_separator", "move", "full_move_separator"}
	stepI := 0
	i := 0
	for i < len(p.s) {
		// Calculate all tokens that match
		var tokenMatches []tokenMatch
		for rx, fs := range parserAlgebraicTransitions[stepOrder[stepI]] {
			matches := rxs[rx].FindStringSubmatch(p.s[i:])
			if matches != nil {
				tokenMatches = append(tokenMatches, fs(matches))
			}
		}

		// Bail if no token matches
		if len(tokenMatches) == 0 {
			err := fmt.Errorf("at index %v [%v] didn't match any token", i, p.s[i:])
			return p.parser.parsedGame.gameSteps, err
		}

		tokenMatch := tokenMatches[0]

		// Move steps will advance the game
		if stepOrder[stepI] == "move" {
			var ok bool
			for _, tm := range tokenMatches {
				if ok = p.parser.next(*tm.ap, tm.match); ok {
					tokenMatch = tm
					break // Many regexes may match the token, but only one should match any actions
				}
			}
			if !ok {
				err := fmt.Errorf("at %v matched token %v but no valid action found for it; options were: %v", i, tokenMatch.match, p.parser.possibleNextActions)
				return p.parser.parsedGame.gameSteps, err
			}
		}

		// Calculate characteristics of the notation as we go through the string.
		// Two things here:
		// 1. For a generic parser, any variation is valid, e.g. 0-0 castling, but if then we find an O-O that's an error.
		// 2. A custom parser can already set that castling has to be e.g. O-O, so that if we find 0-0 that's an error.
		if err := p.updateStepCharacteristics(tokenMatch.ch); err != nil {
			return p.parser.parsedGame.gameSteps, err
		}

		// Advance the parser to the next token
		i += len(tokenMatch.match)

		// Cycle the step
		stepI++
		if stepI >= len(stepOrder) {
			stepI = 0
		}
	}
	return p.parser.parsedGame.gameSteps, nil
}

// TODO human-readable error messages here. Also, lacking some context.
func (p *parserAlgebraic) updateStepCharacteristics(sc characteristics) error {
	if sc.isCheck && !p.parser.parsedGame.currentGame().isCheck {
		return errParsedInexistentCheck
	}
	if sc.isCheckmate && !p.parser.parsedGame.currentGame().isCheckmate {
		return errParsedInexistentCheckmate
	}
	if sc.usesCheckSymbol != nil {
		if p.characteristics.usesCheckSymbol == nil {
			p.characteristics.usesCheckSymbol = sc.usesCheckSymbol
		} else if *p.characteristics.usesCheckSymbol != *sc.usesCheckSymbol {
			return fmt.Errorf("expecting CheckSymbol %v but found %v", *p.characteristics.usesCheckSymbol, *sc.usesCheckSymbol)
		}
	}
	if sc.usesCheckmateSymbol != nil {
		if p.characteristics.usesCheckmateSymbol == nil {
			p.characteristics.usesCheckmateSymbol = sc.usesCheckmateSymbol
		} else if *p.characteristics.usesCheckmateSymbol != *sc.usesCheckmateSymbol {
			return fmt.Errorf("expecting CheckmateSymbol %v but found %v", *p.characteristics.usesCheckmateSymbol, *sc.usesCheckmateSymbol)
		}
	}
	if sc.usesFullMoveDot != nil {
		if p.characteristics.usesFullMoveDot == nil {
			p.characteristics.usesFullMoveDot = sc.usesFullMoveDot
		} else if *p.characteristics.usesFullMoveDot != *sc.usesFullMoveDot {
			return fmt.Errorf("expecting FullMoveDot %v but found %v", *p.characteristics.usesFullMoveDot, *sc.usesFullMoveDot)
		}
	}
	if sc.usesNewlineAsFullMoveSeparator != nil {
		if p.characteristics.usesNewlineAsFullMoveSeparator == nil {
			p.characteristics.usesNewlineAsFullMoveSeparator = sc.usesNewlineAsFullMoveSeparator
		} else if *p.characteristics.usesNewlineAsFullMoveSeparator != *sc.usesNewlineAsFullMoveSeparator {
			return fmt.Errorf("expecting NewlineAsFullMoveSeparator %v but found %v", *p.characteristics.usesNewlineAsFullMoveSeparator, *sc.usesNewlineAsFullMoveSeparator)
		}
	}
	if sc.usesThreatenSymbol != nil {
		if p.characteristics.usesThreatenSymbol == nil {
			p.characteristics.usesThreatenSymbol = sc.usesThreatenSymbol
		} else if *p.characteristics.usesThreatenSymbol != *sc.usesThreatenSymbol {
			return fmt.Errorf("expecting ThreatenSymbol %v but found %v", *p.characteristics.usesThreatenSymbol, *sc.usesThreatenSymbol)
		}
	}
	if sc.usesCaptureSymbol != nil {
		if p.characteristics.usesCaptureSymbol == nil {
			p.characteristics.usesCaptureSymbol = sc.usesCaptureSymbol
		} else if *p.characteristics.usesCaptureSymbol != *sc.usesCaptureSymbol {
			return fmt.Errorf("expecting CaptureSymbol %v but found %v", *p.characteristics.usesCaptureSymbol, *sc.usesCaptureSymbol)
		}
	}
	if sc.usesEndGameSymbol != nil {
		if p.characteristics.usesEndGameSymbol == nil {
			p.characteristics.usesEndGameSymbol = sc.usesEndGameSymbol
		} else if *p.characteristics.usesEndGameSymbol != *sc.usesEndGameSymbol {
			return fmt.Errorf("expecting EndGameSymbol %v but found %v", *p.characteristics.usesEndGameSymbol, *sc.usesEndGameSymbol)
		}
	}
	if sc.usesPromotionSymbol != nil {
		if p.characteristics.usesPromotionSymbol == nil {
			p.characteristics.usesPromotionSymbol = sc.usesPromotionSymbol
		} else if *p.characteristics.usesPromotionSymbol != *sc.usesPromotionSymbol {
			return fmt.Errorf("expecting PromotionSymbol %v but found %v", *p.characteristics.usesPromotionSymbol, *sc.usesPromotionSymbol)
		}
	}
	if sc.usesCastlingSymbol != nil {
		if p.characteristics.usesCastlingSymbol == nil {
			p.characteristics.usesCastlingSymbol = sc.usesCastlingSymbol
		} else if *p.characteristics.usesCastlingSymbol != *sc.usesCastlingSymbol {
			return fmt.Errorf("expecting CastlingSymbol %v but found %v", *p.characteristics.usesCastlingSymbol, *sc.usesCastlingSymbol)
		}
	}
	// TODO full move number
	return nil
}
