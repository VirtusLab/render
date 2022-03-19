package renderer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/Masterminds/sprig/v3"
	crypto "github.com/VirtusLab/crypt/crypto/render"
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
	DirRender(inputDir, outputDir string) error
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
			"render":    r.NestedRender,
			"readFile":  r.ReadFile,
			"writeFile": r.WriteFile,
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

// WithNetFunctions mutates Renderer configuration by merging the custom template functions
func WithNetFunctions() func(*config.Config) {
	return WithMoreFunctions(NetFunctions())
}

// MergeFunctions merges two template.FuncMap instances, overrides if necessary
func MergeFunctions(dst *template.FuncMap, src template.FuncMap) error {
	err := mergo.Merge(dst, src, mergo.WithOverride)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// TODO parametrize
var defaultTemplateExtensions = []string{".tpl", ".tmpl"}

// DirRender is used to render files by directory, see also FileRender
// TODO break up to multiple small functions
func (r *renderer) DirRender(inputDir, outputDir string) error {
	logrus.Infof("Directory mode selected: '%s' -> '%s'", inputDir, outputDir)

	fileEntries, err := dirTree(inputDir)
	if err != nil {
		return errors.Wrapf(err, "can't scan the directory tree: '%s'", inputDir)
	}

	for _, file := range fileEntries {
		logrus.Debugf("Processing '%s'", path.Join(file.path, file.name))

		target := trimExtension(file, defaultTemplateExtensions)

		rel, err := filepath.Rel(inputDir, file.path)
		if err != nil {
			return errors.Wrapf(err, "can't get a relative path for: '%s'", file.path)
		}

		target.path = path.Join(outputDir, rel)

		_, err = os.Stat(target.path)
		if os.IsNotExist(err) {
			err := os.MkdirAll(target.path, os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "can't create the target directory: '%s'", target.path)
			}
			logrus.Infof("Target directory was created: '%s'", target.path)
		} else if err != nil {
			return errors.Wrapf(err, "can't get file information for '%s'", target.path)
		}

		err = r.FileRender(path.Join(file.path, file.name), path.Join(target.path, target.name))
		if err != nil {
			return errors.Wrap(err, "can't render a file")
		}
	}

	return nil
}

// FileRender is used to render files by path, see also DirRender
func (r *renderer) FileRender(inputPath, outputPath string) error {
	inputName := inputPath
	outputName := outputPath
	if inputPath == "" {
		inputName = "stdin"
	}
	if outputPath == "" {
		outputName = "stdout"
	}
	logrus.Infof("Rendering '%s' -> '%s'\n", inputName, outputName)

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

	inputString := string(input)
	logrus.Debugf("%s: \n%s", inputName, inputString)
	result, err := r.NamedRender(templateName, inputString)
	if err != nil {
		return err
	}
	logrus.Debugf("%s: \n%s", outputName, result)

	err = files.WriteOutput(outputPath, []byte(result), 0644)
	if err != nil {
		logrus.Debugf("Can't save the rendered file: %v", err)
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

// ExtraFunctions provides additional template functions
// to the standard (text/template) ones
func ExtraFunctions() template.FuncMap {
	return template.FuncMap{
		"n":        N,
		"toYaml":   ToYAML,
		"fromYaml": FromYAML,
		"fromJson": FromJSON,
		"jsonPath": JSONPath,
		"ungzip":   Ungzip,
		"gzip":     Gzip,
	}
}

// NetFunctions provides additional template functions
// to the standard (text/template) ones
func NetFunctions() template.FuncMap {
	return template.FuncMap{
		"cidrhost":    CidrHost,
		"cidrnetmask": CidrNetmask,
		"cidrsubnet":  CidrSubnet,
		"cidrsubnets": CidrSubnets,
	}
}

// TODO move to files package
type dirEntry struct {
	path      string
	name      string
	extension string
}

// TODO move to files package
func dirTree(input string) (entries []dirEntry, err error) {
	err = filepath.Walk(input, func(path string, info os.FileInfo, dirErr error) error {
		if dirErr != nil {
			logrus.Errorf("error '%v' on path '%s'", dirErr, path)
			return dirErr
		}

		logrus.Debugf("Discovered path: '%s'", path)

		if !info.IsDir() {
			logrus.Tracef("  dir  : '%s'", filepath.Dir(path))
			logrus.Tracef("  name : '%s'", info.Name())
			logrus.Tracef("  ext  : '%s'", filepath.Ext(path))

			entry := dirEntry{
				path:      filepath.Dir(path),
				name:      info.Name(),
				extension: filepath.Ext(path),
			}
			entries = append(entries, entry)
		}
		return nil
	})
	if err != nil {
		return entries, errors.Wrapf(err, "can't walk the directory tree '%s'", input)
	}

	return entries, nil
}

func trimExtension(file dirEntry, extensions []string) (new dirEntry) {
	new = file
	for _, ext := range extensions {
		if file.extension == ext {
			new.name = strings.TrimSuffix(file.name, ext)
			new.extension = filepath.Ext(new.name)
		}
	}
	return
}
