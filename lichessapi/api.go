package lichessapi

import (
    "log"
    "net/http"
)

type LichessApi struct {
    Challenge Challenge
    authtkn string
    gamesInProgress int
}

func NewLichessApi(auth string) *LichessApi {
    return &LichessApi{
        authtkn: auth,
    }
}

func (s *LichessApi) request(method, path string) *http.Response {
	var client = &http.Client{}
    req, err := http.NewRequest(method, "https://lichess.org/api"+path, nil)
	req.Header.Add("Authorization", "Bearer "+s.authtkn)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed request "+method+" "+path, err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Response %d %s %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), method, path)
	}
	return resp
}