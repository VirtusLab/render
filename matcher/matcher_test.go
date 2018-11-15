package matcher

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMatcher_MatchGroups(t *testing.T) {
	type TestCase struct {
		name  string
		expr  string
		pairs map[string]map[string]string
		f     func(tc TestCase)
	}

	standard := func(tc TestCase) {
		m, err := New(tc.expr)
		if err != nil {
			t.Errorf("[%s] unexpected error: %s", tc.name, err.Error())
		}
		for value, expected := range tc.pairs {
			result := m.MatchGroups(value)
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("[%s] got \n%+v\n, expected: \n%+v\n", tc.name, result, expected)
			}
		}
	}

	cases := []TestCase{
		{
			name: "empty expression",
			expr: ``,
			pairs: map[string]map[string]string{
				"": {},
			},
			f: standard,
		}, {
			name: "simple expression",
			expr: `^(?P<name>\S+)=(?P<value>\S*)$`,
			pairs: map[string]map[string]string{
				"test=something": {
					"name":  "test",
					"value": "something",
				},
			},
			f: standard,
		}, {
			name: "git url expression",
			expr: `^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`,
			pairs: map[string]map[string]string{
				"":                              {},
				"invalid":                       {},
				"git@something.com:anorg/arepo": {},
				"git@something.com:anorg/arepo.git": {
					"hostname":     "something.com",
					"organisation": "anorg",
					"name":         "arepo",
				},
			},
			f: standard,
		}, {
			name: "compile fail",
			expr: "<?:[",
			f: func(tc TestCase) {
				_, err := New(tc.expr)
				if err == nil {
					t.Errorf("[%s] expected an error", tc.name)
				}
			},
		},
	}

	logrus.SetLevel(logrus.DebugLevel)

	for i, c := range cases {
		t.Run(fmt.Sprintf("[%d] %s", i, c.name), func(t *testing.T) { c.f(c) })
	}
}
