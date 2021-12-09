package engine

type ChessEngine interface{
    Name() string
    Color() string
    Reset() // reset's the engine's internal state
    SetColor(color string)
    Move(move string)
    NextBestMove() string // takes in a move, returns a new move in response
    GameMoves() string // returns all the game's moves
    IsGameOver() bool
}