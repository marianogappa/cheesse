package parser

import (
	"fmt"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

func NewNotationParserSmith(initialCharacteristics Characteristics) *NotationParser {
	var (
		evolveCharacteristics = func(ch Characteristics, sc Characteristics) (Characteristics, error) {
			if sc.usesNewlineAsFullMoveSeparator != nil {
				if ch.usesNewlineAsFullMoveSeparator == nil {
					ch.usesNewlineAsFullMoveSeparator = sc.usesNewlineAsFullMoveSeparator
				} else if *ch.usesNewlineAsFullMoveSeparator != *sc.usesNewlineAsFullMoveSeparator {
					return ch, fmt.Errorf("expecting newline as full move separator %v but found %v", *ch.usesNewlineAsFullMoveSeparator, *sc.usesNewlineAsFullMoveSeparator)
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
				`([a-h])([1-8])([a-h])([1-8])([pnbrqkEcC])?([QNBR])?`: func(ms []string, g core.Game) []tokenMatch {
					fromFile, fromRank, toFile, toRank, captureOrCastle, promotionPiece := ms[1], ms[2], ms[3], ms[4], ms[5], ms[6]

					ap := actionPattern{
						fromX: fileToPInt(fromFile),
						fromY: rankToPInt(fromRank),
						toX:   fileToPInt(toFile),
						toY:   rankToPInt(toRank),
					}

					switch captureOrCastle {
					case "c":
						ap.isCastle = pBool(true)
						ap.isKingsideCastle = pBool(true)
						ap.isCapture = pBool(false)
					case "C":
						ap.isCastle = pBool(true)
						ap.isQueensideCastle = pBool(true)
						ap.isCapture = pBool(false)
					case "E":
						ap.isCapture = pBool(true)
						ap.isEnPassantCapture = pBool(true)
					case "":
						ap.isCapture = pBool(false)
					default:
						ap.isCapture = pBool(true)
						ap.capturedPieceType = smithCaptureToType(captureOrCastle)
					}

					if promotionPiece != "" {
						ap.isPromotion = pBool(true)
						ap.promotionPieceType = smithPromotionToType(promotionPiece)
					} else {
						ap.isPromotion = pBool(false)
					}

					return []tokenMatch{{ms[0], &ap, Characteristics{}}}
				},
			},
		}
	)
	return newNotationParser(transitions, evolveCharacteristics, initialCharacteristics)
}

func smithCaptureToType(s string) core.PieceType {
	if pt, ok := map[string]core.PieceType{
		"p": core.PiecePawn,
		"n": core.PieceKnight,
		"b": core.PieceBishop,
		"r": core.PieceRook,
		"q": core.PieceQueen,
		"k": core.PieceKing,
	}[s]; ok {
		return pt
	}
	return core.PieceNone
}

func smithPromotionToType(s string) core.PieceType {
	if pt, ok := map[string]core.PieceType{
		"Q": core.PieceQueen,
		"N": core.PieceKnight,
		"B": core.PieceBishop,
		"R": core.PieceRook,
	}[s]; ok {
		return pt
	}
	return core.PieceNone
}
