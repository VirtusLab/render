package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/constants"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer"
	"github.com/VirtusLab/render/renderer/configuration"
	"github.com/VirtusLab/render/version"
	"gopkg.in/urfave/cli.v1"
)

var (
	app         *cli.App
	inputPath   string
	outputPath  string
	configPaths cli.StringSlice
	vars        cli.StringSlice
)

func main() {
	app = cli.NewApp()
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
			Usage:       "the input template file, stdin if empty",
			Destination: &inputPath,
		},
		cli.StringFlag{
			Name:        "out",
			Value:       "",
			Usage:       "the output file, stdout if empty",
			Destination: &outputPath,
		},
		cli.StringSliceFlag{
			Name:  "config",
			Usage: "optional configuration YAML file, can be used multiple times",
			Value: &configPaths,
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

	if len(c.Args()) == 0 {
		return nil
	}

	if c.Args()[0] == "help" {
		return nil
	}

	return nil
}

func action(_ *cli.Context) error {
	config, err := configuration.All(configPaths, vars)
	if err != nil {
		return err
	}

	r := renderer.New(config)
	err = r.RenderFile(inputPath, outputPath)
	if err != nil {
		if err == files.ErrExpectedStdin {
			return fmt.Errorf("expected either stdin or --in parameter, for usage use --help")
		}
		logrus.Debugf("Rendering failed: %s", err)
		return err
	}

	return nil
}
