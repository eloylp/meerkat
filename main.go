package main

import (
	"github.com/eloylp/meerkat/app"
	"github.com/eloylp/meerkat/config"
	"log"
)

func main() {
	cfg := config.C()
	d := app.NewApp(cfg)
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
}
