package main

import (
	"log/slog"

	"github.com/dsha256/packer/internal/server"
)

const restAPIPort = "8080"

func main() {
	bootstrap()
}

func bootstrap() {
	newServer := server.NewServer()
	err := newServer.Serve(restAPIPort)
	if err != nil {
		slog.Error(err.Error())
	}
}
