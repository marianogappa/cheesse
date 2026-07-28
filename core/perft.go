package core

// Perft counts the number of leaf nodes in the move tree at the given depth.
// It is the standard method for validating move generation correctness in chess engines.
// https://www.chessprogramming.org/Perft
func Perft(g Game, depth int) int {
	if depth == 0 {
		return 1
	}
	nodes := 0
	for _, a := range g.Actions {
		if a.IsResign {
			continue
		}
		nodes += Perft(g.DoAction(a), depth-1)
	}
	return nodes
}
