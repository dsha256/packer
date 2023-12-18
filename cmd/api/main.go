package main

import (
	"log/slog"

	"github.com/dsha256/packer/internal/packer"
	"github.com/dsha256/packer/internal/server"
)

const restAPIPort = "8080"

func main() {
	bootstrap()
}

func bootstrap() {
	newSizerSrvc := packer.NewSizerService(packer.SortedSizes)
	newPackerSrvc := packer.NewPacketsService(newSizerSrvc)

	newServer := server.NewServer(newSizerSrvc, newPackerSrvc)
	err := newServer.Serve(restAPIPort)
	if err != nil {
		slog.Error(err.Error())
	}
}
