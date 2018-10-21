package renderer

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
)

// TODO improve error handling fatals -> error

func Render(input, output, config string, extraParams []string) error {
	configuration := ParseConfiguration(config, extraParams)
	parsedTemplate := ParseTemplate(input)
	render(parsedTemplate, configuration, output)

	return nil
}

func render(tmpl *template.Template, configuration map[string]interface{}, output string) {
	stdout := false
	if len(output) == 0 {
		stdout = true
	}

	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, configuration)
	if err != nil {
		log.Fatal("Can't render template file", err)
	}

	if stdout {
		log.Print(buffer.String())
	} else {
		err := ioutil.WriteFile(output, buffer.Bytes(), 0644)
		if err != nil {
			log.Fatal("Can't save rendered file", err)
		}
	}
}

func ParseTemplate(input string) *template.Template {
	if !NotEmptyAndExists(input) {
		log.Fatalf("Template file %v is empty or does not exist", input)
	}

	b, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal("Can't open template file", err)
	}

	extraFunctions := sprig.TxtFuncMap()
	extraFunctions["readFile"] = readFile
	parsed, err := template.New(input).Funcs(extraFunctions).Parse(string(b))
	if err != nil {
		log.Fatal("Can't parse template file", err)
	}

	return parsed
}

func ParseConfiguration(config string, extraParams []string) map[string]interface{} {
	if !NotEmptyAndExists(config) {
		log.Fatalf("Config file %v is empty or does not exist", config)
	}

	var configMap map[string]interface{}
	b, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal("Can't open config file", err)
	}
	err = yaml.Unmarshal(b, &configMap)
	if err != nil {
		log.Fatal("Can't parse config file", err)
	}

	for _, v:= range extraParams {
		if strings.Contains(v, "=") {
			pair := strings.Split(v, "=")
			configMap[pair[0]] = pair[1]
		}
	}

	return configMap
}

func NotEmptyAndExists(file string) bool {
	if len(file) == 0 {
		return false
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}

	if fileInfo.Size() == 0 {
		return false
	}

	return true
}

func readFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}