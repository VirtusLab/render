package matcher

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMatcher_MatchGroups(t *testing.T) {
	type entry struct {
		value  string
		wantOk bool
		want   map[string]string
	}
	type test struct {
		name    string
		expr    string
		entries []entry
		f       func(tc test)
	}

	standard := func(tt test) {
		m, err := New(tt.expr)
		assert.NoError(t, err)
		for _, e := range tt.entries {
			actual, ok := m.MatchGroups(e.value)
			assert.Equal(t, e.wantOk, ok, "test: '%s', entry: '%s'", tt.name, e.value)
			assert.EqualValues(t, e.want, actual, "test: '%s', entry: '%s'", tt.name, e.value)
		}
	}

	tests := []test{
		{
			name: "empty expression",
			expr: ``,
			entries: []entry{
				{
					value:  "",
					wantOk: true,
					want:   map[string]string{},
				},
			},
			f: standard,
		}, {
			name: "simple expression",
			expr: `^(?P<name>\S+)=(?P<value>\S*)$`,
			entries: []entry{
				{
					value:  "test=something",
					wantOk: true,
					want: map[string]string{
						"name":  "test",
						"value": "something",
					},
				},
			},
			f: standard,
		}, {
			name: "git url expression",
			expr: `^git@(?P<hostname>[\w\-\.]+):(?P<organisation>[\w\-]+)\/(?P<name>[\w\-]+)\.git$`,
			entries: []entry{
				{value: "", wantOk: false, want: map[string]string{}},
				{value: "invalid", wantOk: false, want: map[string]string{}},
				{value: "git@something.com:anorg/arepo", wantOk: false, want: map[string]string{}},
				{value: "git@something.com:anorg/arepo.git", wantOk: true, want: map[string]string{
					"hostname":     "something.com",
					"organisation": "anorg",
					"name":         "arepo",
				}},
			},
			f: standard,
		}, {
			name: "compile fail",
			expr: "<?:[",
			f: func(tc test) {
				_, err := New(tc.expr)
				assert.Error(t, err)
			},
		},
	}

	logrus.SetLevel(logrus.DebugLevel)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tt.name), func(t *testing.T) { tt.f(tt) })
	}
}

func Test_matcher_Match(t *testing.T) {
	type fields struct {
		matcher Matcher
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "empty",
			fields: fields{matcher: Must("^[a-z]+[0-9]+")},
			args: args{
				value: "",
			},
			want: false,
		}, {
			name:   "simple match",
			fields: fields{matcher: Must("^[a-z]+[0-9]+")},
			args: args{
				value: "asdf1234",
			},
			want: true,
		},
		{
			name:   "no match",
			fields: fields{matcher: Must("^[a-z]+[0-9]+")},
			args: args{
				value: "1234asdf",
			},
			want: false,
		},
	}

	logrus.SetLevel(logrus.DebugLevel)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.matcher
			got := m.Match(tt.args.value)
			assert.Equal(t, tt.want, got, "[%d] matcher.Match() = %v, want %v", i, got, tt.want)
		})
	}
}
