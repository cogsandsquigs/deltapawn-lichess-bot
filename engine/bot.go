package engine

type ChessEngine interface{
    Name() string
    Color() string
    SetColor(color string)
    Move(move string)
    GameMoves() string // returns all the game's moves
    NextBestMove() string // takes in a move, returns a new move in response
}