package engine

import (
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/opening"
)

type Deltapawn struct {
	color int
	game  *chess.Game
	depth int
	book  opening.Book
}

func NewDeltapawn() *Deltapawn {
	return &Deltapawn{
		game: chess.NewGame(
			chess.UseNotation(chess.UCINotation{}),
		),
		depth: searchDepth,
		book:  opening.NewBookECO(),
	}
}

func (b *Deltapawn) Name() string {
	return "Deltapawn"
}

func (b *Deltapawn) Color() int {
	return b.color
}

func (b *Deltapawn) New() ChessEngine {
	return NewDeltapawn()
}

func (b *Deltapawn) SetColor(color int) {
	if (color != 1 || color != -1) && color > 1 {
		b.color = 1
	} else if (color != 1 || color != -1) && color < 1 {
		b.color = -1
	} else {
		b.color = color
	}
}

func (b *Deltapawn) Move(move string) {
	if move == "" {
		return
	}

	err := b.game.MoveStr(move)
	if err != nil {
		log.Println(err)
	}
}

func (b *Deltapawn) GameMoves() string {
	var stringmoves []string
	for _, m := range b.game.Moves() {
		stringmoves = append(stringmoves, m.String())
	}

	return strings.Join(stringmoves, " ")
}

func (b *Deltapawn) NextBestMove() string {
	rand.Seed(time.Now().Unix())

	moves := b.game.ValidMoves()
	bestmove := moves[0]
	bestscore := -9999999999

	if len(b.game.Moves()) < 4 {
		return b.game.ValidMoves()[rand.Intn(len(b.game.ValidMoves()))].String()
	} else if len(b.game.Moves()) < 8 {
		for _, move := range moves {
			g := b.game.Clone()
			g.Move(move)
			v := pvs(g, math.MinInt64, math.MaxInt64, 0, 1)
			if v > bestscore {
				bestmove = move
			}
		}

		return bestmove.String()
	} else {
		for _, move := range moves {
			g := b.game.Clone()
			g.Move(move)
			v := pvs(g, math.MinInt64, math.MaxInt64, b.depth, 1)
			if v > bestscore {
				bestmove = move
			}
		}

		return bestmove.String()
	}
}

func (b *Deltapawn) IsGameOver() bool {
	return b.game.Outcome() != "*"
}

func (b *Deltapawn) Reset() {
	b.game = NewDeltapawn().game
}

func sortMoves(game *chess.Game, moves []*chess.Move, color int) []*chess.Move {
	sort.SliceStable(moves, func(i, j int) bool {
		g1 := game.Clone()
		g1.Move(moves[i])

		g2 := game.Clone()
		g2.Move(moves[j])

		b1 := g1.Position().Board()
		b2 := g2.Position().Board()

		return evalBoard(b1)*color > evalBoard(b2)*color
	})

	return moves
}

// evaluates the board in white's favor
func evalBoard(board *chess.Board) int {
	// define variables here

	score := 0

	for _, piece := range board.SquareMap() {
		m := -2*int(piece.Color()) + 3 // turns color to white: 1, black: -1
		switch piece.Type() {
		case chess.Queen:
			score += PieceValuation[0] * m
		case chess.Rook:
			score += PieceValuation[1] * m
		case chess.Bishop:
			score += PieceValuation[2] * m
		case chess.Knight:
			score += PieceValuation[3] * m
		case chess.Pawn:
			score += PieceValuation[4] * m
		default:
			score += 0
		}
	}
	return score
}
