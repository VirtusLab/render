package files

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
)

// ReadInput reads bytes from inputPath (if not empty) or stdin
func ReadInput(inputPath string) ([]byte, error) {
	var inputFile *os.File
	if inputPath == "" {
		stdinFileInfo, _ := os.Stdin.Stat()
		if (stdinFileInfo.Mode() & os.ModeNamedPipe) != 0 {
			logrus.Debug("no input path, using piped stdin")
			inputFile = os.Stdin
		} else {
			return nil, errors.New("expected a pipe stdin")
		}
	} else {
		logrus.Debugf("input path: %v", inputPath)
		f, err := os.Open(inputPath)
		if err != nil {
			logrus.Error("cannot open file")
			return nil, err
		}
		defer f.Close()
		inputFile = f
	}
	fileContent, err := ioutil.ReadAll(inputFile)
	if err != nil {
		logrus.Error("cannot read file")
		return nil, err
	}
	return fileContent, nil
}

// WriteOutput writes given bytes into outputPath (if not empty) or stdout
func WriteOutput(outputPath string, outputContent []byte, perm os.FileMode) error {
	if outputPath == "" {
		logrus.Debug("no output path, writing to stdout")
		count, err := os.Stdout.Write(outputContent)
		if err == nil && count < len(outputContent) {
			logrus.Warnf("wrote only %v/%v bytes", count, len(outputContent))
			return io.ErrShortWrite
		}
		if err != nil {
			logrus.Error("error writing to stdout")
			return err
		}
	} else {
		logrus.Debugf("writing to file: %v", outputPath)
		err := ioutil.WriteFile(outputPath, outputContent, perm)
		if err != nil {
			logrus.Error("error writing to file")
			return err
		}
	}
	return nil
}

func IsEmptyOrDoesNotExist(file string) bool {
	if len(file) == 0 {
		logrus.Infof("Configuration file path is empty")
		return true
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		logrus.Infof("Configuration file path does not exist")
		return true
	}

	if fileInfo.Size() == 0 {
		logrus.Infof("Configuration file is empty")
		return true
	}

	return false
}
