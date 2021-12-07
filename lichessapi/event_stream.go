package lichessapi

import (
    "log"
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


func (s *LichessApi) handleChallengeEvent(e *Event) {
	challengeId := e.Challenge.Id
	if s.validChallenge(&e.Challenge) && s.gamesInProgress < 1 {
		log.Println("Accepting challenge", e.Challenge)
		resp := s.request("POST", "challenge/"+challengeId+"/accept")
		resp.Body.Close()
	} else {
		log.Println("Declining challenge", e.Challenge)
		resp := s.request("POST", "challenge/"+challengeId+"/decline")
		resp.Body.Close()
	}
}

func (s *LichessApi) validChallenge(c *Challenge) bool {
    return c.Status == "created" &&
		c.Challenger.Online == true &&
		includes(s.Challenge.Variants, c.Variant.Key) &&
		includes(s.Challenge.Speeds, c.Speed) &&
		(c.Rated == true && includes(s.Challenge.Modes, "rated") ||
			c.Rated == false && includes(s.Challenge.Modes, "casual"))
}