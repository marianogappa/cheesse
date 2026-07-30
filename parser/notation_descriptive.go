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
				// Promotion move: P-R8(Q), P-Q8(N)ch, P-Kt8=Q
				`P-(QR|QN|QKt|QB|Q|KB|KN|KKt|KR|B|N|Kt|R|K)?([1-8])[\(=]?(Q|K|B|N|Kt|R)\)?(\+\+|dbl\.? ?ch|dis\.? ?ch|\+|†|ch|#|mate|‡|≠)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
					toSquareFile, toSquareRank, sPromotionPieceType := ms[1], ms[2], ms[3]
					threatenSymbol, _ := ms[4], ms[5]

					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					ap := actionPattern{
						fromPieceType:      core.PiecePawn,
						toX:                descFileToPInt(toSquareFile),
						toY:                descRankToPInt(toSquareRank, g),
						isCapture:          pBool(false),
						isPromotion:        pBool(true),
						promotionPieceType: descStringToPieceType(sPromotionPieceType),
						isCastle:           pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: pBool(false),
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}
					ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}

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

				// Move with optional disambiguation: KN-K2, QR-Q1, R(R5)-QR5, R(Kt)-Kt6, QB-B4
				// prefix: optional [QK] side + piece, or piece + (file-rank) parenthesized
				`([QK])?([QKBNRPt]t?)(?:\(([QRNKtB]t?)([1-8])?\))?-(QR|QN|QKt|QB|Q|KB|KN|KKt|KR|B|N|Kt|R|K)?([1-8])?(\+\+|dbl\.? ?ch|dis\.? ?ch|\+|†|ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
					disambigSide, sFromPieceType := ms[1], ms[2]
					parenFile, parenRank := ms[3], ms[4]
					toSquareFile, toSquareRank, threatenSymbol, _ := ms[5], ms[6], ms[7], ms[8]

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

					disambigAlts := applyDescDisambig(&ap, disambigSide, parenFile, parenRank, g)

					ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}

					bases := []actionPattern{ap}
					if len(disambigAlts) > 0 {
						bases = nil
						for _, x := range disambigAlts {
							newAP := ap.Clone()
							newAP.fromX = pInt(x)
							bases = append(bases, newAP)
						}
					}

					ambiguousFiles := map[string][]int{
						"R":  {0, 7},
						"N":  {1, 6},
						"Kt": {1, 6},
						"B":  {2, 5},
					}
					if xs, ok := ambiguousFiles[toSquareFile]; ok {
						tokenMatches := []tokenMatch{}
						for _, base := range bases {
							for _, x := range xs {
								tokenMatches = append(tokenMatches, tokenMatch{ms[0], cloneActionPatternWithToX(base, x), ch})
							}
						}
						return tokenMatches
					}

					tokenMatches := []tokenMatch{}
					for _, base := range bases {
						b := base
						tokenMatches = append(tokenMatches, tokenMatch{ms[0], &b, ch})
					}
					return tokenMatches
				},

				// Capture: PxP, BxN, QxP, BxBch, KxP, RxNch, QxKtP, BxQP, etc.
				// Also handles disambiguation: R(R5)xP, QBxP
				`(?:([QK])?\(?([QRNKtB]t?)?([1-8])?\)?)?(Q|K|B|N|Kt|R|P)x(?:([QK])?([QRNKtB]t?)?)?(Q|K|B|N|Kt|R|P)(e\.p\.)?(\+\+|dbl\.? ?ch|dis\.? ?ch|\+|†|ch|#|mate|‡|≠)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
					disambigSide, disambigFile, disambigRank := ms[1], ms[2], ms[3]
					sFromPieceType := ms[4]
					capturedSide, capturedFile, sCapturedPieceType := ms[5], ms[6], ms[7]
					enPassant, threatenSymbol, _ := ms[8], ms[9], ms[10]
					_ = capturedSide

					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					ap := actionPattern{
						fromPieceType:      descStringToPieceType(sFromPieceType),
						capturedPieceType:  descStringToPieceType(sCapturedPieceType),
						isCapture:          pBool(true),
						isPromotion:        pBool(false),
						isCastle:           pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: nilOrTrue(enPassant != ""),
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}

					// N.B. capture handler doesn't fork on from-disambiguation alternatives
				// because the captured piece type + file already constrain the match.
				_ = applyDescDisambig(&ap, disambigSide, disambigFile, disambigRank, g)

					capturedCombined := capturedSide + capturedFile
					if capturedCombined != "" {
						cpx := descFileToPInt(capturedCombined)
						if cpx != nil {
							ap.capturedPieceX = cpx
						} else {
							ambiguousFiles := map[string][]int{
								"R":  {0, 7},
								"N":  {1, 6},
								"Kt": {1, 6},
								"B":  {2, 5},
							}
							if xs, ok := ambiguousFiles[capturedCombined]; ok {
								tokenMatches := []tokenMatch{}
								for _, x := range xs {
									newAP := ap.Clone()
									newAP.capturedPieceX = pInt(x)
									ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
									tokenMatches = append(tokenMatches, tokenMatch{ms[0], &newAP, ch})
								}
								return tokenMatches
							}
						}
					}

					ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
					return []tokenMatch{{ms[0], &ap, ch}}
				},

				// Capture-promotion: PxR(Q)ch, PxP(N), etc.
				`(P)x(Q|K|B|N|Kt|R|P)\(?(Q|K|B|N|Kt|R)\)?(\+\+|dbl\.? ?ch|dis\.? ?ch|\+|†|ch|#|mate|‡|≠)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
					_, sCapturedPieceType, sPromotionPieceType := ms[1], ms[2], ms[3]
					threatenSymbol, _ := ms[4], ms[5]

					isCheck, isCheckmate, usesCheckSymbol, usesCheckmateSymbol := processThreatenSymbol(threatenSymbol)
					ap := actionPattern{
						fromPieceType:      core.PiecePawn,
						capturedPieceType:  descStringToPieceType(sCapturedPieceType),
						promotionPieceType: descStringToPieceType(sPromotionPieceType),
						isCapture:          pBool(true),
						isPromotion:        pBool(true),
						isCastle:           pBool(false),
						isResign:           pBool(false),
						isEnPassantCapture: pBool(false),
						isCheck:            isCheck,
						isCheckmate:        isCheckmate,
					}

					ch := Characteristics{usesCheckSymbol: usesCheckSymbol, usesCheckmateSymbol: usesCheckmateSymbol}
					return []tokenMatch{{ms[0], &ap, ch}}
				},

				// Castling
				`(0-0|0-0-0|O-O|O-O-O)(\+\+|dbl\.? ?ch|dis\.? ?ch|\+|†|ch|#|mate|‡|≠|X|x|×)?(!!|\?\?|!\?|\?!|!|\?)?`: func(ms []string, g core.Game) []tokenMatch {
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
				rxEndOfGame: processEndOfGameToken,
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

	np := newNotationParser(transitions, evolveCharacteristics, initialCharacteristics)
	np.preprocessor = normalizeDescriptive
	return np
}

// normalizeDescriptive pre-processes a descriptive notation string by gluing
// space-separated suffixes (like "mate", "ch", "dblch", "dis.ch") back onto the
// preceding move token, e.g. "P-R8(Q) mate" → "P-R8(Q)mate".
func normalizeDescriptive(s string) string {
	for _, suffix := range []string{"mate", "dblch", "dbl.ch", "dbl ch", "dis.ch", "dis ch", "disch", "ch"} {
		s = strings.ReplaceAll(s, " "+suffix, suffix)
	}
	return s
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

// applyDescDisambig sets fromX/fromY on an action pattern based on descriptive
// disambiguation prefixes. Returns extra alternatives for ambiguous file names (e.g.
// "Kt" could be QN or KN file). If the returned slice is non-empty, the caller
// should use those instead of the original ap.
func applyDescDisambig(ap *actionPattern, side, file, rank string, g core.Game) []int {
	if side != "" && file == "" {
		homeFiles := map[string]map[core.PieceType]int{
			"Q": {core.PieceRook: 0, core.PieceKnight: 1, core.PieceBishop: 2},
			"K": {core.PieceRook: 7, core.PieceKnight: 6, core.PieceBishop: 5},
		}
		if files, ok := homeFiles[side]; ok {
			if x, ok := files[ap.fromPieceType]; ok {
				ap.fromX = pInt(x)
			}
		}
	} else if file != "" {
		combined := side + file
		x := descFileToPInt(combined)
		if x != nil {
			ap.fromX = x
		} else {
			ambiguousFiles := map[string][]int{
				"R":  {0, 7},
				"N":  {1, 6},
				"Kt": {1, 6},
				"B":  {2, 5},
			}
			if xs, ok := ambiguousFiles[combined]; ok {
				if rank != "" {
					ap.fromY = descRankToPInt(rank, g)
				}
				return xs
			}
		}
	}
	if rank != "" {
		ap.fromY = descRankToPInt(rank, g)
	}
	return nil
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
