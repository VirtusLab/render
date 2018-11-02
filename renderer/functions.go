package renderer

import (
	"io/ioutil"

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
func (r *Renderer) ToYaml(marshallable interface{}) (string, error) {
	marshaledYaml, err := yaml.Marshal(marshallable)
	return string(marshaledYaml), err
}

// TODO: gzip, ungzip, encrypt, decrypt
