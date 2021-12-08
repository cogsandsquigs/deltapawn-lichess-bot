package main

import (
    "os"
    "log"
    //"fmt"

    "github.com/ipratt-code/deltapawn-lichess-bot/lichessapi"
    "github.com/ipratt-code/deltapawn-lichess-bot/engine"
)

var preferences = lichessapi.BotPreferences{
    Variants: []string{"standard"},
    Speeds: []string{"rapid", "classical", "unlimited"} ,
    Modes: []string{"rated", "casual"},
}

func main() {
    server := lichessapi.NewLichessApi(os.Getenv("LICHESSAUTH"), preferences)
    bot := engine.NewRandomBot()
    log.Println("bot starting up!...")
    server.StreamEvent(bot)
}