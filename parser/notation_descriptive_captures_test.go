package parser

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescriptiveNotation_CheernevGames(t *testing.T) {
	testCases := []struct {
		name      string
		fen       string
		s         string
		moveCount int
	}{
		{
			name:      "Chernev Game 1: promotion and resign",
			fen:       "8/8/8/8/8/1k5P/8/2K5 w - - 0 1",
			s:         "1. P-R4 K-B5\n2. P-R5 K-Q4\n3. P-R6 K-K3\n4. P-R7 K-B2\n5. P-R8(Q) Resigns",
			moveCount: 10,
		},
		{
			name:      "Chernev Game 2: promotion checkmate",
			fen:       "8/8/8/6K1/8/5k2/P7/8 w - - 0 1",
			s:         "1. K-B5 K-K6\n2. K-K5 K-Q6\n3. K-Q5 K-B6\n4. K-B5 K-Q6\n5. P-R4 K-B6\n6. P-R5 K-Kt6\n7. P-R6 K-R5\n8. P-R7 K-R4\n9. P-R8(Q)mate",
			moveCount: 17,
		},
		{
			name:      "Chernev Game 3: promotion with check",
			fen:       "4k3/8/3K4/4P3/8/8/8/8 w - - 0 1",
			s:         "1. K-K6 K-Q1\n2. K-B7 K-Q2\n3. P-K6ch K-Q1\n4. P-K7ch K-Q2\n5. P-K8(Q)ch",
			moveCount: 9,
		},
		{
			name:      "Chernev Game 5: pawn advance",
			fen:       "3k4/8/1P6/8/4K3/8/8/8 w - - 0 1",
			s:         "1. K-Q5 K-Q2\n2. K-B5 K-Q1\n3. K-Q6 K-B1\n4. K-B6 K-Kt1\n5. P-Kt7 K-R2\n6. K-B7",
			moveCount: 11,
		},
		{
			name:      "Chernev Game 6: king maneuver",
			fen:       "6k1/8/7K/8/2P5/8/8/8 w - - 0 1",
			s:         "1. K-Kt6 K-B1\n2. K-B6 K-K1\n3. K-K6 K-Q1\n4. K-Q6 K-B1\n5. K-B6 K-Kt1\n6. K-Q7 K-Kt2\n7. P-B5 K-Kt1\n8. P-B6 K-R2\n9. P-B7",
			moveCount: 17,
		},
		{
			name:      "Chernev Game 7: promotion",
			fen:       "8/7k/5K2/6P1/8/8/8/8 w - - 0 1",
			s:         "1. K-B7 K-R1\n2. K-Kt6 K-Kt1\n3. K-R6 K-R1\n4. P-Kt6 K-Kt1\n5. P-Kt7 K-B2\n6. K-R7 K-K2\n7. P-Kt8(Q)",
			moveCount: 13,
		},
		{
			name:      "Chernev Game 8: pawn advance with check",
			fen:       "8/8/5k2/8/5K2/8/4P3/8 w - - 0 1",
			s:         "1. K-K4 K-K3\n2. P-K3 K-Q3\n3. K-B5 K-K2\n4. K-K5 K-Q2\n5. K-B6 K-K1\n6. K-K6 K-B1\n7. P-K4 K-K1\n8. P-K5 K-B1\n9. K-Q7 K-B2\n10. P-K6ch K-B1\n11. P-K7ch",
			moveCount: 21,
		},
		{
			name:      "Chernev Game 9: resign",
			fen:       "8/7k/8/8/8/3P4/8/6K1 w - - 0 1",
			s:         "1. K-B2 K-Kt3\n2. K-K3 K-B4\n3. K-Q4 K-K3\n4. K-B5 K-Q2\n5. K-Q5 K-K2\n6. K-B6 K-K3\n7. P-Q4 K-K2\n8. P-Q5 K-Q1\n9. K-Q6 Resigns",
			moveCount: 18,
		},
		{
			name:      "Chernev Game 12: sacrifice promotion",
			fen:       "7k/7P/6P1/8/8/6K1/8/8 w - - 0 1",
			s:         "1. K-B4 K-Kt2\n2. K-B5 K-R1\n3. K-Kt5 K-Kt2\n4. P-R8(Q)ch KxQ\n5. K-B6 K-Kt1\n6. P-Kt7 K-R2\n7. K-B7 Resigns",
			moveCount: 14,
		},
		{
			name:      "Chernev Game 13: sacrifice and advance",
			fen:       "6k1/8/5KP1/6P1/8/8/8/8 w - - 0 1",
			s:         "1. P-Kt7 K-R2\n2. P-Kt8(Q)ch KxQ\n3. K-Kt6 K-R1\n4. K-B7 K-R2\n5. P-Kt6ch K-R1\n6. P-Kt7ch Resigns",
			moveCount: 12,
		},
		{
			name:      "Chernev Game 24: mutual pawn race",
			fen:       "8/6p1/7k/8/1K6/8/1P6/8 w - - 0 1",
			s:         "1. K-B5 P-Kt4\n2. P-Kt4 P-Kt5\n3. K-Q4 P-Kt6\n4. K-K3 K-Kt4\n5. P-Kt5 K-Kt5\n6. P-Kt6 K-R6\n7. P-Kt7 P-Kt7\n8. K-B2 K-R7\n9. P-Kt8(Q)ch",
			moveCount: 17,
		},
		{
			name:      "Chernev Game 85: en passant and pawn captures",
			fen:       "8/2pp2pp/8/2PP1P2/1p5k/8/PP3p2/5K2 w - - 0 1",
			s:         "1. P-KB6 PxP\n2. KxP K-Kt4\n3. P-R4 PxPe.p.\n4. PxP K-B4\n5. P-R4 K-K4\n6. P-Q6 PxP\n7. P-B6 PxP\n8. P-R5",
			moveCount: 15,
		},
		{
			name:      "Chernev Game 87: knight check",
			fen:       "8/8/8/4P3/2k5/6K1/8/2n5 w - - 0 1",
			s:         "1. P-K6 Kt-K7ch\n2. K-R2",
			moveCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, tc.s)
			require.NoError(t, err, "failed to parse")
			require.Len(t, gameSteps, tc.moveCount, "move count mismatch")
		})
	}
}

func TestDescriptiveNotation_ComplexGames(t *testing.T) {
	testCases := []struct {
		name      string
		fen       string
		s         string
		moveCount int
	}{
		{
			name: "French Defense (ostinato six-notation suite)",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s:    "1. P-K4 P-K3\n2. P-Q4 P-Q4\n3. N-QB3 B-N5\n4. B-N5ch B-Q2\n5. BxBch QxB\n6. KN-K2 PxP\n7. 0-0",
			moveCount: 13,
		},
		{
			name: "Evans Gambit (ostinato issue #1 game, 47 half-moves)",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: "1. P-K4 P-K4\n2. N-KB3 N-QB3\n3. B-B4 B-B4\n4. P-QN4 BxNP\n5. P-B3 B-R4\n6. P-Q4 PxP\n7. O-O P-Q6\n8. Q-N3 Q-B3\n9. P-K5 Q-N3\n10. R-K1 KN-K2\n11. B-R3 P-N4\n12. QxP R-QN1\n13. Q-R4 B-N3\n14. QN-Q2 B-N2\n15. N-K4 Q-B4\n16. BxQP Q-R4\n17. N-B6ch PxN\n18. PxP R-N1\n19. QR-Q1 QxN\n20. RxNch NxR\n21. QxPch KxQ\n22. B-B5dblch K-K1\n23. B-Q7ch K-B1\n24. BxNmate",
			moveCount: 47,
		},
		{
			name:      "Chernev Game 294: discovered check and underpromotion",
			fen:       "1q3k2/7K/1P1P1p2/R2P4/8/B5p1/8/8 w - - 0 1",
			s:         "1. P-Kt7 K-B2\n2. R-R8 QxKtP\n3. R-B8ch KxRch\n4. P-Q7dis.ch K-B2\n5. P-Q8(N)ch K-K1ch\n6. KtxQ",
			moveCount: 11,
		},
		{
			name:      "Chernev Game 296: rook disambiguation with parentheses",
			fen:       "RK6/8/1k6/7R/8/8/pp6/8 w - - 0 1",
			s:         "1. R(R5)-QR5 K-B3\n2. K-B8 K-Q3\n3. K-Q8 K-K3\n4. R(R8)-R6ch K-B2\n5. R-B5ch K-Kt2\n6. R-Kt5ch K-B2\n7. R(Kt)-Kt6 P-Kt8(Q)\n8. R(R6)-QB6",
			moveCount: 15,
		},
		{
			name:      "Chernev Game 278: complex tactics with checkmate",
			fen:       "2b4k/8/5Pr1/5N2/8/8/8/K1B5 w - - 0 1",
			s:         "1. P-B7 R-R3ch\n2. B-R3 RxBch\n3. K-Kt2 R-R7ch\n4. K-B1 R-R8ch\n5. K-Q2 R-R7ch\n6. K-K3 R-R6ch\n7. K-B4 R-R5ch\n8. K-Kt5 R-Kt5ch\n9. K-R6 R-Kt1\n10. Kt-K7 B-K3\n11. PxR(Q)ch BxQ\n12. Kt-Kt6mate",
			moveCount: 23,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, tc.s)
			require.NoError(t, err, "failed to parse")
			require.Len(t, gameSteps, tc.moveCount, "move count mismatch")
		})
	}
}

func TestDescriptiveNotation_KtKnightVariant(t *testing.T) {
	g, err := core.NewGameFromFEN("8/1pK5/8/8/8/8/1k3P2/8 w - - 0 1")
	require.NoError(t, err)
	stepsN, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, "1. KxP K-N6\n2. K-B6 K-B5\n3. K-Q6 K-Q5\n4. P-B4 K-K5")
	require.NoError(t, err)
	stepsKt, err := NewNotationParserDescriptive(Characteristics{}).Parse(g, "1. KxP K-Kt6\n2. K-B6 K-B5\n3. K-Q6 K-Q5\n4. P-B4 K-K5")
	require.NoError(t, err)
	require.Len(t, stepsN, 8)
	require.Len(t, stepsKt, 8)
	for i := range stepsN {
		assert.Equal(t, stepsN[i].StepGame.ToFEN(), stepsKt[i].StepGame.ToFEN(),
			"Move %d: FEN mismatch between N and Kt", i+1)
	}
}
