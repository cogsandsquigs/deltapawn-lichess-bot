package engine

import (
	"log"
	"math"

	"github.com/notnil/chess"
)

// implementation of principal valuation search
func pvs(game *chess.Game, alpha, beta, depth, color int) int {
	switch depth {
	case 0:
		log.Printf("evaluating move %s with depth of %d", game.Moves()[len(game.Moves())-1], depth)
		return quiesce(game, math.MinInt64, math.MaxInt64, -color) // evalBoard(game.Position().Board())
	default:
		moves := sortMoves(game, game.ValidMoves(), color)
		bSearchPv := true
		for _, move := range moves {
			score := 0
			g := game.Clone()

			g.Move(move)

			log.Printf("evaluating move %s with depth of %d", move.String(), depth)

			if bSearchPv {
				score = -pvs(g, -beta, -alpha, depth-1, -color)
			} else {
				score = -pvs(g, -alpha-1, -alpha, depth-1, -color)

				// in fail-soft ... && score < beta ) is common
				if score > alpha {
					score = -pvs(g, -beta, -alpha, depth-1, -color) // re-search
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
	}

}
