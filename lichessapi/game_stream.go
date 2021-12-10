package lichessapi

import (
	"log"
	//"time"
	"encoding/json"
	"strings"

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

func (s *LichessApi) gameStreamWrapper(c chan bool, f func(string, engine.ChessEngine), str string, eng engine.ChessEngine) {
	g := make(chan bool)
	go func() {
		f(str, eng)
		g <- true
	}()
	for {
		select {
		case <-c:
		case <-g:
			return
		}
	}

}

func (s *LichessApi) streamGame(gameId string, eng engine.ChessEngine) {
	resp, _ := s.request("GET", "bot/game/stream/"+gameId)
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
			move := moves[len(moves)-1]

			if eng.Color() == "white" && len(moves)%2 == 0 {
				s.runEngine(gameId, move, eng)
			} else if eng.Color() == "black" && len(moves)%2 == 1 {
				s.runEngine(gameId, move, eng)
			}
		} else if eng.Color() == "white" {
			s.runEngine(gameId, "", eng)
		}
	}
}

func (s *LichessApi) runEngine(gameId, move string, eng engine.ChessEngine) {
	log.Println("got move: " + move)
	eng.Move(move)
	if eng.IsGameOver() {
		eng.Reset()
		return
	}
	nextMove := eng.NextBestMove()
	eng.Move(nextMove)
	err := s.makeMove(gameId, nextMove)
	if err != nil {
		_, err = s.request("POST", "bot/game/"+gameId+"/abort")
		if err != nil {
			s.request("POST", "bot/game/"+gameId+"/resign")
		}
	}
}

func (s *LichessApi) makeMove(gameId, move string) error {
	log.Println("REQUEST", "bot/game/"+gameId+"/move/"+move)
	if move == "(none)" {
		return nil
	}
	resp, err := s.request("POST", "bot/game/"+gameId+"/move/"+move)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
