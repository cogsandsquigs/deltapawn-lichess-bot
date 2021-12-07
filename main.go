package main

import (
    "os"
    "fmt"

    "github.com/ipratt-code/deltapawn-lichess-bot/lichessapi"
)

func main() {
    server := lichessapi.NewLichessApi(os.Getenv("LICHESSAUTH"))
    u := server.GetAccountData()
    fmt.Println(u)
}