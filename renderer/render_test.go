package renderer

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/VirtusLab/render/renderer/configuration"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name    string
	f       func(tc TestCase)
	logHook *test.Hook
}

func TestRenderer_Render_Empty(t *testing.T) {
	Run(t, TestCase{
		name: "empty render",
		f: func(tc TestCase) {
			input := ""
			expected := ""
			config := configuration.Configuration{}

			result, err := New(config).Render(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
			assert.Equal(t, 0, len(tc.logHook.Entries))
		},
	})
}

func TestRenderer_Render_Error(t *testing.T) {
	Run(t, TestCase{
		name: "parse error",
		f: func(tc TestCase) {
			input := "{{ wrong+ }}"
			expected := ""
			config := configuration.Configuration{}

			result, err := New(config).Render(tc.name, input)

			assert.Error(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
			assert.Equal(t, 1, len(tc.logHook.Entries))
			assert.Equal(t, logrus.ErrorLevel, tc.logHook.LastEntry().Level)
			assert.Contains(t, tc.logHook.LastEntry().Message, "Can't parse the template")
		},
	})
}

func TestRenderer_Render_Render(t *testing.T) {
	Run(t, TestCase{
		name: "render render",
		f: func(tc TestCase) {
			input := "key: {{ .inner | render }}"
			expected := "key: some"
			config := configuration.Configuration{
				"inner": "{{ .value }}",
				"value": "some",
			}

			result, err := New(config).Render(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
		},
	})
}

func TestRenderer_Render_Func(t *testing.T) {
	Run(t, TestCase{
		name: "parse func",
		f: func(tc TestCase) {
			input := "key: {{ b64enc .value }}"
			expected := "key: c29tZQ=="
			config := configuration.Configuration{
				"value": "some",
			}

			result, err := New(config).Render(tc.name, input)

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
			config := configuration.Configuration{
				"key":   "awe",
				"value": "some",
			}

			result, err := New(config).Render(tc.name, input)

			assert.NoError(t, err, tc.name)
			assert.Equal(t, expected, result, tc.name)
		},
	})
}

func Run(t *testing.T, c TestCase) {
	logrus.SetLevel(logrus.DebugLevel)
	hook := test.NewGlobal()
	c.logHook = hook
	t.Run(c.name, func(t *testing.T) { c.f(c) })
	hook.Reset()
}
