package parser

import (
	"strconv"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

func NewNotationParserICCF(initialCharacteristics Characteristics) *NotationParser {
	var (
		transitions = map[string]map[string]func([]string, core.Game) []tokenMatch{
			"full_move_start": {
				`[\t\f\r ]*([0-9]{0,3}[\t\f\r ]|[0-9]{0,3}\.)?[\t\f\r ]*`: func(ms []string, g core.Game) []tokenMatch {
					var fullMoveNumber *int
					if len(ms[1]) > 0 {
						fmn, _ := strconv.Atoi(ms[1])
						fullMoveNumber = &fmn
					}
					var usesFullMoveDot *bool
					if strings.Contains(ms[0], ".") {
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
				`([1-8])([1-8])([1-8])([1-8])([1-8])?`: func(ms []string, g core.Game) []tokenMatch {
					fromFile, fromRank, toFile, toRank, promotion := ms[1], ms[2], ms[3], ms[4], ms[5]

					ap := actionPattern{
						fromPieceType:      stringToPieceType(""),
						fromX:              iccfFileToPInt(fromFile),
						fromY:              iccfRankToPInt(fromRank),
						toX:                iccfFileToPInt(toFile),
						toY:                iccfRankToPInt(toRank),
						isCapture:          pBool(promotion != ""),
						isPromotion:        nil,
						isCastle:           nil,
						isResign:           nil,
						isEnPassantCapture: nil,
						isCheck:            nil,
						isCheckmate:        nil,
					}
					ch := Characteristics{}
					return []tokenMatch{{ms[0], &ap, ch}}
				},
			},
		}
		evolveCharacteristics = func(ch Characteristics, sc Characteristics) (Characteristics, error) {
			return ch, nil
		}
	)

	return newNotationParser(transitions, evolveCharacteristics, initialCharacteristics)
}

func iccfFileToPInt(file string) *int {
	if file == "" {
		return nil
	}
	num, err := strconv.Atoi(file)
	if err != nil {
		return nil
	}
	if num > 8 || num < 1 {
		return nil
	}
	x := num - 1
	return &x
}

func iccfRankToPInt(rank string) *int {
	if rank == "" {
		return nil
	}
	num, err := strconv.Atoi(rank)
	if err != nil {
		return nil
	}
	if num > 8 || num < 1 {
		return nil
	}
	y := 8 - num
	return &y
}
