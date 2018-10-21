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
	"github.com/imdario/mergo"
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

func NewConfiguration(configPaths []string, extraParams []string) (Configuration, error) {
	var accumulator = make(map[string]interface{})
	err := mergeFiles(accumulator, configPaths)
	if err != nil {
		logrus.Debug("Cannot merge multiple files into the main config: %s", err)
		return nil, err
	}
	logrus.Debugf("Configuration from files: %v", accumulator)

	err = mergeVars(accumulator, extraParams)
	if err != nil {
		logrus.Debug("Cannot merge vars into the main config: %s", err)
		return nil, err
	}
	logrus.Debugf("Configuration from files and vars: %v", accumulator)

	return accumulator, nil
}

func mergeFiles(accumulator Configuration, configPaths []string) error {
	for i, configPath := range configPaths {
		logrus.Debugf("Reading configuration file [%d]: %v", i, configPath)
		if files.IsNotEmptyAndExists(configPath) {
			b, err := ioutil.ReadFile(configPath)
			if err != nil {
				logrus.Errorf("Can't open the configuration file: %v", err)
				return err
			}
			var config map[string]interface{}
			err = yaml.Unmarshal(b, &config)
			if err != nil {
				logrus.Errorf("Can't parse the configuration file: %v", err)
				return err
			}
			MergeConfigurations(&accumulator, config)
		}
	}
	return nil
}

func mergeVars(accumulator Configuration, extraParams []string) error {
	var config = make(Configuration)
	for _, v := range extraParams {
		if varArgRegexp.Match(v) {
			groups := varArgRegexp.MatchGroups(v)
			name := groups["name"]
			value := groups["value"]
			logrus.Debugf("Extra var: %s=%s", name, value)
			config[name] = value
		} else {
			logrus.Error("Expected a valid extra parameter: '%s'", v)
		}
	}
	err := MergeConfigurations(&accumulator, config)
	if err != nil {
		return err
	}
	return nil
}

func MergeConfigurations(dst *Configuration, src Configuration) error {
	err := mergo.Merge(dst, src)
	if err != nil {
		return err
	}
	return nil
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
