package main

import (
	"os"
	"github.com/urfave/cli"
	"log"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "render"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Usage = "Simple go-template files render"

	var in string
	var out string
	var config string
	var extraParams cli.StringSlice

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "in",
			Value:       "",
			Usage:       "the template file",
			Destination: &in,
		},
		cli.StringFlag{
			Name:        "out",
			Value:       "",
			Usage:       "the output file, stdout if empty",
			Destination: &out,
		},
		cli.StringFlag{
			Name:        "config",
			Value:       "",
			Usage:       "the config file",
			Destination: &config,
		},
		cli.StringSliceFlag{
			Name:  "set",
			Usage: "an additional parameters in key=value format",
			Value: &extraParams,
		},
	}

	app.Action = func(c *cli.Context) error {
		err := Render(in, out, config, extraParams)
		if err != nil {
			log.Fatal("Rendering failed", err)
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
