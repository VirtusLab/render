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
	RootKey = "root"
)

var (
	varArgRegexp = matcher.NewMust(`^(?P<name>\S+)=(?P<value>\S*)$`)
)

type Configuration map[string]interface{}

func New(configs ...Configuration) Configuration {
	var accumulator = make(Configuration)
	for _, config := range configs {
		MergeConfigurations(&accumulator, config)
	}

	return accumulator
}

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

func WithVars(extraParams []string) (Configuration, error) {
	var config = make(Configuration)
	for _, v := range extraParams {
		if varArgRegexp.Match(v) {
			groups := varArgRegexp.MatchGroups(v)
			name := groups["name"]
			value := groups["value"]
			logrus.Debugf("Extra var: %s=%s", name, value)
			config[name] = value
		} else {
			logrus.Error("Expected a valid extra parameter: '%s'", v)
		}
	}

	logrus.Debugf("Configuration from vars: %v", config)
	return config, nil
}

func MergeConfigurations(dst *Configuration, src Configuration) error {
	err := mergo.Merge(dst, src)
	if err != nil {
		return err
	}
	return nil
}
