package main

import (
	"fmt"
	"os"

	"github.com/VirtusLab/render/constants"
	"github.com/VirtusLab/render/renderer"
	"github.com/VirtusLab/render/renderer/parameters"
	"github.com/VirtusLab/render/version"

	"github.com/VirtusLab/go-extended/pkg/files"
	"github.com/sirupsen/logrus"
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
		_, err := fmt.Fprintf(cli.ErrWriter, "There is no %q command.\n", command)
		if err != nil {
			logrus.Errorf("Unexpected error: %s", err)
		}
		cli.OsExiter(1)
	}
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			return err
		}

		_, err = fmt.Fprintf(cli.ErrWriter, "WRONG: %v\n", err)
		if err != nil {
			logrus.Errorf("Unexpected error: %s", err)
		}
		return nil
	}
	cli.OsExiter = func(c int) {
		if c != 0 {
			logrus.Debugf("exiting with %d", c)
		}
		os.Exit(c)
	}

	if err := app.Run(os.Args); err != nil {
		_, err := fmt.Fprintf(cli.ErrWriter, "ERROR: %v\n", err)
		if err != nil {
			logrus.Errorf("Unexpected error: %s", err)
		}
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
	params, err := parameters.All(configPaths, vars)
	if err != nil {
		return err
	}

	r := renderer.New(
		renderer.WithParameters(params),
		renderer.WithSprigFunctions(),
		renderer.WithExtraFunctions(),
	)
	err = r.FileRender(inputPath, outputPath)
	if err != nil {
		if err == files.ErrExpectedStdin {
			return fmt.Errorf("expected either stdin or --in parameter, for usage use --help")
		}
		return err
	}

	return nil
}
