package main

import (
	"log"
	"os"
	//"fmt"

	"github.com/ipratt-code/deltapawn-lichess-bot/engine"
	"github.com/ipratt-code/deltapawn-lichess-bot/lichessapi"
	"github.com/joho/godotenv"
)

var preferences = lichessapi.BotPreferences{
	Variants: []string{"standard"},
	Speeds:   []string{"rapid", "classical", "unlimited"},
	Modes:    []string{"rated", "casual"},
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("was not able to read environment variables from .env file: ", err)
	}
	server := lichessapi.NewLichessApi(os.Getenv("LICHESSAUTH"), preferences)
	bot := engine.NewRandomBot()
	log.Println("bot starting up!...")
	server.StreamEvent(bot)
}
