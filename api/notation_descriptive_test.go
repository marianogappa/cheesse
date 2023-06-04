package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringsToDescriptiveSquares(t *testing.T) {
	ts := []struct {
		input      []string
		moveNumber int
		expected   []xy
	}{
		{
			input:    []string{},
			expected: []xy{},
		},
		{
			input:      []string{"QR1"},
			moveNumber: 1,
			expected:   []xy{{0, 0}},
		},
		{
			input:      []string{"QR1"},
			expected:   []xy{{0, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN1"},
			expected:   []xy{{1, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN1"},
			expected:   []xy{{1, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB1"},
			expected:   []xy{{2, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB1"},
			expected:   []xy{{2, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q1"},
			expected:   []xy{{3, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q1"},
			expected:   []xy{{3, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"K1"},
			expected:   []xy{{4, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"K1"},
			expected:   []xy{{4, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB1"},
			expected:   []xy{{5, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB1"},
			expected:   []xy{{5, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN1"},
			expected:   []xy{{6, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN1"},
			expected:   []xy{{6, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR1"},
			expected:   []xy{{7, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR1"},
			expected:   []xy{{7, 7}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR2"},
			expected:   []xy{{0, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR2"},
			expected:   []xy{{0, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN2"},
			expected:   []xy{{1, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN2"},
			expected:   []xy{{1, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB2"},
			expected:   []xy{{2, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB2"},
			expected:   []xy{{2, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q2"},
			expected:   []xy{{3, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q2"},
			expected:   []xy{{3, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"K2"},
			expected:   []xy{{4, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"K2"},
			expected:   []xy{{4, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB2"},
			expected:   []xy{{5, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB2"},
			expected:   []xy{{5, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN2"},
			expected:   []xy{{6, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN2"},
			expected:   []xy{{6, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR2"},
			expected:   []xy{{7, 1}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR2"},
			expected:   []xy{{7, 6}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR3"},
			expected:   []xy{{0, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR3"},
			expected:   []xy{{0, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN3"},
			expected:   []xy{{1, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN3"},
			expected:   []xy{{1, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB3"},
			expected:   []xy{{2, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB3"},
			expected:   []xy{{2, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q3"},
			expected:   []xy{{3, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q3"},
			expected:   []xy{{3, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"K3"},
			expected:   []xy{{4, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"K3"},
			expected:   []xy{{4, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB3"},
			expected:   []xy{{5, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB3"},
			expected:   []xy{{5, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN3"},
			expected:   []xy{{6, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN3"},
			expected:   []xy{{6, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR3"},
			expected:   []xy{{7, 2}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR3"},
			expected:   []xy{{7, 5}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR4"},
			expected:   []xy{{0, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR4"},
			expected:   []xy{{0, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN4"},
			expected:   []xy{{1, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN4"},
			expected:   []xy{{1, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB4"},
			expected:   []xy{{2, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB4"},
			expected:   []xy{{2, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q4"},
			expected:   []xy{{3, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q4"},
			expected:   []xy{{3, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"K4"},
			expected:   []xy{{4, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"K4"},
			expected:   []xy{{4, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB4"},
			expected:   []xy{{5, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB4"},
			expected:   []xy{{5, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN4"},
			expected:   []xy{{6, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN4"},
			expected:   []xy{{6, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR4"},
			expected:   []xy{{7, 3}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR4"},
			expected:   []xy{{7, 4}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR5"},
			expected:   []xy{{0, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR5"},
			expected:   []xy{{0, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN5"},
			expected:   []xy{{1, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN5"},
			expected:   []xy{{1, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB5"},
			expected:   []xy{{2, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB5"},
			expected:   []xy{{2, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q5"},
			expected:   []xy{{3, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q5"},
			expected:   []xy{{3, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"K5"},
			expected:   []xy{{4, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"K5"},
			expected:   []xy{{4, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB5"},
			expected:   []xy{{5, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB5"},
			expected:   []xy{{5, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN5"},
			expected:   []xy{{6, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN5"},
			expected:   []xy{{6, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR5"},
			expected:   []xy{{7, 4}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR5"},
			expected:   []xy{{7, 3}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR6"},
			expected:   []xy{{0, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR6"},
			expected:   []xy{{0, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN6"},
			expected:   []xy{{1, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN6"},
			expected:   []xy{{1, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB6"},
			expected:   []xy{{2, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB6"},
			expected:   []xy{{2, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q6"},
			expected:   []xy{{3, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q6"},
			expected:   []xy{{3, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"K6"},
			expected:   []xy{{4, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"K6"},
			expected:   []xy{{4, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB6"},
			expected:   []xy{{5, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB6"},
			expected:   []xy{{5, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN6"},
			expected:   []xy{{6, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN6"},
			expected:   []xy{{6, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR6"},
			expected:   []xy{{7, 5}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR6"},
			expected:   []xy{{7, 2}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR7"},
			expected:   []xy{{0, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR7"},
			expected:   []xy{{0, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN7"},
			expected:   []xy{{1, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN7"},
			expected:   []xy{{1, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB7"},
			expected:   []xy{{2, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB7"},
			expected:   []xy{{2, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q7"},
			expected:   []xy{{3, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q7"},
			expected:   []xy{{3, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"K7"},
			expected:   []xy{{4, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"K7"},
			expected:   []xy{{4, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB7"},
			expected:   []xy{{5, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB7"},
			expected:   []xy{{5, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN7"},
			expected:   []xy{{6, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN7"},
			expected:   []xy{{6, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR7"},
			expected:   []xy{{7, 6}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR7"},
			expected:   []xy{{7, 1}},
			moveNumber: 0,
		},
		{
			input:      []string{"QR8"},
			expected:   []xy{{0, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"QR8"},
			expected:   []xy{{0, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"QN8"},
			expected:   []xy{{1, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"QN8"},
			expected:   []xy{{1, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"QB8"},
			expected:   []xy{{2, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"QB8"},
			expected:   []xy{{2, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"Q8"},
			expected:   []xy{{3, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"Q8"},
			expected:   []xy{{3, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"K8"},
			expected:   []xy{{4, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"K8"},
			expected:   []xy{{4, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"KB8"},
			expected:   []xy{{5, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"KB8"},
			expected:   []xy{{5, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"KN8"},
			expected:   []xy{{6, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"KN8"},
			expected:   []xy{{6, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR8"},
			expected:   []xy{{7, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR8"},
			expected:   []xy{{7, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"", "KR8"},
			expected:   []xy{{7, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"", "KR8"},
			expected:   []xy{{7, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"KR8", ""},
			expected:   []xy{{7, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"KR8", ""},
			expected:   []xy{{7, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"8"},
			expected:   []xy{{-1, 7}},
			moveNumber: 1,
		},
		{
			input:      []string{"8"},
			expected:   []xy{{-1, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"B8"},
			expected:   []xy{{2, 7}, {5, 7}, {5, 0}},
			moveNumber: 1,
		},
		{
			input:      []string{"B8"},
			expected:   []xy{{2, 0}, {5, 7}, {5, 0}},
			moveNumber: 0,
		},
		{
			input:      []string{"B"},
			expected:   []xy{{2, -1}},
			moveNumber: 1,
		},
		{
			input:      []string{"B"},
			expected:   []xy{{5, -1}},
			moveNumber: 0,
		},
	}
	for _, tc := range ts {
		t.Run(strings.Join(tc.input, ","), func(t *testing.T) {
			assert.ElementsMatch(t, tc.expected, stringsToDescriptiveSquares(tc.input, tc.moveNumber))
		})
	}
}

func TestNotationParserDescriptive(t *testing.T) {
	testCases := []struct {
		fen                   string
		s                     string
		expectedErr           error
		expectedMatchedTokens []string
		expectedBoard         []string
		expectedFEN           string
	}{
		// {
		// 	fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// 	s: `1. P-K4 P-K3
		// 		2. P-Q4 P-Q4
		// 		3. N-QB3 B-N5
		// 		4. B-N5ch B-Q2
		// 		5. BxBch QxB
		// 		6. KN-K2 PxP
		// 		7. 0-0`,
		// 	expectedErr:           nil,
		// 	expectedMatchedTokens: []string{"P-K4", "P-K3", "P-Q4", "P-Q4", "N-QB3", "B-N5", "B-N5ch", "B-Q2", "BxBch", "QxB", "KN-K2", "PxP", "0-0"},
		// 	expectedBoard: []string{
		// 		"♜♞  ♚ ♞♜",
		// 		"♟♟♟♛ ♟♟♟",
		// 		"    ♟   ",
		// 		"        ",
		// 		" ♝ ♙♟   ",
		// 		"  ♘     ",
		// 		"♙♙♙ ♘♙♙♙",
		// 		"♖ ♗♕ ♖♔ ",
		// 	},
		// },
		{
			fen:         "8/8/8/8/8/1k5P/8/2K5 w - - 0 1",
			s:           "1. P-R4    K-B5 2. P-R5    K-Q4 3. P-R6    K-K3 4. P-R7    K-B2 5. P-R8(Q) ",
			expectedErr: nil,
			expectedFEN: "7Q/5k2/8/8/8/8/8/2K5 b - - 0 5",
		},
		// {
		// 	fen:         "8/8/8/6K1/8/5k2/P7/8 w - - 0 1",
		// 	s:           "1. K-B5!  K-K6 2. K-K5!  K-Q6 3. K-Q5!  K-B6 4. K-B5!  K-Q6 5. P-R4    K-B6 6. P-R5    K-Kt6 7. P-R6    K-R5 8. P-R7    K-R4 9. P-R8(Q) mate ",
		// 	expectedErr: nil,
		// 	expectedFEN: "Q7/8/8/k1K5/8/8/8/8 b - - 0 9",
		// },
		// {
		// 	fen:         "4k3/8/3K4/4P3/8/8/8/8 w - - 0 1",
		// 	s:           "1. K-K6! K-Q1 2. K-B7 K-Q2 3. P-K6ch K-Q1 4. P-K7ch K-Q2 5. P-K8(Q)ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "4Q3/3k1K2/8/8/8/8/8/8 b - - 0 5",
		// },
		// {
		// 	fen:         "8/5k2/3P4/8/6K1/8/8/8 w - - 0 1",
		// 	s:           "1. K-B5 K-B1 2. K-B6! K-K1 3. K-K6  K-Q1 4. P-Q7 K-QB2 5. K-K7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/2kPK3/8/8/8/8/8/8 b - - 2 5",
		// },
		// {
		// 	fen:         "3k4/8/1P6/8/4K3/8/8/8 w - - 0 1",
		// 	s:           "1. K-Q5 K-Q2 2. K-B5 K-Q1 3. K-Q6! K-B1 4. K-B6 K-Kt1 5. P-Kt7 K-R2 6. K-B7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/kPK5/8/8/8/8/8/8 b - - 2 6",
		// },
		// {
		// 	fen:         "6k1/8/7K/8/2P5/8/8/8 w - - 0 1",
		// 	s:           "1. K-Kt6! K-B1 2. K-B6 K-K1 3. K-K6 K-Q1 4. K-Q6 K-B1 5. K-B6 K-Kt1 6. K-Q7 K-Kt2 7. P-B5 K-Kt1 8. P-B6 K-R2 9. P-B7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/k1PK4/8/8/8/8/8/8 b - - 0 9",
		// },
		// {
		// 	fen:         "8/7k/5K2/6P1/8/8/8/8 w - - 0 1",
		// 	s:           "1. K-B7! K-R1 2. K-Kt6! K-Kt1 3. K-R6 K-R1 4. P-Kt6 K-Kt1 5. P-Kt7 K-B2 6. K-R7 K-K2 7. P-Kt8(Q) ",
		// 	expectedErr: nil,
		// 	expectedFEN: "6Q1/4k2K/8/8/8/8/8/8 b - - 0 7",
		// },
		// {
		// 	fen:         "8/8/5k2/8/5K2/8/4P3/8 w - - 0 1",
		// 	s:           "1. K-K4 K-K3 2. P-K3! K-Q3 3. K-B5 K-K2 4.  K-K5 K-Q2 5. K-B6 K-K1 6. K-K6 K-B1 7. P-K4   K-K1 8. P-K5   K-B1 9. K-Q7 K-B2 10. P-K6ch K-B1 11. P-K7ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "5k2/3KP3/8/8/8/8/8/8 b - - 0 11",
		// },
		// {
		// 	fen:         "8/7k/8/8/8/3P4/8/6K1 w - - 0 1",
		// 	s:           "1. K-B2 K-Kt3 2. K-K3 K-B4 3. K-Q4! K-K3 4. K-B5 K-Q2 5. K-Q5 K-K2 6. K-B6 K-K3 7. P-Q4 K-K2 8. P-Q5 K-Q1 9. K-Q6 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "3k4/8/3K4/3P4/8/8/8/8 b - - 2 9",
		// },
		// {
		// 	fen:         "8/8/8/1k6/8/K7/6P1/8 w - - 0 1",
		// 	s:           "1. K-Kt3 K-B4 2. K-B3 K-Q4 3. K-Q3 K-K4 4. K-K3 K-B4 5. K-B3 K-Kt4 6. K-Kt3 K-B4 7. K-R4 K-B3 8. K-R5 K-Kt2 9. K-Kt5 K-B2 10. K-R6 K-Kt1 11. K-Kt6 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "6k1/8/6K1/8/8/8/6P1/8 b - - 21 11",
		// },
		// {
		// 	fen:         "3k4/8/K7/8/8/8/P7/8 w - - 0 1",
		// 	s:           "K-Kt7! K-Q2 2. P-R4 K-Q3 3. P-R5 K-B4 4. P-R6 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/1K6/P7/2k5/8/8/8/8 b - - 0 4",
		// },
		// {
		// 	fen:         "7k/7P/6P1/8/8/6K1/8/8 w - - 0 1",
		// 	s:           "1. K-B4 K-Kt2 2. K-B5 K-R1 3. K-Kt5 K-Kt2 4. P-R8(Q)ch! KxQ 5. K-B6 K-Kt1 6. P-Kt7 K-R2 7. K-B7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/5KPk/8/8/8/8/8/8 b - - 2 7",
		// },
		// {
		// 	fen:         "6k1/8/5KP1/6P1/8/8/8/8 w - - 0 1",
		// 	s:           "P-Kt7 K-R2 P-Kt8(Q)ch! KxQ K-Kt6 K-R1 K-B7 K-R2 P-Kt6ch K-R1 P-Kt7ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "7k/5KP1/8/8/8/8/8/8 b - - 0 6",
		// },
		// {
		// 	fen:         "8/8/8/8/2k5/P1P5/7K/8 w - - 0 1",
		// 	s:           "P-R4! K-B4 K-Kt3 K-Kt3 P-B4 K-R4 P-B5 K-R3 K-B3 K-Kt2 P-R5 K-B3 P-R6 K-B2 K-K4 K-B3 K-K5 K-B2 K-Q5 K-Kt1 P-B6 K-B2 P-R7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/P1k5/2P5/3K4/8/8/8/8 b - - 0 12",
		// },
		// {
		// 	fen:         "8/8/5Ppk/8/8/4K3/8/8 w - - 0 1",
		// 	s:           "K-B4 K-R2 K-Kt5 K-R1 K-R6! K-Kt1 KxP K-R1 K-B7 K-R2 K-K8 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "4K3/7k/5P2/8/8/8/8/8 b - - 4 6",
		// },
		// {
		// 	fen:         "8/8/5k2/8/p7/8/1PK5/8 w - - 0 1",
		// 	s:           "1. K-Kt1! P-R6 2. P-Kt3! K-K4 3. K-R2 K-Q4 4. KxP K-B3 5. K-R4! K-Kt3 6. K-Kt4 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/8/1k6/8/1K6/1P6/8/8 b - - 4 6",
		// },
		// {
		// 	fen:         "8/8/2p5/8/8/P4k2/8/2K5 w - - 0 1",
		// 	s:           "1. P-R4 R-K5 2. P-R5 K-Q4 3. P-R6 K-Q3 4. P-R7 K-B2 5. P-R8(Q) ",
		// 	expectedErr: nil,
		// 	expectedFEN: "Q7/2k5/2p5/8/8/8/8/2K5 b - - 0 5",
		// },
		// {
		// 	fen:         "8/p4K2/P7/8/8/8/1k6/8 w - - 0 1",
		// 	s:           "1. K-K6 K-B6 2. K-Q5! K-Kt5 3. K-B6 K-B5 4. K-Kt7 K-Kt4 5. KxP K-B3 6. K-Kt8 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "1K6/8/P1k5/8/8/8/8/8 b - - 2 6",
		// },
		// {
		// 	fen:         "8/8/8/7p/1PK2k2/8/8/8 w - - 0 1",
		// 	s:           "1. P-Kt5 K-K4 2. P-Kt6! K-Q3 3. K-Kt5 P-R5 4. K-R6 P-R6 5. P-Kt7 K-B2 6. K-R7 P-R7 7. P-Kt8(Q)ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "1Q6/K1k5/8/8/8/8/7p/8 b - - 0 7",
		// },
		// {
		// 	fen:         "8/8/8/3KP1k1/6p1/8/8/8 w - - 0 1",
		// 	s:           "1. P-K6          K-B3 2. K-Q6       P-Kt6 3. P-K7       P-Kt7 4. P-K8(Q) P-Kt8(Q) 5.Q-B8ch       K-Kt4 6. Q-Kt8ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "6Q1/8/3K4/6k1/8/8/8/6q1 b - - 3 6",
		// },
		// {
		// 	fen:         "8/1p6/7K/8/7P/8/5k2/8 w - - 0 1",
		// 	s:           "K-Kt5 P-Kt4 K-B4 –K-K7 K-K4 K-Q7 K-Q4 K-B7 K-B5 K-B6 P-R5 P-Kt5 P-R6 P-Kt6 P-R7 P-Kt7 P-R8(Q)ch K-B7 Q-R2ch K-B8 Q-B4ch K-B7 Q-B4ch K-Q7 Q-Kt3 K-B8 Q-B3ch K-Kt8 K-B4 K-R7 Q-R5ch K-Kt8 K-Kt3 K-B8 Q-K1 mate ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/8/8/8/8/1K6/1p6/2k1Q3 b - - 18 18",
		// },
		// {
		// 	fen:         "8/8/5P2/8/6K1/7p/6k1/8 w - - 0 1",
		// 	s:           "1. P-B7 P-R7 2. P-B8(Q) P-R8(Q) 3. Q-B3ch K-Kt8 4. Q-K3ch K-B8 5. Q-B1ch K-Kt7 6. Q-Q2ch K-B8 7. Q-Q1ch K-Kt7 8. Q-K2ch K-Kt8 9. K-Kt3! ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/8/8/8/8/6K1/4Q3/6kq b - - 13 9",
		// },
		// {
		// 	fen:         "8/8/1p6/8/8/6P1/k1K5/8 w - - 0 1",
		// 	s:           "1. K-B3! K-R6 2. K-B4 K-R5 3. P-Kt4 P-Kt4ch 4. K-Q3! K-R6 5. P-Kt5 P-Kt5 6. P-Kt6 P-Kt6 7. P-Kt7 P-Kt7 8. K-B2! K-R7 9. P-Kt8(Q)ch K-R8 10. Q-R8 mate ",
		// 	expectedErr: nil,
		// 	expectedFEN: "Q7/8/8/8/8/8/1pK5/k7 b - - 2 10",
		// },
		// {
		// 	fen:         "8/6p1/7k/8/1K6/8/1P6/8 w - - 0 1",
		// 	s:           "1. K-B5! P-Kt4 2. P-Kt4 P-Kt5 3. K-Q4 P-Kt6 4. K-K3 K-Kt4 5. P-Kt5 K-Kt5 6. P-Kt6 K-R6 7. P-Kt7 P-Kt7 8. K-B2 K-R7 9. P-Kt8(Q)ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "1Q6/8/8/8/8/8/5Kpk/8 b - - 0 9",
		// },
		// {
		// 	fen:         "8/1pK5/8/8/8/8/k4P2/8 w - - 0 1",
		// 	s:           "K-Q6! K-R6 K-B5 K-R5 P–B4 P-Kt4 P-B5 P-Kt5 K-B4! P-Kt6 K-B3! K-R6 P-B6 P-Kt7 P-B7 P-Kt8(Q) P-B8(Q)ch K-R5 Q-R8ch K-Kt4 Q-Kt7ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/1Q6/8/1k6/8/2K5/8/1q6 b - - 4 11",
		// },
		// {
		// 	fen:         "8/2p5/6K1/8/8/5k2/P7/8 w - - 0 1",
		// 	s:           "1. K-B5! K-K6 2. K-K5 P-B3 3. P-R4 K-Q6 4. P-R5 P-B4 5. P-R6 P-B5 6. P-R7 P-B6 7. P-R8(Q) P-B7 8. Q-Q5ch! K-K7 9. Q-R2! K-Q8 10. K-Q4 P-B8(Q) 11. K-Q3 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/8/8/8/8/3K4/Q7/2qk4 b - - 1 11",
		// },
		// {
		// 	fen:         "8/4K1pp/8/8/8/8/k6P/8 w - - 0 1",
		// 	s:           "P-R4! P-R4 K-B8! P-Kt3 K-K7! P-Kt4 PXP P-R5 P-Kt6 P-R6 P-Kt7 P-R7 P-Kt8(Q)ch K-R6 Q-Kt2 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/4K3/8/8/8/k7/6Qp/8 b - - 2 8",
		// },
		// {
		// 	fen:         "8/6k1/1p6/6KP/P7/8/8/8 w - - 0 1",
		// 	s:           "(a) K-B5 K-R3 K-K5 KxP K-Q5 K-Kt3 K-B6 K-B3 KxP  K-K2 K-B7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/2K1k3/8/8/P7/8/8/8 b - - 2 6",
		// },
		// {
		// 	fen:         "8/8/5k2/p7/P3KP2/8/8/8 w - - 0 1",
		// 	s:           "K-Q5 K-B4 K-B5 KxP K-Kt4 K-K4 KxP K-Q3 K-Kt6 K-Q2 K-Kt7 K-Q1 P-R5 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "3k4/1K6/8/P7/8/8/8/8 b - - 0 7",
		// },
		// {
		// 	fen:         "8/2p1k3/2P5/1P2K3/8/8/8/8 w - - 0 1",
		// 	s:           "K-Q5! K-Q1 K-K6 K-B1 K-K7 K-Kt1 K-Q8 K-R2 KxP K-R1 K-Q8 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "k2K4/8/2P5/1P6/8/8/8/8 b - - 2 6",
		// },
		// {
		// 	fen:         "8/kp6/8/PP6/1K6/8/8/8 w - - 0 1",
		// 	s:           "K-B5 K-Kt1 K-Kt6 K-R1 K-B7 K-R2 P-R6 PxP P-Kt6ch K-R1 P-Kt7ch K-R2 P-Kt8(Q) mate ",
		// 	expectedErr: nil,
		// 	expectedFEN: "1Q6/k1K5/p7/8/8/8/8/8 b - - 0 7",
		// },
		// {
		// 	fen:         "8/6k1/3p4/3P1PK1/8/8/8/8 w - - 0 1",
		// 	s:           "(a) P-B6ch K-B1! P-B7! KxP K-B5 K-K2 K-Kt6 K-Q1 K-B6 K-Q2 K-B7 K-Q1 K-K6 K-B2 K-K7 K-B1 KxP K-Q1 K-B6 K-B1 P-Q6 K-Q1 P-Q7 K-K2 K-B7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/2KPk3/8/8/8/8/8/8 b - - 2 13",
		// },
		// {
		// 	fen:         "8/3k4/2pP4/2P1K3/8/8/8/8 w - - 0 1",
		// 	s:           "K-B6 K-Q1! P-Q7! KxP K-B7! K-Q1 K-K6 K-B2 K-K7 K-B1 K-Q6 K-Kt2 K-Q7 K-Kt1 KxP K-B1 K-Q6 K-Q1 P-B6 K-B1 P-B7 K-Kt2 K-Q7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/1kPK4/8/8/8/8/8/8 b - - 2 12",
		// },
		// {
		// 	fen:         "8/1k6/2pK4/8/1P1P4/8/8/8 w - - 0 1",
		// 	s:           "K-Q7 K-Kt3 K-B8 K-R3 K-B7 K-Kt4 K-Kt7 KxP KxP K-B5 P-Q5 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/8/2K5/3P4/2k5/8/8/8 b - - 0 6",
		// },
		// {
		// 	fen:         "2k5/1p6/1P6/2PK4/8/8/8/8 w - - 0 1",
		// 	s:           "K-Q6 K-Q1 K-K6 K-B1 K-K7 K-Kt1 K-Q7 K-R1 P-B6 PxP K-B7 P-B4 P-Kt7ch K-R2 P-Kt8(Q)ch K-R3 Q-Kt6 mate ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/2K5/kQ6/2p5/8/8/8/8 b - - 2 9",
		// },
		// {
		// 	fen:         "8/8/kpK5/8/P1P5/8/8/8 w - - 0 1",
		// 	s:           "1.    K-Q7! K-Kt2 2.    P-R5! PxP 3.    P-B5 P-R5 4.    P-B6ch K-Kt3 5.    P-B7 P-R6 6.    P-B8(Q) P-R7 7.    Q-QR8 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "Q7/3K4/1k6/8/8/8/p7/8 b - - 1 7",
		// },
		// {
		// 	fen:         "3k4/2p5/2K5/1P1P4/8/8/8/8 w - - 0 1",
		// 	s:           "K-Kt7 K-Q2 K-Kt8 K-Q1 P-Q6! PxP P-Kt6 P-Q4 K-R7 P-Q5 P-Kt7 P-Q6 P-Kt8(Q)ch K-Q2 O-Kt5ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/K2k4/8/1Q6/8/3p4/8/8 b - - 2 8",
		// },
		// {
		// 	fen:         "4k3/8/4KP2/7p/8/6P1/8/8 w - - 0 1",
		// 	s:           "P-B7ch K-B1 K-B6 P-R5 P-Kt4 P-R6 P-Kt5 P-R7 P-Kt6 P-R8(Q) P-Kt7 mate ",
		// 	expectedErr: nil,
		// 	expectedFEN: "5k2/5PP1/5K2/8/8/8/8/7q b - - 0 6",
		// },
		// {
		// 	fen:         "6k1/7p/7K/7P/8/8/6P1/8 w - - 0 1",
		// 	s:           "P-Kt3! K-R1 P-Kt4 K-Kt1 P-Kt5 K-R1 P-Kt6 PxP PxP K-Kt1 P-Kt7 K-B2 K-R7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/5kPK/8/8/8/8/8/8 b - - 2 7",
		// },
		// {
		// 	fen:         "8/8/8/5kP1/4p2P/4K3/8/8 w - - 0 1",
		// 	s:           "K-B2 K-Kt3 K-K2 K-B4 K-K3 K-K4 P-Kt6! K-B3 P-R5 K-Kt2 KxP ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/6k1/6P1/7P/4K3/8/8/8 b - - 0 6",
		// },
		// {
		// 	fen:         "8/2k5/p1P5/P1K5/8/8/8/8 w - - 0 1",
		// 	s:           "1.    K-Q5 K-B1 2.    K-Q4! K-Q1 3.    K-B4 K-B1 4.    K-Q5 K-B2 5.    K-B5 K-B1 6.    K-Kt6 K-Kt1 7.    KxP K-B2 8.    K-Kt5 K-B1 9.    K-Kt6 K-Kt1 10.    P-B7ch K-B1 11.    P-R6 K-Q2 12.    K-Kt7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/1KPk4/P7/8/8/8/8/8 b - - 2 12",
		// },
		// {
		// 	fen:         "8/1pk5/8/1K6/1PP5/8/8/8 w - - 0 1",
		// 	s:           "P-B5 K-B1 K-Kt6 K-Kt1 P-B6 K-R1 K-B7 PxP KxP K-R2 P-Kt5 K-Kt1 K-Kt6 K-R1 K-B7 K-R2 P-Kt6ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/k1K5/1P6/8/8/8/8/8 b - - 0 9",
		// },
		// {
		// 	fen:         "8/8/8/8/4k1p1/8/4KPP1/8 w - - 0 1",
		// 	s:           "1. P-Kt3 K-Q5 2. P-B3 PxPch 3. KxP K-K4 4. K-Kt4! K-B3 5. K-R5 K-B4 6. P-Kt4ch K-B3 7. K-R6 K-B2 8. P-Kt5 K-Kt1 9. K-Kt6 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "6k1/8/6K1/6P1/8/8/8/8 b - - 2 9",
		// },
		// {
		// 	fen:         "8/8/1kp5/8/K1PP4/8/8/8 w - - 0 1",
		// 	s:           "K-Kt3 K-B2 K-B3 K-Q3 K-Q3 K-Q2 K-K4 K-K3 P-B5 K-B3 P-Q5! K-K2 P-Q6ch K-Q2 K-K5 K-Q1 P-Q7! KxP K-B6 K-Q1 K-K6 K-B2 K-K7 K-B1 K-Q6 K-Kt2 K-Q7 K-Kt1 KxP K-B1 K-Q6 K-Q1 P-B6 K-B1 P-B7 K-Kt2 K-Q7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/1kPK4/8/8/8/8/8/8 b - - 2 19",
		// },
		// {
		// 	fen:         "8/p7/P7/8/8/4k3/4P3/4K3 w - - 0 1",
		// 	s:           "K-B1 K-Q5 K-B2 K-B4 P-K4! K-Q5 K-B3 K-K4 K-K3 K-K3 K-Q4 K-Q3 P-K5ch K-K3 K-K4 K-K2 K-Q5 K-Q2 P-K6ch K-K2 K-K5 K-K1 K-Q6 K-Q1 K-B6 K-K2 K-Kt7 KxP KxP K-Q2 K-Kt7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/1K1k4/P7/8/8/8/8/8 b - - 2 16",
		// },
		// {
		// 	fen:         "3k4/1K3p2/8/4P3/8/6P1/8/8 w - - 0 1",
		// 	s:           "1.    K-B6 K-K2 2.    K-Q5 K-Q2 3.    K-Q4! K-K2 4.    K-K3 K-Q2 5.    K-B4 K-K3 6.    K-K4 K-Q2 7.    K-B5 K-K2 8.    P-Kt4 K-K1! 9.    K-B6! K-B1 10.    P-Kt5 K-K1 11.    K-Kt7! K-K2 12.    K-Kt8 K-K1 13.    P-K6! PxP 14.    P-Kt6 P-K4 15.    P-Kt7 P-K5 16.    K-R7 P-K6 17.    P-Kt8(Q)ch ",
		// 	expectedErr: nil,
		// 	expectedFEN: "4k1Q1/7K/8/8/8/4p3/8/8 b - - 0 17",
		// },
		// {
		// 	fen:         "8/6p1/8/5P2/5P2/5K2/8/6k1 w - - 0 1",
		// 	s:           "P-B6! PxP P-B5 K-R7 K-Kt4 K-Kt7 K-R5 K-B6 K-Kt6 K-Kt5 KxP K-R4 K-Kt7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/6K1/8/5P1k/8/8/8/8 b - - 2 7",
		// },
		// {
		// 	fen:         "8/6p1/8/K6P/7P/1k6/8/8 w - - 0 1",
		// 	s:           "K-Kt5 K-B6 K-B5 K-Q6 K-Q5 K-K6 K-K5 K-B6 K-B5 K-Kt6 P-R6! PxP P-R5 K-R5 K-Kt6 K-Kt5 KxP K-B4 K-Kt7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/6K1/8/5k1P/8/8/8/8 b - - 2 10",
		// },
		// {
		// 	fen:         "8/6p1/8/4K3/7P/4k2P/8/8 w - - 0 1",
		// 	s:           "1. K-K6! K-B5 2. K-B7  K-Kt6 3. P-R5  K-R5 4. K-Kt6! KxP(R6) 5. KxP ",
		// 	expectedErr: nil,
		// 	expectedFEN: "8/6K1/8/7P/8/7k/8/8 b - - 0 5",
		// },
		// {
		// 	fen:         "5k2/5p2/8/4P3/2K1P3/8/8/8 w - - 0 1",
		// 	s:           "P-K6! PxP P-K5! K-K2 K-B5 K-Q2 K-Kt6 K-Q1 K-B6 K-K2 K-B7 K-K1 K-Q6 K-B2 K-Q7 K-B1 KxP K-K1 K-B6 K-B1 P-K6 K-K1 P-K7 ",
		// 	expectedErr: nil,
		// 	expectedFEN: "4k3/4P3/5K2/8/8/8/8/8 b - - 0 12",
		// },
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test notation parser descriptive %v", i), func(t *testing.T) {
			g, err := newGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := newNotationParserDescriptive(characteristics{}).parse(g, tc.s)
			require.Equal(t, tc.expectedErr, err)
			if tc.expectedErr != nil {
				return
			}
			for _, gs := range gameSteps {
				fmt.Println(gs.a)
			}
			if tc.expectedMatchedTokens != nil {
				require.Len(t, tc.expectedMatchedTokens, len(gameSteps))
				for i, gameStep := range gameSteps {
					assert.Equal(t, tc.expectedMatchedTokens[i], gameStep.s)
				}
			}
			if tc.expectedBoard != nil {
				assert.Equal(t, tc.expectedBoard, gameSteps[len(gameSteps)-1].g.toBoard().board)
			}
			if tc.expectedFEN != "" {
				assert.Equal(t, tc.expectedFEN, gameSteps[len(gameSteps)-1].g.toFEN())
			}
		})
	}
}
