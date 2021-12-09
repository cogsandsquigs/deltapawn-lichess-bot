package engine

import (
    "log"
    "sort"
    "time"
    "strings"
    "math/rand"

    "github.com/notnil/chess"
)

type DummyBot struct {
    color string
    game  *chess.Game
    currentMoveList []*chess.Move
}

var dummyBotPieceValuation = []float64{9999, 9, 5, 3.5, 3, 1} // from king to pawn

func NewDummyBot() *DummyBot {
    return &DummyBot{
        game: chess.NewGame(
            chess.UseNotation(chess.UCINotation{}),
        ),
    }
}

func (b *DummyBot) Name() string {
    return "DummyBot"
}

func (b *DummyBot) Color() string {
    return b.color
}

func (b *DummyBot) SetColor(color string) {
    b.color = color
}

func (b *DummyBot) Move(move string) {
    if move == "" {
        return
    } 
    
    err := b.game.MoveStr(move)
    if err != nil {
        log.Println(err)
    }
}

func (b *DummyBot) GameMoves() string {
    var stringmoves []string
    for _, m := range b.game.Moves() {
        stringmoves = append(stringmoves, m.String())
    }

    return strings.Join(stringmoves, " ")
}

func (b *DummyBot) NextBestMove() string {
    var move *chess.Move
    b.currentMoveList = b.game.ValidMoves()
    sort.Sort(b)
    if evalPiece(b.game.Position().Board().Piece(b.currentMoveList[0].S2())) == 0 {
        rand.Seed(time.Now().Unix())
        move = b.currentMoveList[rand.Intn(len(b.currentMoveList))]
    } else {
        move = b.currentMoveList[0]
    }
    
    return move.String()
}

func (b *DummyBot) IsGameOver() bool {
    return b.game.Outcome() != "*"
}

func (b *DummyBot) Reset() {
    b.game = NewDummyBot().game
    b.currentMoveList = []*chess.Move{}
}


func (b *DummyBot) Len() int {
    m := b.currentMoveList
    return len(m)
}

func (b *DummyBot) Less(i, j int) bool {
    m := b.currentMoveList
    board := b.game.Position().Board()
    m1, m2 := m[i], m[j]
    p1, p2 := board.Piece(m1.S2()), board.Piece(m2.S2())
    return evalPiece(p1) > evalPiece(p2)

}

func (b *DummyBot) Swap(i, j int) {
    
    b.currentMoveList[i], b.currentMoveList[j] = b.currentMoveList[j], b.currentMoveList[i]
}

func evalPiece(piece chess.Piece) float64 {
    switch piece.Type() {
    case chess.King:
        return dummyBotPieceValuation[0]
    case chess.Queen:
        return dummyBotPieceValuation[1]
    case chess.Rook:
        return dummyBotPieceValuation[2]
    case chess.Bishop:
        return dummyBotPieceValuation[3]
    case chess.Knight:
        return dummyBotPieceValuation[4]
    case chess.Pawn:
        return dummyBotPieceValuation[5]
    default:
        return 0
    }
}