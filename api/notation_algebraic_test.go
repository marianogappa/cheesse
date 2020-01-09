package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserAlgebraic(t *testing.T) {
	testCases := []struct {
		fen                   string
		s                     string
		expectedErr           error
		expectedMatchedTokens []string
		expectedBoard         []string
		expectedFEN           string
	}{
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. e4 e6
				2. d4 d5
				3. Nc3 Bb4
				4. Bb5+ Bd7
				5. Bxd7+ Qxd7
				6. Ne2 dxe4
				7. 0-0`,
			expectedErr:           nil,
			expectedMatchedTokens: []string{"e4", "e6", "d4", "d5", "Nc3", "Bb4", "Bb5+", "Bd7", "Bxd7+", "Qxd7", "Ne2", "dxe4", "0-0"},
			expectedBoard: []string{
				"♜♞  ♚ ♞♜",
				"♟♟♟♛ ♟♟♟",
				"    ♟   ",
				"        ",
				" ♝ ♙♟   ",
				"  ♘     ",
				"♙♙♙ ♘♙♙♙",
				"♖ ♗♕ ♖♔ ",
			},
		},
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			s: `1. e4 e6
				2. d4 d5
				3. Nc3 Bb4
				4. Bb5+ Bd7
				5. Bxd7† Qxd7
				6. Ne2 dxe4
				7. 0-0`,
			expectedErr: fmt.Errorf("expecting CheckSymbol + but found †"),
		},
		{
			fen: "8/8/8/8/8/1k5P/8/2K5 w - - 0 1",
			s: `1. h4 Kc4
				2. h5 Kd5
				3. h6 Ke6
				4. h7 Kf7
				5. h8=Q`,
			expectedErr:           nil,
			expectedMatchedTokens: []string{"h4", "Kc4", "h5", "Kd5", "h6", "Ke6", "h7", "Kf7", "h8=Q"},
			expectedBoard: []string{
				"       ♕",
				"     ♚  ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"  ♔     ",
			},
		},
		{
			fen: "8/8/8/8/8/1k5P/8/2K5 w - - 0 1",
			s: `1. h4 Kc4
				2. h5 Kd5
				3. h6 Ke6
				4. h7 Kf7
				5. h8=Q`,
			expectedFEN: "7Q/5k2/8/8/8/8/8/2K5 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/6K1/8/5k2/P7/8 w - - 0 1",
			s: `1. Kf5! Ke3
				2. Ke5! Kd3
				3. Kd5! Kc3
				4. Kc5! Kd3
				5. a4 Kc3
				6. a5 Kb3
				7. a6 Ka4
				8. a7 Ka5
				9. a8=Q# `,
			expectedFEN: "Q7/8/8/k1K5/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "4k3/8/3K4/4P3/8/8/8/8 w - - 0 1",
			s: `1. Ke6! Kd8
				2. Kf7  Kd7
				3. e6+  Kd8
				4. e7+  Kd7
				5. e8=Q+ `,
			expectedFEN: "4Q3/3k1K2/8/8/8/8/8/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "8/5k2/3P4/8/6K1/8/8/8 w - - 0 1",
			s: `1. Kf5  Kf8
				2. Kf6! Ke8
				3. Ke6  Kd8
				4. d7 Kc7
				5. Ke7 `,
			expectedFEN: "8/2kPK3/8/8/8/8/8/8 b - - 2 5",
			expectedErr: nil,
		},
		{
			fen: "3k4/8/1P6/8/4K3/8/8/8 w - - 0 1",
			s: `1. Kd5  Kd7
				2. Kc5  Kd8
				3. Kd6! Kc8
				4. Kc6  Kb8
				5. b7 Ka7
				6. Kc7`,
			expectedFEN: "8/kPK5/8/8/8/8/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "6k1/8/7K/8/2P5/8/8/8 w - - 0 1",
			s: `1. Kg6! Kf8
				2. Kf6  Ke8
				3. Ke6  Kd8
				4. Kd6  Kc8
				5. Kc6  Kb8
				6. Kd7  Kb7
				7. c5 Kb8
				8. c6 Ka7
				9. c7 `,
			expectedFEN: "8/k1PK4/8/8/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/7k/5K2/6P1/8/8/8/8 w - - 0 1",
			s: `1. Kf7! Kh8
				2. Kg6! Kg8
				3. Kh6  Kh8
				4. g6 Kg8
				5. g7 Kf7
				6. Kh7  Ke7
				7. g8=Q `,
			expectedFEN: "6Q1/4k2K/8/8/8/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/5k2/8/5K2/8/4P3/8 w - - 0 1",
			s: `1. Ke4  Ke6
				2. e3!  Kd6
				3. Kf5  Ke7
				4. Ke5  Kd7
				5. Kf6  Ke8
				6. Ke6  Kf8
				7. e4 Ke8
				8. e5 Kf8
				9. Kd7  Kf7
				10. e6+ Kf8
				11. e7+`,
			expectedFEN: "5k2/3KP3/8/8/8/8/8/8 b - - 0 11",
			expectedErr: nil,
		},
		{
			fen: "8/7k/8/8/8/3P4/8/6K1 w - - 0 1",
			s: `1. Kf2  Kg6
				2. Ke3  Kf5
				3. Kd4! Ke6
				4. Kc5  Kd7
				5. Kd5  Ke7
				6. Kc6  Ke6
				7. d4 Ke7
				8. d5 Kd8
				9. Kd6`,
			expectedFEN: "3k4/8/3K4/3P4/8/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/1k6/8/K7/6P1/8 w - - 0 1",
			s: `1. Kb3  Kc5
				2. Kc3  Kd5
				3. Kd3  Ke5
				4. Ke3  Kf5
				5. Kf3  Kg5
				6. Kg3  Kf5
				7. Kh4  Kf6
				8. Kh5  Kg7
				9. Kg5  Kf7
				10. Kh6 Kg8
				11. Kg6`,
			expectedFEN: "6k1/8/6K1/8/8/8/6P1/8 b - - 21 11",
			expectedErr: nil,
		},
		{
			fen: "3k4/8/K7/8/8/8/P7/8 w - - 0 1",
			s: `1. Kb7! Kd7
				2. a4 Kd6
				3. a5 Kc5
				4. a6`,
			expectedFEN: "8/1K6/P7/2k5/8/8/8/8 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "7k/7P/6P1/8/8/6K1/8/8 w - - 0 1",
			s: `1. Kf4      Kg7
				2. Kf5      Kh8
				3. Kg5      Kg7
				4. h8=Q+! Kxh8
				5. Kf6      Kg8
				6. g7     Kh7
				7. Kf7 `,
			expectedFEN: "8/5KPk/8/8/8/8/8/8 b - - 2 7",
			expectedErr: nil,
		},
		{
			fen: "6k1/8/5KP1/6P1/8/8/8/8 w - - 0 1",
			s: `1. g7 Kh7
				2. g8=Q+! Kxg8
				3. Kg6  Kh8
				4. Kf7  Kh7
				5. g6+  Kh8
				6. g7+`,
			expectedFEN: "7k/5KP1/8/8/8/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/2k5/P1P5/7K/8 w - - 0 1",
			s: `1. a4!  Kc5
				2. Kg3  Kb6
				3. c4 Ka5
				4. c5 Ka6
				5. Kf3  Kb7
				6. a5 Kc6
				7. a6 Kc7
				8. Ke4  Kc6
				9. Ke5  Kc7
				10. Kd5 Kb8
				11. c6  Kc7
				12. a7 `,
			expectedFEN: "8/P1k5/2P5/3K4/8/8/8/8 b - - 0 12",
			expectedErr: nil,
		},
		{
			fen: "8/8/5Ppk/8/8/4K3/8/8 w - - 0 1",
			s: `1. Kf4  Kh7
				2. Kg5  Kh8
				3. Kh6! Kg8
				4. Kxg6 Kh8
				5. Kf7  Kh7
				6. Ke8`,
			expectedFEN: "4K3/7k/5P2/8/8/8/8/8 b - - 4 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/5k2/8/p7/8/1PK5/8 w - - 0 1",
			s: `1. Kb1! a3
				2. b3!   Ke5
				3. Ka2   Kd5
				4. Kxa3 Kc6
				5. Ka4!  Kb6
				6. Kb4`,
			expectedFEN: "8/8/1k6/8/1K6/1P6/8/8 b - - 4 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/2p5/8/8/P4k2/8/2K5 w - - 0 1",
			s: `1. a4 Ke4
				2. a5 Kd5
				3. a6 Kd6
				4. a7 Kc7
				5. a8=Q`,
			expectedFEN: "Q7/2k5/2p5/8/8/8/8/2K5 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "8/p4K2/P7/8/8/8/1k6/8 w - - 0 1",
			s: `1. Ke6    Kc3
				2. Kd5!  Kb4
				3. Kc6    Kc4
				4. Kb7    Kb5
				5. Kxa7 Kc6
				6. Kb8 `,
			expectedFEN: "1K6/8/P1k5/8/8/8/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/7p/1PK2k2/8/8/8 w - - 0 1",
			s: `1. b5 Ke5
				2. b6!  Kd6
				3. Kb5  h4
				4. Ka6  h3
				5. b7 Kc7
				6. Ka7  h2
				7. b8=Q+`,
			expectedFEN: "1Q6/K1k5/8/8/8/8/7p/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/3KP1k1/6p1/8/8/8 w - - 0 1",
			s: `1. e6 Kf6
				2. Kd6 g3
				3. e7 g2
				4. e8=Q g1=Q
				5. Qf8+ Kg5
				6. Qg8+`,
			expectedFEN: "6Q1/8/3K4/6k1/8/8/8/6q1 b - - 3 6",
			expectedErr: nil,
		},
		{
			fen: "8/1p6/7K/8/7P/8/5k2/8 w - - 0 1",
			s: `1. Kg5       b5
				2. Kf4       Ke2
				3. Ke4       Kd2
				4. Kd4       Kc2
				5. Kc5       Kc3
				6. h5      b4
				7. h6      b3
				8. h7      b2
				9. h8=Q+ Kc2
				10. Qh2+   Kc1
				11. Qf4+    Kc2
				12. Qc4+    Kd2
				13. Qb3      Kc1
				14. Qc3+    Kb1
				15. Kc4      Ka2
				16. Qa5+    Kb1
				17. Kb3      Kc1
				18. Qe1# `,
			expectedFEN: "8/8/8/8/8/1K6/1p6/2k1Q3 b - - 18 18",
			expectedErr: nil,
		},
		{
			fen: "8/8/5P2/8/6K1/7p/6k1/8 w - - 0 1",
			s: `1. f7   h2
				2. f8=Q h1=Q
				3. Qf3+ Kg1
				4. Qe3+ Kf1
				5. Qc1+ Kg2
				6. Qd2+ Kf1
				7. Qd1+ Kg2
				8. Qe2+ Kg1
				9. Kg3!`,
			expectedFEN: "8/8/8/8/8/6K1/4Q3/6kq b - - 13 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/1p6/8/8/6P1/k1K5/8 w - - 0 1",
			s: `1. Kc3!      Ka3
				2. Kc4       Ka4
				3. g4      b5+
				4. Kd3!     Ka3
				5. g5      b4
				6. g6      b3
				7. g7      b2
				8. Kc2!      Ka2
				9. g8=Q+ Ka1
				10. Qa8# `,
			expectedFEN: "Q7/8/8/8/8/8/1pK5/k7 b - - 2 10",
			expectedErr: nil,
		},
		{
			fen: "8/6p1/7k/8/1K6/8/1P6/8 w - - 0 1",
			s: `1. Kc5! g5
				2. b4  g4
				3. Kd4   g3
				4. Ke3   Kg5
				5. b5  Kg4
				6. b6  Kh3
				7. b7  g2
				8. Kf2   Kh2
				9. b8=Q+ `,
			expectedFEN: "1Q6/8/8/8/8/8/5Kpk/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/1pK5/8/8/8/8/k4P2/8 w - - 0 1",
			s: `1. Kd6!    Ka3
				2. Kc5      Ka4
				3. f4     b5
				4. f5     b4
				5. Kc4!     b3
				6. Kc3!     Ka3
				7. f6     b2
				8. f7     b1=Q
				9. f8=Q+ Ka4
				10. Qa8+  Kb5
				11. Qb7+ `,
			expectedFEN: "8/1Q6/8/1k6/8/2K5/8/1q6 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "8/2p5/6K1/8/8/5k2/P7/8 w - - 0 1",
			s: `1. Kf5!     Ke3
				2. Ke5      c6
				3. a4     Kd3
				4. a5     c5
				5. a6     c4
				6. a7     c3
				7. a8=Q   c2
				8. Qd5+! Ke2
				9. Qa2!     Kd1
				10. Kd4     c1=Q
				11. Kd3 `,
			expectedFEN: "8/8/8/8/8/3K4/Q7/2qk4 b - - 1 11",
			expectedErr: nil,
		},
		{
			fen: "8/4K1pp/8/8/8/8/k6P/8 w - - 0 1",
			s: `1. h4!       h5
				2. Kf8!     g6
				3. Ke7!      g5
				4. hxg5     h4
				5. g6      h3
				6. g7      h2
				7. g8=Q+ Ka3
				8. Qg2 `,
			expectedFEN: "8/4K3/8/8/8/k7/6Qp/8 b - - 2 8",
			expectedErr: nil,
		},
		{
			fen: "8/6k1/1p6/6KP/P7/8/8/8 w - - 0 1",
			s: `1. Kf5    Kh6
				2. Ke5    Kxh5
				3. Kd5    Kg6
				4. Kc6    Kf6
				5. Kxb6 Ke7
				6. Kc7 `,
			expectedFEN: "8/2K1k3/8/8/P7/8/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/5k2/p7/P3KP2/8/8/8 w - - 0 1",
			s: `1. Kd5    Kf5
				2. Kc5    Kxf4
				3. Kb5    Ke5
				4. Kxa5 Kd6
				5. Kb6    Kd7
				6. Kb7    Kd8
				7. a5 `,
			expectedFEN: "3k4/1K6/8/P7/8/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/2p1k3/2P5/1P2K3/8/8/8/8 w - - 0 1",
			s: `1. Kd5!  Kd8
				2. Ke6    Kc8
				3. Ke7    Kb8
				4. Kd8    Ka7
				5. Kxc7 Ka8
				6. Kd8`,
			expectedFEN: "k2K4/8/2P5/1P6/8/8/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "8/kp6/8/PP6/1K6/8/8/8 w - - 0 1",
			s: `1. Kc5  Kb8
				2. Kb6  Ka8
				3. Kc7  Ka7
				4. a6 bxa6
				5. b6+  Ka8
				6. b7+  Ka7
				7. b8=Q# `,
			expectedFEN: "1Q6/k1K5/p7/8/8/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/6k1/3p4/3P1PK1/8/8/8/8 w - - 0 1",
			s: `1. f6+    Kf8!
				2. f7!    Kxf7
				3. Kf5    Ke7
				4. Kg6    Kd8
				5. Kf6    Kd7
				6. Kf7    Kd8
				7. Ke6    Kc7
				8. Ke7    Kc8
				9. Kxd6 Kd8
				10. Kc6   Kc8
				11. d6    Kd8
				12. d7    Ke7
				13. Kc7`,
			expectedFEN: "8/2KPk3/8/8/8/8/8/8 b - - 2 13",
			expectedErr: nil,
		},
		{
			fen: "8/3k4/2pP4/2P1K3/8/8/8/8 w - - 0 1",
			s: `1. Kf6    Kd8!
				2. d7!    Kxd7
				3. Kf7!   Kd8
				4. Ke6    Kc7
				5. Ke7    Kc8
				6. Kd6    Kb7
				7. Kd7    Kb8
				8. Kxc6 Kc8
				9. Kd6    Kd8
				10. c6    Kc8
				11. c7    Kb7
				12. Kd7 `,
			expectedFEN: "8/1kPK4/8/8/8/8/8/8 b - - 2 12",
			expectedErr: nil,
		},
		{
			fen: "8/1k6/2pK4/8/1P1P4/8/8/8 w - - 0 1",
			s: `1. Kd7    Kb6
				2. Kc8    Ka6
				3. Kc7    Kb5
				4. Kb7    Kxb4
				5. Kxc6 Kc4
				6. d5`,
			expectedFEN: "8/8/2K5/3P4/2k5/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "2k5/1p6/1P6/2PK4/8/8/8/8 w - - 0 1",
			s: `1. Kd6      Kd8
				2. Ke6      Kc8
				3. Ke7      Kb8
				4. Kd7      Ka8
				5. c6     bxc6
				6. Kc7      c5
				7. b7+      Ka7
				8. b8=Q+ Ka6
				9. Qb6#`,
			expectedFEN: "8/2K5/kQ6/2p5/8/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/kpK5/8/P1P5/8/8/8 w - - 0 1",
			s: `1. Kd7!  Kb7
				2. a5!    bxa5
				3. c5   a4
				4. c6+    Kb6
				5. c7   a3
				6. c8=Q a2
				7. Qa8`,
			expectedFEN: "Q7/3K4/1k6/8/8/8/p7/8 b - - 1 7",
			expectedErr: nil,
		},
		{
			fen: "3k4/2p5/2K5/1P1P4/8/8/8/8 w - - 0 1",
			s: `1. Kb7      Kd7
				2. Kb8      Kd8
				3. d6!      cxd6
				4. b6     d5
				5. Ka7      d4
				6. b7     d3
				7. b8=Q+ Kd7
				8. Qb5+`,
			expectedFEN: "8/K2k4/8/1Q6/8/3p4/8/8 b - - 2 8",
			expectedErr: nil,
		},
		{
			fen: "4k3/8/4KP2/7p/8/6P1/8/8 w - - 0 1",
			s: `1. f7+  Kf8
				2. Kf6  h4
				3. g4 h3
				4. g5 h2
				5. g6 h1=Q
				6. g7# `,
			expectedFEN: "5k2/5PP1/5K2/8/8/8/8/7q b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "6k1/7p/7K/7P/8/8/6P1/8 w - - 0 1",
			s: `1. g3!    Kh8
				2. g4     Kg8
				3. g5     Kh8
				4. g6  hxg6
				5. hxg6 Kg8
				6. g7  Kf7
				7. Kh7 `,
			expectedFEN: "8/5kPK/8/8/8/8/8/8 b - - 2 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/5kP1/4p2P/4K3/8/8 w - - 0 1",
			s: `1. Kf2  Kg6
				2. Ke2  Kf5
				3. Ke3  Ke5
				4. g6!  Kf6
				5. h5 Kg7
				6. Kxe4 `,
			expectedFEN: "8/6k1/6P1/7P/4K3/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/2k5/p1P5/P1K5/8/8/8/8 w - - 0 1",
			s: `1. Kd5    Kc8
				2. Kd4!  Kd8
				3. Kc4    Kc8
				4. Kd5    Kc7
				5. Kc5    Kc8
				6. Kb6    Kb8
				7. Kxa6 Kc7
				8. Kb5    Kc8
				9. Kb6    Kb8
				10. c7+   Kc8
				11. a6    Kd7
				12. Kb7 `,
			expectedFEN: "8/1KPk4/P7/8/8/8/8/8 b - - 2 12",
			expectedErr: nil,
		},
		{
			fen: "8/1pk5/8/1K6/1PP5/8/8/8 w - - 0 1",
			s: `1. c5   Kc8
				2. Kb6    Kb8
				3. c6   Ka8
				4. Kc7    bxc6
				5. Kxc6 Ka7
				6. b5   Kb8
				7. Kb6    Ka8
				8. Kc7    Ka7
				9. b6+ `,
			expectedFEN: "8/k1K5/1P6/8/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/4k1p1/8/4KPP1/8 w - - 0 1",
			s: `1. g3   Kd4
				2. f3   gxf3+
				3. Kxf3  Ke5
				4. Kg4!  Kf6
				5. Kh5    Kf5
				6. g4+    Kf6
				7. Kh6    Kf7
				8. g5   Kg8
				9. Kg6`,
			expectedFEN: "6k1/8/6K1/6P1/8/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/1kp5/8/K1PP4/8/8/8 w - - 0 1",
			s: `1. Kb3    Kc7
				2. Kc3    Kd6
				3. Kd3    Kd7
				4. Ke4    Ke6
				5. c5   Kf6
				6. d5!    Ke7
				7. d6+    Kd7
				8. Ke5    Kd8
				9. d7!    Kxd7
				10. Kf6   Kd8
				11. Ke6   Kc7
				12. Ke7   Kc8
				13. Kd6   Kb7
				14. Kd7   Kb8
				15. Kxc6 Kc8
				16. Kd6   Kd8
				17. c6    Kc8
				18. c7    Kb7
				19. Kd7 `,
			expectedFEN: "8/1kPK4/8/8/8/8/8/8 b - - 2 19",
			expectedErr: nil,
		},
		{
			fen: "8/p7/P7/8/8/4k3/4P3/4K3 w - - 0 1",
			s: `1. Kf1    Kd4
				2. Kf2    Kc5
				3. e4!    Kd4
				4. Kf3    Ke5
				5. Ke3    Ke6
				6. Kd4    Kd6
				7. e5+    Ke6
				8. Ke4    Ke7
				9. Kd5    Kd7
				10. e6+   Ke7
				11. Ke5   Ke8
				12. Kd6   Kd8
				13. Kc6   Ke7
				14. Kb7   Kxe6
				15. Kxa7 Kd7
				16. Kb7 `,
			expectedFEN: "8/1K1k4/P7/8/8/8/8/8 b - - 2 16",
			expectedErr: nil,
		},
		{
			fen: "3k4/1K3p2/8/4P3/8/6P1/8/8 w - - 0 1",
			s: `1. Kc6    Ke7
				2. Kd5    Kd7
				3. Kd4!  Ke7
				4. Ke3    Kd7
				5. Kf4    Ke6
				6. Ke4    Kd7
				7. Kf5    Ke7
				8. g4   Ke8!
				9. Kf6!   Kf8
				10. g5    Ke8
				11. Kg7! Ke7
				12. Kg8   Ke8
				13. e6!   fxe6
				14. g6    e5
				15. g7    e4
				16. Kh7   e3
				17. g8=Q+ `,
			expectedFEN: "4k1Q1/7K/8/8/8/4p3/8/8 b - - 0 17",
			expectedErr: nil,
		},
		{
			fen: "8/6p1/8/5P2/5P2/5K2/8/6k1 w - - 0 1",
			s: `1. f6!   gxf6
				2. f5  Kh2
				3. Kg4   Kg2
				4. Kh5   Kf3
				5. Kg6   Kg4
				6. Kxf6 Kh5
				7. Kg7`,
			expectedFEN: "8/6K1/8/5P1k/8/8/8/8 b - - 2 7",
			expectedErr: nil,
		},
		{
			fen: "8/6p1/8/K6P/7P/1k6/8/8 w - - 0 1",
			s: `1. Kb5    Kc3
				2. Kc5    Kd3
				3. Kd5    Ke3
				4. Ke5    Kf3
				5. Kf5    Kg3
				6. h6!    gxh6
				7. h5   Kh4
				8. Kg6    Kg4
				9. Kxh6 Kf5
				10. Kg7 `,
			expectedFEN: "8/6K1/8/5k1P/8/8/8/8 b - - 2 10",
			expectedErr: nil,
		},
		{
			fen: "8/6p1/8/4K3/7P/4k2P/8/8 w - - 0 1",
			s: `1. Ke6!  Kf4
				2. Kf7   Kg3
				3. h5  Kh4
				4. Kg6! Kxh3
				5. Kxg7 `,
			expectedFEN: "8/6K1/8/7P/8/7k/8/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "5k2/5p2/8/4P3/2K1P3/8/8/8 w - - 0 1",
			s: `1. e6!   fxe6
				2. e5!   Ke7
				3. Kc5   Kd7
				4. Kb6   Kd8
				5. Kc6   Ke7
				6. Kc7   Ke8
				7. Kd6   Kf7
				8. Kd7   Kf8
				9. Kxe6 Ke8
				10. Kf6  Kf8
				11. e6   Ke8
				12. e7 `,
			expectedFEN: "4k3/4P3/5K2/8/8/8/8/8 b - - 0 12",
			expectedErr: nil,
		},
		{
			fen: "4K3/8/2p5/8/P2k4/8/P7/8 w - - 0 1",
			s: `1. a5  Kc5
				2. a4  Kd6
				3. Kd8   c5
				4. a6  Kc6
				5. a5  c4
				6. Kc8  c3
				7. a7  c2
				8. a8=Q+ Kd6
				9. Qh1`,
			expectedFEN: "2K5/8/3k4/P7/8/8/2p5/7Q b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/1k6/4p1P1/6P1/6K1 w - - 0 1",
			s: `1. Kf1  Kc3
				2. Ke2  Kd4
				3. g4 Ke4
				4. g3 Kd4
				5. g5 Ke5
				6. g4 Ke6
				7. Kxe3 Kf7
				8. Ke4  Kg6
				9. Kf4  Kf7
				10. Kf5 Kg7
				11. g6  Kg8
				12. Kf6 Kf8
				13. g7+ Kg8
				14. g5  Kh7
				15. g8=Q+ Kxg8
				16. Kg6`,
			expectedFEN: "6k1/8/6K1/6P1/8/8/8/8 b - - 1 16",
			expectedErr: nil,
		},
		{
			fen: "8/6p1/8/8/7K/7P/4k2P/8 w - - 0 1",
			s: `1. Kg3! Ke3
				2. h4 Ke4
				3. Kg4  Ke5
				4. Kg5  Ke4
				5. h5 Kf3
				6. Kf5! Kg2
				7. Kg6`,
			expectedFEN: "8/6p1/6K1/7P/8/8/6kP/8 b - - 4 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/3k4/2p2p2/P1K2P2/8/8/8 w - - 0 1",
			s: `1. a5 Kc6
				2. a6 Kb6
				3. a7 Kxa7
				4. Kxc5 Kb7
				5. Kd6  Kc8
				6. Ke6  Kd8
				7. Kxf5 Ke7
				8. Kg6  Ke6
				9. f5+  Ke7
				10. Kg7`,
			expectedFEN: "8/4k1K1/8/5P2/8/8/8/8 b - - 2 10",
			expectedErr: nil,
		},
		{
			fen: "8/p5kp/8/1P6/8/1K6/P7/8 w - - 0 1",
			s: `1. a4 h5
				2. a5 h4
				3. b6 axb6
				4. a6!  h3
				5. a7 h2
				6. a8=Q`,
			expectedFEN: "Q7/6k1/1p6/8/8/1K6/7p/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/2k5/2Pp3p/1P6/8/5K2/8/8 w - - 0 1",
			s: `1. Kf4  Kb6
				2. Kf5  Kc7
				3. Kf6  Kb6
				4. Ke6  Kc7
				5. Kd5  h5
				6. b6+  Kxb6
				7. Kxd6 h4
				8. c7 Kb7
				9. Kd7`,
			expectedFEN: "8/1kPK4/8/8/7p/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/7k/4K2p/6p1/6P1/7P/8/8 w - - 0 1",
			s: `1. Kf7! h5
				2. h4!  Kh6
				3. Kf6! hxg4
				4. hxg5+  Kh5
				5. g6 g3
				6. g7 Kg4
				7. g8=Q+ Kf3
				8. Kg5  g2
				9. Kh4  Kf2
				10. Qg3+  Kf1
				11. Qf3+  Kg1
				12. Kh3 Kh1
				13. Qxg2#`,
			expectedFEN: "8/8/8/8/8/7K/6Q1/7k b - - 0 13",
			expectedErr: nil,
		},
		{
			fen: "7k/6p1/6P1/8/8/p5K1/P7/8 w - - 0 1",
			s: `1. Kf4! Kg8
				2. Ke5  Kf8
				3. Kd6! Ke8
				4. Ke6  Kf8
				5. Kd7  Kg8
				6. Ke7  Kh8
				7. Kd6  Kg8
				8. Kc5  Kf8
				9. Kb4  Ke7
				10. Kxa3  Kd6
				11. Kb4 Kc6
				12. a4  Kb6
				13. Kc4 Kc6
				14. a5  Kc7
				15. Kd5`,
			expectedFEN: "8/2k3p1/6P1/P2K4/8/8/8/8 b - - 2 15",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/7p/5Ppk/4KP2/8/8 w - - 0 1",
			s: `1. Kf2  Kh3
				2. Kg1  Kh4
				3. Kg2  g3
				4. Kg1! g2
				5. Kh2! g1=Q+
				6. Kxg1 Kg3
				7. f5`,
			expectedFEN: "8/8/8/5P1p/8/5Pk1/8/6K1 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/p2K1p2/3P1Pk1/8/8/8 w - - 0 1",
			s: `1. Kc4! Kxf4
				2. d5 Ke5
				3. Kc5  f4
				4. d6 f3
				5. d7 f2
				6. d8=Q f1=Q
				7. Qe8+ Kf5
				8. Qf8+`,
			expectedFEN: "5Q2/8/8/p1K2k2/8/8/8/5q2 b - - 3 8",
			expectedErr: nil,
		},
		{
			fen: "8/5k2/5p2/p7/2P1K1P1/8/8/8 w - - 0 1",
			s: `1. Kd4! Ke6
				2. Kc5  Ke5
				3. Kb5  Kd4
				4. g5!  fxg5
				5. c5 a4
				6. c6 a3
				7. c7 a2
				8. c8=Q a1=Q
				9. Qh8+`,
			expectedFEN: "7Q/8/8/1K4p1/3k4/8/8/q7 b - - 1 9",
			expectedErr: nil,
		},
		{
			fen: "2K5/8/5p2/3k2p1/P5P1/8/8/8 w - - 0 1",
			s: `1. a5!  Kc6
				2. Kb8! Kb5
				3. Kb7  Kxa5
				4. Kc6  f5
				5. gxf5 g4
				6. f6 g3
				7. f7 g2
				8. f8=Q g1=Q
				9. Qa3#!`,
			expectedFEN: "8/8/2K5/k7/8/Q7/8/6q1 b - - 1 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/7p/p7/k1P2K2/8/P7/8 w - - 0 1",
			s: `1. a3!  h5
				2. Kg3  h4+
				3. Kh3! Kxa3
				4. c5 a4
				5. c6 Kb2
				6. c7 a3
				7. c8=Q a2
				8. Qb7+ Kc2
				9. Qc6+ Kb2
				10. Qb5+  Kc2
				11. Qa4+  Kb2
				12. Qb4+  Kc2
				13. Qa3 Kb1
				14. Qb3+  Ka1
				15. Kg4!  h3
				16. Qc2 h2
				17. Qc1#`,
			expectedFEN: "8/8/8/8/6K1/8/p6p/k1Q5 b - - 1 17",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/KP6/1p6/k4p2/5P2/8 w - - 0 1",
			s: `1. b6 b3
				2. b7 b2
				3. b8=R!  Ka2
				4. Ka4  b1=Q
				5. Rxb1 Kxb1
				6. Kb3`,
			expectedFEN: "8/8/8/8/8/1K3p2/5P2/1k6 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/4k3/1p5p/1P5P/4K3/8/8 w - - 0 1",
			s: `1. Ke4  Kd6
				2. Kd4  Ke6
				3. Kc5  Kf5
				4. Kxb5 Kg4
				5. Kc4  Kxh4
				6. b5 Kg3
				7. b6 h4
				8. b7 h3
				9. b8=Q+ Kg2
				10. Qg8+  Kf2
				11. Qd5 Kg1
				12. Qg5+  Kf2
				13. Qh4+  Kg2
				14. Qg4+  Kh2
				15. Kd4 Kh1
				16. Qxh3+ `,
			expectedFEN: "8/8/8/8/3K4/7Q/8/7k b - - 0 16",
			expectedErr: nil,
		},
		{
			fen: "4K3/8/8/1p5p/1P5P/8/8/4k3 w - - 0 1",
			s: `1. Ke7! Ke2
				2. Ke6! Ke3
				3. Ke5! Ke2
				4. Ke4! Ke1
				5. Ke3! Kf1
				6. Kf4  Ke2
				7. Kg5  Kd3
				8. Kxh5 Kc4
				9. Kg4  Kxb4
				10. h5  Ka3
				11. h6   b4
				12. h7  b3
				13. h8=Q  b2
				14. Qc3+  Ka2
				15. Qc2 Ka1
				16. Qa4+  Kb1
				17. Kf3 Kc1
				18. Qc4+  Kd2
				19. Qb3 Kc1
				20. Qc3+  Kb1
				21. Ke3 Ka2
				22. Qc2 Ka1
				23. Qa4+  Kb1
				24. Kd3 Kc1
				25. Qc2#`,
			expectedFEN: "8/8/8/8/8/3K4/1pQ5/2k5 b - - 23 25",
			expectedErr: nil,
		},
		{
			fen: "8/4p2p/8/1K1k4/8/5P2/P7/8 w - - 0 1",
			s: `1. a4 Kd6
				2. Kb6  Kd7
				3. Kb7  h5
				4. a5 h4
				5. a6 h3
				6. a7 h2
				7. a8=Q h1=Q
				8. Qc8+ Kd6
				9. Qc6+ Ke5
				10. f4`,
			expectedFEN: "8/1K2p3/2Q5/4k3/5P2/8/8/7q b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "8/8/5p2/6p1/P2k2P1/8/8/K7 w - - 0 1",
			s: `1. a5 Kc5
				2. Kb2! f5
				3. gxf5 g4
				4. f6!  Kd6
				5. a6 g3
				6. f7 Ke7
				7. a7 g2
				8. f8=Q+   Kxf8
				9. a8=Q+ Ke7
				10. Qxg2`,
			expectedFEN: "8/4k3/8/8/8/8/1K4Q1/8 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "8/p2p4/8/8/8/k7/5P1P/7K w - - 0 1",
			s: `1. f4 Kb4
				2. h4 d5
				3. f5!  Kc5
				4. h5 d4
				5. f6 Kd6
				6. h6 d3
				7. f7!  Ke7
				8. h7 d2
				9. f8=Q+   Kxf8
				10. h8=Q+ Ke7
				11. Qd4`,
			expectedFEN: "8/p3k3/8/8/3Q4/8/3p4/7K b - - 2 11",
			expectedErr: nil,
		},
		{
			fen: "8/8/2k2ppp/8/P2K2PP/8/8/8 w - - 0 1",
			s: `1. a5 Kb5
				2. Kd5  Kxa5
				3. Ke6  f5
				4. gxf5 gxf5
				5. Kxf5 Kb6
				6. Kg6  Kc7
				7. Kxh6 Kd7
				8. h5 Ke7
				9. Kg7`,
			expectedFEN: "8/4k1K1/8/7P/8/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/p6p/4k3/8/4K3/8/P5PP/8 w - - 0 1",
			s: `1. g4 a5
				2. a4 Kf6
				3. h4 Ke6
				4. g5 Kf7
				5. Kf5  Kg7
				6. h5 Kf7
				7. Ke5  Ke7
				8. g6 hxg6
				9. hxg6 Ke8
				10. Kd5 Kf8
				11. Kc5`,
			expectedFEN: "5k2/8/6P1/p1K5/P7/8/8/8 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/1p3k1p/1P5P/5KP1/8/8 w - - 0 1",
			s: `1. g4+  hxg4+
				2. Kg3  Kg6
				3. Kxg4 Kh6
				4. Kf5  Kh5
				5. Ke5  Kxh4
				6. Kd5  Kg5
				7. Kc5  Kf6
				8. Kxb5 Ke7
				9. Kc6  Kd8
				10. Kb7 Kd7
				11. b5`,
			expectedFEN: "8/1K1k4/8/1P6/8/8/8/8 b - - 0 11",
			expectedErr: nil,
		},
		{
			fen: "8/3pp2p/8/1kPP4/5P2/8/1K6/8 w - - 0 1",
			s: `1. c6!  Kb6
				2. d6!  exd6
				3. f5 Kc7
				4. f6 Kd8
				5. c7+  Kxc7
				6. f7`,
			expectedFEN: "8/2kp1P1p/3p4/8/8/8/1K6/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/4k1pp/8/4KPPP/8/8/8/8 w - - 0 1",
			s: `1. g6 h6
				2. Kd5  Kf6
				3. Ke4  Kg5
				4. Ke5  Kxh5
				5. Ke6  Kg5
				6. f6!  gxf6
				7. g7`,
			expectedFEN: "8/6P1/4Kp1p/6k1/8/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/p7/Pp6/1P6/8/K2k4/P7/8 w - - 0 1",
			s: `1. Kb4  Kd4
				2. a4 Kd5
				3. a5 bxa5+
				4. Ka4! Kc5
				5. Kxa5 Kd6
				6. b6 axb6+
				7. Kxb6`,
			expectedFEN: "8/8/PK1k4/8/8/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "k7/P7/1P6/4p3/4Pp2/5K2/8/8 w - - 0 1",
			s: `1. Ke2  Kb7
				2. Kd3  Ka8
				3. Kc4  Kb7
				4. Kc5  f3
				5. Kd6  f2
				6. a8=Q+ Kxa8
				7. Kc7  f1=Q
				8. b7+  Ka7
				9. b8=Q+ Ka6
				10. Qb6#`,
			expectedFEN: "8/2K5/kQ6/4p3/4P3/8/8/5q2 b - - 2 10",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/2p5/2Pp4/3K2Pk/7P/8 w - - 0 1",
			s: `1. Ke4  Kg4
				2. h4 Kh5
				3. Kf4  Kh6
				4. g4 Kg6
				5. h5+  Kh6
				6. Ke4  Kg5
				7. Kf3  Kh6
				8. Kf4  Kh7
				9. g5 Kg7
				10. g6  Kh6
				11. Kg4 Kg7
				12. Kg5!  d3
				13. h6+ Kg8
				14. Kf6 d2
				15. h7+ Kh8
				16. Kf7 d1=Q
				17. g7+ Kxh7
				18. g8=Q+ Kh6
				19. Qg6#`,
			expectedFEN: "8/5K2/6Qk/2p5/2P5/8/8/3q4 b - - 2 19",
			expectedErr: nil,
		},
		{
			fen: "K7/5p2/k7/6p1/P5P1/8/P7/8 w - - 0 1",
			s: `1. a3 f5
				2. gxf5 g4
				3. f6 g3
				4. f7 g2
				5. f8=Q g1=Q
				6. Qc8+ Ka5
				7. Qc3+ Kb6
				8. a5+  Kb5
				9. Qb4+ Ka6
				10. Qb7+  Kxa5
				11. Qb4+  Ka6
				12. Qa4+  Kb6
				13. Qa7+`,
			expectedFEN: "K7/Q7/1k6/8/8/P7/8/6q1 b - - 5 13",
			expectedErr: nil,
		},
		{
			fen: "8/1p6/p3k3/8/6p1/PP6/4KPP1/8 w - - 0 1",
			s: `1. f3 gxf3+
				2. Kxf3 Kf5
				3. g4+  Kg5
				4. b4 b6
				5. a4 Kg6
				6. Kf4  Kf6
				7. g5+  Kg6
				8. b5 a5
				9. Kg4  Kg7
				10. Kf5 Kf7
				11. Ke5 Ke7
				12. g6  Ke8
				13. Kd6 Kf8
				14. Kc6`,
			expectedFEN: "5k2/8/1pK3P1/pP6/P7/8/8/8 b - - 4 14",
			expectedErr: nil,
		},
		{
			fen: "8/1pp5/p2p4/P2P4/1PP5/7k/8/7K w - - 0 1",
			s: `1. c5 Kg4
				2. b5!  dxc5
				3. b6 cxb6
				4. d6`,
			expectedFEN: "8/1p6/pp1P4/P1p5/6k1/8/8/7K b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "3k4/2p5/1pKp4/p2P4/2P5/P7/1P6/8 w - - 0 1",
			s: `1. c5!  bxc5
				2. Kb5  Kd7
				3. a4 Kc8
				4. Kxa5 Kb7
				5. Kb5  Ka7
				6. Kc6  Kb8
				7. a5 Kc8
				8. a6 Kb8
				9. a7+  Kxa7
				10. Kxc7  Ka6
				11. Kxd6  Kb5
				12. b3  Kb4
				13. Kc6`,
			expectedFEN: "8/8/2K5/2pP4/1k6/1P6/8/8 b - - 2 13",
			expectedErr: nil,
		},
		{
			fen: "7k/8/5PpK/Pp1P2pp/3P4/8/5p2/8 w - - 0 1",
			s: `1. a6 f1=Q
				2. a7 Qa1
				3. f7 Qa3
				4. d6 Qf3
				5. d5!  Qxf7
				6. a8=Q+ Qg8
				7. Qa1+ Qg7+
				8. Qxg7#`,
			expectedFEN: "7k/6Q1/3P2pK/1p1P2pp/8/8/8/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "k7/2p1pp2/2P3p1/4P1P1/5P2/p7/Kp3P2/8 w - - 0 1",
			s: `1. f5 e6
				2. fxe6 fxe6
				3. f4 Kb8
				4. f5!  exf5
				5. e6 Kc8
				6. e7`,
			expectedFEN: "2k5/2p1P3/2P3p1/5pP1/8/p7/Kp6/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/2p5/1pPp4/1P1Pp3/4Pp1k/5P2/5KP1/8 w - - 0 1",
			s: `1. g3+  fxg3+
				2. Kg2  Kh5
				3. Kxg3 Kg5
				4. f4+  exf4+
				5. Kf3  Kg6
				6. Kxf4 Kf6
				7. e5+  dxe5+
				8. Ke4  Kf7
				9. Kxe5 Ke7
				10. d6+ cxd6+
				11. Kd5 Ke8
				12. Kxd6 Kd8
				13. c7+ Kc8
				14. Ke6 Kxc7
				15. Ke7 Kc8
				16. Kd6 Kb7
				17. Kd7 Kb8
				18. Kc6 Ka7
				19. Kc7 Ka8
				20. Kxb6 Kb8
				21. Ka6 Ka8
				22. b6  Kb8
				23. b7  Kc7
				24. Ka7`,
			expectedFEN: "8/KPk5/8/8/8/8/8/8 b - - 2 24",
			expectedErr: nil,
		},
		{
			fen: "8/2pp2pp/8/2PP1P2/1p5k/8/PP4p1/6K1 w - - 0 1",
			s: `1. f6!  gxf6
				2. Kxg2 Kg5
				3. a4 bxa3e.p.
				4. bxa3 Kf5
				5. a4 Ke5
				6. d6!  cxd6
				7. c6!  dxc6
				8. a5`,
			expectedFEN: "8/7p/2pp1p2/P3k3/8/8/6K1/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "3K4/8/8/1n5P/5k2/8/8/8 w - - 0 1",
			s: `1. h6 Nd6
				2. h7 Nf7+
				3. Ke7  Nh8
				4. Kf6!`,
			expectedFEN: "7n/7P/5K2/8/5k2/8/8/8 b - - 4 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/4P3/2k5/6K1/8/2n5 w - - 0 1",
			s: `1. e6!  Ne2+
				2. Kh2!!`,
			expectedFEN: "8/8/4P3/8/2k5/8/4n2K/8 b - - 2 2",
			expectedErr: nil,
		},
		{
			fen: "8/8/1P6/8/2K2kn1/8/6P1/8 w - - 0 1",
			s: `1. Kd5  Ne5
				2. g3+  Kf5
				3. g4+  Kf6
				4. g5+  Kf5
				5. g6 Kf6
				6. g7`,
			expectedFEN: "8/6P1/1P3k2/3Kn3/8/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/5n2/8/3P2PK/3k2p1/8/8/8 w - - 0 1",
			s: `1. g6 Nd6
				2. g7 Ne8
				3. g8=Q Nf6+
				4. Kg5! Nxg8
				5. d6!  g3
				6. d7 g2
				7. d8=Q`,
			expectedFEN: "3Q2n1/8/8/6K1/3k4/8/6p1/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/7K/7P/1P6/2n4k/4n3/1P6/8 w - - 0 1",
			s: `1. Kg6  Ne5+
				2. Kf6  Ng4+
				3. Ke6  Nxh6
				4. b6 Nf7
				5. Kxf7 Nc4
				6. b7 Nd6+
				7. Ke7  Nxb7
				8. b4 Kg5
				9. Kd7  Kf6
				10. Kc7`,
			expectedFEN: "8/1nK5/5k2/8/1P6/8/8/8 b - - 4 10",
			expectedErr: nil,
		},
		{
			fen: "8/6b1/5k2/8/P3K1P1/8/8/8 w - - 0 1",
			s: `1. a5 Bf8
				2. Kd5  Bh6
				3. g5+! Bxg5
				4. Ke4! Bh4
				5. Kf3`,
			expectedFEN: "8/8/5k2/P7/7b/5K2/8/8 b - - 3 5",
			expectedErr: nil,
		},
		{
			fen: "8/8/1P6/8/2K2k2/8/1b4P1/8 w - - 0 1",
			s: `1. Kd5! Be5
				2. g3+  Kf5
				3. g4+  Kf6
				4. g5+  Kf5
				5. g6 Kf6
				6. g7`,
			expectedFEN: "8/6P1/1P3k2/3Kb3/8/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/2b5/7P/1p6/1P4K1/8/8/5k2 w - - 0 1",
			s: `1. Kf5  Bb6
				2. Ke4  Bd8
				3. Ke5! Bg5
				4. h7 Bc1
				5. Kd5  Bb2
				6. Kc5  Kf2
				7. Kxb5 Kf3
				8. Kc6  Kf4
				9. b5 Kf5
				10. b6  Kg6
				11. b7  Be5
				12. b8=Q Bxb8
				13. h8=Q`,
			expectedFEN: "1b5Q/8/2K3k1/8/8/8/8/8 b - - 0 13",
			expectedErr: nil,
		},
		{
			fen: "6b1/8/2k4K/4P3/7P/8/8/8 w - - 0 1",
			s: `1. Kg7  Bb3
				2. h5 Kd7
				3. h6 Bc2
				4. Kf7! Bb3+
				5. e6+  Bxe6+
				6. Kf6  Bg8
				7. Kg7`,
			expectedFEN: "6b1/3k2K1/7P/8/8/8/8/8 b - - 3 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/5b2/5P2/2P1K1k1/8/8/8 w - - 0 1",
			s: `1. c5 Kg5
				2. c6 Bd8
				3. Ke5  Kh6
				4. f6!  Kh7
				5. Ke6  Kg8
				6. Kd7  Ba5
				7. Ke8`,
			expectedFEN: "4K1k1/8/2P2P2/b7/8/8/8/8 b - - 6 7",
			expectedErr: nil,
		},
		{
			fen: "4k3/8/2K5/8/PP6/8/8/4b3 w - - 0 1",
			s: `1. a5!  Kd8
				2. a6 Bf2
				3. Kb7  Kd7
				4. b5 Kd8
				5. a7 Bxa7
				6. Kxa7`,
			expectedFEN: "3k4/K7/8/1P6/8/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "K7/8/kPP5/1b6/1P6/8/8/8 w - - 0 1",
			s: `1. c7 Bc6+
				2. b7!  Bxb7+
				3. Kb8  Kb6
				4. b5!`,
			expectedFEN: "1K6/1bP5/1k6/1P6/8/8/8/8 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "8/5b2/7p/7P/5PP1/8/1K4k1/8 w - - 0 1",
			s: `1. f5 Kg3
				2. g5!  hxg5
				3. h6 Bg8
				4. f6 Kf4
				5. h7 Bxh7
				6. f7`,
			expectedFEN: "8/5P1b/8/6p1/5k2/8/1K6/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "4b1k1/4P1P1/5P1K/8/8/8/8/8 w - - 0 1",
			s: `1. Kg5  Kf7
				2. Kf5  Bd7+
				3. Ke5  Ba4
				4. Kd6  Bb5
				5. Kc7  Ba4
				6. Kd8  Bb5
				7. g8=Q+ Kxg8
				8. e8=Q+ Bxe8
				9. Kxe8`,
			expectedFEN: "4K1k1/8/5P2/8/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "6k1/3p4/P4P2/2P5/1K6/7b/8/8 w - - 0 1",
			s: `1. c6 dxc6
				2. Kc5  Bc8
				3. a7 Bb7
				4. Kd6  Kf7
				5. Kc7  Ba8
				6. Kb8`,
			expectedFEN: "bK6/P4k2/2p2P2/8/8/8/8/8 b - - 6 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/1KP5/3r4/8/8/8/k7 w - - 0 1",
			s: `1. c7 Rd6+
				2. Kb5  Rd5+
				3. Kb4  Rd4+
				4. Kb3  Rd3+
				5. Kc2  Rd4
				6. c8=R!  Ra4
				7. Kb3`,
			expectedFEN: "2R5/8/8/8/r7/1K6/8/k7 b - - 2 7",
			expectedErr: nil,
		},
		{
			fen: "6K1/1k4P1/6P1/8/8/8/r7/8 w - - 0 1",
			s: `1. Kf7  Rf2+
				2. Ke6  Re2+
				3. Kf5  Rf2+
				4. Ke4  Re2+
				5. Kf4! Re8
				6. Kg5  Kc7
				7. Kh6  Kd7
				8. Kh7  Ke7
				9. g8=Q`,
			expectedFEN: "4r1Q1/4k2K/6P1/8/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/pN6/8/2K5/k7 w - - 0 1",
			s: `1. Kc1  a3
				2. Nc2+ Ka2
				3. Nd4  Ka1
				4. Kc2  Ka2
				5. Ne2  Ka1
				6. Nc1  a2
				7. Nb3#`,
			expectedFEN: "8/8/8/8/8/1N6/p1K5/k7 b - - 1 7",
			expectedErr: nil,
		},
		{
			fen: "8/3N4/8/8/8/2P1k3/8/7K w - - 0 1",
			s: `1. Nc5  Kd2
				2. Na4! Kd3
				3. Kg2  Kc4
				4. Kf3  Kb3
				5. Ke4  Kxa4
				6. Kd5  Kb5
				7. c4+  Kb6
				8. Kd6  Kb7
				9. c5 Kc8
				10. Kc6!`,
			expectedFEN: "2k5/8/2K5/2P5/8/8/8/8 b - - 2 10",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/2K5/p1N5/Pk6/8 w - - 0 1",
			s: `1. Kd3  Ka1
				2. Na4! Kb1
				3. Kd2  Ka1
				4. Kc1  Kxa2
				5. Kc2  Ka1
				6. Nc5  Ka2
				7. Nd3  Ka1
				8. Nc1  a2
				9. Nb3#`,
			expectedFEN: "8/8/8/8/8/1N6/p1K5/k7 b - - 1 9",
			expectedErr: nil,
		},
		{
			fen: "8/5k2/7P/6K1/1N6/p7/8/8 w - - 0 1",
			s: `1. Na2! Kf8
				2. Kf6! Kg8
				3. Kg6  Kh8
				4. Nb4  Kg8
				5. h7+  Kh8
				6. Nc6  a2
				7. Nd8  a1=Q
				8. Nf7#`,
			expectedFEN: "7k/5N1P/6K1/8/8/8/8/q7 b - - 1 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/4N3/8/7p/3K1k2/7P/8 w - - 0 1",
			s: `1. h3 Kg3
				2. Ng5  Kf4
				3. Ne4  Kf3
				4. Kd4  Kf4
				5. Kd5  Kf5
				6. Nc3! Kf4
				7. Ne2+ Kf3
				8. Ng1+ Kg2
				9. Ke4  Kxg1
				10. Kf3 Kh2
				11. Kg4 Kg2
				12. Kxh4  Kf3
				13. Kg5 Ke4
				14. h4  Ke5
				15. h5  Ke6
				16. Kg6`,
			expectedFEN: "8/8/4k1K1/7P/8/8/8/8 b - - 2 16",
			expectedErr: nil,
		},
		{
			fen: "8/8/7p/8/6pp/4N1Pk/5K2/8 w - - 0 1",
			s: `1. Ng2! hxg3+
				2. Kg1  h5
				3. Kh1  h4
				4. Nf4#`,
			expectedFEN: "8/8/8/8/5Npp/6pk/8/7K b - - 1 4",
			expectedErr: nil,
		},
		{
			fen: "b3K3/8/3P4/4k1N1/8/8/8/8 w - - 0 1",
			s: `1. Kd7! Kd5
				2. Kc7  Bc6
				3. Ne4! Kxe4
				4. Kxc6`,
			expectedFEN: "8/8/2KP4/8/4k3/8/8/8 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "8/Kb6/1Pk5/8/7N/8/8/8 w - - 0 1",
			s: `1. Nf5  Ba8
				2. Nd4+ Kc5
				3. Ne6+ Kc6
				4. Nc7  Bb7
				5. Nd5!`,
			expectedFEN: "8/Kb6/1Pk5/3N4/8/8/8/8 b - - 9 5",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/1PK5/3N4/8/5kb1/8 w - - 0 1",
			s: `1. Nc6  Bf1
				2. b6 Ba6
				3. Kd6  Bb7
				4. Kc7  Ba8
				5. Nd8  Ke3
				6. Nb7  Kd4
				7. Kb8  Kd5
				8. Kxa8 Kc6
				9. Ka7`,
			expectedFEN: "8/KN6/1Pk5/8/8/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/P1K1k3/8/8/2N5/5b2/8/8 w - - 0 1",
			s: `1. Na5  Ba8
				2. Kc8! Kd6
				3. Kb8  Kd7
				4. Nb7  Kc6
				5. Kxa8 Kc7
				6. Nd6! Kxd6
				7. Kb7`,
			expectedFEN: "8/PK6/3k4/8/8/8/8/8 b - - 1 7",
			expectedErr: nil,
		},
		{
			fen: "8/3PK3/8/N1k5/5b2/8/8/8 w - - 0 1",
			s: `1. Nb7+ Kb4
				2. Nd6  Be3
				3. Nc8  Bf4
				4. Nb6! Kc5
				5. Nd5!`,
			expectedFEN: "8/3PK3/8/2kN4/5b2/8/8/8 b - - 9 5",
			expectedErr: nil,
		},
		{
			fen: "4k3/8/4N3/5n1P/8/8/8/K7 w - - 0 1",
			s: `1. Ng7+!   Nxg7
				2. h6   Kf7
				3. h7`,
			expectedFEN: "8/5knP/8/8/8/8/8/K7 b - - 0 3",
			expectedErr: nil,
		},
		{
			fen: "8/7k/8/2K5/1P2N3/n7/8/8 w - - 0 1",
			s: `1. Nd2  Kg7
				2. Nc4  Nb1
				3. Kd4! Kf7
				4. b5 Ke7
				5. b6 Kd7
				6. Kc5  Nc3
				7. Ne5+ Kc8
				8. Kc6  Ne2
				9. b7+  Kb8
				10. Kb6`,
			expectedFEN: "1k6/1P6/1K6/4N3/8/8/4n3/8 b - - 2 10",
			expectedErr: nil,
		},
		{
			fen: "K2n4/8/P7/8/2k5/8/2N5/8 w - - 0 1",
			s: `1. Ka7! Kb5
				2. Nb4! Ka5
				3. Kb8  Nc6+
				4. Kb7  Nd8+
				5. Kc7  Ne6+
				6. Kb8  Nc5
				7. a7 Nd7+
				8. Kc7  Nb6
				9. Kb7  Kb5
				10. Nd5`,
			expectedFEN: "8/PK6/1n6/1k1N4/8/8/8/8 b - - 6 10",
			expectedErr: nil,
		},
		{
			fen: "8/KP1n4/2k5/8/5N2/8/8/8 w - - 0 1",
			s: `1. Ne6  Kd5
				2. Nf8  Ne5
				3. Ka8  Nc6
				4. Nd7  Ke6
				5. Nb6  Kd6
				6. Nc8+ Kc7
				7. Na7  Nb8
				8. Nb5+ Kb6
				9. Kxb8`,
			expectedFEN: "1K6/1P6/1k6/1N6/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "1K6/P7/1nk5/8/4N3/8/8/8 w - - 0 1",
			s: `1. Nf6  Na8!
				2. Nd5! Kd7
				3. Kb7  Kd8
				4. Nb6  Nc7
				5. Kc6`,
			expectedFEN: "3k4/P1n5/1NK5/8/8/8/8/8 b - - 9 5",
			expectedErr: nil,
		},
		{
			fen: "k4n2/8/8/2N1P3/5K2/8/8/8 w - - 0 1",
			s: `1. Kg5! Ka7
				2. Kf5! Kb6
				3. Nd7+!  Nxd7
				4. e6 Nc5
				5. e7 Nb7
				6. Ke5`,
			expectedFEN: "8/1n2P3/1k6/4K3/8/8/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/5P1p/3N3k/8/8/3n3K w - - 0 1",
			s: `1. f6 Nf2+
				2. Kg2  Nd3
				3. f7 Nf4+
				4. Kh2  Ng6
				5. Nf3+ Kg4
				6. Ne5+!  Nxe5
				7. f8=Q`,
			expectedFEN: "5Q2/8/8/4n2p/6k1/8/7K/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/2K5/kp2p3/p4P2/3N4/8/3n4 w - - 0 1",
			s: `1. f5 Ne3
				2. f6 Nf5
				3. f7 Ne7+
				4. Kb7  Ng6
				5. Nxe5 Nf8
				6. Nc6#!`,
			expectedFEN: "5n2/1K3P2/2N5/kp6/p7/8/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "8/pK6/pP6/8/8/r6N/8/5k2 w - - 0 1",
			s: `1. bxa7 Rb3+
				2. Kc7  Rc3+
				3. Kd7  Rd3+
				4. Ke7  Re3+
				5. Kf7  Rf3+
				6. Kg7  Rg3+
				7. Ng5! Rxg5+
				8. Kf7  Rf5+
				9. Ke7  Re5+
				10. Kd7 Rd5+
				11. Kc7 Rc5+
				12. Kb7 Rb5+
				13. Kc6!`,
			expectedFEN: "8/P7/p1K5/1r6/8/8/8/5k2 b - - 11 13",
			expectedErr: nil,
		},
		{
			fen: "7k/4K1p1/6Pp/5N1P/8/8/8/8 w - - 0 1",
			s: `1. Nd6  Kg8
				2. Ne8  Kh8
				3. Nf6! gxf6
				4. Kf7  f5
				5. g7+  Kh7
				6. g8=Q#`,
			expectedFEN: "6Q1/5K1k/7p/5p1P/8/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "1K6/5p2/k3P3/p7/8/3p4/P1N5/8 w - - 0 1",
			s: `1. Nb4+!   axb4
				2. e7  d2
				3. e8=Q  d1=Q
				4. Qc6+    Ka5
				5. Kb7!  Qd3
				6. Qb6+  Ka4
				7. Qa7+  Kb5
				8. Qa6+`,
			expectedFEN: "8/1K3p2/Q7/1k6/1p6/3q4/P7/8 b - - 9 8",
			expectedErr: nil,
		},
		{
			fen: "2n5/8/1PN1Pk2/8/4K3/8/8/8 w - - 0 1",
			s: `1. b7 Nd6+
				2. Kd4! Nxb7
				3. Kd5  Kg7
				4. Nd8! Nxd8
				5. e7`,
			expectedFEN: "3n4/4P1k1/8/3K4/8/8/8/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "5n2/8/3K4/2N2k1P/8/6P1/8/8 w - - 0 1",
			s: `1. Ke7  Nh7
				2. Kf7  Kg5
				3. Ne4+ Kxh5
				4. Kg7  Ng5
				5. Nf6#`,
			expectedFEN: "8/6K1/5N2/6nk/8/6P1/8/8 b - - 3 5",
			expectedErr: nil,
		},
		{
			fen: "5n2/8/7P/4k1P1/p4N2/8/5K2/8 w - - 0 1",
			s: `1. g6 Kf6
				2. g7 Kf7
				3. Nd5  Ne6
				4. Ne7  Nxg7
				5. h7!`,
			expectedFEN: "8/4NknP/8/8/p7/8/5K2/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "8/8/2k4b/P7/8/8/2N2PKp/8 w - - 0 1",
			s: `1. Nd4+ Kb7
				2. Kxh2 Ka6
				3. Nb3    Bf4+
				4. Kh3    Kb5
				5. Kg4    Bb8
				6. f4   Kb4
				7. f5   Kxb3
				8. f6   Kb4
				9. f7   Bd6
				10. a6`,
			expectedFEN: "8/5P2/P2b4/8/1k4K1/8/8/8 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "4k3/3p1N2/8/P1P4b/8/8/7K/8 w - - 0 1",
			s: `1. c6 dxc6
				2. a6 Bf3
				3. Ng5  Bd5
				4. Ne6  c5
				5. Nc7+ Kd7
				6. Nxd5 Kc8
				7. Nb6+ Kb8
				8. Nd7+ Ka7
				9. Nxc5`,
			expectedFEN: "8/k7/P7/2N5/8/8/7K/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/1pN2p2/1P6/8/2P5/5k2/b7/6K1 w - - 0 1",
			s: `1. c5 Bb1
				2. Ne6! fxe6
				3. c6 Be4
				4. c7`,
			expectedFEN: "8/1pP5/1P2p3/8/4b3/5k2/8/6K1 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/P4PK1/8/6N1/8/8/r6k w - - 0 1",
			s: `1. f7 Rxa6+
				2. Nf6  Ra8
				3. Ne8  Ra6+
				4. Kg5  Ra5+
				5. Kg4  Ra4+
				6. Kg3  Ra3+
				7. Kf2  Ra2+
				8. Ke3  Ra3+
				9. Ke4  Ra4+
				10. Ke5 Ra5+
				11. Ke6 Ra6+
				12. Kd7 Ra7+
				13. Nc7`,
			expectedFEN: "8/r1NK1P2/8/8/8/8/8/7k b - - 23 13",
			expectedErr: nil,
		},
		{
			fen: "8/7p/1P5K/3N1r2/8/3k4/6P1/8 w - - 0 1",
			s: `1. b7 Rf8
				2. Nb4+ Ke4
				3. Nc6  Kf4
				4. g4!  Kxg4
				5. Kg7  Re8
				6. Kf7  Rh8
				7. Ke7! Rg8
				8. Nd8  Rg7+
				9. Nf7  Rg8
				10. Nh6+`,
			expectedFEN: "6r1/1P2K2p/7N/8/6k1/8/8/8 b - - 11 10",
			expectedErr: nil,
		},
		{
			fen: "1k2N2b/2p4P/2p2p2/2P2P2/8/5K2/8/8 w - - 0 1",
			s: `1. Kg4     Kc8
				2. Kh5     Kd8
				3. Ng7!    Bxg7
				4. h8=Q+ Bxh8
				5. Kg6     Ke7
				6. Kh7     Kf7
				7. Kxh8    Kf8
				8. Kh7     Kf7
				9. Kh6     Kf8
				10. Kg6    Ke7
				11. Kg7    Ke8
				12. Kxf6     Kd8
				13. Ke6    Ke8
				14. f6     Kf8
				15. f7     Kg7
				16. Ke7`,
			expectedFEN: "8/2p1KPk1/2p5/2P5/8/8/8/8 b - - 2 16",
			expectedErr: nil,
		},
		{
			fen: "2K5/p7/k1B5/p7/p7/8/8/8 w - - 0 1",
			s: `1. Kc7  a3
				2. Ba4  a2
				3. Kc6  a1=Q
				4. Bb5#`,
			expectedFEN: "8/p7/k1K5/pB6/8/8/8/q7 b - - 1 4",
			expectedErr: nil,
		},
		{
			fen: "8/4k3/8/7P/4B3/5K2/8/8 w - - 0 1",
			s: `1. h6 Kf7
				2. Bh7! Kf6
				3. Kf4  Kf7
				4. Kf5  Kf8
				5. Kf6  Ke8
				6. Bf5  Kf8
				7. h7 Ke8
				8. h8=Q#`,
			expectedFEN: "4k2Q/8/5K2/5B2/8/8/8/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "4k3/8/8/7P/8/4K2B/8/8 w - - 0 1",
			s: `1. Be6! Ke7
				2. h6 Kf6
				3. Bf5! Kf7
				4. Bh7! Kf6
				5. Kf4  Kf7
				6. Kf5  Kf8
				7. Kf6  Ke8
				8. Bf5  Kf8
				9. h7 Ke8
				10. h8=Q#`,
			expectedFEN: "4k2Q/8/5K2/5B2/8/8/8/8 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "6K1/8/7p/8/k7/2B5/1P6/8 w - - 0 1",
			s: `1. Kf7  h5
				2. Ke6  h4
				3. Kd5  h3
				4. Kc4! h2
				5. Bb4  h1=Q
				6. b3#!`,
			expectedFEN: "8/8/8/8/kBK5/1P6/8/7q b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/3k4/P7/K7/1B6/1p6/8/8 w - - 0 1",
			s: `1. a7 b2
				2. a8=Q b1=Q
				3. Qb7+ Ke6
				4. Qe7+ Kd5
				5. Qd6+ Kc4
				6. Qc5+ Kb3
				7. Qc3+ Ka2
				8. Qa3#`,
			expectedFEN: "8/8/8/K7/1B6/Q7/k7/1q6 b - - 11 8",
			expectedErr: nil,
		},
		{
			fen: "6k1/3p4/8/8/8/B7/P7/7K w - - 0 1",
			s: `1. Bb4! Kf7
				2. a4 Ke6
				3. a5 Kd5
				4. a6 Kc6
				5. Ba5! d5
				6. Kg2  d4
				7. Kf3  d3
				8. Ke3  Kd7
				9. a7`,
			expectedFEN: "8/P2k4/8/B7/8/3pK3/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/k1K5/8/pB6/P7/8/8 w - - 0 1",
			s: `1. Bc5  Ka5
				2. Kb7  Kb5
				3. Bb6! Kc4
				4. Kc6  Kb3
				5. Bc5  Kc4
				6. Bd6  Kd4
				7. Kb5  Kd5
				8. Bh2  Ke6
				9. Kxa4 Kd7
				10. Kb5 Kc8
				11. Kc6`,
			expectedFEN: "2k5/8/2K5/8/8/P7/7B/8 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "8/2P5/1K6/3k4/8/p7/6p1/B7 w - - 0 1",
			s: `1. Bd4! g1=Q
				2. Bxg1 a2
				3. Bd4! Kxd4
				4. c8=Q a1=Q
				5. Qh8+`,
			expectedFEN: "7Q/8/1K6/8/3k4/8/8/q7 b - - 1 5",
			expectedErr: nil,
		},
		{
			fen: "8/3p4/8/7p/1KPk4/8/B7/8 w - - 0 1",
			s: `1. c5!  h4
				2. Be6! dxe6
				3. c6 h3
				4. c7 h2
				5. c8=Q h1=Q
				6. Qc3+ Kd5
				7. Qc5+ Ke4
				8. Qc6+`,
			expectedFEN: "8/8/2Q1p3/8/1K2k3/8/8/7q b - - 5 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/1pk1K3/2p5/8/1P5B/8 w - - 0 1",
			s: `1. Bg1+ Kb4
				2. Bd4  Kb3
				3. Bc3  b4
				4. Kd4! bxc3
				5. bxc3 Ka4
				6. Kxc4`,
			expectedFEN: "8/8/8/8/k1K5/2P5/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/k7/8/1pK5/p4B2/P7/8/8 w - - 0 1",
			s: `1. Kc6  Ka8
				2. Kb6! b4
				3. axb4 a3
				4. b5 a2
				5. Be5  a1=Q
				6. Bxa1 Kb8
				7. Be5+ Ka8
				8. Kc7  Ka7
				9. b6+`,
			expectedFEN: "8/k1K5/1P6/4B3/8/8/8/8 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/3p4/4p2B/4P3/8/4K1kp/8 w - - 0 1",
			s: `1. Bf3+  Kg1
				2. Bh1!  Kxh1
				3. Kf1   d5
				4. exd5  e4
				5. d6  e3
				6. d7  e2+
				7. Kxe2  Kg1
				8. d8=Q  h1=Q
				9. Qd4+  Kh2
				10. Qh4+ Kg2
				11. Qg4+ Kh2
				12. Kf2`,
			expectedFEN: "8/8/8/8/6Q1/8/5K1k/7q b - - 7 12",
			expectedErr: nil,
		},
		{
			fen: "8/3K4/3P1p2/p4p2/k3B3/p7/8/8 w - - 0 1",
			s: `1. Bb1!   f4
				2. Kc6    f3
				3. Kc5!   Kb3
				4. d7   f2
				5. d8=Q   f1=Q
				6. Qd5+   Kc3
				7. Qd4+   Kb3
				8. Qa4+! Kxa4
				9. Bc2#`,
			expectedFEN: "8/8/5p2/p1K5/k7/p7/2B5/5q2 b - - 1 9",
			expectedErr: nil,
		},
		{
			fen: "3n4/2K5/1P6/k7/2B5/8/8/8 w - - 0 1",
			s: `1. Bf7  Kb5
				2. Be8+ Ka5
				3. Bd7  Ka6
				4. Bh3  Ka5
				5. Bg4  Kb5
				6. Be2+ Ka5
				7. Bc4`,
			expectedFEN: "3n4/2K5/1P6/k7/2B5/8/8/8 b - - 13 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/3B4/8/7K/5k1P/8/3n4 w - - 0 1",
			s: `1. Kg5  Nf2
				2. h4!  Ne4+
				3. Kg6  Nxd6
				4. h5 Nc4
				5. h6 Ne5+
				6. Kg7`,
			expectedFEN: "8/6K1/7P/4n3/8/5k2/8/8 b - - 2 6",
			expectedErr: nil,
		},
		{
			fen: "6K1/5P2/6k1/2b5/8/8/3B4/8 w - - 0 1",
			s: `1. Bc3  Ba3
				2. Bg7  Bb4
				3. Bf8  Bd2
				4. Bc5  Bh6
				5. Bd4  Kf5
				6. Bg7`,
			expectedFEN: "6K1/5PB1/7b/5k2/8/8/8/8 b - - 11 6",
			expectedErr: nil,
		},
		{
			fen: "3B4/5K2/4P3/2bk4/8/8/8/8 w - - 0 1",
			s: `1. Be7  Be3
				2. Bf8  Bg5
				3. Bg7  Kd6
				4. Bf6  Bd2
				5. e7`,
			expectedFEN: "8/4PK2/3k1B2/8/8/8/3b4/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "8/4B1b1/6P1/5K1k/8/8/8/8 w - - 0 1",
			s: `1. Bg5  Bf8
				2. Kf6! Be7+
				3. Kf7  Bf8
				4. Be3`,
			expectedFEN: "5b2/5K2/6P1/7k/8/4B3/8/8 b - - 7 4",
			expectedErr: nil,
		},
		{
			fen: "5k2/8/6PK/8/8/2b3B1/8/8 w - - 0 1",
			s: `1. Kh7! Bb2
				2. Bf4  Bd4
				3. Bh6+ Ke8
				4. Bg7  Bc5
				5. Be5  Bf8
				6. Bd6! Bxd6
				7. g7`,
			expectedFEN: "4k3/6PK/3b4/8/8/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/6K1/4B2P/6k1/4b3/8/8 w - - 0 1",
			s: `1. Bg7  Bg5
				2. Bh6  Bf6
				3. Be3  Bg7
				4. Bg5  Bf8
				5. Bf6  Kf4
				6. Bg7`,
			expectedFEN: "5b2/6B1/6K1/7P/5k2/8/8/8 b - - 11 6",
			expectedErr: nil,
		},
		{
			fen: "2KB4/1P6/2k5/8/8/8/7b/8 w - - 0 1",
			s: `1. Bh4  Kb6
				2. Bf2+ Ka6
				3. Bc5! Bg3
				4. Be7  Kb6
				5. Bd8+ Kc6
				6. Bh4! Bh2
				7. Bf2  Kb5
				8. Ba7  Ka6
				9. Bb8  Bg1
				10. Bf4 Ba7
				11. Be3`,
			expectedFEN: "2K5/bP6/k7/8/8/4B3/8/8 b - - 21 11",
			expectedErr: nil,
		},
		{
			fen: "2k5/P1p5/1n2K3/8/4B3/8/8/8 w - - 0 1",
			s: `1. Bc6! Kd8
				2. Kf5! Ke7
				3. Ke5! Kf7
				4. Kd4  Ke6
				5. Kc5  Ke7
				6. Bf3  Kd7
				7. Bg4+ Ke7
				8. Kc6  Kd8
				9. Kb7`,
			expectedFEN: "3k4/PKp5/1n6/8/6B1/8/8/8 b - - 17 9",
			expectedErr: nil,
		},
		{
			fen: "2B5/8/5k1K/4p3/1b4P1/8/8/8 w - - 0 1",
			s: `1. g5+  Ke7
				2. g6 Ke8
				3. Bd7+ Kxd7
				4. g7`,
			expectedFEN: "8/3k2P1/7K/4p3/1b6/8/8/8 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "1B6/8/7P/4p3/3b3k/8/8/2K5 w - - 0 1",
			s: `1. Ba7! Ba1
				2. Kb1  Bc3
				3. Kc2  Ba1
				4. Bd4! Bxd4
				5. Kd3  Ba1
				6. Ke4`,
			expectedFEN: "8/8/7P/4p3/4K2k/8/8/b7 b - - 3 6",
			expectedErr: nil,
		},
		{
			fen: "6b1/5p2/8/2B2K2/7k/8/6P1/8 w - - 0 1",
			s: `1. Bf2+ Kh5
				2. g4+  Kh6
				3. Kf6  Kh7
				4. g5 Kh8
				5. Bd4  Kh7
				6. Ba1  Kh8
				7. g6!  fxg6
				8. Kxg6#`,
			expectedFEN: "6bk/8/6K1/8/8/8/8/B7 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "b7/6K1/2p5/8/4Pp2/8/2k3B1/8 w - - 0 1",
			s: `1. e5 Bb7
				2. e6 Bc8
				3. e7 Bd7
				4. Bh3! Be8
				5. Kf8  Bh5
				6. Bg4  Bg6
				7. Bf5+ Bxf5
				8. e8=Q`,
			expectedFEN: "4QK2/8/2p5/5b2/5p2/8/2k5/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "6k1/8/p4K2/2P5/6p1/1b6/8/5B2 w - - 0 1",
			s: `1. c6 Ba4
				2. c7 Bd7
				3. Ke7  Bc8
				4. Kd8  Bf5
				5. Bd3  Be6
				6. Bc4  Kf7
				7. c8=Q`,
			expectedFEN: "2QK4/5k2/p3b3/8/2B3p1/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "k7/8/3P2K1/B7/8/8/7r/8 w - - 0 1",
			s: `1. d7 Rh8
				2. Kg7  Rb8
				3. Bc7`,
			expectedFEN: "kr6/2BP2K1/8/8/8/8/8/8 b - - 4 3",
			expectedErr: nil,
		},
		{
			fen: "8/8/1P6/3K4/1k6/5rB1/8/8 w - - 0 1",
			s: `1. b7 Rd3+
				2. Ke6! Rd8
				3. Bc7  Rh8
				4. Be5  Rd8
				5. Ke7  Rg8
				6. Kf7  Rd8
				7. Bc7  Rh8
				8. Bd6+ Ka5
				9. Bf8  Rh7+
				10. Bg7`,
			expectedFEN: "8/1P3KBr/8/k7/8/8/8/8 b - - 18 10",
			expectedErr: nil,
		},
		{
			fen: "8/8/2P5/4K3/8/1B5k/6r1/8 w - - 0 1",
			s: `1. c7 Re2+
				2. Kf6  Re8
				3. Ba4! Rg8
				4. Kf7  Rh8
				5. Kg7  Ra8
				6. Bc6  Ra7
				7. Bd7+!  Kh4
				8. c8=Q`,
			expectedFEN: "2Q5/r2B2K1/8/8/7k/8/8/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/P7/3p3r/3k3B/3P4/5K2/8 w - - 0 1",
			s: `1. a7 Rf5+
				2. Ke2  Re5+
				3. Kd2  Re8
				4. Bf2+ Ke5
				5. Bg3+ Kf5
				6. Bb8`,
			expectedFEN: "1B2r3/P7/8/3p1k2/8/3P4/3K4/8 b - - 10 6",
			expectedErr: nil,
		},
		{
			fen: "4k3/8/1P2b3/4K3/4B3/8/8/r7 w - - 0 1",
			s: `1. b7 Ra5+
				2. Kd6! Rb5!
				3. Bc6+ Kd8
				4. Bxb5 Bc8!
				5. b8=B!  Bg4
				6. Bc7+ Kc8
				7. Ba6#`,
			expectedFEN: "2k5/2B5/B2K4/8/6b1/8/8/8 b - - 4 7",
			expectedErr: nil,
		},
		{
			fen: "8/3b1K1B/6Pk/8/8/8/7P/8 w - - 0 1",
			s: `1. Kg8! Be6+
				2. Kh8  Bf5
				3. g7!  Bxh7
				4. h3! Kg6
				5. h4 Kh6
				6. h5`,
			expectedFEN: "7K/6Pb/7k/7P/8/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "5k2/8/5K2/6P1/5P2/1bB5/8/8 w - - 0 1",
			s: `1. g6 Bc2
				2. f5 Bb3
				3. Kg5  Bc2
				4. f6 Bb3
				5. Bb4+ Kg8
				6. Kf4  Bc4
				7. Ke5  Bb3
				8. Kd6  Kf8
				9. Kd7+ Kg8
				10. Ke7 Bc4
				11. f7+`,
			expectedFEN: "6k1/4KP2/6P1/8/1Bb5/8/8/8 b - - 0 11",
			expectedErr: nil,
		},
		{
			fen: "5b2/5k2/8/4P3/5P2/5K2/4B3/8 w - - 0 1",
			s: `1. Bc4+ Ke7
				2. Ke4! Bg7
				3. Kf5  Bh6
				4. Kg4! Bf8
				5. Kg5  Bg7
				6. Kg6  Bf8
				7. f5 Ke8
				8. f6 Bc5
				9. Kg7  Bf8+
				10. Kg8 Bc5
				11. e6  Bb4
				12. Bb5+  Kd8
				13. Kf7 Bc5
				14. e7+`,
			expectedFEN: "3k4/4PK2/5P2/1Bb5/8/8/8/8 b - - 0 14",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/4bk2/8/8/3B1PP1/5K2 w - - 0 1",
			s: `1. Ke2  Kg4
				2. Be1  Bd6
				3. f3+  Kf4
				4. g3+  Kf5
				5. g4+  Ke6
				6. Kd3  Kd5
				7. Bd2  Bc7
				8. f4 Bb6
				9. Bc3  Bc5
				10. g5  Bb6
				11. g6  Ke6
				12. Ke4 Bd8
				13. f5+ Ke7
				14. Kd5 Kf8
				15. Ke6 Bg5
				16. f6  Be3
				17. Bb4+`,
			expectedFEN: "5k2/8/4KPP1/8/1B6/4b3/8/8 b - - 2 17",
			expectedErr: nil,
		},
		{
			fen: "8/8/7k/8/8/6Pb/4B2P/6K1 w - - 0 1",
			s: `1. Bf1  Bg4
				2. h4 Bf5
				3. Kf2  Bg4
				4. Ke3  Be6
				5. Kf4  Bd7
				6. Bd3  Bh3
				7. Bf5  Bf1
				8. g4 Be2
				9. g5+  Kh5
				10. Kg3!  Bd1
				11. Be4 Bb3
				12. Bf3+  Kg6
				13. Kf4 Bf7
				14. h5+ Kg7
				15. Ke5 Bb3
				16. Be4 Bf7
				17. h6+ Kh8
				18. Kf6 Bh5
				19. Bd5 Kh7
				20. Bf7 Be2
				21. g6+`,
			expectedFEN: "8/5B1k/5KPP/8/8/8/4b3/8 b - - 0 21",
			expectedErr: nil,
		},
		{
			fen: "8/5b2/8/4k3/4B3/3PKP2/8/8 w - - 0 1",
			s: `1. f4+   Kd6
				2. f5  Ke5
				3. d4+   Kf6
				4. Kf4   Bb3
				5. Bc6   Bc2
				6. Bd7   Bb3
				7. Ke4   Bc4
				8. d5  Bb3
				9. Be6   Bc4
				10. Kd4  Be2
				11. d6   Bb5
				12. d7   Ke7
				13. f6+  Kd8
				14. f7   Ke7
				15. f8=Q+ Kxf8
				16. d8=Q+`,
			expectedFEN: "3Q1k2/8/4B3/1b6/3K4/8/8/8 b - - 0 16",
			expectedErr: nil,
		},
		{
			fen: "4k3/4P3/8/4Pp2/7K/8/4B2n/8 w - - 0 1",
			s: `1. Kg5  Ng4
				2. Kxf5 Nxe5
				3. Ke6! Nd7
				4. Bh5#`,
			expectedFEN: "4k3/3nP3/4K3/7B/8/8/8/8 b - - 3 4",
			expectedErr: nil,
		},
		{
			fen: "4b2k/7p/1p6/6K1/2P1B1P1/8/8/8 w - - 0 1",
			s: `1. Kh6  Bf7
				2. Bd3  Be6
				3. g5 Bg8
				4. Bxh7!  Bxh7
				5. g6 Bxg6
				6. Kxg6 Kg8
				7. Kf6  Kf8
				8. Ke6  Ke8
				9. Kd6  Kd8
				10. Kc6 Kc8
				11. Kxb6  Kb8
				12. c5  Kc8
				13. Kc6!`,
			expectedFEN: "2k5/8/2K5/2P5/8/8/8/8 b - - 2 13",
			expectedErr: nil,
		},
		{
			fen: "k2K4/8/6P1/2B5/5P2/p7/4p2b/8 w - - 0 1",
			s: `1. g7 e1=Q
				2. g8=Q Kb7
				3. Qb3+ Kc6
				4. Qb6+ Kd5
				5. Qb5! Qc1
				6. Be3+`,
			expectedFEN: "3K4/8/8/1Q1k4/5P2/p3B3/7b/2q5 b - - 8 6",
			expectedErr: nil,
		},
		{
			fen: "K7/8/Pp6/kp6/1n6/1P6/8/B7 w - - 0 1",
			s: `1. a7 Nc6
				2. Kb7  Nxa7
				3. Bc3+ b4
				4. Bd4  Nb5
				5. Bxb6#`,
			expectedFEN: "8/1K6/1B6/kn6/1p6/1P6/8/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "6n1/7k/2p2p2/8/P2P4/8/8/2B4K w - - 0 1",
			s: `1. Ba3  f5
				2. d5!  cxd5
				3. a5 Nf6
				4. a6 Ne8
				5. Bd6! Nxd6
				6. a7`,
			expectedFEN: "8/P6k/3n4/3p1p2/8/8/8/7K b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/2r3pP/6P1/8/K7/7k/3B4 w - - 0 1",
			s: `1. Bh5! Kg3
				2. Bxg6!  Kf4
				3. h7 Rc8
				4. Be8! Rxe8
				5. g6`,
			expectedFEN: "4r3/7P/6P1/8/5k2/K7/8/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "1kr5/4P1Pp/8/8/8/B7/7K/8 w - - 0 1",
			s: `1. e8=Q Rxe8
				2. Bf8  Re2+
				3. Kh3! Re3+
				4. Kh4  Re4+
				5. Kh5  Re5+
				6. Kh6  Re1
				7. Bc5! Re8
				8. Kxh7 Rd8
				9. Be7! Rc8
				10. Bf8 Rc7
				11. Bd6!`,
			expectedFEN: "1k6/2r3PK/3B4/8/8/8/8/8 b - - 6 11",
			expectedErr: nil,
		},
		{
			fen: "4k3/4p3/P2p4/8/2bP4/p7/2P5/K2B4 w - - 0 1",
			s: `1. a7 Bd5
				2. c4 Bb7
				3. Bf3  Bxf3
				4. d5`,
			expectedFEN: "4k3/P3p3/3p4/3P4/2P5/p4b2/8/K7 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "8/2p2b2/Pp1p4/4pp2/k7/2P5/1P2BPK1/8 w - - 0 1",
			s: `1. b3+  Kxb3
				2. c4 Be8
				3. Bf3  e4
				4. Bh5! Bc6
				5. Bf7  Ba8
				6. Bd5  c6
				7. c5+! cxd5
				8. cxb6`,
			expectedFEN: "b7/8/PP1p4/3p1p2/4p3/1k6/5PK1/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "8/1K6/8/8/3pk3/8/8/R7 w - - 0 1",
			s: `1. Kc6  d3
				2. Kc5  Ke3
				3. Kc4  d2
				4. Kc3  Ke2
				5. Ra2`,
			expectedFEN: "8/8/8/8/8/2K5/R2pk3/8 b - - 3 5",
			expectedErr: nil,
		},
		{
			fen: "K6R/8/7p/6k1/8/8/8/8 w - - 0 1",
			s: `1. Kb7  h5
				2. Kc6  h4
				3. Kd5  Kg4
				4. Ke4  Kg3
				5. Ke3  h3
				6. Rg8+ Kh2
				7. Kf2  Kh1
				8. Ra8  Kh2
				9. Rh8  Kh1
				10. Rxh3#`,
			expectedFEN: "8/8/8/8/8/7R/5K2/7k b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "K5R1/8/6p1/5k2/8/8/8/8 w - - 0 1",
			s: `1. Kb7  g5
				2. Kc6  g4
				3. Kd5  Kf4
				4. Kd4  Kf3
				5. Kd3  g3
				6. Rf8+ Kg2
				7. Ke2  Kg1
				8. Rg8  g2
				9. Kf3  Kh1
				10. Kf2`,
			expectedFEN: "6R1/8/8/8/8/8/5Kp1/7k b - - 3 10",
			expectedErr: nil,
		},
		{
			fen: "8/7p/8/8/6k1/8/8/KR6 w - - 0 1",
			s: `1. Rg1+ Kf5
				2. Rh1  Kg6
				3. Kb2  h5
				4. Kc3  Kg5
				5. Kd2  h4
				6. Ke2  Kg4
				7. Kf2  h3
				8. Rh2  Kh4
				9. Kf3`,
			expectedFEN: "8/8/8/8/7k/5K1p/7R/8 b - - 3 9",
			expectedErr: nil,
		},
		{
			fen: "1R6/8/8/7p/6k1/8/8/K7 w - - 0 1",
			s: `1. Rg8+ Kf3
				2. Rh8  Kg4
				3. Kb2  h4
				4. Kc2  h3
				5. Kd2  Kg3
				6. Ke2! Kg2
				7. Rg8+ Kh1
				8. Kf3  h2
				9. Kg3! Kg1
				10. Kh3+  Kh1
				11. Rd8 Kg1
				12. Rd1+`,
			expectedFEN: "8/8/8/8/8/7K/7p/3R2k1 b - - 7 12",
			expectedErr: nil,
		},
		{
			fen: "8/5K2/R7/8/p7/k7/8/8 w - - 0 1",
			s: `1. Rb6! Ka2
				2. Ke6  a3
				3. Kd5  Ka1
				4. Kc4  a2
				5. Kb3! Kb1
				6. Ka3+ Ka1
				7. Re6  Kb1
				8. Re1+`,
			expectedFEN: "8/8/8/8/8/K7/p7/1k2R3 b - - 7 8",
			expectedErr: nil,
		},
		{
			fen: "3R4/4K3/8/5kp1/8/8/8/8 w - - 0 1",
			s: `1. Rf8+ Kg6
				2. Rf6+ Kh5
				3. Ke6  g4
				4. Kf5  g3
				5. Rg6  Kh4
				6. Rh6#`,
			expectedFEN: "8/8/7R/5K2/7k/6p1/8/8 b - - 3 6",
			expectedErr: nil,
		},
		{
			fen: "2R5/7k/5K2/8/p7/8/1p6/8 w - - 0 1",
			s: `1. Rc7+ Kg8
				2. Rg7+ Kh8
				3. Rb7  a3
				4. Kg6  b1=Q
				5. Rxb1 a2
				6. Rb8#     `,
			expectedFEN: "1R5k/8/6K1/8/8/8/p7/8 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "K7/2k5/R7/8/5p2/6p1/8/8 w - - 0 1",
			s: `1. Rg6! Kd7
				2. Rg4! g2
				3. Rxg2 Ke6
				4. Rg5! Kf6
				5. Rc5  Ke6
				6. Kb7  Kf6
				7. Kb6  Ke6
				8. Kb5  Kf6
				9. Kc4  Ke6
				10. Kd3 Kf6
				11. Ke4`,
			expectedFEN: "8/8/5k2/2R5/4Kp2/8/8/8 b - - 16 11",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/3R4/2K5/6pp/k7/8 w - - 0 1",
			s: `1. Rd2+ Kb1
				2. Kc3  Kc1
				3. Ra2  Kd1
				4. Kd3  Kc1
				5. Ke3  h2
				6. Ra1+ Kb2
				7. Rh1  g2
				8. Rxh2`,
			expectedFEN: "8/8/8/8/8/4K3/1k4pR/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/p4K1R/8/p7/8/8/k7 w - - 0 1",
			s: `1. Rh2  a3
				2. Ke5  a2
				3. Kd4  Kb1
				4. Kc3  a1=N
				5. Rb2+ Kc1
				6. Ra2  Kb1
				7. Rxa6 Nc2
				8. Re6! Na3
				9. Re1+ Ka2
				10. Re2+  Ka1
				11. Kb3 Nb1
				12. Ra2#`,
			expectedFEN: "8/8/8/8/8/1K6/R7/kn6 b - - 10 12",
			expectedErr: nil,
		},
		{
			fen: "6k1/8/4K3/3R4/3p4/8/7b/8 w - - 0 1",
			s: `1. Rg5+ Kf8
				2. Rh5  Bc7
				3. Kd7  Bb6
				4. Rb5  Ba7
				5. Ra5  Bb6
				6. Ra8+ Kf7
				7. Kc6`,
			expectedFEN: "R7/5k2/1bK5/8/3p4/8/8/8 b - - 13 7",
			expectedErr: nil,
		},
		{
			fen: "8/p7/8/1P6/8/3K4/p3R3/1k6 w - - 0 1",
			s: `1. Re1+ Kb2
				2. Ra1! Kb3
				3. Kd2  Kb2
				4. Kd1  Kxa1
				5. Kc2  a5
				6. b6 a4
				7. b7 a3
				8. Kc3! Kb1
				9. b8=Q+  Kc1
				10. Qf4+  Kd1
				11. Qf1#`,
			expectedFEN: "8/8/8/8/8/p1K5/p7/3k1Q2 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "8/kp6/4R3/1P2K3/8/3pp3/8/8 w - - 0 1",
			s: `1. Kd6! d2
				2. Kc7  d1=Q
				3. Ra6+ bxa6
				4. b6+  Ka8
				5. b7+  Ka7
				6. b8=Q#`,
			expectedFEN: "1Q6/k1K5/p7/8/8/4p3/8/3q4 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/pp5R/k6p/2K5/1P4p1/8 w - - 0 1",
			s: `1. Rg5  h3
				2. Rg4+ b4+
				3. Kc4  h2
				4. Rg3  g1=Q
				5. Ra3+!  bxa3
				6. b3#`,
			expectedFEN: "8/8/8/p7/k1K5/pP6/7p/6q1 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "5Rn1/7b/8/8/k5P1/8/3K4/8 w - - 0 1",
			s: `1. Rf7  Bb1
				2. g5 Ba2
				3. Ra7+ Kb3
				4. Kc1`,
			expectedFEN: "6n1/R7/8/6P1/8/1k6/b7/2K5 b - - 4 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/4P2R/8/1r6/k7/8/K7 w - - 0 1",
			s: `1. e7 Re4
				2. Rh3+ Kb4
				3. Rh4! Rxh4
				4. e8=Q`,
			expectedFEN: "4Q3/8/8/8/1k5r/8/8/K7 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/4kr1P/7K/8/8/7R w - - 0 1",
			s: `1. h6  Kf6
				2. h7  Kg7
				3. h8=Q+! Kxh8
				4. Kg4+`,
			expectedFEN: "7k/8/8/5r2/6K1/8/8/7R b - - 1 4",
			expectedErr: nil,
		},
		{
			fen: "3r4/8/8/8/k2P4/3K4/8/1R6 w - - 0 1",
			s: `1. d5!  Ka5
				2. Kd4  Ka6
				3. Ke5  Re8+
				4. Kf6  Rd8
				5. Ke6  Re8+
				6. Kd7  Rh8
				7. d6 Rh7+
				8. Kc6  Ka5
				9. d7 Rh8
				10. Kc7 Rh7
				11. Kc8`,
			expectedFEN: "2K5/3P3r/8/k7/8/8/8/1R6 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "K7/P5R1/2k5/8/8/8/2r5/8 w - - 0 1",
			s: `1. Kb8  Rb2+
				2. Kc8  Ra2
				3. Rg6+ Kc5
				4. Kb7  Rb2+
				5. Kc7  Ra2
				6. Rg5+ Kc4
				7. Kb7  Rb2+
				8. Kc6  Ra2
				9. Rg4+ Kc3
				10. Kb6 Rb2+
				11. Kc5 Ra2
				12. Rg3+  Kc2
				13. Rg2+  Kb3
				14. Rxa2  Kxa2
				15. a8=Q+`,
			expectedFEN: "Q7/8/8/2K5/8/8/k7/8 b - - 0 15",
			expectedErr: nil,
		},
		{
			fen: "R7/P7/8/8/6K1/8/6k1/r7 w - - 0 1",
			s: `1. Kf4  Kf2
				2. Ke4  Ke2
				3. Kd4  Kd2
				4. Kc5  Kc3
				5. Rc8! Rxa7
				6. Kb6+`,
			expectedFEN: "2R5/r7/1K6/8/8/2k5/8/8 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "K7/P3k3/8/8/8/8/2R5/1r6 w - - 0 1",
			s: `1. Rc8  Kd6
				2. Rb8  Rh1
				3. Kb7  Rb1+
				4. Kc8  Rc1+
				5. Kd8  Rh1
				6. Rb6+ Kc5
				7. Rc6+!  Kb5
				8. Rc8  Rh8+
				9. Kc7  Rh7+
				10. Kb8`,
			expectedFEN: "1KR5/P6r/8/1k6/8/8/8/8 b - - 19 10",
			expectedErr: nil,
		},
		{
			fen: "R7/1K1k4/P7/8/8/8/8/2r5 w - - 0 1",
			s: `1. a7 Rb1+
				2. Ka6  Ra1+
				3. Kb6  Rb1+
				4. Kc5  Ra1
				5. Rh8  Rxa7
				6. Rh7+`,
			expectedFEN: "8/r2k3R/8/2K5/8/8/8/8 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "r7/8/4k3/P7/K7/3R4/8/8 w - - 0 1",
			s: `1. Kb5  Rb8+
				2. Kc6  Rc8+
				3. Kb7  Rc1
				4. a6 Rb1+
				5. Kc6  Rc1+
				6. Kb5  Rb1+
				7. Ka4  Rb8
				8. Ka5  Ra8
				9. Kb6  Rb8+
				10. Kc7 Rb1
				11. Ra3 Rc1+
				12. Kb6 Rb1+
				13. Ka5 Rb8
				14. a7  Ra8
				15. Kb6`,
			expectedFEN: "r7/P7/1K2k3/8/8/R7/8/8 b - - 2 15",
			expectedErr: nil,
		},
		{
			fen: "8/r4k2/8/8/P7/K3R3/8/8 w - - 0 1",
			s: `1. Kb4  Re7
				2. Ra3  Ke8
				3. a5 Kd8
				4. a6 Ra7
				5. Kb5  Kc8
				6. Rh3  Kb8
				7. Rh8+ Kc7
				8. Rg8  Kd6
				9. Kb6  Rf7
				10. a7`,
			expectedFEN: "6R1/P4r2/1K1k4/8/8/8/8/8 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "6K1/6R1/r5P1/8/7k/8/8/8 w - - 0 1",
			s: `1. Rh7+ Kg5
				2. g7 Ra8+
				3. Kf7  Ra7+
				4. Ke6  Ra6+
				5. Kd5! Rg6
				6. Ke5! Kg4
				7. Rh1  Kg5
				8. Rg1+ Kh6
				9. Rxg6+  Kxg6
				10. g8=Q+`,
			expectedFEN: "6Q1/8/6k1/4K3/8/8/8/8 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "7r/6k1/1P6/8/8/5R2/8/K7 w - - 0 1",
			s: `1. Kb2   Rb8
				2. Rb3   Kf7
				3. Ka3!  Ke7
				4. Ka4   Kd7
				5. Kb5   Kc8
				6. Rc3+ Kb7
				7. Rc7+ Ka8
				8. Ra7#!`,
			expectedFEN: "kr6/R7/1P6/1K6/8/8/8/8 b - - 15 8",
			expectedErr: nil,
		},
		{
			fen: "1r6/4k3/8/8/1P6/1K6/8/3R4 w - - 0 1",
			s: `1. Rd4! Ke6
				2. Kc4! Rc8+
				3. Kb5  Rb8+
				4. Kc6  Rc8+
				5. Kb7  Rf8
				6. b5`,
			expectedFEN: "5r2/1K6/4k3/1P6/3R4/8/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "1r6/8/8/5k2/8/1P6/1K6/4R3 w - - 0 1",
			s: `1. Kc3  Rc8+
				2. Kd4  Rb8
				3. Kc4  Rc8+
				4. Kd5  Rb8
				5. Rb1  Kf6
				6. b4 Ke7
				7. Kc6! Kd8
				8. Rd1+ Kc8
				9. Rh1  Kd8
				10. Rh8+`,
			expectedFEN: "1r1k3R/8/2K5/8/1P6/8/8/8 b - - 8 10",
			expectedErr: nil,
		},
		{
			fen: "8/k7/8/3P4/7r/2K5/1R6/8 w - - 0 1",
			s: `1. Kd3! Rg4
				2. Ke3  Rh4
				3. d6 Rh6
				4. Rd2! Rh8
				5. d7 Rd8
				6. Ke4  Kb7
				7. Ke5  Kc7
				8. Ke6  Rh8
				9. Rc2+ Kb7
				10. Rh2!  Rg8
				11. Kf7 Ra8
				12. Ke7 Kc7
				13. d8=Q+ Rxd8
				14. Rc2+  Kb7
				15. Kxd8`,
			expectedFEN: "3K4/1k6/8/8/8/8/2R5/8 b - - 0 15",
			expectedErr: nil,
		},
		{
			fen: "3K4/3P2k1/8/8/8/8/2r5/5R2 w - - 0 1",
			s: `1. Ra1 Kf7
				2. Ra8 Rc1
				3. Rc8 Rd1
				4. Kc7 Rc1+
				5. Kb6 Rb1+
				6. Kc5`,
			expectedFEN: "2R5/3P1k2/8/2K5/8/8/8/1r6 b - - 11 6",
			expectedErr: nil,
		},
		{
			fen: "6k1/r2KP3/8/8/8/8/8/5R2 w - - 0 1",
			s: `1. Ke6  Ra6+
				2. Ke5! Ra5+
				3. Kf6  Ra6+
				4. Kg5  Ra5+
				5. Kg6  Ra6+
				6. Rf6  Ra8
				7. Rd6`,
			expectedFEN: "r5k1/4P3/3R2K1/8/8/8/8/8 b - - 13 7",
			expectedErr: nil,
		},
		{
			fen: "8/1r2K3/4P1k1/8/8/8/8/R7 w - - 0 1",
			s: `1. Kd8  Rb8+
				2. Kc7  Rb2
				3. Re1  Rc2+
				4. Kd7  Rd2+
				5. Ke8  Ra2
				6. e7 Ra8+
				7. Kd7  Ra7+
				8. Kc6`,
			expectedFEN: "8/r3P3/2K3k1/8/8/8/8/4R3 b - - 4 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/7k/4P3/r7/4K3/8/6R1 w - - 0 1",
			s: `1. Kd3! Rb4
				2. e6!  Rb6
				3. Re1  Rb8
				4. e7 Re8
				5. Kd4  Kg7
				6. Kd5  Kf7
				7. Kd6  Ra8
				8. Rf1+  Ke8
				9. Rf8#`,
			expectedFEN: "r3kR2/4P3/3K4/8/8/8/8/8 b - - 10 9",
			expectedErr: nil,
		},
		{
			fen: "4r3/8/7k/8/8/4P3/4K3/6R1 w - - 0 1",
			s: `1. Kd3! Rd8+
				2. Kc4  Rc8+
				3. Kd5  Rd8+
				4. Ke5  Re8+
				5. Kf6! Rf8+
				6. Ke7  Rf5
				7. e4!  Re5+
				8. Kf6  Rh5
				9. e5 Rh2
				10. Rf1 Kh7
				11. Kf7`,
			expectedFEN: "8/5K1k/8/4P3/8/8/7r/5R2 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "1K6/1P1k4/1P6/8/8/r7/2R5/8 w - - 0 1",
			s: `1. Rd2+ Ke7
				2. Rd6! Rc3
				3. Rc6! Rxc6
				4. Ka7`,
			expectedFEN: "8/KP2k3/1Pr5/8/8/8/8/8 b - - 1 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/P7/8/3K4/r7/P2R4/k7 w - - 0 1",
			s: `1. Kc4  Rxa6
				2. a4!  Rxa4+
				3. Kb3!`,
			expectedFEN: "8/8/8/8/r7/1K6/3R4/k7 b - - 1 3",
			expectedErr: nil,
		},
		{
			fen: "2R5/2P5/1k6/8/7r/1K5P/1P6/8 w - - 0 1",
			s: `1. Rd8    Rxh3+
				2. Rd3!   Rxd3+
				3. Kc2!   Rd6!
				4. c8=N+! Kc6
				5. Nxd6   Kxd6
				6. Kb3    Kc6
				7. Ka4    Kb6
				8. Kb4    Kc6
				9. Ka5`,
			expectedFEN: "8/8/2k5/K7/8/8/1P6/8 b - - 7 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/P5rk/8/8/KpR5/8/8 w - - 0 1",
			s: `1. Rh3+ Kg7
				2. Rg3! Rxg3
				3. a7 Rg1
				4. Kb2  Rg2+
				5. Kxb3 Rg3+
				6. Kb4  Rg4+
				7. Kb5  Rg5+
				8. Kb6  Rg6+
				9. Kb7`,
			expectedFEN: "8/PK4k1/6r1/8/8/8/8/8 b - - 8 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/6P1/8/5p1k/r7/6K1/2R5 w - - 0 1",
			s: `1. Rg1  Ra8
				2. g7 Kh5
				3. Kf3  Rg8
				4. Kxf4 Kh6
				5. Kf5  Kh7
				6. Kf6  Ra8
				7. Rh1+ Kg8
				8. Rh8#`,
			expectedFEN: "r5kR/6P1/5K2/8/8/8/8/8 b - - 8 8",
			expectedErr: nil,
		},
		{
			fen: "6k1/4K3/8/r4p2/4P3/8/8/5R2 w - - 0 1",
			s: `1. Rg1+ Kh7
				2. e5!  Rxe5+
				3. Kf7  Kh6
				4. Kf6`,
			expectedFEN: "8/8/5K1k/4rp2/8/8/8/6R1 b - - 3 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/2P2K2/8/6pr/6k1/R7/8 w - - 0 1",
			s: `1. Kg5! Kh3
				2. Rh2+!  Kxh2
				3. Kxh4 g3
				4. c7 g2
				5. c8=Q g1=Q
				6. Qh3#`,
			expectedFEN: "8/8/8/8/7K/7Q/7k/6q1 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "4k3/6K1/1P6/8/1p6/8/r7/6R1 w - - 0 1",
			s: `1. b7 Ra7
				2. Re1+ Kd8
				3. Re7! Kxe7
				4. b8=Q`,
			expectedFEN: "1Q6/r3k1K1/8/8/1p6/8/8/8 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "3r4/P4Kp1/k7/8/8/8/1R6/8 w - - 0 1",
			s: `1. Ke7! Ra8
				2. Kd7! Rf8
				3. Rf2! Ra8
				4. Kc7 Rxa7+
				5. Kc6 Ka5
				6. Ra2+`,
			expectedFEN: "8/r5p1/2K5/k7/8/8/R7/8 b - - 3 6",
			expectedErr: nil,
		},
		{
			fen: "2K5/2P2R2/k7/8/8/8/2r2p2/8 w - - 0 1",
			s: `1. Kb8  Rb2+
				2. Ka8  Rc2
				3. Rf6+ Ka5
				4. Kb8  Rb2+
				5. Ka7  Rc2
				6. Rf5+ Ka4
				7. Kb7  Rb2+
				8. Ka6  Rc2
				9. Rf4+ Ka3
				10. Kb6 Rb2+
				11. Ka5 Rc2
				12. Rf3+  Kb2
				13. Rxf2  Rxf2
				14. c8=Q`,
			expectedFEN: "2Q5/8/8/K7/8/8/1k3r2/8 b - - 0 14",
			expectedErr: nil,
		},
		{
			fen: "R6K/7P/p6k/r7/8/8/8/8 w - - 0 1",
			s: `1. Rf8  Rg5
				2. Rf6+ Kh5
				3. Rf5! Rxf5
				4. Kg7  Rg5+
				5. Kf7  Rf5+
				6. Ke7  Re5+
				7. Kd7  Rd5+
				8. Kc7  Rc5+
				9. Kb7  Rb5+
				10. Ka7`,
			expectedFEN: "8/K6P/p7/1r5k/8/8/8/8 b - - 13 10",
			expectedErr: nil,
		},
		{
			fen: "8/8/1P6/8/2k3pr/8/3RK3/8 w - - 0 1",
			s: `1. b7 Rh2+
				2. Ke3  Rh3+
				3. Ke4! Rb3
				4. Rc2+ Kb5
				5. b8=Q`,
			expectedFEN: "1Q6/8/8/1k6/4K1p1/1r6/2R5/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "8/6pP/6R1/7r/8/8/5K2/7k w - - 0 1",
			s: `1. Rg1+ Kh2
				2. Rg2+ Kh1
				3. Kg3! Rh6
				4. h8=Q!  Rxh8
				5. Ra2`,
			expectedFEN: "7r/6p1/8/8/8/6K1/R7/7k b - - 1 5",
			expectedErr: nil,
		},
		{
			fen: "8/5R2/4PK2/2r5/8/2k1p3/8/8 w - - 0 1",
			s: `1. e7 Rc8
				2. Rf8  Rc6+
				3. Kf5  Re6!
				4. Kxe6 e2
				5. Rf3+ Kd4
				6. Rf4+ Kd3
				7. Re4! Kxe4
				8. e8=Q e1=Q
				9. Kf6+`,
			expectedFEN: "4Q3/8/5K2/8/4k3/8/8/4q3 b - - 1 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/2R3P1/8/7p/4k3/3r4/7K w - - 0 1",
			s: `1. g7  Rd8
				2. Re6+  Kf4
				3. Rf6+  Kg3
				4. Rg6+  Kh3
				5. Kg1   Rg8
				6. Kf2   Kh2
				7. Rg4!  Kh3
				8. Kf3   Kh2
				9. Rxh4+ Kg1
				10. Rh7  Rd8
				11. Rh8`,
			expectedFEN: "3r3R/6P1/8/8/8/5K2/8/6k1 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "8/8/1R4P1/5p2/3K4/k3p3/4r3/8 w - - 0 1",
			s: `1. Kc3  Ka4
				2. Rb4+ Ka5
				3. Rg4! fxg4
				4. g7`,
			expectedFEN: "8/6P1/8/k7/6p1/2K1p3/4r3/8 b - - 0 4",
			expectedErr: nil,
		},
		{
			fen: "K7/2rp4/8/2p3P1/8/8/5R2/7k w - - 0 1",
			s: `1. Kb8  Rc6
				2. Rf6  c4
				3. Rh6+ Kg2
				4. Rxc6 dxc6
				5. g6 c3
				6. g7 c2
				7. g8=Q+`,
			expectedFEN: "1K4Q1/8/2p5/8/8/8/2p3k1/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/R1p5/5pP1/8/8/7K/1k6/2r5 w - - 0 1",
			s: `1. Rb7+   Ka3
				2. g7     Rg1
				3. Rxc7   Kb4!
				4. Kh4     Rg5
				5. Re7     Kc5
				6. Re5+! Rxe5
				7. g8=Q`,
			expectedFEN: "6Q1/8/5p2/2k1r3/7K/8/8/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "7k/8/4K1Rp/6pP/5p1r/8/8/8 w - - 0 1",
			s: `1. Kf7    Rxh5
				2. Rg8+ Kh7
				3. Rg7+ Kh8
				4. Kg6    g4
				5. Ra7    Rg5+
				6. Kxh6 g3
				7. Kxg5 g2
				8. Ra1    f3
				9. Kg6!`,
			expectedFEN: "7k/8/6K1/8/8/5p2/6p1/R7 b - - 1 9",
			expectedErr: nil,
		},
		{
			fen: "k7/p2K1p2/7R/r1p1P3/1p6/8/8/8 w - - 0 1",
			s: `1. e6!    fxe6
				2. Kc6    a6
				3. Rh8+ Ka7
				4. Rh7+ Ka8
				5. Kb6    Rb5+
				6. Kxa6 Rb8
				7. Ra7#`,
			expectedFEN: "kr6/R7/K3p3/2p5/1p6/8/8/8 b - - 2 7",
			expectedErr: nil,
		},
		{
			fen: "8/2k5/2P3Pr/8/8/8/6p1/2R3K1 w - - 0 1",
			s: `1. g7 Rg6
				2. Ra1  Kb8
				3. c7+  Kxc7
				4. Ra8  Rxg7
				5. Ra7+`,
			expectedFEN: "8/R1k3r1/8/8/8/8/6p1/6K1 b - - 1 5",
			expectedErr: nil,
		},
		{
			fen: "1R6/2pk4/1P6/8/3r4/K7/P7/8 w - - 0 1",
			s: `1. Rd8+! Kxd8
				2. b7    Rb4!
				3. Kxb4  c5+
				4. Kb5!   Kc7
				5. Ka6     Kb8
				6. Kb6     c4
				7. a4    c3
				8. a5    c2
				9. a6    c1=Q
				10. a7#!`,
			expectedFEN: "1k6/PP6/1K6/8/8/8/8/2q5 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "2r5/6k1/6p1/P7/3P4/8/6K1/3R4 w - - 0 1",
			s: `1. d5   Ra8
				2. d6   Kf7
				3. d7   Ke7
				4. a6   Kd8
				5. Rd6   Rb8
				6. Kf3    Rb4
				7. a7   Ra4
				8. Rxg6 Kxd7
				9. Rg8    Rxa7
				10. Rg7+`,
			expectedFEN: "8/r2k2R1/8/8/8/5K2/8/8 b - - 1 10",
			expectedErr: nil,
		},
		{
			fen: "6R1/3r3p/2p2pP1/8/p1P5/8/5K1k/8 w - - 0 1",
			s: `1. Rh8    Rd2+
				2. Kf1    Rd1+
				3. Ke2    Rg1
				4. Rxh7+  Kg3
				5. Rh1! Rg2+
				6. Ke3    Kg4
				7. Rh2! Rg3+
				8. Kf2    Rf3+
				9. Kg1! Rg3+
				10. Rg2`,
			expectedFEN: "8/8/2p2pP1/8/p1P3k1/6r1/6R1/6K1 b - - 12 10",
			expectedErr: nil,
		},
		{
			fen: "5R2/1n6/P6k/8/1P5r/8/3K1p2/8 w - - 0 1",
			s: `1. a7   f1=Q!
				2. Rxf1  Rh2+
				3. Kc1    Ra2
				4. b5!    Ra1+
				5. Kc2    Ra2+
				6. Kb1    Rxa7
				7. Rf6+  Kg5
				8. Ra6    Nd6
				9. Rxa7  Nxb5
				10. Ra5`,
			expectedFEN: "8/8/8/Rn4k1/8/8/8/1K6 b - - 1 10",
			expectedErr: nil,
		},
		{
			fen: "5r1k/R5p1/6Pp/5P1K/7P/1p1p4/8/8 w - - 0 1",
			s: `1. f6!     Rg8
				2. Rf7     d2
				3. fxg7+ Rxg7
				4. Kxh6  d1=Q
				5. Rf8+   Rg8
				6. g7#`,
			expectedFEN: "5Rrk/6P1/7K/8/7P/1p6/8/3q4 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "8/r4kp1/5p1p/3R1K1P/4PPP1/8/8/8 w - - 0 1",
			s: `1. e5 fxe5
				2. fxe5 Ke7
				3. e6 Ra4
				4. g5!  hxg5
				5. Rd7+ Kf8
				6. Rf7+ Kg8
				7. Kg6  g4
				8. h6!  gxh6
				9. e7 Ra8
				10. Rf6 Re8
				11. Rd6`,
			expectedFEN: "4r1k1/4P3/3R2Kp/8/6p1/8/8/8 b - - 4 11",
			expectedErr: nil,
		},
		{
			fen: "8/Q5K1/8/8/8/8/3kp3/8 w - - 0 1",
			s: `1. Qf2    Kd1
				2. Qd4+ Kc2
				3. Qe3    Kd1
				4. Qd3+ Ke1
				5. Kf6    Kf2
				6. Qd2    Kf1
				7. Qf4+ Kg2
				8. Qe3    Kf1
				9. Qf3+ Ke1
				10. Ke5   Kd2
				11. Qf2   Kd1
				12. Qd4+  Kc2
				13. Qe3   Kd1
				14. Qd3+  Ke1
				15. Ke4   Kf2
				16. Qf3+  Ke1
				17. Kd3   Kd1
				18. Qxe2+ Kc1
				19. Qc2#`,
			expectedFEN: "8/8/8/8/8/3K4/2Q5/2k5 b - - 2 19",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/4K3/8/8/3Q1p2/6k1 w - - 0 1",
			s: `1. Kf4    f1=Q+
				2. Kg3 Qf5
				3. Qg2#`,
			expectedFEN: "8/8/8/5q2/8/6K1/6Q1/6k1 b - - 3 3",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/6K1/6Q1/2p5/3k4 w - - 0 1",
			s: `1. Qb3    Kd2
				2. Qb2    Kd1
				3. Kf3!   Kd2
				4. Kf2    Kd1
				5. Qd4+ Kc1
				6. Qb4    Kd1
				7. Qe1#`,
			expectedFEN: "8/8/8/8/8/8/2p2K2/3kQ3 b - - 13 7",
			expectedErr: nil,
		},
		{
			fen: "5Q2/8/8/8/K7/8/2p5/1k6 w - - 0 1",
			s: `1. Qf5    Ka1
				2. Kb3    c1=Q
				3. Qa5+ Kb1
				4. Qa2#`,
			expectedFEN: "8/8/8/8/8/1K6/Q7/1kq5 b - - 3 4",
			expectedErr: nil,
		},
		{
			fen: "Q7/8/8/8/8/8/pk2K3/8 w - - 0 1",
			s: `1. Qb7+ Ka1
				2. Qe4    Kb2
				3. Qd4+ Kb1
				4. Qd1+ Kb2
				5. Qd2+ Kb1
				6. Kd1!  a1=Q
				7. Qc2#`,
			expectedFEN: "8/8/8/8/8/8/2Q5/qk1K4 b - - 1 7",
			expectedErr: nil,
		},
		{
			fen: "Q7/8/8/3K4/8/8/pk6/8 w - - 0 1",
			s: `1. Qh8+ Kb1
				2. Qh1+ Kb2
				3. Qh2+ Kb1
				4. Kc4!   a1=Q
				5. Kb3`,
			expectedFEN: "8/8/8/8/8/1K6/7Q/qk6 b - - 1 5",
			expectedErr: nil,
		},
		{
			fen: "8/KQ6/8/8/8/8/p7/k7 w - - 0 1",
			s: `1. Kb6!   Kb2
				2. Ka5+  Kc1
				3. Qh1+  Kb2
				4. Qg2+  Kb1
				5. Ka4     a1=Q+
				6. Kb3`,
			expectedFEN: "8/8/8/8/8/1K6/6Q1/qk6 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "7Q/8/1K6/8/p7/8/7p/6k1 w - - 0 1",
			s: `1. Qg8+ Kf2
				2. Qh7    Kg3
				3. Qd3+ Kg2
				4. Qe4+ Kg3
				5. Kc5    a3
				6. Kd4    a2
				7. Qh1    a1=Q+
				8. Qxa1 Kg2
				9. Qb2+ Kg1
				10. Ke3 h1=Q
				11. Qf2#`,
			expectedFEN: "8/8/8/8/8/4K3/5Q2/6kq b - - 1 11",
			expectedErr: nil,
		},
		{
			fen: "3Q4/kr6/2K5/8/8/8/8/8 w - - 0 1",
			s: `1. Qd4+ Ka8
				2. Qh8+ Ka7
				3. Qd8 Ka6
				4. Qa8+ Ra7
				5. Qb8`,
			expectedFEN: "1Q6/r7/k1K5/8/8/8/8/8 b - - 9 5",
			expectedErr: nil,
		},
		{
			fen: "7Q/8/8/8/8/8/6rp/4K2k w - - 0 1",
			s: `1. Qa8  Kg1
				2. Qa7+ Kh1
				3. Qb7  Kg1
				4. Qb6+ Kh1
				5. Qc6  Kg1
				6. Qc5+ Kh1
				7. Qd5  Kg1
				8. Qd4+ Kh1
				9. Qe4  Kg1
				10. Qe3+  Kh1
				11. Qf3 Kg1
				12. Qf1#`,
			expectedFEN: "8/8/8/8/8/8/6rp/4KQk1 b - - 23 12",
			expectedErr: nil,
		},
		{
			fen: "8/P7/R7/8/2q5/5K2/7k/8 w - - 0 1",
			s: `1. Rh6+ Kg1
				2. Rh1+ Kxh1
				3. a8=Q Kh2
				4. Qh8+ Kg1
				5. Qg7+ Kf1
				6. Qa1+ Qc1
				7. Qxc1#`,
			expectedFEN: "8/8/8/8/8/5K2/8/2Q2k2 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/5P1k/5Q2/7q/8/6K1/8/8 w - - 0 1",
			s: `1. f8=N+! Kg8
				2. Ne6  Qf7
				3. Qd8+ Kh7
				4. Ng5+`,
			expectedFEN: "3Q4/5q1k/8/6N1/8/6K1/8/8 b - - 6 4",
			expectedErr: nil,
		},
		{
			fen: "8/3KP1q1/8/8/8/4Q3/k7/8 w - - 0 1",
			s: `1. Kc8  Qg4+
				2. Kb8  Qb4+
				3. Ka8`,
			expectedFEN: "K7/4P3/8/8/1q6/4Q3/k7/8 b - - 5 3",
			expectedErr: nil,
		},
		{
			fen: "6K1/5P2/6Q1/3q4/1k6/8/8/8 w - - 0 1",
			s: `1. Kh7  Qh1+
				2. Kg7! Qa1+
				3. Kg8  Qa2
				4. Qb6+ Kc3
				5. Kg7! Qg2+
				6. Qg6  Qb7
				7. Kg8  Qd5
				8. Kh7  Qh1+
				9. Qh6  Qe4+
				10. Kh8 Qe5+
				11. Qg7`,
			expectedFEN: "7K/5PQ1/8/4q3/8/2k5/8/8 b - - 21 11",
			expectedErr: nil,
		},
		{
			fen: "K7/1P6/k1q5/8/8/1Q6/8/8 w - - 0 1",
			s: `1. Qb4  Qh1
				2. Qa3+ Kb6
				3. Qb2+ Kc5
				4. Ka7  Qh7
				5. Qb6+ Kd5
				6. Ka6`,
			expectedFEN: "8/1P5q/KQ6/3k4/8/8/8/8 b - - 11 6",
			expectedErr: nil,
		},
		{
			fen: "6k1/5q2/3Pp2K/8/4Q3/8/8/8 w - - 0 1",
			s: `1. Qg2+ Kf8
				2. Qa8+ Qe8
				3. Qb7  Qd8
				4. Kg6  Qe8+
				5. Kf6  Qd8+
				6. Kxe6 Qe8+
				7. Qe7+ Qxe7+
				8. dxe7+  Ke8
				9. Kf6`,
			expectedFEN: "4k3/4P3/5K2/8/8/8/8/8 b - - 2 9",
			expectedErr: nil,
		},
		{
			fen: "8/7q/2K2p2/4p3/2k1P2p/8/3P4/7Q w - - 0 1",
			s: `1. Qb1  Kd4
				2. Qb3! Qxe4+
				3. Kd6  Qa8
				4. Qe3+ Kc4
				5. Qc3+ Kb5
				6. Qb3+ Ka6
				7. Qa4+ Kb7
				8. Qb5+ Ka7
				9. Kc7!`,
			expectedFEN: "q7/k1K5/5p2/1Q2p3/7p/8/3P4/8 b - - 13 9",
			expectedErr: nil,
		},
		{
			fen: "2b4q/4K3/4p3/p2p1pk1/8/1P3P2/2P5/1Q6 w - - 0 1",
			s: `1. Qc1+ f4
				2. Qg1+ Kf5
				3. Qg4+ Ke5
				4. Qg5+ Kd4
				5. Qg1+ Ke5
				6. Qa1+ d4
				7. Qxa5#`,
			expectedFEN: "2b4q/4K3/4p3/Q3k3/3p1p2/1P3P2/2P5/8 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/5k2/3q1N2/Q4K2/8 w - - 0 1",
			s: `1. Qf7+ Qf5
				2. Qc4+ Qe4
				3. Qc7+ Kf5
				4. Qf7+ Kg4
				5. Qg7+ Kf5
				6. Nd4+ Kf4
				7. Qg3#`,
			expectedFEN: "8/8/8/8/3Nqk2/6Q1/5K2/8 b - - 13 7",
			expectedErr: nil,
		},
		{
			fen: "8/5B2/q7/8/4k3/8/5K2/3Q4 w - - 0 1",
			s: `1. Qg4+ Ke5
				2. Qg5+ Ke4
				3. Bg6+ Kd4
				4. Qe3+ Kd5
				5. Be4+!`,
			expectedFEN: "8/8/q7/3k4/4B3/4Q3/5K2/8 b - - 9 5",
			expectedErr: nil,
		},
		{
			fen: "1QK3kq/6p1/8/6B1/8/8/8/8 w - - 0 1",
			s: `1. Kb7+!  Kh7
				2. Qh2+ Kg8
				3. Qa2+ Kh7
				4. Qf7! Qg8
				5. Qh5#!`,
			expectedFEN: "6q1/1K4pk/8/6BQ/8/8/8/8 b - - 9 5",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/1q6/3Q3K/p7/k1N5 w - - 0 1",
			s: `1. Qf1  Qb1
				2. Qf6+ Qb2
				3. Nb3+ Kb1
				4. Qf1+ Kc2
				5. Na1+ Kd2
				6. Qf2+ Kc3
				7. Qf6+`,
			expectedFEN: "8/8/5Q2/8/8/2k4K/pq6/N7 b - - 13 7",
			expectedErr: nil,
		},
		{
			fen: "8/2K5/2N5/pk6/q7/4Q3/8/8 w - - 0 1",
			s: `1. Qb6+ Kc4
				2. Qd4+ Kb3
				3. Qd1+ Ka3
				4. Qa1+ Kb3
				5. Nxa5+  Kb4
				6. Qd4+ Ka3
				7. Nc4+ Kb3
				8. Nd2+ Ka3
				9. Qa1+ Kb4
				10. Qb2+  Kc5
				11. Qb6+  Kd5
				12. Qd6#`,
			expectedFEN: "8/2K5/3Q4/3k4/q7/8/3N4/8 b - - 14 12",
			expectedErr: nil,
		},
		{
			fen: "K7/8/8/5q2/3Q4/k6p/4N3/8 w - - 0 1",
			s: `1. Qa1+ Kb4
				2. Qb2+ Ka4
				3. Nc3+ Ka5
				4. Qa3+ Kb6
				5. Qd6+ Ka5
				6. Kb8!!  h2
				7. Qa3+ Kb6
				8. Qa7+ Kc6
				9. Qc7#`,
			expectedFEN: "1K6/2Q5/2k5/5q2/8/2N5/7p/8 b - - 5 9",
			expectedErr: nil,
		},
		{
			fen: "8/8/1q1p4/1ppp4/8/5N1K/5k2/3Q4 w - - 0 1",
			s: `1. Nh2  Ke3
				2. Ng4+ Kf4
				3. Qf1+ Ke4
				4. Nf6+ Kd4
				5. Qd1+ Kc4
				6. Qxd5+  Kc3
				7. Qa8! Kd4
				8. Nd5`,
			expectedFEN: "Q7/8/1q1p4/1ppN4/3k4/7K/8/8 b - - 4 8",
			expectedErr: nil,
		},
		{
			fen: "1R6/8/8/6K1/8/2Q5/qr6/k7 w - - 0 1",
			s: `1. Kh4! Qa4+
				2. Kg3  Qa2
				3. Qe1+ Rb1
				4. Qe5+ Rb2
				5. Qc3  Kb1
				6. Qe1+ Kc2
				7. Rc8+ Kd3
				8. Rc3+ Kd4
				9. Qe3+ Kd5
				10. Rc5+  Kd6
				11. Qe5+  Kd7
				12. Rc7+  Kd8
				13. Qe7#`,
			expectedFEN: "3k4/2R1Q3/8/8/8/6K1/qr6/8 b - - 25 13",
			expectedErr: nil,
		},
		{
			fen: "3q3b/2pp4/8/2B5/8/1k6/3P3Q/3K4 w - - 0 1",
			s: `1. d4   Kc4
				2. Qa2+   Kb5
				3. Qb3+   Kc6
				4. Be7!   Qxe7
				5. d5+    Kd6
				6. Qb4+   c5
				7. dxc6e.p.+  Ke6
				8. Qxe7+    Kxe7
				9. c7`,
			expectedFEN: "7b/2Ppk3/8/8/8/8/8/3K4 b - - 0 9",
			expectedErr: nil,
		},
		{
			fen: "3q2r1/6P1/p7/8/8/k3N1Qp/8/2K5 w - - 0 1",
			s: `1. Nc2+ Ka4
				2. Qg4+ Kb5
				3. Nd4+ Kb6
				4. Qg6+ Kb7
				5. Qe4+!  Kc8
				6. Qa8+ Kd7
				7. Qc6+ Ke7
				8. Qe6#!`,
			expectedFEN: "3q2r1/4k1P1/p3Q3/8/3N4/7p/8/2K5 b - - 15 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/2Q4K/5p2/1k1qn1N1/1P3p2/5P2/8 w - - 0 1",
			s: `1. Qa4+ Kc5
				2. Qa7+ Kd5
				3. Ne3+ Ke5
				4. Qg7+ Kf4!
				5. Ng2+!  fxg2
				6. Qxd4   g1=Q
				7. Qe3+ Ke5
				8. f4+`,
			expectedFEN: "8/8/7K/4kp2/4nP2/1P2Q3/8/6q1 b - f3 0 8",
			expectedErr: nil,
		},
		{
			fen: "6q1/7k/1r2NQ2/1p4p1/5PP1/1K3p2/8/8 w - - 0 1",
			s: `1. Qf5+   Kh8
				2. Qe5+   Kh7
				3. Qe4+   Kh8
				4. Qd4+   Kh7
				5. Qd3+   Kh8
				6. Qc3+   Kh7
				7. Qc2+   Kh8
				8. Qb2+   Kh7
				9. Qb1+ Kh8
				10. Qh1+  Qh7
				11. Qa1+  Kg8
				12. Qa8+  Kf7
				13. Nxg5+`,
			expectedFEN: "Q7/5k1q/1r6/1p4N1/5PP1/1K3p2/8/8 b - - 0 13",
			expectedErr: nil,
		},
		{
			fen: "8/K6N/8/2N5/1n6/6Q1/6pn/7k w - - 0 1",
			s: `1. Ne4 Nd3
				2. Qf2! Nxf2
				3. Ng3+ Kg1
				4. Ng5 Ng4
				5. Nf3#`,
			expectedFEN: "8/K7/8/8/6n1/5NN1/5np1/6k1 b - - 5 5",
			expectedErr: nil,
		},
		{
			fen: "3KNB1k/7p/8/8/8/8/3n4/8 w - - 0 1",
			s: `1. Bh6! Ne4
				2. Kd7! Nf6+
				3. Ke7  Ng8+
				4. Kf8! Nxh6
				5. Nd6  Nf5
				6. Nf7#`,
			expectedFEN: "5K1k/5N1p/8/5n2/8/8/8/8 b - - 3 6",
			expectedErr: nil,
		},
		{
			fen: "3K4/1p5p/1p6/8/8/1B6/1p6/b1k3B1 w - - 0 1",
			s: `1. Be3+ Kb1
				2. Bh6  b5
				3. Ke7  b4
				4. Kf6  b5
				5. Kg5  Kc1
				6. Kf5+ Kb1
				7. Kf4  Kc1
				8. Kf3+ Kb1
				9. Ke3  Kc1
				10. Ke2+  Kb1
				11. Bd2 h5
				12. Kd1 h4
				13. Bxb4  h3
				14. Bd5 h2
				15. Kd2 h1=Q
				16. Bxh1  Ka2
				17. Bd5+  Kb1
				18. Ba3 b4
				19. Bb3 bxa3
				20. Bg8 a2
				21. Bh7#`,
			expectedFEN: "8/7B/8/8/8/8/pp1K4/bk6 b - - 1 21",
			expectedErr: nil,
		},
		{
			fen: "8/8/4P3/k1N4r/8/K3N3/8/8 w - - 0 1",
			s: `1. e7 Rh8
				2. Ne6  Re8
				3. Nc7! Rxe7
				4. Nc4#!`,
			expectedFEN: "8/2N1r3/8/k7/2N5/K7/8/8 b - - 1 4",
			expectedErr: nil,
		},
		{
			fen: "8/5B2/p7/q5B1/k7/3K4/1P6/8 w - - 0 1",
			s: `1. b3+  Kb5
				2. Be8+ Kc5
				3. b4+  Qxb4
				4. Be7+`,
			expectedFEN: "4B3/4B3/p7/2k5/1q6/3K4/8/8 b - - 1 4",
			expectedErr: nil,
		},
		{
			fen: "2b4k/8/5Pr1/5N2/8/8/8/K1B5 w - - 0 1",
			s: `1. f7 Ra6+
				2. Ba3! Rxa3+!
				3. Kb2 Ra2+!
				4. Kc1! Ra1+
				5. Kd2 Ra2+
				6. Ke3 Ra3+
				7. Kf4 Ra4+
				8. Kg5 Rg4+!
				9. Kh6! Rg8!
				10. Ne7 Be6
				11. fxg8=Q+ Bxg8
				12. Ng6#!`,
			expectedFEN: "6bk/8/6NK/8/8/8/8/8 b - - 1 12",
			expectedErr: nil,
		},
		{
			fen: "N7/1p1kqP2/8/3P4/KP6/4B3/8/8 w - - 0 1",
			s: `1. Nb6+ Kc7
				2. f8=Q!  Qxf8
				3. d6+  Kc6
				4. b5+  Kxd6
				5. Bc5+ Kxc5
				6. Nd7+`,
			expectedFEN: "5q2/1p1N4/8/1Pk5/K7/8/8/8 b - - 1 6",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/5N2/p1RK4/pk6/8 w - - 0 1",
			s: `1. Rc2+ Kb3
				2. Rc1  a1=Q!
				3. Rxa1 Kb2
				4. Rf1! a2
				5. Kc4! a1=Q
				6. Nd3+ Ka2
				7. Nb4+ Kb2
				8. Rf2+ Kb1
				9. Kb3`,
			expectedFEN: "8/8/8/8/1N6/1K6/5R2/qk6 b - - 7 9",
			expectedErr: nil,
		},
		{
			fen: "2k5/6R1/N7/8/8/6K1/7p/5b2 w - - 0 1",
			s: `1. Rg8+ Kb7
				2. Nc5+ Kb6
				3. Na4+ Kb5
				4. Nc3+ Kb4
				5. Na2+ Kb3
				6. Nc1+ Kb2
				7. Kxh2 Kxc1
				8. Rg1`,
			expectedFEN: "8/8/8/8/8/8/7K/2k2bR1 b - - 1 8",
			expectedErr: nil,
		},
		{
			fen: "8/2k3Nr/8/3K4/1R6/8/8/8 w - - 0 1",
			s: `1. Ne8+ Kc8
				2. Nd6+ Kd8
				3. Rb8+ Kd7
				4. Rb7+ Kd8
				5. Nf7+ Kc8
				6. Kc6  Rg7
				7. Nd6+ Kd8
				8. Rb8+ Ke7
				9. Nf5+ Kf6
				10. Nxg7`,
			expectedFEN: "1R6/6N1/2K2k2/8/8/8/8/8 b - - 0 10",
			expectedErr: nil,
		},
		{
			fen: "3Nk3/r7/7R/8/3K4/8/8/8 w - - 0 1",
			s: `1. Rh8+ Kd7
				2. Rh7+ Kd6!
				3. Nf7+ Kc7
				4. Ne5+!  Kb6!
				5. Nc4+ Ka6
				6. Rh6+ Kb7!
				7. Nd6+ Kb8
				8. Rh8+ Kc7
				9. Nb5+`,
			expectedFEN: "7R/r1k5/8/1N6/3K4/8/8/8 b - - 17 9",
			expectedErr: nil,
		},
		{
			fen: "R7/8/8/8/4K2p/8/8/3k1N1r w - - 0 1",
			s: `1. Ra1+ Ke2
				2. Ng3+ hxg3
				3. Ra2+!  Kd1
				4. Kd3  Ke1
				5. Ke3  Kf1
				6. Kf3  Kg1
				7. Kxg3 Kf1
				8. Ra1+`,
			expectedFEN: "8/8/8/8/8/6K1/8/R4k1r b - - 2 8",
			expectedErr: nil,
		},
		{
			fen: "8/b7/8/8/7K/1R6/2n5/N2k4 w - - 0 1",
			s: `1. Rb1+ Kd2
				2. Rb2  Kc1
				3. Rxc2+  Kb1
				4. Rc4! Bf2+
				5. Kg4  Kxa1
				6. Kf3!  Bh4
				7. Rxh4`,
			expectedFEN: "8/8/8/8/7R/5K2/8/k7 b - - 0 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/1K2k1n1/1n6/5R2/8/4N3 w - - 0 1",
			s: `1. Re3+ Kd4
				2. Rg3  Ne4
				3. Ra3  Nd5
				4. Nf3#!`,
			expectedFEN: "8/8/8/1K1n4/3kn3/R4N2/8/8 b - - 7 4",
			expectedErr: nil,
		},
		{
			fen: "8/8/7n/3k4/B4p2/8/8/b2K2R1 w - - 0 1",
			s: `1. Bb3+ Ke5
				2. Kc2  Bd4
				3. Rg6  Nf5
				4. Re6#!`,
			expectedFEN: "8/8/4R3/4kn2/3b1p2/1B6/2K5/8 b - - 7 4",
			expectedErr: nil,
		},
		{
			fen: "7k/7p/B7/1K6/8/1n3b2/8/5R2 w - - 0 1",
			s: `1. Kb4  Nd4
				2. Kc3  Be2
				3. Bxe2 Nxe2+
				4. Kd3  Ng3
				5. Rf3  Nh5
				6. Rf5  Ng7
				7. Rf8#!`,
			expectedFEN: "5R1k/6np/8/8/8/3K4/8/8 b - - 7 7",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/5p2/8/3R4/P7/1k1K1B1r w - - 0 1",
			s: `1. a3!  Rxf1+
				2. Ke2  Rf4
				3. Rb3+ Kc2
				4. Rb4! Rxb4
				5. axb4`,
			expectedFEN: "8/8/8/5p2/1P6/8/2k1K3/8 b - - 0 5",
			expectedErr: nil,
		},
		{
			fen: "3N4/4R3/6k1/8/5K1P/8/5p2/8 w - - 0 1",
			s: `1. h5+  Kh6
				2. Nf7+ Kxh5
				3. Re5+ Kh4
				4. Ng5! f1=Q+
				5. Nf3+ Kh3
				6. Rh5+ Kg2
				7. Rh2#`,
			expectedFEN: "8/8/8/8/5K2/5N2/6kR/5q2 b - - 5 7",
			expectedErr: nil,
		},
		{
			fen: "q3N1R1/p7/8/8/4P2k/8/4K2P/8 w - - 0 1",
			s: `1. Rh8+ Kg5
				2. Nd6  Qc6
				3. Nf7+ Kf4
				4. Rh4#!`,
			expectedFEN: "8/p4N2/2q5/8/4Pk1R/8/4K2P/8 b - - 7 4",
			expectedErr: nil,
		},
		{
			fen: "3N4/2p5/8/3q4/3k4/8/3PKP2/R7 w - - 0 1",
			s: `1. Ra4+ Ke5
				2. Ra5  c5
				3. Rxc5 Qxc5
				4. d4+  Kxd4
				5. Ne6+`,
			expectedFEN: "8/8/4N3/2q5/3k4/8/4KP2/8 b - - 1 5",
			expectedErr: nil,
		},
		{
			fen: "5q2/2R5/5p2/1k3N2/2p1P3/2P5/2K5/8 w - - 0 1",
			s: `1. Rc8  Qa3
				2. Nd4+   Kb6
				3. Rb8+   Kc5
				4. Rb5+   Kd6
				5. Rd5+   Ke7
				6. Ra5! Qxa5
				7. Nc6+`,
			expectedFEN: "8/4k3/2N2p2/q7/2p1P3/2P5/2K5/8 b - - 1 7",
			expectedErr: nil,
		},
		{
			fen: "1q3k2/7K/1P1P1p2/R2P4/8/B5b1/8/8 w - - 0 1",
			s: `1. b7!   Kf7
				2. Ra8   Qxb7
				3. Rf8+!   Kxf8+
				4. d7+   Kf7
				5. d8=N+! Ke8+
				6. Nxb7`,
			expectedFEN: "4k3/1N5K/5p2/3P4/8/B5b1/8/8 b - - 0 6",
			expectedErr: nil,
		},
		{
			fen: "2K5/4p3/1p1p1p2/1p2k3/1P3q2/8/1PBPPP2/7R w - - 0 1",
			s: `1. Rh5+ Ke6!
				2. Bb3+ d5
				3. Rxd5 Qxf2
				4. d4!  f5
				5. e4 Qxb2
				6. Rd6+!  Kxd6
				7. e5+  Kc6
				8. d5#!`,
			expectedFEN: "2K5/4p3/1pk5/1p1PPp2/1P6/1B6/1q6/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "RK6/8/1k6/7R/8/8/pp6/8 w - - 0 1",
			s: `1. Ra5  Kc6
				2. Kc8  Kd6
				3. Kd8  Ke6
				4. Ra6+   Kf7
				5. Rf5+   Kg7
				6. Rg5+   Kf7
				7. Rg6  b1=Q
				8. Rf6# `,
			expectedFEN: "3K4/5k2/5RR1/8/8/8/p7/1q6 b - - 1 8",
			expectedErr: nil,
		},
		{
			fen: "8/8/8/8/7k/6pq/R1R5/4K3 w - - 0 1",
			s: `1. Ra4+   Kg5
				2. Rc5+   Kf6
				3. Ra6+   Ke7
				4. Rc7+   Kd8
				5. Rh7  Qg2
				6. Ra8+ Qxa8
				7. Rh8+`,
			expectedFEN: "q2k3R/8/8/8/8/6p1/8/4K3 b - - 1 7",
			expectedErr: nil,
		},
		{
			fen: "8/4r3/1P5K/8/1p2kr2/8/R7/R7 w - - 0 1",
			s: `1. Re1+   Kf5
				2. Rxe7   Kf6
				3. Rf7+!  Kxf7
				4. b7   Rf6+
				5. Kh7  Rb6
				6. Ra7  Ke6
				7. Ra6!   Rxa6
				8. b8=Q`,
			expectedFEN: "1Q6/7K/r3k3/8/1p6/8/8/8 b - - 0 8",
			expectedErr: nil,
		},
		{
			fen: "5kr1/7R/3K4/7R/2p5/2P1p1r1/8/8 w - - 0 1",
			s: `1. Rf5+ Ke8
				2. Re7+ Kd8
				3. Rd7+ Kc8
				4. Rc5+ Kb8
				5. Rb5+ Ka8
				6. Kc7  Rg3g5
				7. Rd8+ Ka7
				8. Rb7+ Ka6
				9. Rd6+ Ka5
				10. Ra7+  Kb5
				11. Rb6+  Kc5
				12. Ra5#`,
			expectedFEN: "6r1/2K5/1R6/R1k3r1/2p5/2P1p3/8/8 b - - 23 12",
			expectedErr: nil,
		},
		{
			fen: "1B6/8/5kP1/8/8/7K/8/1r1B2N1 w - - 0 1",
			s: `1. Be5+!  Kxe5
				2. g7 Rb8
				3. Bb3! Rxb3+
				4. Nf3+ Rxf3+
				5. Kg2`,
			expectedFEN: "8/6P1/8/4k3/8/5r2/6K1/8 b - - 1 5",
			expectedErr: nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test parser algebraic %v", i), func(t *testing.T) {
			g, err := newGameFromFEN(tc.fen)
			require.NoError(t, err)
			gameSteps, err := newParserAlgebraic(g, tc.s).parse()
			require.Equal(t, tc.expectedErr, err)
			if tc.expectedErr != nil {
				return
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
