package renderer

import (
	"fmt"
	"testing"

	"github.com/VirtusLab/render/renderer/configuration"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestRenderer_Render(t *testing.T) {
	type TestCase struct {
		name    string
		f       func(TestCase)
		logHook *test.Hook
	}

	when := func(config configuration.Configuration, templateName, rawTemplate string) (string, error) {
		r := New(config)
		return r.Render("test-"+templateName, rawTemplate)
	}

	cases := []TestCase{
		{
			name: "empty render",
			f: func(tc TestCase) {
				input := ""
				expected := ""
				config := configuration.Configuration{}

				result, err := when(config, tc.name, input)

				assert.NoError(t, err, tc.name)
				assert.Equal(t, expected, result, tc.name)
				assert.Equal(t, 0, len(tc.logHook.Entries))
			},
		}, {
			name: "parse error",
			f: func(tc TestCase) {
				input := "{{ wrong+ }}"
				expected := ""
				config := configuration.Configuration{}

				result, err := when(config, tc.name, input)

				assert.Error(t, err, tc.name)
				assert.Equal(t, expected, result, tc.name)
				assert.Equal(t, 1, len(tc.logHook.Entries))
				assert.Equal(t, logrus.ErrorLevel, tc.logHook.LastEntry().Level)
				assert.Contains(t, tc.logHook.LastEntry().Message, "Can't parse the template file:")
			},
		}, {
			name: "render render",
			f: func(tc TestCase) {
				input := "key: {{ .inner | render }}"
				expected := "key: some"
				config := configuration.Configuration{
					"inner": "{{ .value }}",
					"value": "some",
				}

				result, err := when(config, tc.name, input)

				assert.NoError(t, err, tc.name)
				assert.Equal(t, expected, result, tc.name)
			},
		},
	}

	logrus.SetLevel(logrus.DebugLevel)
	hook := test.NewGlobal()

	for i, c := range cases {
		c.logHook = hook
		t.Run(fmt.Sprintf("[%d] %s", i, c.name), func(t *testing.T) { c.f(c) })
		hook.Reset()
	}
}
