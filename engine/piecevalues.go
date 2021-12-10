package engine

import "github.com/notnil/chess"

var PieceValuation = map[chess.PieceType]int{
	chess.King:   20000,
	chess.Queen:  900,
	chess.Rook:   500,
	chess.Bishop: 330,
	chess.Knight: 320,
	chess.Pawn:   100,
}
