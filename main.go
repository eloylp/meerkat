package main

import (
	"fmt"
	"github.com/eloylp/meerkat/app"
	"github.com/eloylp/meerkat/config"
	"log"
)

var version string

func main() {
	fmt.Println(fmt.Sprintf("Meerkat %s", version))
	cfg := config.C()
	d := app.NewApp(cfg)
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
}
