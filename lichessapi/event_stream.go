package lichessapi

import (
	"log"
	//"time"
	"encoding/json"

	"github.com/ipratt-code/deltapawn-lichess-bot/engine"
)

type Event struct {
	Type      string    `json:"type"`
	Challenge Challenge `json:"challenge"`
	Game      Game      `json:"game"`
}

type Game struct {
	Id string `json:"id"`
}

type Challenge struct {
	Id         string `json:"id"`
	Status     string `json:"status"`
	Challenger struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Rating      int    `json:"rating"`
		Provisional bool   `json:"provisional"`
		Online      bool   `json:"online"`
		Lag         int    `json:"lag"`
	}
	DestUser struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Rating      int    `json:"rating"`
		Provisional bool   `json:"provisional"`
		Online      bool   `json:"online"`
		Lag         int    `json:"lag"`
	}
	Variant struct {
		Key string `json:"key"`
	}
	Rated bool   `json:"rated"`
	Speed string `json:"speed"`
	Color string `json:"color"`
}

func (s *LichessApi) StreamEvent(eng engine.ChessEngine) {
	resp, _ := s.request("GET", "stream/event")
	dec := json.NewDecoder(resp.Body)

	for dec.More() {
		var e Event
		err := dec.Decode(&e)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("got event: " + e.Type)
			s.handleEvent(&e, eng)
		}
		//time.Sleep(time.Second)
	}
}

func (s *LichessApi) handleEvent(e *Event, eng engine.ChessEngine) {
	switch e.Type {
	case "challenge":
		s.handleChallengeEvent(e)
	case "gameStart":
		log.Println("starting game")
		c := make(chan bool)
		go s.gameStreamWrapper(c, s.streamGame, e.Game.Id, eng.New())
		s.gameKillChans[e.Game.Id] = c
		s.gamesInProgress += 1
	case "gameFinish":
        s.gamesInProgress -= 1
		log.Println("ending game against")
	default:
		log.Printf("Unhandled Event %v\n", e.Type)
	}
}

func (s *LichessApi) handleChallengeEvent(e *Event) {
	challengeId := e.Challenge.Id
	if s.validChallenge(&e.Challenge) && s.gamesInProgress < 1 {
		log.Println("Accepting challenge", "https://lichess.org/"+e.Challenge.Id)
		resp, _ := s.request("POST", "challenge/"+challengeId+"/accept")
		resp.Body.Close()
	} else {
		log.Println("Declining challenge", "https://lichess.org/"+e.Challenge.Id)
		resp, _ := s.request("POST", "challenge/"+challengeId+"/decline")
		resp.Body.Close()
	}
}

func (s *LichessApi) validChallenge(c *Challenge) bool {
	return c.Status == "created" &&
		c.Challenger.Online &&
		includes(s.Challenge.Variants, c.Variant.Key) &&
		includes(s.Challenge.Speeds, c.Speed) &&
		(c.Rated && includes(s.Challenge.Modes, "rated") ||
			!c.Rated && includes(s.Challenge.Modes, "casual"))
}
