package main

import (
	"log/slog"

	"github.com/PCoelho07/gochat/server"
)

func main() {
	err := server.Start(":7007")
	if err != nil {
		slog.Error(err.Error())
	}
}
