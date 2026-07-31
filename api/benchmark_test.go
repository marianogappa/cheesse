package api

import (
	"testing"
)

func BenchmarkParseNotation_50MoveGame(b *testing.B) {
	game := `1. d4 Nf6 2. c4 c5 3. d5 b5 4. cxb5 a6 5. e3 axb5 6. Bxb5 Qa5+ 7. Nc3 Bb7 8.
Nge2 Bxd5 9. O-O Bc6 10. a4 e6 11. Ng3 d5 12. Bd2 Qd8 13. e4 d4 14. Bxc6+ Nxc6
15. Nb5 Be7 16. Qc2 O-O 17. Rfc1 Qb6 18. Na3 Nd7 19. Nc4 Qa6 20. a5 Nde5 21.
Nb6 Ra7 22. b3 Qb5 23. Nc4 Rfa8 24. Rcb1 Bd8 25. Ra4`

	a := New()
	ig := InputGame{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.ParseNotation(ig, game)
	}
}

func BenchmarkConvertNotation_ToICCF(b *testing.B) {
	game := `1. e4 e6 2. d4 d5 3. Nc3 Bb4 4. Bb5+ Bd7 5. Bxd7+ Qxd7 6. Ne2 dxe4 7. 0-0`
	a := New()
	ig := InputGame{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.ConvertNotation(ig, game, "ICCF")
	}
}
