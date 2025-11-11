package main

import (
	"log"

	"app/config"
	"app/internal"
)

func main() {
	cfg := config.Load(".")
	s := internal.CreateApp(cfg)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
