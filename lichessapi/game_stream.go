package lichessapi

import (
    "log"
    //"time"
    "strings"
    "encoding/json"

    "github.com/ipratt-code/deltapawn-lichess-bot/engine"
)

type gameState struct {
	Type    string
	Id      string
	Rated   bool
	Variant struct {
		Key string
	}
	Clock struct {
		Initial   int
		Increment int
	}
	Speed      string
	InitialFen string
	State      struct {
		Type  string
		Moves string
		Wtime int
		Btime int
		Winc  int
		Binc  int
	}
	White struct {
		Id string
	}
	Black struct {
		Id string
	}
	Moves string
	Wtime int
	Btime int
	Winc  int
	Binc  int
}

func (s *LichessApi) streamGame(gameId string, eng engine.ChessEngine) {
    resp := s.request("GET", "bot/game/stream/"+gameId)
	dec := json.NewDecoder(resp.Body)

    for dec.More() {
        var gS gameState
		
        err := dec.Decode(&gS)

		
        if err != nil {
			log.Println(err)
		}

        if gS.Type == "gameFull" {
            if gS.White.Id == "deltapawn" {
                eng.SetColor("white")
            } else {
                eng.SetColor("black")
            }
        }

        moves := strings.Split(gS.Moves, " ")
        
        if moves[0] != "" {
            move := moves[len(moves) - 1]
            
            log.Println("got move: " + move)            
            log.Println()
            if eng.Color() == "white" && len(moves) % 2 == 0 {
                s.runEngine(gameId, move, eng)
            
            } else if eng.Color() == "black" && len(moves) % 2 == 1 {
                s.runEngine(gameId, move, eng)
            
            } else {
                //eng.Move(move)
            }
        } else if eng.Color() == "white" {
            s.runEngine(gameId, "", eng)
        }
    }
}

func (s *LichessApi) runEngine(gameId, move string, eng engine.ChessEngine) {
    eng.Move(move)
    nextMove := eng.NextBestMove()
    eng.Move(move)
    s.makeMove(gameId, nextMove)
}

func (s *LichessApi) makeMove(gameId, move string) {
	log.Println("REQUEST", "bot/game/"+gameId+"/move/"+move)
	if move == "(none)" {
		return
	}
	resp := s.request("POST", "bot/game/"+gameId+"/move/"+move)
	resp.Body.Close()
}