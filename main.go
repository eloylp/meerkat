package main

import (
	"fmt"
	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/factory"
	"log"
)

var version string

func main() {
	fmt.Println(fmt.Sprintf("Meerkat %s", version))
	cfg := config.C()
	d := factory.NewApp(cfg)
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}
}
