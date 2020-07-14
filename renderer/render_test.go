package renderer

import (
	"strings"
	"testing"

	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/Masterminds/sprig/v3"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestRenderer_NamedRender_Empty(t *testing.T) {
	Run(t, Test{
		name: "empty render",
		f: func(tt Test) {
			input := ""
			expected := ""

			result, err := New().NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_NamedRender_Simple(t *testing.T) {
	Run(t, Test{
		name: "simple render",
		f: func(tt Test) {
			input := `key: {{ .value }}
something:
  nested: {{ .something.nested }}`

			expected := `key: some
something:
  nested: val`

			fromVars, err := parameters.FromVars([]string{
				"something.nested=val",
			})
			if err != nil {
				t.Fatal(err)
			}

			params, err := parameters.Merge(
				parameters.Parameters{
					"value": "some",
				},
				fromVars,
			)
			if err != nil {
				t.Fatal(err)
			}

			result, err := New(WithParameters(params)).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 3, len(tt.logHook.Entries))
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_Render_Error(t *testing.T) {
	Run(t, Test{
		name: "parse error",
		f: func(tt Test) {
			input := "{{ wrong+ }}"
			expected := ""

			result, err := New().NamedRender(tt.name, input)

			assert.Error(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
		},
	})
}

func TestRenderer_NamedRender_Render(t *testing.T) {
	Run(t, Test{
		name: "render render",
		f: func(tt Test) {
			input := "key: {{ .inner | render }}"
			expected := "key: some"
			params := parameters.Parameters{
				"inner": "{{ .value }}",
				"value": "some",
			}

			result, err := New(WithParameters(params)).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_NamedRender_RenderOverride(t *testing.T) {
	Run(t, Test{
		name: "render render with override",
		f: func(tt Test) {
			input := `first: {{ .value }}
second: {{ .inner | render .override }}`
			expected := `first: some
second: other`
			params := parameters.Parameters{
				"inner": "{{ .value }}",
				"value": "some",
				"override": map[string]interface{}{
					"value": "other",
				},
			}

			result, err := New(WithParameters(params)).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_NamedRender_RenderTooManyArgs(t *testing.T) {
	Run(t, Test{
		name: "render render with override and too many args",
		f: func(tt Test) {
			input := `{{ .inner | render .override "too-many" }}`
			params := parameters.Parameters{
				"inner": "{{ .value }}",
				"value": "some",
				"override": map[string]interface{}{
					"value": "other",
				},
			}

			_, err := New(WithParameters(params)).NamedRender(tt.name, input)

			assert.Error(t, err, tt.name)
			assert.True(t, strings.Contains(err.Error(), "expected 1 or 2 parameters, got: 3"))
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_NamedRender_RenderOverride_Twice(t *testing.T) {
	Run(t, Test{
		name: "render render with override",
		f: func(tt Test) {
			input := `first: {{ .inner | render }}
second: {{ .inner | render .override }}
third: {{ .inner | render .overrider }}`
			expected := `first: some
second: other
third: otherer`
			params := parameters.Parameters{
				"inner": "{{ .value }}",
				"value": "some",
				"override": map[string]interface{}{
					"value": "other",
				},
				"overrider": map[string]interface{}{
					"value": "otherer",
				},
			}

			result, err := New(WithParameters(params)).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_NamedRender_Func(t *testing.T) {
	Run(t, Test{
		name: "parse func",
		f: func(tt Test) {
			input := "key: {{ b64enc .value }}"
			expected := "key: c29tZQ=="
			params := parameters.Parameters{
				"value": "some",
			}

			result, err := New(
				WithParameters(params),
				WithMoreFunctions(sprig.TxtFuncMap()),
			).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_Render_Pipe(t *testing.T) {
	Run(t, Test{
		name: "parse func",
		f: func(tt Test) {
			input := "{{ .key }}: {{ b64enc .value | b64dec }}"
			expected := "awe: some"
			params := parameters.Parameters{
				"key":   "awe",
				"value": "some",
			}

			result, err := New(
				WithParameters(params),
				WithMoreFunctions(sprig.TxtFuncMap()),
			).NamedRender(tt.name, input)

			assert.NoError(t, err, tt.name)
			assert.Equal(t, expected, result, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

func TestRenderer_FuzzFalsePositive1(t *testing.T) {
	t.Run("go-fuzz crash 1", func(t *testing.T) {
		input := "{{range $,$ =.}}{{end}}"
		output, err := New(
			WithSprigFunctions(),
			WithExtraFunctions(),
		).Render(input)
		assert.NoError(t, err)
		assert.Equal(t, "", output)
	})
}

func TestRenderer_Render_Validate_Default(t *testing.T) {
	Run(t, Test{
		name: "validation",
		f: func(tt Test) {
			err := New().Validate()
			assert.NoError(t, err, tt.name)
			assert.Equal(t, 0, CountProblems(tt.logHook))
		},
	})
}

type Test struct {
	name    string
	f       func(tt Test)
	logHook *test.Hook
}

func Run(t *testing.T, tt Test) {
	logrus.SetLevel(logrus.DebugLevel)
	hook := test.NewGlobal()
	tt.logHook = hook
	t.Run(tt.name, func(t *testing.T) { tt.f(tt) })
	hook.Reset()
}

func FilterEntries(filterOut []logrus.Level, entries []*logrus.Entry) []*logrus.Entry {
	var errors []*logrus.Entry
	for _, entry := range entries {
		skip := false
		for _, level := range filterOut {
			if entry.Level == level {
				skip = true
			}
		}
		if !skip {
			logrus.Debugf("Log: %+v", *entry)
			errors = append(errors, entry)
		}
	}
	return errors
}

func CountProblems(hook *test.Hook) int {
	return len(FilterEntries([]logrus.Level{logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}, hook.AllEntries()))
}
