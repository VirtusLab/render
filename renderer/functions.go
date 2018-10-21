package renderer

import "io/ioutil"

// ReadFile provides a custom template function for in-template file opening
func ReadFile(file string) (string, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// TODO: gzip, ungzip, encrypt, decrypt
