package renderer

import (
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/VirtusLab/go-extended/pkg/files"
	base "github.com/VirtusLab/go-extended/pkg/renderer"
	"github.com/VirtusLab/go-extended/pkg/renderer/config"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Renderer allows for parameterised text template rendering
type Renderer interface {
	base.Renderer

	FileRender(inputPath, outputPath string) error
}

type renderer struct {
	base.Renderer
}

// New creates a new renderer with the specified parameters and zero or more options
func New(configurators ...func(*config.Config)) Renderer {
	r := &renderer{
		Renderer: base.New(configurators...),
	}
	r.Reconfigure(
		WithMoreFunctions(template.FuncMap{
			"render":   r.Render,
			"readFile": r.ReadFile,
		}),
	)
	return r
}

// WithParameters mutates Renderer configuration with new template parameters
func WithParameters(parameters map[string]interface{}) func(*config.Config) {
	return base.WithParameters(parameters)
}

// WithOptions mutates Renderer configuration with new template functions
func WithOptions(options ...string) func(*config.Config) {
	return base.WithOptions(options...)
}

// WithDelim mutates Renderer configuration with new left and right delimiters
func WithDelim(left, right string) func(*config.Config) {
	return base.WithDelim(left, right)
}

// WithFunctions mutates Renderer configuration with new template functions
func WithFunctions(extraFunctions template.FuncMap) func(*config.Config) {
	return base.WithFunctions(extraFunctions)
}

// WithMoreFunctions mutates Renderer with new template functions,
func WithMoreFunctions(moreFunctions template.FuncMap) func(*config.Config) {
	return func(c *config.Config) {
		allFunctions := c.ExtraFunctions
		err := MergeFunctions(&allFunctions, moreFunctions)
		if err != nil {
			logrus.Panicf("unexpected problem merging extra functions")
		}
		c.ExtraFunctions = allFunctions
	}
}

// WithExtraFunctions mutates Renderer configuration with the custom template functions
func WithExtraFunctions() func(*config.Config) {
	return WithMoreFunctions(ExtraFunctions())
}

// WithSprigFunctions mutates Renderer configuration with the Sprig template functions
func WithSprigFunctions() func(*config.Config) {
	return WithMoreFunctions(sprig.TxtFuncMap())
}

// MergeFunctions merges two template.FuncMap instances, overrides if necessary
func MergeFunctions(dst *template.FuncMap, src template.FuncMap) error {
	err := mergo.Merge(dst, src, mergo.WithOverride)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// TODO DirRender

// FileRender is used to render files by path, see also Render
func (r *renderer) FileRender(inputPath, outputPath string) error {
	input, err := files.ReadInput(inputPath)
	if err != nil {
		logrus.Debugf("Can't open the template: %v", err)
		return err
	}

	var templateName string
	if inputPath == "" {
		templateName = "stdin"
	} else {
		templateName = inputPath
	}

	result, err := r.NamedRender(templateName, string(input))
	if err != nil {
		return err
	}

	err = files.WriteOutput(outputPath, []byte(result), 0644)
	if err != nil {
		logrus.Debugf("Can't save the rendered: %v", err)
		return err
	}

	return nil
}

/*
ExtraFunctions provides additional template functions to the standard (text/template) ones:

  - toYaml - provides a configuration data structure fragment as a YAML format
  - gzip - use gzip compression inside the templates, for best results use with b64enc
  - ungzip - use gzip extraction inside the templates, for best results use with b64dec
  - encryptAWS - encrypts the data from inside of the template using AWS KMS, for best results use with gzip and b64enc
  - decryptAWS - decrypts the data from inside of the template using AWS KMS, for best results use with ungzip and b64dec
  - encryptGCP - encrypts the data from inside of the template using GCP KMS, for best results use with gzip and b64enc
  - decryptGCP - decrypts the data from inside of the template using GCP KMS, for best results use with ungzip and b64dec
  - encryptAzure - encrypts the data from inside of the template using Azure Key Vault, for best results use with gzip and b64enc
  - decryptAzure - decrypts the data from inside of the template using Azure Key Vault, for best results use with ungzip and b64dec

*/
func ExtraFunctions() template.FuncMap {
	return template.FuncMap{
		"toYaml":       ToYaml,
		"ungzip":       Ungzip,
		"gzip":         Gzip,
		"encryptAWS":   EncryptAWS,
		"decryptAWS":   DecryptAWS,
		"encryptGCP":   EncryptGCP,
		"decryptGCP":   DecryptGCP,
		"encryptAzure": EncryptAzure,
		"decryptAzure": DecryptAzure,
	}
}
