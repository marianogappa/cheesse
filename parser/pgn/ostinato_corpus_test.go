package pgn

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
	"github.com/marianogappa/cheesse/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOstinatoCorpus_PGN ports ostinato's "should parse pgn notation": the
// six-notation-suite French Defense game as a PGN with tag pairs.
func TestOstinatoCorpus_PGN(t *testing.T) {
	pgnString := `[Event "Ostinato Testing"]
[Site "Buenos Aires, Argentina"]
[Date "2015.??.??"]
[Round "1"]
[Result "1/2-1/2"]
[White "Fake Player 1"]
[Black "Fake Player 2"]

1. e4 e6 2. d4 d5 3. Nc3 Bb4 4. Bb5+ Bd7 5. Bxd7+ Qxd7 6. Nge2
dxe4 7. 0-0
`
	g := core.NewDefaultGame()
	parsed, err := parser.NewGenericNotationParser(NewVariantPGN()).Parse(g, pgnString)
	require.NoError(t, err)
	require.Len(t, parsed.GameSteps, 13)

	expectedFEN := "rn2k1nr/pppq1ppp/4p3/8/1b1Pp3/2N5/PPP1NPPP/R1BQ1RK1 b kq - 1 7"
	assert.Equal(t, expectedFEN, parsed.GameSteps[len(parsed.GameSteps)-1].StepGame.ToFEN())
	assert.Equal(t, "Ostinato Testing", parsed.Metadata["Event"])
	assert.Equal(t, "Buenos Aires, Argentina", parsed.Metadata["Site"])
}

// TestOstinatoCorpus_EverittWang ports ostinato's "should parse final actions": the
// complete 54-move Everitt-Wang 1997 correspondence game with 12 PGN headers,
// ending 0-1 (White's loss).
func TestOstinatoCorpus_EverittWang(t *testing.T) {
	pgnString := `[Event "1997 NAPZ/ MPP(F) M-01"]
[Site "ICCF"]
[Date "1997.??.??"]
[Round "?"]
[White "Everitt, Gordon T. (USA)"]
[Black "Wang, Mong Lin (SIN)"]
[Result "0-1"]
[ECO "A57"]
[WhiteElo "2336"]
[BlackElo "2428"]
[PlyCount "108"]
[EventDate "1997.??.??"]

1. d4 Nf6 2. c4 c5 3. d5 b5 4. cxb5 a6 5. e3 axb5 6. Bxb5 Qa5+ 7. Nc3 Bb7 8.
Nge2 Bxd5 9. O-O Bc6 10. a4 e6 11. Ng3 d5 12. Bd2 Qd8 13. e4 d4 14. Bxc6+ Nxc6
15. Nb5 Be7 16. Qc2 O-O 17. Rfc1 Qb6 18. Na3 Nd7 19. Nc4 Qa6 20. a5 Nde5 21.
Nb6 Ra7 22. b3 Qb5 23. Nc4 Rfa8 24. Rcb1 Bd8 25. Ra4 Nd7 26. Na3 Qb7 27. b4
cxb4 28. Bxb4 Qa6 29. Bd2 Rc8 30. Qc4 d3 31. Qxa6 Rxa6 32. Nc4 Nc5 33. Ra3 Be7
34. e5 Nd4 35. f3 Ncb3 36. Raxb3 Nxb3 37. Rxb3 Rxc4 38. Rb8+ Bf8 39. Ne4 Rac6
40. Rd8 Rc2 41. h3 Rc8 42. Rxd3 Ra2 43. Kh2 h6 44. Rd7 g5 45. Nf6+ Kg7 46. h4
Rcc2 47. Ne4 gxh4 48. Kh3 Kg6 49. Rd3 Be7 50. Rd4 Bg5 51. f4 Be7 52. Rd3 Kf5
53. Rd4 Bc5 54. Nxc5 Rxd2 0-1
`
	g := core.NewDefaultGame()
	parsed, err := parser.NewGenericNotationParser(NewVariantPGN()).Parse(g, pgnString)
	require.NoError(t, err)

	// 108 half-moves plus the result marker
	require.Len(t, parsed.GameSteps, 109)
	assert.Equal(t, "0-1", parsed.GameSteps[len(parsed.GameSteps)-1].StepString)
	assert.Equal(t, "Everitt, Gordon T. (USA)", parsed.Metadata["White"])
	assert.Equal(t, "Wang, Mong Lin (SIN)", parsed.Metadata["Black"])
	assert.Equal(t, "0-1", parsed.Metadata["Result"])
}
