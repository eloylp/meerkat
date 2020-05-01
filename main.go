package main

import (
	"fmt"
	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/factory"
	"log"
)

var (
	Name      = "meerkat"
	Version   = "vx.y.z"
	Build     = "a1234"
	BuildTime = "a1234"
)

func main() {

	fmt.Printf("%s %s - Build: %s at %s \n", Name, Version, Build, BuildTime)
	cfg := config.FromArguments()
	d, err := factory.NewHTTPServedApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("gracefully shutted down server")
}
