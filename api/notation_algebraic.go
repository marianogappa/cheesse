package api

import (
	"fmt"
	"strconv"
	"strings"
)

func newNotationParserAlgebraic(initialCharacteristics characteristics) *notationParser {
	var (
		transitions = map[string]map[string]func([]string) tokenMatch{
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
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
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
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Capture
				`([QKBNR])([a-h])?([1-8])?(x|:)?([a-h])([1-8])(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
					sFromPieceType, fromSquareFile, fromSquareRank, _, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
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
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Capture with colon at the end
				`([QKBNR])([a-h])?([1-8])?([a-h])([1-8]):(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
					sFromPieceType, fromSquareFile, fromSquareRank, toSquareFile, toSquareRank, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
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
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Capture with pawn, potentially without rank
				`([a-h])(x|:)?([a-h])([1-8]?)( ?e.p.)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
					fromSquareFile, _, toSquareFile, toSquareRank, enPassantCapture, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
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
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Capture and promotion with pawn, potentially without rank
				`([a-h])(x|:)?([a-h])([1-8]?)([=\(])([QBNR])\)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
					fromSquareFile, _, toSquareFile, toSquareRank, promotionSymbol, sPromotionPieceType, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
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
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := characteristics{
						usesCheckSymbol:     usesCheckSymbol,
						usesCheckmateSymbol: usesCheckmateSymbol,
						usesPromotionSymbol: &promotionSymbol,
					}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Promotion
				`([a-h])([1-8])([=\(])([QBNR])\)?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
					toSquareFile, toSquareRank, promotionSymbol, sPromotionPieceType, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
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
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := characteristics{
						usesCheckSymbol:     usesCheckSymbol,
						usesCheckmateSymbol: usesCheckmateSymbol,
						usesPromotionSymbol: &promotionSymbol,
					}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Castling
				`(0-0|0-0-0|O-O|O-O-O)(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string) tokenMatch {
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
					ch := characteristics{
						usesCheckSymbol:     usesCheckSymbol,
						usesCheckmateSymbol: usesCheckmateSymbol,
						usesCastlingSymbol:  &cs,
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

		// TODO human-readable error messages here. Also, lacking some context.
		evolveCharacteristics = func(ch characteristics, sc characteristics) (characteristics, error) {
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

func processThreatenSymbol(threatenSymbol string) (isCheck *bool, isCheckmate *bool, usesCheckSymbol *string, usesCheckmateSymbol *string) {
	switch threatenSymbol {
	case "+", "†", "ch", "++", "dblch", "dbl ch", "dbl.ch", "disch", "dis ch", "dis.ch":
		isCheck = pBool(true)
		isCheckmate = nil
		usesCheckSymbol = &threatenSymbol
		usesCheckmateSymbol = nil
	case "#", "mate", "‡", "≠", "X", "x", "×":
		isCheck = nil
		isCheckmate = pBool(true)
		usesCheckSymbol = nil
		usesCheckmateSymbol = &threatenSymbol
	default:
		isCheck = nil
		isCheckmate = nil
		usesCheckSymbol = nil
		usesCheckmateSymbol = nil
	}
	return
}
