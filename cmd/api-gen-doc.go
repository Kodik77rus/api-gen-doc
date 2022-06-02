package main

import (
	"log"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
	"github.com/Kodik77rus/api-gen-doc/internal/server"
)

func main() {
	appConf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.New(&appConf.Server)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())
}
