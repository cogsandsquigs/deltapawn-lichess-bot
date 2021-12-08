package engine

import (
    "time"
    "strings"
    "math/rand"

    "github.com/notnil/chess"
)

type RandomBot struct {
    color string
    game  *chess.Game
}

func NewRandomBot() *RandomBot {
    return &RandomBot{
        game: chess.NewGame(),
    }
}

func (b *RandomBot) Name() string {
    return "RandomBot"
}

func (b *RandomBot) Color() string {
    return b.color
}

func (b *RandomBot) SetColor(color string) {
    b.color = color
}

func (b *RandomBot) Move(move string) {
    b.game.MoveStr(move)
}

func (b *RandomBot) GameMoves() string {
    var stringmoves []string
    for _, m := range b.game.Moves() {
        stringmoves = append(stringmoves, m.String())
    }

    return strings.Join(stringmoves, " ")
}

func (b *RandomBot) NextBestMove() string {
    rand.Seed(time.Now().Unix())
    moves := b.game.ValidMoves()
    randmove := moves[rand.Intn(len(moves))]
    return randmove.String()
}