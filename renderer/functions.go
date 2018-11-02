package renderer

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"

	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer/configuration"
	"github.com/ghodss/yaml"
)

func (r *Renderer) root() (string, error) {
	if value, ok := r.configuration[configuration.RootKey].(string); ok {
		return value, nil
	}
	return files.Pwd()
}

// ReadFile is a template function that allows for an in-template file opening
func (r *Renderer) ReadFile(file string) (string, error) {
	root, err := r.root()
	if err != nil {
		return "", err
	}
	absPath, err := files.ToAbsPath(file, root)
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// ToYaml is a template function, it turns a marshallable structure into a YAML fragment
func ToYaml(marshallable interface{}) (string, error) {
	marshaledYaml, err := yaml.Marshal(marshallable)
	return string(marshaledYaml), err
}

// Gzip compresses the input using gzip algorithm
func Gzip(input interface{}) (string, error) {
	inputAsBytes, err := asBytes(input)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	_, err = w.Write(inputAsBytes)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// Ungzip un-compresses the input using gzip algorithm
func Ungzip(input interface{}) (string, error) {
	inputAsBytes, err := asBytes(input)
	if err != nil {
		return "", err
	}

	in := bytes.NewBuffer(inputAsBytes)
	r, err := gzip.NewReader(in)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func asBytes(input interface{}) ([]byte, error) {
	switch input.(type) {
	case []byte:
		return input.([]byte), nil
	case string:
		return []byte(input.(string)), nil
	default:
		return nil, errors.Errorf("expected []byte or string, got: '%v'", reflect.TypeOf(input))
	}
}

// TODO: encrypt, decrypt
