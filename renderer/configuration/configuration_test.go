package configuration

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		configs []Configuration
	}
	type test struct {
		name    string
		args    args
		want    Configuration
		wantErr *error
		f       func(tt test)
	}

	standard := func(tt test) {
		got, err := New(tt.args.configs...)
		if tt.wantErr != nil {
			assert.EqualError(t, err, (*tt.wantErr).Error())
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tt.want, got)
	}

	mustWithVars := func(vars []string) Configuration {
		c, e := WithVars(vars)
		if e != nil {
			t.Fatal("invalid test input")
		}
		return c
	}

	tests := []test{
		{
			name: "empty",
			args: args{
				[]Configuration{},
			},
			want: Configuration{},
			f:    standard,
		}, {
			name: "empty merge",
			args: args{
				[]Configuration{
					{},
					mustWithVars([]string{}),
				},
			},
			want: Configuration{},
			f:    standard,
		}, {
			name: "single",
			args: args{
				[]Configuration{
					{"akey": "avalue"},
				},
			},
			want: Configuration{
				"akey": "avalue",
			},
			f: standard,
		}, {
			name: "two",
			args: args{
				[]Configuration{
					{"akey": "avalue"},
					{"another": "entry"},
				},
			},
			want: Configuration{
				"akey":    "avalue",
				"another": "entry",
			},
			f: standard,
		}, {
			name: "merge",
			args: args{
				[]Configuration{
					{"akey": "avalue"},
					{"akey": "overriden"},
				},
			},
			want: Configuration{
				"akey": "overriden",
			},
			f: standard,
		}, {
			name: "merge with vars",
			args: args{
				[]Configuration{
					{"akey": "avalue"},
					{"akey": "overriden"},
					mustWithVars([]string{}),
				},
			},
			want: Configuration{
				"akey": "overriden",
			},
			f: standard,
		},
	}

	logrus.SetLevel(logrus.DebugLevel)

	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tt.name), func(t *testing.T) { tt.f(tt) })
	}
}
