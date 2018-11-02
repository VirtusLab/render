package files

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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

// IsNotEmptyAndExists checks the given file exists and is not empty
func IsNotEmptyAndExists(file string) bool {
	if len(file) == 0 {
		logrus.Infof("Configuration file path is empty")
		return false
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		logrus.Infof("Configuration file path does not exist")
		return false
	}

	if fileInfo.Size() == 0 {
		logrus.Infof("Configuration file is empty")
		return false
	}

	return true
}

// ToAbsPath turns a relative path into an absolute path with the given root path, absolute paths are ignored
func ToAbsPath(path, root string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(root, path), nil
}

// Pwd returns the process working directory
func Pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}
