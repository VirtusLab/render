package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/VirtusLab/render/constants"
	"github.com/VirtusLab/render/renderer"
	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/VirtusLab/go-extended/pkg/files"
	"github.com/VirtusLab/go-extended/pkg/renderer/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

var (
	app                     *cli.App
	inputFile               string
	outputFile              string
	inputDir                string
	outputDir               string
	configPaths             cli.StringSlice
	vars                    cli.StringSlice
	unsafeIgnoreMissingKeys bool
)

func main() {
	app = cli.NewApp()
	app.Name = constants.Name
	app.Usage = constants.Description
	app.Author = constants.Author
	app.Version = constants.Version()
	app.Before = preload
	app.Action = action

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "run in debug mode",
		},
		cli.StringFlag{
			Name:        "indir",
			Value:       "",
			Usage:       "the input directory, can't be used with --out",
			Destination: &inputDir,
		},
		cli.StringFlag{
			Name:        "outdir",
			Value:       "",
			Usage:       "the output directory, the same as --outdir if empty, can't be used with --in",
			Destination: &outputDir,
		},
		cli.StringFlag{
			Name:        "in",
			Value:       "",
			Usage:       "the input template file, stdin if empty, can't be used with --outdir",
			Destination: &inputFile,
		},
		cli.StringFlag{
			Name:        "out",
			Value:       "",
			Usage:       "the output file, stdout if empty, can't be used with --indir",
			Destination: &outputFile,
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
		cli.BoolFlag{
			Name:        "unsafe-ignore-missing-keys",
			Usage:       "do not fail on missing map key and print '<no value>' ('missingkey=invalid')",
			Destination: &unsafeIgnoreMissingKeys,
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		logrus.Errorf("Command not found: '%s'", command)
		cli.OsExiter(1)
	}
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			return err
		}

		logrus.Errorf("Usage error: %s", err)
		return nil
	}
	cli.OsExiter = func(c int) {
		if c != 0 {
			logrus.Debugf("Exiting with code %d", c)
		}
		os.Exit(c)
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("Unexpected error: %v", err)
		cli.OsExiter(1)
	}
}

func preload(c *cli.Context) error {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Infof("Version %s", app.Version)

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

func action(c *cli.Context) error {
	if c.NArg() > 0 {
		return fmt.Errorf("have not expected any arguments, got %d", c.NArg())
	}

	opts := []string{config.MissingKeyErrorOption}
	if unsafeIgnoreMissingKeys {
		logrus.Warnf("You are using '--unsafe-ignore-missing-keys' and %s will use option '%s'",
			app.Name, config.MissingKeyInvalidOption)
		opts = []string{config.MissingKeyInvalidOption}
	}

	if len(configPaths) > 0 {
		logrus.Infof("Configurations:\n\t%s", strings.Join(configPaths, "\n\t"))
	}
	if len(vars) > 0 {
		logrus.Infof("Variables:\n\t%s", strings.Join(vars, "\n\t"))
	}
	params, err := parameters.All(configPaths, vars)
	if err != nil {
		return err
	}

	r := renderer.New(
		renderer.WithOptions(opts...),
		renderer.WithParameters(params),
		renderer.WithSprigFunctions(),
		renderer.WithExtraFunctions(),
		renderer.WithCryptFunctions(),
	)

	if len(inputDir) > 0 {
		if len(inputFile) > 0 {
			return fmt.Errorf("conflict, --in can't be used with --indir or --outdir")
		}
		if len(outputFile) > 0 {
			return fmt.Errorf("conflict, --out can't be used with --indir or --outdir")
		}
		if len(outputDir) == 0 {
			outputDir = inputDir
		}

		err = r.DirRender(inputDir, outputDir)
		switch err.(type) {
		case nil:
			return nil
		default:
			return err
		}
	}

	if len(inputDir) > 0 {
		return fmt.Errorf("conflict, --indir can't be used with --in or --out")
	}
	if len(outputDir) > 0 {
		return fmt.Errorf("conflict, --outdir can't be used with --in or --out")
	}
	err = r.FileRender(inputFile, outputFile)
	switch err.(type) {
	case nil:
		return nil
	case *files.ErrExpectedStdin:
		return fmt.Errorf("expected either stdin, --indir or --in parameter, for usage use --help")
	default:
		return err
	}
}
