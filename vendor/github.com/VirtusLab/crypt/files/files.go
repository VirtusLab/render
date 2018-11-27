package files

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// ErrExpectedStdin indicates that an stdin pipe was expected but not present
	ErrExpectedStdin = errors.New("expected a pipe stdin")
)

// ReadInput reads bytes from inputPath (if not empty) or stdin
func ReadInput(inputPath string) ([]byte, error) {
	var inputFile *os.File
	if inputPath == "" {
		stdinFileInfo, _ := os.Stdin.Stat()
		if (stdinFileInfo.Mode() & os.ModeNamedPipe) != 0 {
			logrus.Debug("No input path, using piped stdin")
			inputFile = os.Stdin
		} else {
			return nil, ErrExpectedStdin
		}
	} else {
		logrus.Debugf("input path: %v", inputPath)
		f, err := os.Open(inputPath)
		if err != nil {
			logrus.Debugf("Cannot open file: '%s'; %v", inputPath, err)
			return nil, err
		}
		defer f.Close()
		inputFile = f
	}
	fileContent, err := ioutil.ReadAll(inputFile)
	if err != nil {
		logrus.Debugf("Cannot read file: '%s'; %v", inputPath, err)
		return nil, err
	}
	return fileContent, nil
}

// WriteOutput writes given bytes into outputPath (if not empty) or stdout
func WriteOutput(outputPath string, outputContent []byte, perm os.FileMode) error {
	if outputPath == "" {
		logrus.Debug("No output path, writing to stdout")
		count, err := os.Stdout.Write(outputContent)
		if err == nil && count < len(outputContent) {
			logrus.Warnf("Wrote only %v/%v bytes", count, len(outputContent))
			return io.ErrShortWrite
		}
		if err != nil {
			logrus.Debugf("Error writing to stdout; %v", err)
			return err
		}
	} else {
		logrus.Debugf("Writing to file: %v", outputPath)
		err := ioutil.WriteFile(outputPath, outputContent, perm)
		if err != nil {
			logrus.Debugf("Error writing to file: '%s'; %v", outputPath, err)
			return err
		}
	}
	return nil
}
