package renderer

import (
	"fmt"
	"text/template"

	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/Masterminds/sprig"
	"github.com/VirtusLab/crypt/crypto"
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
	Clone(configurators ...func(*config.Config)) Renderer
	FileRender(inputPath, outputPath string) error

	NestedRender(args ...interface{}) (string, error)
	ReadFile(file string) (string, error)
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
			"render":   r.NestedRender,
			"readFile": r.ReadFile,
		}),
	)
	return r
}

// WithParameters mutates Renderer configuration by replacing all template parameters
func WithParameters(parameters map[string]interface{}) func(*config.Config) {
	return base.WithParameters(parameters)
}

// WithMoreParameters mutates Renderer configuration by merging the given template parameters
func WithMoreParameters(extraParams ...map[string]interface{}) func(*config.Config) {
	return func(c *config.Config) {
		var err error
		for _, extra := range extraParams {
			c.Parameters, err = parameters.Merge(c.Parameters, extra)
		}
		if err != nil {
			logrus.Panicf("unexpected problem merging extra functions")
		}
	}
}

// WithOptions mutates Renderer configuration by replacing the template functions
func WithOptions(options ...string) func(*config.Config) {
	return base.WithOptions(options...)
}

// WithDelim mutates Renderer configuration by replacing the left and right delimiters
func WithDelim(left, right string) func(*config.Config) {
	return base.WithDelim(left, right)
}

// WithFunctions mutates Renderer configuration by replacing the template functions
func WithFunctions(extraFunctions template.FuncMap) func(*config.Config) {
	return base.WithFunctions(extraFunctions)
}

// WithMoreFunctions mutates Renderer configuration by merging the given template functions,
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

// WithExtraFunctions mutates Renderer configuration by merging the custom template functions
func WithExtraFunctions() func(*config.Config) {
	return WithMoreFunctions(ExtraFunctions())
}

// WithSprigFunctions mutates Renderer configuration by merging the Sprig template functions
func WithSprigFunctions() func(*config.Config) {
	return WithMoreFunctions(sprig.TxtFuncMap())
}

// WithCryptFunctions mutates Renderer configuration by merging the Crypt template functions
func WithCryptFunctions() func(*config.Config) {
	return WithMoreFunctions(crypto.TemplateFunctions())
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

// Clone returns a new copy of the renderer modified with the optional configurators
func (r *renderer) Clone(configurators ...func(*config.Config)) Renderer {
	clone := &renderer{
		Renderer: base.NewWithConfig(r.Configuration()),
	}
	clone.Reconfigure(configurators...)
	logrus.Debugf("cloned renderer: %+v", clone.String())
	return clone
}

func (r *renderer) String() string {
	return fmt.Sprintf("%+v", r.Renderer.Configuration())
}

/*
ExtraFunctions provides additional template functions to the standard (text/template) ones:

  - toYaml - provides a configuration data structure fragment as a YAML format
  - gzip - use gzip compression inside the templates, for best results use with b64enc
  - ungzip - use gzip extraction inside the templates, for best results use with b64dec

*/
func ExtraFunctions() template.FuncMap {
	return template.FuncMap{
		"toYaml": ToYaml,
		"ungzip": Ungzip,
		"gzip":   Gzip,
	}
}
