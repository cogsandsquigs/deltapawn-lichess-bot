package engine

import (
    "log"
    "time"
    "sort"
    "math"
    "strings"
    "math/rand"

    "github.com/notnil/chess"
)

type LookaheadBot struct {
    color string
    game  *chess.Game
}

func NewLookaheadBot() *LookaheadBot {
    return &LookaheadBot{
        game: chess.NewGame(
            chess.UseNotation(chess.UCINotation{}),
        ),
    }
}

func (b *LookaheadBot) Name() string {
    return "LookaheadBot"
}

func (b *LookaheadBot) Color() string {
    return b.color
}

func (b *LookaheadBot) SetColor(color string) {
    b.color = color
}

func (b *LookaheadBot) Move(move string) {
    if move == "" {
        return
    } 
    
    err := b.game.MoveStr(move)
    if err != nil {
        log.Println(err)
    }
}

func (b *LookaheadBot) GameMoves() string {
    var stringmoves []string
    for _, m := range b.game.Moves() {
        stringmoves = append(stringmoves, m.String())
    }

    return strings.Join(stringmoves, " ")
}

func (b *LookaheadBot) NextBestMove() string {
    rand.Seed(time.Now().Unix())
    
    moves := b.game.ValidMoves()
    bestmove := moves[0]
    bestscore := -9999999999.0
    for _, move := range moves {
        g := b.game.Clone()
        g.Move(move)
        v := b.negamax(g, 0, 1)
        if v > bestscore {
            bestmove = move
        }
    }
    
    return bestmove.String()
}

func (b *LookaheadBot) IsGameOver() bool {
    return b.game.Outcome() != "*"
}

func (b *LookaheadBot) Reset() {
    b.game = NewLookaheadBot().game
}

func (b *LookaheadBot) negamax(game *chess.Game, depth, color int) float64 {
    switch depth {
    case 0:
        return b.evalBoard(game.Position().Board())
    default:
        moves := b.sortMoves(game, game.ValidMoves(), color)
        best := -99999999999.0
        for _, move := range moves {
            g := game.Clone()
            g.Move(move)
            log.Printf("evaluating move %s with depth of %d", move.String(), depth)
            best = math.Max(best, b.negamax(g, depth - 1, -color))
        }

        return best
    }
    
}

func (b *LookaheadBot) sortMoves(game *chess.Game, moves []*chess.Move, color int) []*chess.Move {
    sort.SliceStable(moves, func(i, j int) bool {
        g1 := game.Clone()
        g1.Move(moves[i])
        
        g2 := game.Clone()
        g2.Move(moves[j])
        
        b1 := g1.Position().Board()
        b2 := g2.Position().Board()

        return b.evalBoard(b1) * float64(color) > b.evalBoard(b2) * float64(color)
    })

    return moves
} 

func (b *LookaheadBot) evalBoard(board *chess.Board) float64 {
    score := 0.0
    for _, piece := range board.SquareMap() {
        var m float64
        if piece.Color() == 1 {
            m = 1
        } else {
            m = -1
        }
        switch piece.Type() {
        case chess.King:
            score += dummyBotPieceValuation[0] * m
        case chess.Queen:
            score += dummyBotPieceValuation[1] * m
        case chess.Rook:
            score += dummyBotPieceValuation[2] * m
        case chess.Bishop:
            score += dummyBotPieceValuation[3] * m
        case chess.Knight:
            score += dummyBotPieceValuation[4] * m
        case chess.Pawn:
            score += dummyBotPieceValuation[5] * m
        default:
            score += 0
        }
    }
    return score
}