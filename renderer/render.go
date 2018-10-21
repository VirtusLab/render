package renderer

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/matcher"
	"github.com/ghodss/yaml"
)

var (
	varArgRegexp = matcher.NewMust(`^(?P<name>\S+)=(?P<value>\S*)$`)
)

type Configuration map[string]interface{}

type Renderer struct {
	configuration Configuration
}

func New(configuration Configuration) *Renderer {
	return &Renderer{
		configuration: configuration,
	}
}

func (r *Renderer) RenderFile(inputPath, outputPath string) error {
	input, err := files.ReadInput(inputPath)
	if err != nil {
		logrus.Errorf("Can't open the template file: %v", err)
		return err
	}

	result, err := r.Render(outputPath, string(input))
	if err != nil {
		logrus.Errorf("Can't render the template: %v", err)
		return err
	}

	err = files.WriteOutput(outputPath, []byte(result), 0644)
	if err != nil {
		logrus.Errorf("Can't save the rendered file: %v", err)
		return err
	}

	return nil
}

func NewConfiguration(configPath string, extraParams []string) (Configuration, error) {
	var configMap = make(map[string]interface{})
	if files.IsNotEmptyAndExists(configPath) {
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			logrus.Errorf("Can't open the configuration file: %v", err)
			return nil, err
		}
		err = yaml.Unmarshal(b, &configMap)
		if err != nil {
			logrus.Errorf("Can't parse the configuration file: %v", err)
			return nil, err
		}
	}
	logrus.Debugf("Configuration from files: %v", configMap)

	for _, v := range extraParams {
		if varArgRegexp.Match(v) {
			groups := varArgRegexp.MatchGroups(v)
			name := groups["name"]
			value := groups["value"]
			logrus.Debugf("Extra var: %s=%s", name, value)
			configMap[name] = value
		} else {
			logrus.Error("Expected a valid extra parameter: '%s'", v)
		}
	}
	logrus.Debugf("Configuration from files and vars: %v", configMap)

	return configMap, nil
}

func (r *Renderer) Render(templateName, rawTemplate string) (string, error) {
	extraFunctions := sprig.TxtFuncMap()
	extraFunctions["render"] = r.render
	extraFunctions["readFile"] = ReadFile
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
