package main

import (
	"os"
	"log"
	"text/template"
	"io/ioutil"
	"bytes"
	"github.com/ghodss/yaml"
)

func Render(input, output, config string, extraParams []string) error {
	if !NotEmptyAndExists(input) {
		log.Fatalf("File %v is empty or does not exist", input)
	}
	if !NotEmptyAndExists(config) {
		log.Fatalf("File %v is empty or does not exist", config)
	}

	stdout := false
	if len(output) == 0 {
		stdout = true
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

	b, err = ioutil.ReadFile(input)
	if err != nil {
		log.Fatal("Can't open template file", err)
	}

	parsed, err := template.New(input).Parse(string(b))
	if err != nil {
		log.Fatal("Can't parse template file", err)
	}

	var buffer bytes.Buffer
	err = parsed.Execute(&buffer, configMap)
	if err != nil {
		log.Fatal("Can't render template file", err)
	}

	rendered := buffer.String()
	if stdout {
		log.Printf("Rendering %v", input)
		log.Print(rendered)
	} else {
		err := ioutil.WriteFile(output, buffer.Bytes(), 0644)
		log.Printf("Rendering %v and saving %v", input, output)
		if err != nil {
			log.Fatal("Can't save rendered file", err)
		}
	}

	return nil
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
