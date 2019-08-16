package main

import (
	"log"
	"meerkat/app"
	"meerkat/config"
)

func main() {
	cfg := config.C()
	d := app.NewApp(cfg)
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
}
