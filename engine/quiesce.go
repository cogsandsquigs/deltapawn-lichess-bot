package engine

import (
	"log"

	"github.com/notnil/chess"
)

// implementation of quiescence search
func quiesce(game *chess.Game, alpha, beta, color int) int {
	curmove := game.Moves()[len(game.Moves())-1]
	if curmove.HasTag(chess.Capture) || curmove.HasTag(chess.Check) {
		moves := sortMoves(game, game.ValidMoves(), color)
		if len(moves) > 5 {
			history := game.Moves()
			historystr := ""

			for _, move := range history {
				historystr += move.String() + " "
			}

			log.Printf("quiescing move path %s", historystr)
			return evalBoard(game.Position().Board())
		}
		bSearchPv := true
		for _, move := range moves {
			score := 0
			g := game.Clone()

			g.Move(move)

			history := g.Moves()
			historystr := ""

			for _, move := range history {
				historystr += move.String() + " "
			}

			log.Printf("quiescing move path %s", historystr)

			if bSearchPv {
				score = -quiesce(g, -beta, -alpha, -color)
			} else {
				score = -quiesce(g, -alpha-1, -alpha, -color)

				// in fail-soft ... && score < beta ) is common
				if score > alpha {
					score = -quiesce(g, -beta, -alpha, -color) // re-search
				}
			}

			if score >= beta {
				return beta // fail-hard beta-cutoff
			}

			if score > alpha {
				alpha = score     // alpha acts like max in MiniMax
				bSearchPv = false // *1)
			}
		}

		return alpha
	} else {
		history := game.Moves()
		historystr := ""

		for _, move := range history {
			historystr += move.String() + " "
		}

		log.Printf("quiescing move path %s", historystr)
		return evalBoard(game.Position().Board())
	}

}
