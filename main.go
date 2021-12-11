package main

import (
    "fmt"
	"log"
	"net/http"
	"os"

	"github.com/ipratt-code/deltapawn-lichess-bot/engine"
	"github.com/ipratt-code/deltapawn-lichess-bot/lichessapi"
	"github.com/joho/godotenv"
)

var preferences = lichessapi.BotPreferences{
	Variants: []string{"standard"},
	Speeds:   []string{"rapid", "classical", "correspondence"},
	Modes:    []string{"rated", "casual"},
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("was not able to read environment variables from .env file: ", err)
	}
	server := lichessapi.NewLichessApi(os.Getenv("LICHESSAUTH"), preferences)
	bot := engine.NewDeltapawn()
	log.Println("bot starting up!...")
	go server.StreamEvent(bot)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("pinged!")
        fmt.Fprintf(w, "hello world!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
