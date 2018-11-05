package renderer

import (
	"bytes"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer/configuration"
	"github.com/pkg/errors"
)

const (
	// MissingKeyInvalidOption is the renderer option to continue execution on missing key and print "<no value>"
	MissingKeyInvalidOption = "missingkey=invalid"
	// MissingKeyErrorOption is the renderer option to stops execution immediately with an error on missing key
	MissingKeyErrorOption = "missingkey=error"
)

// Renderer structure holds configuration and options
type Renderer struct {
	configuration configuration.Configuration
	options       []string
}

// New creates a new renderer with the specified configuration and zero or more options
func New(configuration configuration.Configuration, opts ...string) *Renderer {
	return &Renderer{
		configuration: configuration,
		options:       opts,
	}
}

// RenderFile is used to render files by path, see also Render
func (r *Renderer) RenderFile(inputPath, outputPath string) error {
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

	result, err := r.Render(templateName, string(input))
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

// Render is the main rendering function, see also SimpleRender, RenderWith, Configuration and ExtraFunctions
func (r *Renderer) Render(templateName, rawTemplate string) (string, error) {
	return r.RenderWith(templateName, rawTemplate, r.ExtraFunctions())
}

// SimpleRender is a simple rendering function, also used as a custom template function
// to allow in-template recursive rendering, see also Render, RenderWith
func (r *Renderer) SimpleRender(rawTemplate string) (string, error) {
	return r.Render("nameless", rawTemplate)
}

// RenderWith is the basic rendering function that takes extraFunctions as an argument
func (r *Renderer) RenderWith(templateName, rawTemplate string, extraFunctions template.FuncMap) (string, error) {
	tmpl, err := template.New(templateName).Funcs(extraFunctions).Option(r.options...).Parse(rawTemplate)
	if err != nil {
		logrus.Errorf("Can't parse the template; %v", err)
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, r.configuration)
	if err != nil {
		retErr := err
		logrus.Debugf("(%v): %v", reflect.TypeOf(err), err)
		if e, ok := err.(template.ExecError); ok {
			retErr = errors.Wrapf(err,
				"Error evaluating the template named: '%s'", e.Name)
		} else {
			retErr = errors.Wrap(err, "Can't render the template")
		}
		return "", retErr
	}
	return buffer.String(), nil
}

/*
ExtraFunctions provides additional template functions to the standard (text/template) ones,
it adds sprig functions and custom functions:

  - render - calls the render from inside of the template, making the renderer recursive
  - readFile - reads a file from a given path, relative paths are translated to absolute
          paths, based on root function
  - root - the root path for rendering, used relative to absolute path translation
          in any file based operations
  - toYaml - provides a configuration data structure fragment as a YAML format
  - gzip - use gzip compression inside the templates, for best results use with b64enc
  - ungzip - use gzip extraction inside the templates, for best results use with b64dec

*/
func (r *Renderer) ExtraFunctions() template.FuncMap {
	extraFunctions := sprig.TxtFuncMap()
	extraFunctions["render"] = r.SimpleRender
	extraFunctions["readFile"] = r.ReadFile
	extraFunctions["toYaml"] = ToYaml
	extraFunctions["ungzip"] = Ungzip
	extraFunctions["gzip"] = Gzip
	return extraFunctions
}
