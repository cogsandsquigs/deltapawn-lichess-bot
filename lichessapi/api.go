package lichessapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type BotPreferences struct {
	Variants []string // done
	Speeds   []string // done
	Modes    []string // done
}

type LichessApi struct {
	Challenge       BotPreferences
	authtkn         string
	gamesInProgress int
	gameKillChans   map[string]chan bool
}

func NewLichessApi(auth string, botpreferences BotPreferences) *LichessApi {
	return &LichessApi{
		Challenge:       botpreferences,
		authtkn:         auth,
		gamesInProgress: 0,
		gameKillChans:   make(map[string]chan bool),
	}
}

func (s *LichessApi) request(method, path string) (*http.Response, error) {
	var client = &http.Client{}
	req, err := http.NewRequest(method, "https://lichess.org/api/"+path, nil)
	req.Header.Add("Authorization", "Bearer "+s.authtkn)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed request "+method+" "+path, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		s, _ := io.ReadAll(resp.Body)
		log.Printf("Response %d %s %s %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), s, method, path)
		var ret map[string]string
		json.Unmarshal(s, &ret)
		return resp, fmt.Errorf(ret["error"])
	}
	return resp, nil
}
