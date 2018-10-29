package renderer

import (
	"io/ioutil"

	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer/configuration"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

func (r *Renderer) root() (string, error) {
	if value, ok := r.configuration[configuration.RootKey].(string); ok {
		return value, nil
	}
	return "", errors.Errorf("can't get '%s' as a string key from the configuration", configuration.RootKey)
}

// ReadFile provides a custom template function for in-template file opening
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

func (r *Renderer) ToYaml(yamlSnippet interface{}) (string, error) {
	marshaledYaml, err := yaml.Marshal(yamlSnippet)
	return string(marshaledYaml), err
}

// TODO: gzip, ungzip, encrypt, decrypt
