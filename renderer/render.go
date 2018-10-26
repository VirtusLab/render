package renderer

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer/configuration"
)

type Renderer struct {
	configuration configuration.Configuration
}

func New(configuration configuration.Configuration) *Renderer {
	return &Renderer{
		configuration: configuration,
	}
}

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

func (r *Renderer) Render(templateName, rawTemplate string) (string, error) {
	extraFunctions := sprig.TxtFuncMap()
	extraFunctions["render"] = r.render
	extraFunctions["readFile"] = r.ReadFile
	tmpl, err := template.New(templateName).Funcs(extraFunctions).Parse(rawTemplate)
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

func (r *Renderer) render(rawTemplate string) (string, error) {
	return r.Render("inner", rawTemplate)
}
