package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/constants"
	"github.com/VirtusLab/render/render"
	"github.com/VirtusLab/render/version"
	"github.com/urfave/cli"
)

var (
	in     string
	out    string
	config string
	vars   cli.StringSlice
)

func main() {
	app := cli.NewApp()
	app.Name = constants.Name
	app.Usage = constants.Description
	app.Author = constants.Author
	app.Version = fmt.Sprintf("%s-%s", version.VERSION, version.GITCOMMIT)
	app.Before = preload
	app.Action = action

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "run in debug mode",
		},
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
			Name:  "set, var",
			Usage: "additional parameters in key=value format, can be used multiple times",
			Value: &vars,
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(cli.ErrWriter, "There is no %q command.\n", command)
		cli.OsExiter(1)
	}
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			return err
		}

		fmt.Fprintf(cli.ErrWriter, "WRONG: %v\n", err)
		return nil
	}
	cli.OsExiter = func(c int) {
		if c != 0 {
			logrus.Debugf("exiting with %d", c)
		}
		os.Exit(c)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(cli.ErrWriter, "ERROR: %v\n", err)
		cli.OsExiter(1)
	}
}

func preload(c *cli.Context) error {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	if c.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		logrus.Debug("Debug logging enabled")
	}

	return nil
}

func action(_ *cli.Context) error {
	err := render.Render(in, out, config, vars)
	if err != nil {
		logrus.Fatal("Rendering failed", err)
		return err
	}

	return nil
}