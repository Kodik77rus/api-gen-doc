package main

import (
	"github.com/Kodik77rus/api-gen-doc/internal/config"
	"github.com/Kodik77rus/api-gen-doc/internal/server"
	"log"
)

func main() {
	appConf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.New(&appConf.Server)
	if err != nil {
		log.Fatal(err)
	}

	server.Start(appConf.TemplateBuilder.TemplateFolder)
}
