package main

import (
	"log/slog"

	"github.com/PCoelho07/gochat/client"
)

func main() {
    err := client.Start(":7007")
    if err != nil { 
        slog.Error(err.Error())
    }
}
