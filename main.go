package main

import (
	"go-sentinel/app"
	"go-sentinel/config"
	"log"
)

func main() {
	cfg := config.C()
	d := app.NewApp(cfg)
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
}
