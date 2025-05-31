package main

import (
	"github.com/urashav/event-hub/configs"
	"github.com/urashav/event-hub/internal/app"
	"log"
)

func main() {
	cfg, err := configs.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.App(cfg)
}
