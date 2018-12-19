package parameters

import (
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/matcher"

	"github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// RootKey is an special configuration key key used by e.g. the Base and Root functions
	RootKey = "root"
)

var (
	// VarArgRegexp defines the extra variable parameter format
	VarArgRegexp = matcher.Must(`^(?P<name>\S+)=(?P<value>\S*)$`)
)

// Parameters is a map used to render the templates with
type Parameters map[string]interface{}

// Validate checks the internal state and returns error if necessary
func (parameters Parameters) Validate() error {
	// TODO(pprazak): could/should we do some validation here?
	return nil
}

// Merge creates a new parameters from one or more parameter sets, to be used with other helper functions
func Merge(parameters ...Parameters) (Parameters, error) {
	var accumulator = make(Parameters)
	for _, config := range parameters {
		err := merge(&accumulator, config)
		if err != nil {
			return nil, err
		}
	}

	return accumulator, nil
}

// All creates a configuration from one or more configuration file paths
// and one or more extra variables in addition to base configuration
func All(configPaths, vars []string) (Parameters, error) {
	baseConfig, err := Base()
	if err != nil {
		return nil, errors.Wrap(err, "can't create base configuration")
	}

	filesConfig, err := FromFiles(configPaths)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse configuration filse")
	}

	varsConfig, err := FromVars(vars)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse extra configuration variables")
	}

	return Merge(baseConfig, filesConfig, varsConfig)
}

// Base creates a basic configuration required for some of the functions, it is recommended to use it
func Base() (Parameters, error) {
	pwd, err := files.Pwd()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get PWD")
	}
	c := Parameters{
		RootKey: pwd,
	}
	logrus.Debugf("Base configuration: %v", c)
	return c, nil
}

// FromFiles creates a configuration from one or more configuration file paths
func FromFiles(configPaths []string) (Parameters, error) {
	var accumulator = make(Parameters)
	for i, configPath := range configPaths {
		logrus.Debugf("Reading configuration file [%d]: %v", i, configPath)
		if files.IsNotEmptyAndExists(configPath) {
			b, err := ioutil.ReadFile(configPath)
			if err != nil {
				logrus.Errorf("Can't open the configuration file: %v", err)
				return nil, errors.WithStack(err)
			}
			var config map[string]interface{}
			err = yaml.Unmarshal(b, &config)
			if err != nil {
				logrus.Errorf("Can't parse the configuration file: %v", err)
				return nil, errors.WithStack(err)
			}
			err = merge(&accumulator, config)
			if err != nil {
				return nil, err
			}
		}
	}
	logrus.Debugf("Parameters from files: %v", accumulator)

	return accumulator, nil
}

// FromVars creates a configuration from one or more extra variables (key=value), see also VarArgRegexp
func FromVars(extraParams []string) (Parameters, error) {
	var config = &Parameters{}
	for _, v := range extraParams {
		groups, ok := VarArgRegexp.MatchGroups(v)
		if !ok {
			logrus.Errorf("Expected a valid extra parameter: '%s'", v)
			return nil, errors.Errorf("invalid parameter: '%s'", v)
		}
		name := groups["name"]
		value := groups["value"]
		logrus.Debugf("Extra var: %s=%s", name, value)
		isNested := strings.Contains(name, ".")
		if isNested {
			logrus.Debugf("Extra var key is nested: %s", name)
			var err error
			config, err = appendNested(config, name, value)
			if err != nil {
				return nil, err
			}
		} else {
			(*config)[name] = value
		}
	}

	logrus.Debugf("Parameters from vars: %v", *config)
	return *config, nil
}

func appendNested(parameters *Parameters, nestedKey string, nestedValue interface{}) (*Parameters, error) {
	if parameters == nil {
		return nil, errors.New("unexpected nil parameters")
	}
	if len(nestedKey) == 0 {
		return parameters, errors.New("unexpected empty nestedKey")
	}
	keys := strings.Split(nestedKey, ".")
	lastIndex := len(keys) - 1

	var current = parameters
	for i, key := range keys {
		// assign nested value to the last key
		if i == lastIndex {
			(*current)[key] = nestedValue
			continue
		}
		// get or create value for current key[i]
		value, ok := (*current)[key]
		if !ok {
			(*current)[key] = Parameters{}
		}
		// get value as a map
		newCurrent, ok := (*current)[key].(Parameters)
		if value != nil && !ok {
			return nil, errors.Errorf(
				"key conflict: key '%s' already exists and is not a map, it has type: '%s'",
				key, reflect.TypeOf(value))
		}
		// assign nested map as current
		current = &newCurrent
	}
	return parameters, nil
}

func merge(dst *Parameters, src Parameters) error {
	err := mergo.Merge(dst, src, mergo.WithOverride)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
