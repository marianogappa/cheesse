package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

func NewNotationParserDescriptive(initialCharacteristics Characteristics) *NotationParser {
	var (
		transitions = map[string]map[string]func([]string, core.Game) []tokenMatch{
			"full_move_start": {
				`[\t\f\r ]*([0-9]+)?(\.)?[\t\f\r ]*`: func(ms []string, g core.Game) []tokenMatch {
					var fullMoveNumber *int
					if len(ms[1]) > 0 {
						fmn, _ := strconv.Atoi(ms[1])
						fullMoveNumber = &fmn
					}
					var usesFullMoveDot *bool
					if len(ms[2]) == 1 {
						usesFullMoveDot = pBool(true)
					}
					return []tokenMatch{{ms[0], nil, Characteristics{fullMoveNumber: fullMoveNumber, usesFullMoveDot: usesFullMoveDot}}}
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
				// Move
				`(Q|K|B|N|Kt|R|P)-(QR|QN|QKt|QB|Q|KB|KN|KKt|KR|B|N|Kt|R|K)?([1-8])?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
					sFromPieceType, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5]

					// At least one of file/rank must be present
					if toSquareFile == "" && toSquareRank == "" {
						return []tokenMatch{}
					}

					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					ap := actionPattern{
						fromPieceType:      descStringToPieceType(sFromPieceType),
						toX:                descFileToPInt(toSquareFile),
						toY:                descRankToPInt(toSquareRank, g),
						isCapture:          pBool(false),
						isPromotion:        pBool(false),
						isCastle:           pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: pBool(false),
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}

					// In Descriptive notation, sometimes the file is ambiguous. For example, if the move is "R-N5" it
					// could be either "R-QN5" or "R-KN5". In this case, both patterns are plausible. We return both.
					ambiguousFiles := map[string][]int{
						"R":  {0, 7},
						"N":  {1, 6},
						"Kt": {1, 6},
						"B":  {2, 5},
					}
					if _, ok := ambiguousFiles[toSquareFile]; ok {
						tokenMatches := []tokenMatch{}
						for _, x := range ambiguousFiles[toSquareFile] {
							tokenMatches = append(tokenMatches, tokenMatch{ms[0], cloneActionPatternWithToX(ap, x), ch})
						}
						return tokenMatches
					}

					return []tokenMatch{{ms[0], &ap, ch}}
				},

				// Capture
				// `(Q|K|B|N|Kt|R)(x|:)?([a-h])([1-8])(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
				// 	sFromPieceType, _, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6]
				// 	isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
				// 	ap := actionPattern{
				// 		fromPieceType:      descStringToPieceType(sFromPieceType),
				// 		toX:                fileToPInt(toSquareFile),
				// 		toY:                descRankToPInt(toSquareRank, g),
				// 		isCapture:          pBool(true),
				// 		capturedPieceX:     fileToPInt(toSquareFile),
				// 		capturedPieceY:     descRankToPInt(toSquareRank, g),
				// 		isPromotion:        pBool(false),
				// 		isCastle:           pBool(false),
				// 		isResign:           pBool(false),
				// 		isEnPassantCapture: pBool(false),
				// 		isCheck:            isCheck,
				// 		isCheckmate:        isCheckmate,
				// 	}
				// 	ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
				// 	return []tokenMatch{{ms[0], &ap, ch}}
				// },

				// Capture with colon at the end
				// `([QKBNR])([a-h])?([1-8])?([a-h])([1-8]):(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				// 	sFromPieceType, fromSquareFile, fromSquareRank, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
				// 	isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
				// 	ap := actionPattern{
				// 		fromPieceType:      stringToPieceType(sFromPieceType),
				// 		fromX:              fileToPInt(fromSquareFile),
				// 		fromY:              descRankToPInt(fromSquareRank),
				// 		toX:                fileToPInt(toSquareFile),
				// 		toY:                descRankToPInt(toSquareRank),
				// 		isCapture:          pBool(true),
				// 		capturedPieceX:     fileToPInt(toSquareFile),
				// 		capturedPieceY:     descRankToPInt(toSquareRank),
				// 		isPromotion:        pBool(false),
				// 		isCastle:           pBool(false),
				// 		isResign:           pBool(false),
				// 		isEnPassantCapture: pBool(false),
				// 		isCheck:            isCheck,
				// 		isCheckmate:        isCheckmate,
				// 	}
				// 	ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
				// 	return []tokenMatch{{ms[0], &ap, ch}
				// },

				// Capture with pawn, potentially without rank
				// `([a-h])(x|:)?([a-h])([1-8]?)( ?e.p.)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				// 	fromSquareFile, _, toSquareFile, toSquareRank, enPassantCapture, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
				// 	isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
				// 	ap := actionPattern{
				// 		fromPieceType:      stringToPieceType(""),
				// 		fromX:              fileToPInt(fromSquareFile),
				// 		toX:                fileToPInt(toSquareFile),
				// 		toY:                descRankToPInt(toSquareRank),
				// 		isEnPassantCapture: pBool(strings.HasSuffix(enPassantCapture, "e.p.")),
				// 		isCapture:          pBool(true),
				// 		isPromotion:        pBool(false),
				// 		isCastle:           pBool(false),
				// 		isResign:           pBool(false),
				// 		isCheck:            isCheck,
				// 		isCheckmate:        isCheckmate,
				// 	}
				// 	ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
				// 	return []tokenMatch{{ms[0], &ap, ch}
				// },

				// Capture and promotion with pawn, potentially without rank
				// `([a-h])(x|:)?([a-h])([1-8]?)([=\(])([QBNR])\)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
				// 	fromSquareFile, _, toSquareFile, toSquareRank, promotionSymbol, sPromotionPieceType, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8]
				// 	isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
				// 	ap := actionPattern{
				// 		fromPieceType:      stringToPieceType(""),
				// 		fromX:              fileToPInt(fromSquareFile),
				// 		toX:                fileToPInt(toSquareFile),
				// 		toY:                descRankToPInt(toSquareRank),
				// 		isEnPassantCapture: pBool(false),
				// 		isCapture:          pBool(true),
				// 		isPromotion:        pBool(true),
				// 		isCastle:           pBool(false),
				// 		isResign:           pBool(false),
				// 		promotionPieceType: stringToPieceType(sPromotionPieceType),
				// 		isCheck:            isCheck,
				// 		isCheckmate:        isCheckmate,
				// 	}
				// 	ch := Characteristics{
				// 		usesCheckSymbol:     usesCheckSymbol,
				// 		usesCheckmateSymbol: usesCheckmateSymbol,
				// 		usesPromotionSymbol: &promotionSymbol,
				// 	}
				// 	return []tokenMatch{{ms[0], &ap, ch}
				// },

				// Promotion
				// `([a-h])([1-8])([=\(])([QBNR])\)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
				// 	toSquareFile, toSquareRank, promotionSymbol, sPromotionPieceType, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6]
				// 	isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
				// 	ap := actionPattern{
				// 		fromPieceType:      stringToPieceType(""),
				// 		toX:                fileToPInt(toSquareFile),
				// 		toY:                descRankToPInt(toSquareRank, g),
				// 		isCapture:          pBool(false),
				// 		isPromotion:        pBool(true),
				// 		isCastle:           pBool(false),
				// 		isResign:           pBool(false),
				// 		isEnPassantCapture: pBool(false),
				// 		promotionPieceType: stringToPieceType(sPromotionPieceType),
				// 		isCheck:            isCheck,
				// 		isCheckmate:        isCheckmate,
				// 	}
				// 	ch := Characteristics{
				// 		usesCheckSymbol:     usesCheckSymbol,
				// 		usesCheckmateSymbol: usesCheckmateSymbol,
				// 		usesPromotionSymbol: &promotionSymbol,
				// 	}
				// 	return []tokenMatch{{ms[0], &ap, ch}}
				// },

				// Castling
				`(0-0|0-0-0|O-O|O-O-O)(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
					castlingSymbol, threatenSymbol, _ := ms[1], ms[2], ms[3]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					ap := actionPattern{
						isCastle:           pBool(true),
						isQueensideCastle:  pBool(castlingSymbol == "0-0-0" || castlingSymbol == "O-O-O"),
						isKingsideCastle:   pBool(castlingSymbol == "0-0" || castlingSymbol == "O-O"),
						isCapture:          pBool(false),
						isPromotion:        pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: pBool(false),
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					cs := string(castlingSymbol[0])
					ch := Characteristics{
						usesCheckSymbol:     usesCheckSymbol,
						usesCheckmateSymbol: usesCheckmateSymbol,
						usesCastlingSymbol:  &cs,
					}
					return []tokenMatch{{ms[0], &ap, ch}}
				},

				// End of game
				`(1–0|0–1|½–½|resigns|White resigns|Black resigns|Resigns)`: func(ms []string, g core.Game) []tokenMatch {
					var usesEndGameSymbol string
					switch ms[1] {
					case "1-0", "0-1", "½–½":
						usesEndGameSymbol = "numbers"
					case "resigns":
						usesEndGameSymbol = "resigns"
					case "Resigns":
						usesEndGameSymbol = "Resigns"
					case "White resigns", "Black resigns":
						usesEndGameSymbol = "color resigns"
					}
					ap := actionPattern{
						isResign:           pBool(true),
						isPromotion:        pBool(false),
						isCastle:           pBool(false),
						isCapture:          pBool(false),
						isEnPassantCapture: pBool(false),
					}
					ch := Characteristics{usesEndGameSymbol: &usesEndGameSymbol}
					if ch.isCheck {
						ap.isCheck = pBool(true)
					}
					if ch.isCheckmate {
						ap.isCheckmate = pBool(true)
					}
					return []tokenMatch{{ms[0], &ap, ch}}
				},
			},
		}

		// TODO human-readable error messages here. Also, lacking some context.
		evolveCharacteristics = func(ch Characteristics, sc Characteristics) (Characteristics, error) {
			if sc.usesCheckSymbol != nil {
				if ch.usesCheckSymbol == nil {
					ch.usesCheckSymbol = sc.usesCheckSymbol
				} else if *ch.usesCheckSymbol != *sc.usesCheckSymbol {
					return ch, fmt.Errorf("expecting CheckSymbol %v but found %v", *ch.usesCheckSymbol, *sc.usesCheckSymbol)
				}
			}
			if sc.usesCheckmateSymbol != nil {
				if ch.usesCheckmateSymbol == nil {
					ch.usesCheckmateSymbol = sc.usesCheckmateSymbol
				} else if *ch.usesCheckmateSymbol != *sc.usesCheckmateSymbol {
					return ch, fmt.Errorf("expecting CheckmateSymbol %v but found %v", *ch.usesCheckmateSymbol, *sc.usesCheckmateSymbol)
				}
			}
			if sc.usesFullMoveDot != nil {
				if ch.usesFullMoveDot == nil {
					ch.usesFullMoveDot = sc.usesFullMoveDot
				} else if *ch.usesFullMoveDot != *sc.usesFullMoveDot {
					return ch, fmt.Errorf("expecting FullMoveDot %v but found %v", *ch.usesFullMoveDot, *sc.usesFullMoveDot)
				}
			}
			if sc.usesNewlineAsFullMoveSeparator != nil {
				if ch.usesNewlineAsFullMoveSeparator == nil {
					ch.usesNewlineAsFullMoveSeparator = sc.usesNewlineAsFullMoveSeparator
				} else if *ch.usesNewlineAsFullMoveSeparator != *sc.usesNewlineAsFullMoveSeparator {
					return ch, fmt.Errorf("expecting NewlineAsFullMoveSeparator %v but found %v", *ch.usesNewlineAsFullMoveSeparator, *sc.usesNewlineAsFullMoveSeparator)
				}
			}
			if sc.usesThreatenSymbol != nil {
				if ch.usesThreatenSymbol == nil {
					ch.usesThreatenSymbol = sc.usesThreatenSymbol
				} else if *ch.usesThreatenSymbol != *sc.usesThreatenSymbol {
					return ch, fmt.Errorf("expecting ThreatenSymbol %v but found %v", *ch.usesThreatenSymbol, *sc.usesThreatenSymbol)
				}
			}
			if sc.usesCaptureSymbol != nil {
				if ch.usesCaptureSymbol == nil {
					ch.usesCaptureSymbol = sc.usesCaptureSymbol
				} else if *ch.usesCaptureSymbol != *sc.usesCaptureSymbol {
					return ch, fmt.Errorf("expecting CaptureSymbol %v but found %v", *ch.usesCaptureSymbol, *sc.usesCaptureSymbol)
				}
			}
			if sc.usesEndGameSymbol != nil {
				if ch.usesEndGameSymbol == nil {
					ch.usesEndGameSymbol = sc.usesEndGameSymbol
				} else if *ch.usesEndGameSymbol != *sc.usesEndGameSymbol {
					return ch, fmt.Errorf("expecting EndGameSymbol %v but found %v", *ch.usesEndGameSymbol, *sc.usesEndGameSymbol)
				}
			}
			if sc.usesPromotionSymbol != nil {
				if ch.usesPromotionSymbol == nil {
					ch.usesPromotionSymbol = sc.usesPromotionSymbol
				} else if *ch.usesPromotionSymbol != *sc.usesPromotionSymbol {
					return ch, fmt.Errorf("expecting PromotionSymbol %v but found %v", *ch.usesPromotionSymbol, *sc.usesPromotionSymbol)
				}
			}
			if sc.usesCastlingSymbol != nil {
				if ch.usesCastlingSymbol == nil {
					ch.usesCastlingSymbol = sc.usesCastlingSymbol
				} else if *ch.usesCastlingSymbol != *sc.usesCastlingSymbol {
					return ch, fmt.Errorf("expecting CastlingSymbol %v but found %v", *ch.usesCastlingSymbol, *sc.usesCastlingSymbol)
				}
			}
			// TODO full move number
			return ch, nil
		}
	)

	return newNotationParser(transitions, evolveCharacteristics, initialCharacteristics)
}

func descStringToPieceType(s string) core.PieceType {
	return map[string]core.PieceType{
		"Q":  core.PieceQueen,
		"K":  core.PieceKing,
		"B":  core.PieceBishop,
		"N":  core.PieceKnight,
		"Kt": core.PieceKnight,
		"R":  core.PieceRook,
		"P":  core.PiecePawn,
	}[s]
}

func descFileToPInt(file string) *int {
	if file == "" || file == "B" || file == "N" || file == "R" || file == "Kt" {
		return nil
	}
	x := map[string]int{
		"QR":  0,
		"QN":  1,
		"QKt": 1,
		"QB":  2,
		"Q":   3,
		"K":   4,
		"KB":  5,
		"KN":  6,
		"KKt": 6,
		"KR":  7,
	}[file]
	return &x
}

func descRankToPInt(rank string, g core.Game) *int {
	if rank == "" {
		return nil
	}
	if g.Turn() == core.ColorWhite {
		v := (8 - int(rank[0]-'0'))
		return &v
	}
	v := int(rank[0]-'0') - 1
	return &v
}

func pInt(i int) *int {
	return &i
}

func cloneActionPatternWithToX(ap actionPattern, toX int) *actionPattern {
	newAP := ap.Clone()
	newAP.toX = pInt(toX)
	return &newAP
}

// func pBool(b bool) *bool {
// 	return &b
// }

// func processThreatenSymbol(threatenSymbol string) (isCheck *bool, isCheckmate *bool, usesCheckSymbol *string, usesCheckmateSymbol *string) {
// 	switch threatenSymbol {
// 	case "+", "†", "ch", "++", "dblch", "dbl ch", "dbl.ch", "disch", "dis ch", "dis.ch":
// 		isCheck = pBool(true)
// 		isCheckmate = nil
// 		usesCheckSymbol = &threatenSymbol
// 		usesCheckmateSymbol = nil
// 	case "#", "mate", "‡", "≠", "X", "x", "×":
// 		isCheck = nil
// 		isCheckmate = pBool(true)
// 		usesCheckSymbol = nil
// 		usesCheckmateSymbol = &threatenSymbol
// 	default:
// 		isCheck = nil
// 		isCheckmate = nil
// 		usesCheckSymbol = nil
// 		usesCheckmateSymbol = nil
// 	}
// 	return
// }
