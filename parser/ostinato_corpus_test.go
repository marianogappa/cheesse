package parser

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/printer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// descriptiveToAlgebraic parses a descriptive-notation game and renders each move in
// canonical SAN, mirroring ostinato's descriptiveToAlgebraic acceptance helper.
func descriptiveToAlgebraic(t *testing.T, fen, s string) []string {
	t.Helper()
	g, err := core.NewGameFromFEN(fen)
	require.NoError(t, err)
	steps, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, s)
	require.NoError(t, err)

	gc := printer.SANCharacteristics()
	out := make([]string, len(steps))
	for i, step := range steps {
		if step.StepAction.IsResign {
			out[i] = "resigns"
			continue
		}
		san, err := printer.AlgebraicPrinter{}.PrintAction(step, gc)
		require.NoError(t, err)
		out[i] = san
	}
	return out
}

// TestOstinatoCorpus_DescriptiveToAlgebraic ports ostinato's move-by-move canonical
// SAN assertions for the Chernev games.
func TestOstinatoCorpus_DescriptiveToAlgebraic(t *testing.T) {
	testCases := []struct {
		name     string
		fen      string
		s        string
		expected []string
	}{
		{
			name:     "Chernev Game 1",
			fen:      "8/8/8/8/8/1k5P/8/2K5 w - - 0 1",
			s:        "1. P-R4 K-B5\n2. P-R5 K-Q4\n3. P-R6 K-K3\n4. P-R7 K-B2\n5. P-R8(Q) Resigns",
			expected: []string{"h4", "Kc4", "h5", "Kd5", "h6", "Ke6", "h7", "Kf7", "h8=Q", "resigns"},
		},
		{
			name:     "Chernev Game 2",
			fen:      "8/8/8/6K1/8/5k2/P7/8 w - - 0 1",
			s:        "1. K-B5 K-K6\n2. K-K5 K-Q6\n3. K-Q5 K-B6\n4. K-B5 K-Q6\n5. P-R4 K-B6\n6. P-R5 K-Kt6\n7. P-R6 K-R5\n8. P-R7 K-R4\n9. P-R8(Q)mate",
			expected: []string{"Kf5", "Ke3", "Ke5", "Kd3", "Kd5", "Kc3", "Kc5", "Kd3", "a4", "Kc3", "a5", "Kb3", "a6", "Ka4", "a7", "Ka5", "a8=Q#"},
		},
		{
			name:     "Chernev Game 3",
			fen:      "4k3/8/3K4/4P3/8/8/8/8 w - - 0 1",
			s:        "1. K-K6 K-Q1\n2. K-B7 K-Q2\n3. P-K6ch K-Q1\n4. P-K7ch K-Q2\n5. P-K8(Q)ch",
			expected: []string{"Ke6", "Kd8", "Kf7", "Kd7", "e6+", "Kd8", "e7+", "Kd7", "e8=Q+"},
		},
		{
			name:     "Chernev Game 5",
			fen:      "3k4/8/1P6/8/4K3/8/8/8 w - - 0 1",
			s:        "1. K-Q5 K-Q2\n2. K-B5 K-Q1\n3. K-Q6 K-B1\n4. K-B6 K-Kt1\n5. P-Kt7 K-R2\n6. K-B7",
			expected: []string{"Kd5", "Kd7", "Kc5", "Kd8", "Kd6", "Kc8", "Kc6", "Kb8", "b7", "Ka7", "Kc7"},
		},
		{
			name:     "Chernev Game 85 with en passant",
			fen:      "8/2pp2pp/8/2PP1P2/1p5k/8/PP3p2/5K2 w - - 0 1",
			s:        "1. P-KB6 PxP\n2. KxP K-Kt4\n3. P-R4 PxPe.p.\n4. PxP K-B4\n5. P-R4 K-K4\n6. P-Q6 PxP\n7. P-B6 PxP\n8. P-R5",
			expected: []string{"f6", "gxf6", "Kxf2", "Kg5", "a4", "bxa3 e.p.", "bxa3", "Kf5", "a4", "Ke5", "d6", "cxd6", "c6", "dxc6", "a5"},
		},
		{
			name:     "Chernev Game 294 with underpromotion and discovered check",
			fen:      "1q3k2/7K/1P1P1p2/R2P4/8/B5p1/8/8 w - - 0 1",
			s:        "1. P-Kt7 K-B2\n2. R-R8 QxKtP\n3. R-B8ch KxRch\n4. P-Q7dis.ch K-B2\n5. P-Q8(N)ch K-K1ch\n6. KtxQ",
			expected: []string{"b7", "Kf7", "Ra8", "Qxb7", "Rf8+", "Kxf8+", "d7+", "Kf7", "d8=N+", "Ke8+", "Nxb7"},
		},
		{
			name:     "Chernev Game 6",
			fen:      "6k1/8/7K/8/2P5/8/8/8 w - - 0 1",
			s:        "1. K-Kt6 K-B1\n2. K-B6 K-K1\n3. K-K6 K-Q1\n4. K-Q6 K-B1\n5. K-B6 K-Kt1\n6. K-Q7 K-Kt2\n7. P-B5 K-Kt1\n8. P-B6 K-R2\n9. P-B7",
			expected: []string{"Kg6", "Kf8", "Kf6", "Ke8", "Ke6", "Kd8", "Kd6", "Kc8", "Kc6", "Kb8", "Kd7", "Kb7", "c5", "Kb8", "c6", "Ka7", "c7"},
		},
		{
			name:     "Chernev Game 7",
			fen:      "8/7k/5K2/6P1/8/8/8/8 w - - 0 1",
			s:        "1. K-B7 K-R1\n2. K-Kt6 K-Kt1\n3. K-R6 K-R1\n4. P-Kt6 K-Kt1\n5. P-Kt7 K-B2\n6. K-R7 K-K2\n7. P-Kt8(Q)",
			expected: []string{"Kf7", "Kh8", "Kg6", "Kg8", "Kh6", "Kh8", "g6", "Kg8", "g7", "Kf7", "Kh7", "Ke7", "g8=Q"},
		},
		{
			name:     "Chernev Game 8",
			fen:      "8/8/5k2/8/5K2/8/4P3/8 w - - 0 1",
			s:        "1. K-K4 K-K3\n2. P-K3 K-Q3\n3. K-B5 K-K2\n4. K-K5 K-Q2\n5. K-B6 K-K1\n6. K-K6 K-B1\n7. P-K4 K-K1\n8. P-K5 K-B1\n9. K-Q7 K-B2\n10. P-K6ch K-B1\n11. P-K7ch",
			expected: []string{"Ke4", "Ke6", "e3", "Kd6", "Kf5", "Ke7", "Ke5", "Kd7", "Kf6", "Ke8", "Ke6", "Kf8", "e4", "Ke8", "e5", "Kf8", "Kd7", "Kf7", "e6+", "Kf8", "e7+"},
		},
		{
			name:     "Chernev Game 9",
			fen:      "8/7k/8/8/8/3P4/8/6K1 w - - 0 1",
			s:        "1. K-B2 K-Kt3\n2. K-K3 K-B4\n3. K-Q4 K-K3\n4. K-B5 K-Q2\n5. K-Q5 K-K2\n6. K-B6 K-K3\n7. P-Q4 K-K2\n8. P-Q5 K-Q1\n9. K-Q6 Resigns",
			expected: []string{"Kf2", "Kg6", "Ke3", "Kf5", "Kd4", "Ke6", "Kc5", "Kd7", "Kd5", "Ke7", "Kc6", "Ke6", "d4", "Ke7", "d5", "Kd8", "Kd6", "resigns"},
		},
		{
			name:     "Chernev Game 12",
			fen:      "7k/7P/6P1/8/8/6K1/8/8 w - - 0 1",
			s:        "1. K-B4 K-Kt2\n2. K-B5 K-R1\n3. K-Kt5 K-Kt2\n4. P-R8(Q)ch KxQ\n5. K-B6 K-Kt1\n6. P-Kt7 K-R2\n7. K-B7 Resigns",
			expected: []string{"Kf4", "Kg7", "Kf5", "Kh8", "Kg5", "Kg7", "h8=Q+", "Kxh8", "Kf6", "Kg8", "g7", "Kh7", "Kf7", "resigns"},
		},
		{
			name:     "Chernev Game 13",
			fen:      "6k1/8/5KP1/6P1/8/8/8/8 w - - 0 1",
			s:        "1. P-Kt7 K-R2\n2. P-Kt8(Q)ch KxQ\n3. K-Kt6 K-R1\n4. K-B7 K-R2\n5. P-Kt6ch K-R1\n6. P-Kt7ch Resigns",
			expected: []string{"g7", "Kh7", "g8=Q+", "Kxg8", "Kg6", "Kh8", "Kf7", "Kh7", "g6+", "Kh8", "g7+", "resigns"},
		},
		{
			name:     "Chernev Game 24",
			fen:      "8/6p1/7k/8/1K6/8/1P6/8 w - - 0 1",
			s:        "1. K-B5 P-Kt4\n2. P-Kt4 P-Kt5\n3. K-Q4 P-Kt6\n4. K-K3 K-Kt4\n5. P-Kt5 K-Kt5\n6. P-Kt6 K-R6\n7. P-Kt7 P-Kt7\n8. K-B2 K-R7\n9. P-Kt8(Q)ch",
			expected: []string{"Kc5", "g5", "b4", "g4", "Kd4", "g3", "Ke3", "Kg5", "b5", "Kg4", "b6", "Kh3", "b7", "g2", "Kf2", "Kh2", "b8=Q+"},
		},
		{
			name:     "Chernev Game 87",
			fen:      "8/8/8/4P3/2k5/6K1/8/2n5 w - - 0 1",
			s:        "1. P-K6 Kt-K7ch\n2. K-R2",
			expected: []string{"e6", "Ne2+", "Kh2"},
		},
		{
			name:     "Chernev Game 296 with rook disambiguation",
			fen:      "RK6/8/1k6/7R/8/8/pp6/8 w - - 0 1",
			s:        "1. R(R5)-QR5 K-B3\n2. K-B8 K-Q3\n3. K-Q8 K-K3\n4. R(R8)-R6ch K-B2\n5. R-B5ch K-Kt2\n6. R-Kt5ch K-B2\n7. R(Kt)-Kt6 P-Kt8(Q)\n8. R(R6)-QB6",
			// N.B. ostinato expected Ra6+ and Rc6 here, but both rooks share the a-file
			// after move 1 (a8/a5) and both can reach c6 at the end, so strict SAN
			// requires R8a6+ and Rac6.
			expected: []string{"Rha5", "Kc6", "Kc8", "Kd6", "Kd8", "Ke6", "R8a6+", "Kf7", "Rf5+", "Kg7", "Rg5+", "Kf7", "Rgg6", "b1=Q", "Rac6"},
		},
		{
			name:     "Chernev Game 278 with capture-promotion and checkmate",
			fen:      "2b4k/8/5Pr1/5N2/8/8/8/K1B5 w - - 0 1",
			s:        "1. P-B7 R-R3ch\n2. B-R3 RxBch\n3. K-Kt2 R-R7ch\n4. K-B1 R-R8ch\n5. K-Q2 R-R7ch\n6. K-K3 R-R6ch\n7. K-B4 R-R5ch\n8. K-Kt5 R-Kt5ch\n9. K-R6 R-Kt1\n10. Kt-K7 B-K3\n11. PxR(Q)ch BxQ\n12. Kt-Kt6mate",
			expected: []string{"f7", "Ra6+", "Ba3", "Rxa3+", "Kb2", "Ra2+", "Kc1", "Ra1+", "Kd2", "Ra2+", "Ke3", "Ra3+", "Kf4", "Ra4+", "Kg5", "Rg4+", "Kh6", "Rg8", "Ne7", "Be6", "fxg8=Q+", "Bxg8", "Ng6#"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := descriptiveToAlgebraic(t, tc.fen, tc.s)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// TestOstinatoCorpus_DescriptiveInitialBoard2 ports "should parse descriptive notation
// with initial board 2": P-Kt8=Q equals-promotion and 1-0 result marker.
func TestOstinatoCorpus_DescriptiveInitialBoard2(t *testing.T) {
	g, err := core.NewGameFromFEN("8/1pK5/8/8/8/8/k4P2/8 w - - 0 1")
	require.NoError(t, err)
	s := `1. K-Q6      K-R6
		2. K-B5       K-R5
		3. P-B4       P-Kt4
		4. P-B5       P-Kt5
		5. K-B4      P-Kt6
		6. K-B3      K-R6
		7. P-B6       P-Kt7
		8. P-B7       P-Kt8=Q
		9. P-B8(Q)ch  K-R5
		10. Q-R8ch    K-Kt4
		11. Q-Kt7ch`
	steps, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, s)
	require.NoError(t, err)
	require.Len(t, steps, 21)

	// Both sides promoted to queens
	lastFEN := steps[len(steps)-1].StepGame.ToFEN()
	assert.Contains(t, lastFEN, "Q") // White queen on the board
	assert.Contains(t, lastFEN, "q") // Black queen on the board
	assert.True(t, steps[len(steps)-1].StepGame.IsCheck)
}

