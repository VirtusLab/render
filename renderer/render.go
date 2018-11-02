package renderer

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer/configuration"
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
		logrus.Debugf("Can't open the template file: %v", err)
		return err
	}

	result, err := r.Render(outputPath, string(input))
	if err != nil {
		logrus.Debugf("Can't render the template: %v", err)
		return err
	}

	err = files.WriteOutput(outputPath, []byte(result), 0644)
	if err != nil {
		logrus.Debugf("Can't save the rendered file: %v", err)
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
		logrus.Errorf("Can't parse the template file: %v", err)
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, r.configuration)
	if err != nil {
		logrus.Errorf("Can't render the template file: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

// ExtraFunctions provides additional template functions to the text/template ones,
// it adds sprig functions and custom functions
func (r *Renderer) ExtraFunctions() template.FuncMap {
	extraFunctions := sprig.TxtFuncMap()
	extraFunctions["render"] = r.SimpleRender
	extraFunctions["readFile"] = r.ReadFile
	extraFunctions["toYaml"] = r.ToYaml
	return extraFunctions
}
