package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/eloylp/meerkat/config"
	"github.com/eloylp/meerkat/factory"
)

var (
	Name      = "meerkat"
	Version   = "vx.y.z"
	Build     = "a1234"
	BuildTime = "a1234"
)

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	_, err := fmt.Fprintf(stdout, "%s %s - Build: %s at %s \n", Name, Version, Build, BuildTime)
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(stderr)
	log.SetOutput(stdout)
	resources := fs.String("u", "", "The comma separated URLS to recover frames from")
	pollInterval := fs.Int("i", 1, "The interval to fill the frame buffer")
	listenAddress := fs.String("l", "0.0.0.0:3000", "Pass the http server listen address for serving results")
	if err := fs.Parse(args[1:]); err != nil { //nolint:govet
		return err
	}
	cfg := config.Config{
		Resources:         config.ParseResources(*resources),
		PollInterval:      *pollInterval,
		HTTPListenAddress: *listenAddress,
	}
	if err = config.Validate(cfg); err != nil {
		fs.PrintDefaults()
		return err
	}
	d, err := factory.NewHTTPServedApp(cfg)
	if err != nil {
		return err
	}
	if err := d.Start(); err != nil {
		return err
	}
	log.Println("gracefully shut down server")
	return nil
}
