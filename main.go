package main

import (
	"fmt"
	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/factory"
	"log"
)

var version string

func main() {
	fmt.Printf("Meerkat %s \n", version)
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
