package configuration

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/matcher"
	"github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

const (
	// RootKey is an special configuration key key used by e.g. the Base and Root functions
	RootKey = "root"
)

var (
	// VarArgRegexp defines the extra variable parameter format
	VarArgRegexp = matcher.NewMust(`^(?P<name>\S+)=(?P<value>\S*)$`)
)

// Configuration is a map used to render the templates with
type Configuration map[string]interface{}

// Validate checks the internal state and returns error if necessary
func (configurations Configuration) Validate() error {
	// TODO(pprazak): could/should we do some validation here?
	return nil
}

// New creates a new configuration from one or more configurations, to be used with other helper functions
func New(configs ...Configuration) Configuration {
	var accumulator = make(Configuration)
	for _, config := range configs {
		MergeConfigurations(&accumulator, config)
	}

	return accumulator
}

// All creates a configuration from one or more configuration file paths
// and one or more extra variables in addition to base configuration
func All(configPaths, vars []string) (Configuration, error) {
	baseConfig, err := Base()
	if err != nil {
		return nil, errors.Wrap(err, "can't create base configuration")
	}

	filesConfig, err := WithFiles(configPaths)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse configuration filse")
	}

	varsConfig, err := WithVars(vars)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse extra configuration variables")
	}

	return New(baseConfig, filesConfig, varsConfig), nil
}

// Base creates a basic configuration required for some of the functions, it is recommended to use it
func Base() (Configuration, error) {
	pwd, err := files.Pwd()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get PWD")
	}
	c := Configuration{
		RootKey: pwd,
	}
	logrus.Debugf("Base configuration: %v", c)
	return c, nil
}

// WithFiles creates a configuration from one or more configuration file paths
func WithFiles(configPaths []string) (Configuration, error) {
	var accumulator = make(Configuration)
	for i, configPath := range configPaths {
		logrus.Debugf("Reading configuration file [%d]: %v", i, configPath)
		if files.IsNotEmptyAndExists(configPath) {
			b, err := ioutil.ReadFile(configPath)
			if err != nil {
				logrus.Errorf("Can't open the configuration file: %v", err)
				return nil, err
			}
			var config map[string]interface{}
			err = yaml.Unmarshal(b, &config)
			if err != nil {
				logrus.Errorf("Can't parse the configuration file: %v", err)
				return nil, err
			}
			MergeConfigurations(&accumulator, config)
		}
	}
	logrus.Debugf("Configuration from files: %v", accumulator)

	return accumulator, nil
}

// WithVars creates a configuration from one or more extra variables (key=value), see also VarArgRegexp
func WithVars(extraParams []string) (Configuration, error) {
	var config = make(Configuration)
	for _, v := range extraParams {
		if VarArgRegexp.Match(v) {
			groups := VarArgRegexp.MatchGroups(v)
			name := groups["name"]
			value := groups["value"]
			logrus.Debugf("Extra var: %s=%s", name, value)
			config[name] = value
		} else {
			logrus.Errorf("Expected a valid extra parameter: '%s'", v)
		}
	}

	logrus.Debugf("Configuration from vars: %v", config)
	return config, nil
}

// MergeConfigurations merges two configurations into one, any existing values will be overridden
func MergeConfigurations(dst *Configuration, src Configuration) error {
	err := mergo.Merge(dst, src)
	if err != nil {
		return err
	}
	return nil
}
