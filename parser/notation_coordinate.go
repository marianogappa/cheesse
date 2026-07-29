package parser

import (
	"fmt"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

func NewNotationParserCoordinate(initialCharacteristics Characteristics) *NotationParser {
	var (
		evolveCharacteristics = func(ch Characteristics, sc Characteristics) (Characteristics, error) {
			if sc.usesNewlineAsFullMoveSeparator != nil {
				if ch.usesNewlineAsFullMoveSeparator == nil {
					ch.usesNewlineAsFullMoveSeparator = sc.usesNewlineAsFullMoveSeparator
				} else if *ch.usesNewlineAsFullMoveSeparator != *sc.usesNewlineAsFullMoveSeparator {
					return ch, fmt.Errorf("expecting newline as full move separator %v but found %v", *ch.usesNewlineAsFullMoveSeparator, *sc.usesNewlineAsFullMoveSeparator)
				}
			}
			if sc.usesCastlingSymbol != nil {
				if ch.usesCastlingSymbol == nil {
					ch.usesCastlingSymbol = sc.usesCastlingSymbol
				} else if *ch.usesCastlingSymbol != *sc.usesCastlingSymbol {
					return ch, fmt.Errorf("expecting CastlingSymbol %v but found %v", *ch.usesCastlingSymbol, *sc.usesCastlingSymbol)
				}
			}
			return ch, nil
		}
		transitions = map[string]map[string]func([]string, core.Game) []tokenMatch{
			"full_move_start": {
				`[\t\f\r ]*([0-9]+)?(\.)?[\t\f\r ]*`: func(ms []string, g core.Game) []tokenMatch {
					return []tokenMatch{{ms[0], nil, Characteristics{}}}
				},
			},
			"half_move_separator": {
				`[\t\f\r ]+`: func(ms []string, g core.Game) []tokenMatch {
					return []tokenMatch{{ms[0], nil, Characteristics{}}}
				},
			},
			"full_move_separator": {
				`([\t\f\r ]*?\n|[\t\f\r ]+)`: func(ms []string, g core.Game) []tokenMatch {
					var usesNewlineAsFullMoveSeparator *bool
					if strings.Contains(ms[0], "\n") {
						usesNewlineAsFullMoveSeparator = pBool(true)
					}
					return []tokenMatch{{ms[0], nil, Characteristics{usesNewlineAsFullMoveSeparator: usesNewlineAsFullMoveSeparator}}}
				},
			},
			"move": {
				// Move or capture: e2-e4, b5xd7, e2e4; optional promotion (Q, =Q, (Q), /Q) and check suffix (+, ch, ++, mate)
				`([a-h])([1-8])(-|x|:)?([a-h])([1-8])(?:([=\(/])?([QBNR])\)?)?(\+\+|mate|#|\+|ch)?`: func(ms []string, g core.Game) []tokenMatch {
					fromFile, fromRank, delimiter, toFile, toRank, promotionSymbol, promotionPiece, threatenSymbol := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8]

					var isCheck, isCheckmate *bool
					switch threatenSymbol {
					case "+", "ch":
						isCheck = pBool(true)
					case "mate", "#":
						isCheckmate = pBool(true)
					case "++":
						// "++" is ambiguous: checkmate in classic usage, double check in modern
						// usage. Accept either by only requiring the move to give check.
						isCheck = pBool(true)
					}

					ap := actionPattern{
						fromX:       fileToPInt(fromFile),
						fromY:       rankToPInt(fromRank),
						toX:         fileToPInt(toFile),
						toY:         rankToPInt(toRank),
						isCastle:    pBool(false),
						isResign:    pBool(false),
						isCheck:     isCheck,
						isCheckmate: isCheckmate,
					}

					// "x" or ":" explicitly marks a capture; "-" or nothing leaves it open
					if delimiter == "x" || delimiter == ":" {
						ap.isCapture = pBool(true)
					}

					if promotionPiece != "" {
						ap.isPromotion = pBool(true)
						ap.promotionPieceType = stringToPieceType(promotionPiece)
					} else {
						ap.isPromotion = pBool(false)
					}

					var ch Characteristics
					if promotionSymbol != "" {
						ch.usesPromotionSymbol = &promotionSymbol
					}
					return []tokenMatch{{ms[0], &ap, ch}}
				},

				// Castling
				`(0-0-0|0-0|O-O-O|O-O)(\+\+|mate|#|\+|ch)?`: func(ms []string, g core.Game) []tokenMatch {
					castlingSymbol, threatenSymbol := ms[1], ms[2]

					var isCheck, isCheckmate *bool
					switch threatenSymbol {
					case "+", "ch", "++":
						isCheck = pBool(true)
					case "mate", "#":
						isCheckmate = pBool(true)
					}

					ap := actionPattern{
						isCastle:          pBool(true),
						isQueensideCastle: pBool(castlingSymbol == "0-0-0" || castlingSymbol == "O-O-O"),
						isKingsideCastle:  pBool(castlingSymbol == "0-0" || castlingSymbol == "O-O"),
						isCapture:         pBool(false),
						isPromotion:       pBool(false),
						isResign:          pBool(false),
						isCheck:           isCheck,
						isCheckmate:       isCheckmate,
					}
					cs := string(castlingSymbol[0])
					ch := Characteristics{usesCastlingSymbol: &cs}
					return []tokenMatch{{ms[0], &ap, ch}}
				},
			},
		}
	)
	return newNotationParser(transitions, evolveCharacteristics, initialCharacteristics)
}
