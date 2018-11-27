package renderer

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/VirtusLab/crypt/aws"
	"github.com/VirtusLab/crypt/azure"
	"github.com/VirtusLab/crypt/gcp"
	"github.com/VirtusLab/render/files"
	"github.com/VirtusLab/render/renderer/configuration"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

func (r *Renderer) root() (string, error) {
	if value, ok := r.configuration[configuration.RootKey].(string); ok {
		return value, nil
	}
	return files.Pwd()
}

// ReadFile is a template function that allows for an in-template file opening
func (r *Renderer) ReadFile(file string) (string, error) {
	root, err := r.root()
	if err != nil {
		return "", err
	}
	absPath, err := files.ToAbsPath(file, root)
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// ToYaml is a template function, it turns a marshallable structure into a YAML fragment
func ToYaml(marshallable interface{}) (string, error) {
	marshaledYaml, err := yaml.Marshal(marshallable)
	return string(marshaledYaml), err
}

// Gzip compresses the input using gzip algorithm
func Gzip(input interface{}) (string, error) {
	inputAsBytes, err := asBytes(input)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	_, err = w.Write(inputAsBytes)
	if err != nil {
		return "", err
	}

	err = w.Flush()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// Ungzip un-compresses the input using gzip algorithm
func Ungzip(input interface{}) (string, error) {
	inputAsBytes, err := asBytes(input)
	if err != nil {
		return "", err
	}

	in := bytes.NewBuffer(inputAsBytes)
	r, err := gzip.NewReader(in)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func asBytes(input interface{}) ([]byte, error) {
	switch input.(type) {
	case []byte:
		return input.([]byte), nil
	case string:
		return []byte(input.(string)), nil
	default:
		return nil, errors.Errorf("expected []byte or string, got: '%v'", reflect.TypeOf(input))
	}
}

// EncryptAWS encrypts plaintext using AWS KMS
func EncryptAWS(awsKms, awsRegion, plaintext string) ([]byte, error) {
	amazon := aws.New(awsKms, awsRegion)
	result, err := amazon.Encrypt([]byte(plaintext))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecryptAWS decrypts ciphertext using AWS KMS
func DecryptAWS(awsRegion, ciphertext string) (string, error) {
	amazon := aws.New("" /* not needed for decryption */, awsRegion)
	result, err := amazon.Decrypt([]byte(ciphertext))
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// EncryptGCP encrypts plaintext using GCP KMS
func EncryptGCP(gcpProject, gcpLocation, gcpKeyring, gcpKey, plaintext string) ([]byte, error) {
	googleKms := gcp.New(gcpProject, gcpLocation, gcpKeyring, gcpKey)
	result, err := googleKms.Encrypt([]byte(plaintext))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecryptGCP decrypts ciphertext using GCP KMS
func DecryptGCP(gcpProject, gcpLocation, gcpKeyring, gcpKey, ciphertext string) (string, error) {
	googleKms := gcp.New(gcpProject, gcpLocation, gcpKeyring, gcpKey)
	result, err := googleKms.Decrypt([]byte(ciphertext))
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// EncryptAzure encrypts plaintext using Azure Key Vault
func EncryptAzure(azureVaultURL, azureKey, azureKeyVersion, plaintext string) ([]byte, error) {
	azr := azure.New(azureVaultURL, azureKey, azureKeyVersion)
	result, err := azr.Encrypt([]byte(plaintext))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecryptAzure decrypts ciphertext using Azure Key Vault
func DecryptAzure(azureVaultURL, azureKey, azureKeyVersion, ciphertext string) (string, error) {
	azr := azure.New(azureVaultURL, azureKey, azureKeyVersion)
	result, err := azr.Decrypt([]byte(ciphertext))
	if err != nil {
		return "", err
	}
	return string(result), nil
}
