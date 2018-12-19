package renderer

import (
	"testing"

	"github.com/VirtusLab/render/renderer/parameters"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestRenderer_NamedRender_Empty(t *testing.T) {
	Run(t, TestCase{
		name: "empty render",
		f: func(tc TestCase) {
			input := ""
			expected := ""

			result, err := New().NamedRender(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
			assert.Equal(t, 0, len(tc.logHook.Entries))
		},
	})
}

func TestRenderer_NamedRender_Simple(t *testing.T) {
	Run(t, TestCase{
		name: "simple render",
		f: func(tc TestCase) {
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

			result, err := New().Parameters(params).NamedRender(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
			assert.Equal(t, 3, len(tc.logHook.Entries))
		},
	})
}

func TestRenderer_Render_Error(t *testing.T) {
	Run(t, TestCase{
		name: "parse error",
		f: func(tc TestCase) {
			input := "{{ wrong+ }}"
			expected := ""

			result, err := New().NamedRender(tc.name, input)

			assert.Error(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
			assert.Equal(t, 1, len(tc.logHook.Entries))
			assert.Equal(t, logrus.ErrorLevel, tc.logHook.LastEntry().Level)
			assert.Contains(t, tc.logHook.LastEntry().Message, "Can't parse the template")
		},
	})
}

func TestRenderer_NamedRender_Render(t *testing.T) {
	Run(t, TestCase{
		name: "render render",
		f: func(tc TestCase) {
			input := "key: {{ .inner | render }}"
			expected := "key: some"
			params := parameters.Parameters{
				"inner": "{{ .value }}",
				"value": "some",
			}

			result, err := New().Parameters(params).NamedRender(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
		},
	})
}

func TestRenderer_NamedRender_Func(t *testing.T) {
	Run(t, TestCase{
		name: "parse func",
		f: func(tc TestCase) {
			input := "key: {{ b64enc .value }}"
			expected := "key: c29tZQ=="
			params := parameters.Parameters{
				"value": "some",
			}

			result, err := New().Parameters(params).NamedRender(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
		},
	})
}

func TestRenderer_Render_Pipe(t *testing.T) {
	Run(t, TestCase{
		name: "parse func",
		f: func(tc TestCase) {
			input := "{{ .key }}: {{ b64enc .value | b64dec }}"
			expected := "awe: some"
			params := parameters.Parameters{
				"key":   "awe",
				"value": "some",
			}

			result, err := New().Parameters(params).NamedRender(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
		},
	})
}

func TestRenderer_Render_Validate_Default(t *testing.T) {
	Run(t, TestCase{
		name: "validation",
		f: func(tc TestCase) {
			err := New().Validate()
			assert.NoError(t, err, tc.name)
		},
	})
}

type TestCase struct {
	name    string
	f       func(tc TestCase)
	logHook *test.Hook
}

func Run(t *testing.T, c TestCase) {
	logrus.SetLevel(logrus.DebugLevel)
	hook := test.NewGlobal()
	c.logHook = hook
	t.Run(c.name, func(t *testing.T) { c.f(c) })
	hook.Reset()
}
