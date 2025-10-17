package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

func NewNotationParserICCF(initialCharacteristics Characteristics) *NotationParser {
	var (
		evolveCharacteristics = func(ch Characteristics, sc Characteristics) (Characteristics, error) {
			// ICCF notation doesn't have complex characteristics evolution
			// Just merge the characteristics
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
				`([1-8])([1-8])([1-8])([1-8])([1-8])?`: func(ms []string, g core.Game) []tokenMatch {
					fromFile, fromRank, toFile, toRank, promotionPieceType := ms[1], ms[2], ms[3], ms[4], ms[5]

					ap := actionPattern{
						fromPieceType:      core.PieceNone, // ICCF doesn't specify piece type, let parser infer it
						fromX:              iccfFileToPInt(fromFile),
						fromY:              iccfRankToPInt(fromRank),
						toX:                iccfFileToPInt(toFile),
						toY:                iccfRankToPInt(toRank),
						isCapture:          nil, // ??
						isPromotion:        pBool(promotionPieceType != ""),
						isCastle:           nil,
						isResign:           nil,
						isEnPassantCapture: nil,
						isCheck:            nil,
						isCheckmate:        nil,
						promotionPieceType: iccfStringToPieceType(promotionPieceType),
					}

					return []tokenMatch{{ms[0], &ap, Characteristics{}}}
				},
			},
		}
	)
	return newNotationParser(transitions, evolveCharacteristics, initialCharacteristics)
}

func iccfFileToPInt(file string) *int {
	if file == "" {
		return nil
	}
	iccfFile, err := strconv.Atoi(file)
	if err != nil {
		return nil
	}
	if iccfFile < 1 || iccfFile > 8 {
		return nil
	}
	x := iccfFile - 1
	return &x
}

func iccfRankToPInt(rank string) *int {
	if rank == "" {
		return nil
	}
	iccfRank, err := strconv.Atoi(rank)
	if err != nil {
		return nil
	}
	if iccfRank < 1 || iccfRank > 8 {
		return nil
	}
	y := 8 - iccfRank
	return &y
}

func iccfStringToPieceType(s string) core.PieceType {
	if pt, ok := map[string]core.PieceType{
		"1": core.PieceQueen,
		"2": core.PieceRook,
		"3": core.PieceBishop,
		"4": core.PieceKnight,
	}[s]; ok {
		return pt
	}
	return core.PieceNone
}
