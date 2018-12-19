package renderer

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"

	"github.com/VirtusLab/render/files"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// MissingKeyInvalidOption is the renderer option to continue execution on missing key and print "<no value>"
	MissingKeyInvalidOption = "missingkey=invalid"
	// MissingKeyErrorOption is the renderer option to stops execution immediately with an error on missing key
	MissingKeyErrorOption = "missingkey=error"
	// LeftDelim is the default left template delimiter
	LeftDelim = "{{"
	// RightDelim is the default right template delimiter
	RightDelim = "}}"
)

// Renderer structure holds parameters and options
type Renderer struct {
	parameters     map[string]interface{}
	options        []string
	leftDelim      string
	rightDelim     string
	extraFunctions template.FuncMap
}

// New creates a new renderer with the specified parameters and zero or more options
func New() *Renderer {
	r := &Renderer{
		parameters:     map[string]interface{}{},
		options:        []string{MissingKeyErrorOption},
		leftDelim:      LeftDelim,
		rightDelim:     RightDelim,
		extraFunctions: template.FuncMap{},
	}
	r.Functions(r.ExtraFunctions())
	return r
}

// Delim mutates Renderer with new left and right delimiters
func (r *Renderer) Delim(left, right string) *Renderer {
	r.leftDelim = left
	r.rightDelim = right
	return r
}

// Functions mutates Renderer with new template functions
func (r *Renderer) Functions(extraFunctions template.FuncMap) *Renderer {
	r.extraFunctions = extraFunctions
	return r
}

// Options mutates Renderer with new template functions
func (r *Renderer) Options(options ...string) *Renderer {
	r.options = options
	return r
}

// Parameters mutates Renderer with new template parameters
func (r *Renderer) Parameters(parameters map[string]interface{}) *Renderer {
	r.parameters = parameters
	return r
}

// Render is a simple rendering function, also used as a custom template function
// to allow in-template recursive rendering, see also NamedRender
func (r *Renderer) Render(rawTemplate string) (string, error) {
	return r.NamedRender("nameless", rawTemplate)
}

// TODO DirRender

// FileRender is used to render files by path, see also Render
func (r *Renderer) FileRender(inputPath, outputPath string) error {
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

// NamedRender is the main rendering function, see also Render, Parameters and ExtraFunctions
func (r *Renderer) NamedRender(templateName, rawTemplate string) (string, error) {
	err := r.Validate()
	if err != nil {
		logrus.Errorf("Invalid state; %v", err)
		return "", err
	}
	t, err := r.Parse(templateName, rawTemplate, r.extraFunctions)
	if err != nil {
		logrus.Errorf("Can't parse the template; %v", err)
		return "", err
	}
	out, err := r.Execute(t)
	if err != nil {
		logrus.Errorf("Can't execute the template; %v", err)
		return "", err
	}
	return out, nil
}

// Validate checks the internal state and returns error if necessary
func (r *Renderer) Validate() error {
	if r.parameters == nil {
		return errors.New("unexpected 'nil' parameters")
	}

	if len(r.leftDelim) == 0 {
		return errors.New("unexpected empty leftDelim")
	}
	if len(r.rightDelim) == 0 {
		return errors.New("unexpected empty rightDelim")
	}

	for _, o := range r.options {
		switch o {
		case MissingKeyErrorOption:
		case MissingKeyInvalidOption:
		default:
			return errors.Errorf("unexpected option: '%s', option must be in: '%s'",
				o, strings.Join([]string{MissingKeyInvalidOption, MissingKeyErrorOption}, ", "))
		}
	}
	return nil
}

// Parse is a basic template parsing function
func (r *Renderer) Parse(templateName, rawTemplate string, extraFunctions template.FuncMap) (*template.Template, error) {
	return template.New(templateName).
		Delims(r.leftDelim, r.rightDelim).
		Funcs(extraFunctions).
		Option(r.options...).
		Parse(rawTemplate)
}

// Execute is a basic template execution function
func (r *Renderer) Execute(t *template.Template) (string, error) {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, r.parameters)
	if err != nil {
		retErr := err
		logrus.Debugf("(%v): %v", reflect.TypeOf(err), err)
		if e, ok := err.(template.ExecError); ok {
			retErr = errors.Wrapf(err,
				"Error evaluating the template named: '%s'", e.Name)
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
  - encryptAWS - encrypts the data from inside of the template using AWS KMS, for best results use with gzip and b64enc
  - decryptAWS - decrypts the data from inside of the template using AWS KMS, for best results use with ungzip and b64dec
  - encryptGCP - encrypts the data from inside of the template using GCP KMS, for best results use with gzip and b64enc
  - decryptGCP - decrypts the data from inside of the template using GCP KMS, for best results use with ungzip and b64dec
  - encryptAzure - encrypts the data from inside of the template using Azure Key Vault, for best results use with gzip and b64enc
  - decryptAzure - decrypts the data from inside of the template using Azure Key Vault, for best results use with ungzip and b64dec

*/
func (r *Renderer) ExtraFunctions() template.FuncMap {
	extraFunctions := sprig.TxtFuncMap()
	extraFunctions["render"] = r.Render
	extraFunctions["readFile"] = r.ReadFile
	extraFunctions["toYaml"] = ToYaml
	extraFunctions["ungzip"] = Ungzip
	extraFunctions["gzip"] = Gzip
	extraFunctions["encryptAWS"] = EncryptAWS
	extraFunctions["decryptAWS"] = DecryptAWS
	extraFunctions["encryptGCP"] = EncryptGCP
	extraFunctions["decryptGCP"] = DecryptGCP
	extraFunctions["encryptAzure"] = EncryptAzure
	extraFunctions["decryptAzure"] = DecryptAzure
	return extraFunctions
}
