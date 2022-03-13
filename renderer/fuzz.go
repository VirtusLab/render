//go:build gofuzz
// +build gofuzz

package renderer

import (
	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/sirupsen/logrus"
)

// fuzzing function for go-fuzz
func Fuzz(data []byte) int {
	logrus.SetLevel(logrus.PanicLevel) // speed up the fuzzing
	params := parameters.Parameters{
		"app_name":          "render",
		"embedded":          "{{ .value }}",
		"embedded_override": "{{ .inner | render .override }}",
		"value":             "some",
		"override": map[string]interface{}{
			"value": "other",
		},
		"nested": map[string]interface{}{
			"things": []string{
				"one",
				"two",
				"three",
			},
		},
		"inner": map[string]interface{}{
			"path": "corpus/inner.yaml.tmpl",
		},
		"compressed": map[string]interface{}{
			"encoded": "H4sIAEJ93FsAA0vOzy0oSi0uTk1RSM7PK0nNKwEAmEAxaRIAAAA=",
		},
		"embed": "content to be embedded",
	}
	r := New(
		WithParameters(params),
		WithSprigFunctions(),
		WithExtraFunctions(),
	)
	input := string(data)
	output, err := r.Render(input)
	if err != nil {
		if len(output) > 0 {
			panic("len(output) > 0")
		}
		return 0
	}
	if len(input) > 0 && len(output) == 0 {
		panic("len(input) > 0 && len(output) == 0")
	}
	return 1
}
