package renderer

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/VirtusLab/go-extended/pkg/files"
	json2 "github.com/VirtusLab/go-extended/pkg/json"
	"github.com/VirtusLab/go-extended/pkg/jsonpath"
	yaml2 "github.com/VirtusLab/go-extended/pkg/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (r *renderer) root() (string, error) {
	params := r.Configuration().Parameters
	if value, ok := params[parameters.RootKey].(string); ok {
		return value, nil
	}
	return files.Pwd()
}

// NestedRender template function allows for recursive use of the renderer
// Accepts 1 or 2 arguments:
// - NestedRender(template string)
// - NestedRender(extraParams map[string]interface{}, template string)
// Returns an error when 0 or more than 2 arguments are passed.
func (r *renderer) NestedRender(args ...interface{}) (string, error) {
	argN := len(args)

	logrus.Debugf("Nested render called with %d arguments", argN)
	for i, a := range args {
		logrus.Debugf("[%d] type: '%T', value: '%+v'", i, a, a)
	}

	var template string
	var extraParams map[string]interface{}
	switch argN {
	case 1:
		var ok bool
		template, ok = args[0].(string)
		if !ok {
			return "", errors.Errorf(
				"expected the only parameter to be a 'string', got: '%T'", args[0])
		}
	case 2:
		var ok bool
		extraParams, ok = args[0].(map[string]interface{})
		if !ok {
			return "", errors.Errorf(
				"expected the first parameter to be 'map[string]interface{}', got: '%T'", args[0])
		}
		template, ok = args[1].(string)
		if !ok {
			return "", errors.Errorf(
				"expected the second parameter to be 'string', got: '%T'", args[1])
		}
	default:
		return "", errors.Errorf("expected 1 or 2 parameters, got: %d", argN)
	}
	return r.Clone(
		WithMoreParameters(extraParams),
	).Render(template)
}

// ReadFile is a template function that allows for an in-template file reading.
// It takes a file path argument, the path can be absolute
// or relative to the process working directory.
// The relative path root can be changed with a parameter parameter.RootKey
func (r *renderer) ReadFile(file string) (string, error) {
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

// WriteFile is a template function that allows for an in-template file writing.
// It takes a file path and content arguments, the path can be absolute
// or relative to the process working directory.
// The relative path root can be changed with a parameter parameter.RootKey
func (r *renderer) WriteFile(file string, content string) (string, error) {
	root, err := r.root()
	if err != nil {
		return file, err
	}
	absPath, err := files.ToAbsPath(file, root)
	if err != nil {
		return file, err
	}
	err = os.MkdirAll(filepath.Dir(absPath), 0755)
	if err != nil {
		return file, err
	}
	return file, ioutil.WriteFile(absPath, []byte(content), 0644)
}

// ToYAML is a template function, it turns a marshallable structure into a YAML fragment
func ToYAML(marshallable interface{}) (string, error) {
	logrus.Debug("marshallable: ", marshallable)
	marshaledYaml, err := yaml.Marshal(marshallable)
	return string(marshaledYaml), err
}

// FromYAML is a template function, that unmarshalls YAML string to a map
func FromYAML(unmarshallable string) (interface{}, error) {
	logrus.Debug("unmarshallable: ", unmarshallable)
	result, err := yaml2.ToInterface(strings.NewReader(unmarshallable))
	logrus.Debugf("result: %s (type: %s)", result, reflect.TypeOf(result))
	return result, err
}

// FromJSON is a template function, that unmarshalls JSON string to a map
func FromJSON(unmarshallable string) (interface{}, error) {
	logrus.Debug("unmarshallable: ", unmarshallable)
	result, err := json2.ToInterface(strings.NewReader(unmarshallable))
	logrus.Debug("result: ", result)
	return result, err
}

// JSONPath is a template function, that evaluates JSONPath expression
// against a data structure and returns a list of results
func JSONPath(expression string, marshallable interface{}) (interface{}, error) {
	logrus.Debug("expression: ", expression)
	logrus.Debugf("marshallable: %s (type: %s, kind: %s)",
		marshallable, reflect.TypeOf(marshallable), reflect.ValueOf(marshallable).Kind())

	final, err := jsonpath.New(expression).ExecuteToInterface(marshallable)
	logrus.Debugf("final: %s (type: %s, kind: %s)",
		final, reflect.TypeOf(final), reflect.ValueOf(final).Kind())
	return final, err
}

// Gzip compresses the input using gzip algorithm
func Gzip(input interface{}) (string, error) {
	inputAsBytes, err := asBytes(input)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer func() { _ = w.Close() }()

	_, err = w.Write(inputAsBytes)
	if err != nil {
		return "", err
	}

	err = w.Flush()
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
	defer func() { _ = r.Close() }()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// asBytes ensures input will be []byte if is string
func asBytes(input interface{}) ([]byte, error) {
	switch input := input.(type) {
	case []byte:
		return input, nil
	case string:
		return []byte(input), nil
	default:
		return nil, errors.Errorf("expected []byte or string, got: '%v'", reflect.TypeOf(input))
	}
}
