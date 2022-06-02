package main

import (
	"fmt"
	"log"

	"github.com/Kodik77rus/api-gen-doc/internal/config"
)

func main() {
	appConf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(appConf)
}
