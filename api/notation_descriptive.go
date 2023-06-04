package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var descriptiveStringToPieceType = map[string]pieceType{
	"Q":    pieceQueen,
	"K":    pieceKing,
	"B":    pieceBishop,
	"N":    pieceKnight,
	"Kt":   pieceKnight,
	"R":    pieceRook,
	"P":    piecePawn,
	"QB":   pieceBishop,
	"QN":   pieceKnight,
	"QKt":  pieceKnight,
	"QR":   pieceRook,
	"QP":   piecePawn,
	"KB":   pieceBishop,
	"KN":   pieceKnight,
	"KKt":  pieceKnight,
	"KR":   pieceRook,
	"KP":   piecePawn,
	"BP":   piecePawn,
	"NP":   piecePawn,
	"KtP":  piecePawn,
	"RP":   piecePawn,
	"QBP":  piecePawn,
	"QNP":  piecePawn,
	"QKtP": piecePawn,
	"QRP":  piecePawn,
	"KBP":  piecePawn,
	"KNP":  piecePawn,
	"KKtP": piecePawn,
	"KRP":  piecePawn,
}

// Unfortunately, Descriptive notation has a lot of ambiguity. In the case of squares,
// without knowing if the move is by Black or White, any rank number can be two
// different ranks, and if the file is e.g. "B" we don't know if it's KB or QB so
// it can be both files. And both rank and file can be absent too.
//
// This method compiles all possible xys based on the available information.
func stringsToDescriptiveSquares(ss []string, moveNumber int) []xy {
	rx := regexp.MustCompile(`^(KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)?([1-8]?)$`)
	xys := []xy{}
	xs := []int{}
	ys := []int{}
	for _, s := range ss {
		matches := rx.FindStringSubmatch(s)
		if matches == nil {
			continue
		}
		squarePieceType, squareRank := matches[1], matches[2]
		if squarePieceType != "" {
			switch squarePieceType {
			case "KKt":
				xs = append(xs, 6)
			case "QKt":
				xs = append(xs, 1)
			case "Kt":
				xs = append(xs, 1, 6)
			case "KB":
				xs = append(xs, 5)
			case "KN":
				xs = append(xs, 6)
			case "KR":
				xs = append(xs, 7)
			case "QB":
				xs = append(xs, 2)
			case "QN":
				xs = append(xs, 1)
			case "QR":
				xs = append(xs, 0)
			case "Q":
				xs = append(xs, 3)
			case "K":
				xs = append(xs, 4)
			case "B":
				xs = append(xs, 2, 5)
			case "N":
				xs = append(xs, 1, 6)
			case "R":
				xs = append(xs, 0, 7)
			}
		}
		if squareRank != "" {
			squareRankInt, _ := strconv.Atoi(squareRank)
			if moveNumber%2 == 0 {
				ys = append(ys, 8-squareRankInt)
			} else {
				ys = append(ys, squareRankInt-1)
			}
		}
	}
	if len(xs) == 0 && len(ys) == 0 {
		return []xy{}
	}
	if len(xs) == 0 {
		xs = append(xs, -1)
	}
	if len(ys) == 0 {
		ys = append(ys, -1)
	}
	for _, y := range ys {
		for _, x := range xs {
			xys = append(xys, xy{x, y})
		}
	}
	return xys
}

func newNotationParserDescriptive(initialCharacteristics characteristics) *notationParser {
	var (
		transitions = map[string]map[string]func([]string, int) tokenMatch{
			"full_move_start": {
				`[\t\f\r ]*([0-9]+)?(\.)?[\t\f\r ]*`: func(ms []string, moveNumber int) tokenMatch {
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
				`[\t\f\r ]+`: func(ms []string, moveNumber int) tokenMatch {
					return tokenMatch{ms[0], nil, characteristics{}}
				},
			},
			"full_move_separator": {
				`([\t\f\r ]*?\n|[\t\f\r ]+)`: func(ms []string, moveNumber int) tokenMatch {
					var usesNewlineAsFullMoveSeparator *bool
					if strings.Contains(ms[0], "\n") {
						usesNewlineAsFullMoveSeparator = pBool(true)
					}
					return tokenMatch{ms[0], nil, characteristics{usesNewlineAsFullMoveSeparator: usesNewlineAsFullMoveSeparator}}
				},
			},
			"move": {
				// Move
				`(KKtP|QKtP|QKt|KKt|KtP|QBP|QNP|QRP|KBP|KNP|KRP|KR|KP|QP|BP|NP|RP|QR|KB|KN|QB|QN|Q|B|N|Kt|K|R|P)(\(((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8])\))?-((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8])(\/((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8]))?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, moveNumber int) tokenMatch {
					sFromPieceType, fromPieceSquare1, _, _, toPieceSquare, _, fromPieceSquare2, _, _, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8], ms[9], ms[10], ms[11]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					fromSquare := stringsToDescriptiveSquares([]string{fromPieceSquare1, fromPieceSquare2}, moveNumber)
					toSquare := stringsToDescriptiveSquares([]string{toPieceSquare}, moveNumber)
					ap := actionPattern{
						fromPieceType:      descriptiveStringToPieceType[sFromPieceType],
						fromX:              xySliceToPSInt(fromSquare, 'x'),
						fromY:              xySliceToPSInt(fromSquare, 'y'),
						toX:                xySliceToPSInt(toSquare, 'x'),
						toY:                xySliceToPSInt(toSquare, 'y'),
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
				`(KKtP|QKtP|QKt|KKt|KtP|QBP|QNP|QRP|KBP|KNP|KRP|KR|KP|QP|BP|NP|RP|QR|KB|KN|QB|QN|Q|B|N|Kt|K|R|P)(\(((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8])\))?x(KKtP|QKtP|QKt|KKt|KtP|QBP|QNP|QRP|KBP|KNP|KRP|KR|KP|QP|BP|NP|RP|QR|KB|KN|QB|QN|Q|B|N|Kt|K|R|P)(\/((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8]))?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, moveNumber int) tokenMatch {
					sFromPieceType, fromPieceSquare, _, _, sToPieceType, _, toPieceSquare, _, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8], ms[9], ms[10]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					fromSquare := stringsToDescriptiveSquares([]string{fromPieceSquare}, moveNumber)
					toSquare := stringsToDescriptiveSquares([]string{toPieceSquare}, moveNumber)
					ap := actionPattern{
						fromPieceType:      descriptiveStringToPieceType[sFromPieceType],
						capturedPieceType:  descriptiveStringToPieceType[sToPieceType],
						fromX:              xySliceToPSInt(fromSquare, 'x'),
						fromY:              xySliceToPSInt(fromSquare, 'y'),
						toX:                xySliceToPSInt(toSquare, 'x'),
						toY:                xySliceToPSInt(toSquare, 'y'),
						isCapture:          pBool(true),
						capturedPieceX:     xySliceToPSInt(toSquare, 'x'),
						capturedPieceY:     xySliceToPSInt(toSquare, 'y'),
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

				// Move Promotion
				`(KKtP|QKtP|QKt|KKt|KtP|QBP|QNP|QRP|KBP|KNP|KRP|KR|KP|QP|BP|NP|RP|QR|KB|KN|QB|QN|Q|B|N|Kt|K|R|P)(\(((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8])\))?-((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8])(\((Q|B|N|Kt|R)\)|=(Q|B|N|Kt|R))(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, moveNumber int) tokenMatch {
					sFromPieceType, fromPieceSquare, _, _, toPieceSquare, _, _, promotionPieceType1, promotionPieceType2, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8], ms[9], ms[10], ms[11]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					fromSquare := stringsToDescriptiveSquares([]string{fromPieceSquare}, moveNumber)
					toSquare := stringsToDescriptiveSquares([]string{toPieceSquare}, moveNumber)
					ap := actionPattern{
						fromPieceType:      descriptiveStringToPieceType[sFromPieceType],
						fromX:              xySliceToPSInt(fromSquare, 'x'),
						fromY:              xySliceToPSInt(fromSquare, 'y'),
						toX:                xySliceToPSInt(toSquare, 'x'),
						toY:                xySliceToPSInt(toSquare, 'y'),
						isCapture:          pBool(false),
						isPromotion:        pBool(true),
						isCastle:           pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: pBool(false),
						promotionPieceType: descriptiveStringToPieceType[promotionPieceType1+promotionPieceType2],
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					// TODO add promotion symbol to matching pattern to add characteristic
					ch := characteristics{
						usesCheckSymbol:     usesCheckSymbol,
						usesCheckmateSymbol: usesCheckmateSymbol,
					}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Capture Promotion
				`(KKtP|QKtP|QKt|KKt|KtP|QBP|QNP|QRP|KBP|KNP|KRP|KR|KP|QP|BP|NP|RP|QR|KB|KN|QB|QN|Q|B|N|Kt|K|R|P)(\(((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8])\))?x(KKtP|QKtP|QKt|KKt|KtP|QBP|QNP|QRP|KBP|KNP|KRP|KR|KP|QP|BP|NP|RP|QR|KB|KN|QB|QN|Q|B|N|Kt|K|R|P)(\((Q|B|N|Kt|R)\)|=(Q|B|N|Kt|R))(\/((KKt|QKt|Kt|KB|KN|KR|QB|QN|QR|Q|K|B|N|R)[1-8]?|[1-8]))?(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, moveNumber int) tokenMatch {
					sFromPieceType, fromPieceSquare, _, _, sCapturedPieceType, _, promotionPieceType1, promotionPieceType2, _, toPieceSquare, _, threatenSymbol, _ := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6], ms[7], ms[8], ms[9], ms[10], ms[11], ms[12], ms[13]
					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					fromSquare := stringsToDescriptiveSquares([]string{fromPieceSquare}, moveNumber)
					toSquare := stringsToDescriptiveSquares([]string{toPieceSquare}, moveNumber)
					ap := actionPattern{
						fromPieceType:      descriptiveStringToPieceType[sFromPieceType],
						capturedPieceType:  descriptiveStringToPieceType[sCapturedPieceType],
						fromX:              xySliceToPSInt(fromSquare, 'x'),
						fromY:              xySliceToPSInt(fromSquare, 'y'),
						toX:                xySliceToPSInt(toSquare, 'x'),
						toY:                xySliceToPSInt(toSquare, 'y'),
						isCapture:          pBool(false),
						isPromotion:        pBool(true),
						isCastle:           pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: pBool(false),
						promotionPieceType: descriptiveStringToPieceType[promotionPieceType1+promotionPieceType2],
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					// TODO add promotion symbol to matching pattern to add characteristic
					ch := characteristics{
						usesCheckSymbol:     usesCheckSymbol,
						usesCheckmateSymbol: usesCheckmateSymbol,
					}
					return tokenMatch{ms[0], &ap, ch}
				},

				// Castling
				`(0-0|0-0-0|O-O|O-O-O)(\+|†|ch|dbl\.? ?ch|\+\+|dis\.? ?ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, moveNumber int) tokenMatch {
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
				`(1–0|0–1|½–½|resigns|White resigns|Black resigns|mate)`: func(ms []string, moveNumber int) tokenMatch {
					var usesEndGameSymbol string
					switch ms[1] {
					case "1-0", "0-1", "½–½":
						usesEndGameSymbol = "numbers"
					case "mate":
						usesEndGameSymbol = "mate"
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
