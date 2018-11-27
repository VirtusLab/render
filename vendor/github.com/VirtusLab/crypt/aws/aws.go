package aws

import (
	"errors"

	"github.com/VirtusLab/crypt/crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/sirupsen/logrus"
)

var (
	// ErrKmsMissing - this is the custom error, returned when name, alias or arn is missing
	ErrKmsMissing = errors.New("kms is empty or missing")
	// ErrRegionMissing - this is the custom error, returned when the region is missing
	ErrRegionMissing = errors.New("region is empty or missing")
)

// KMS struct represents AWS Key Management Service
type KMS struct {
	region string
	key    string
}

// New creates a AWS KMS provider
func New(key, region string) crypto.KMS {
	return &KMS{
		key:    key,
		region: region,
	}
}

// Encrypt is responsible for encrypting plaintext and returning ciphertext in bytes using AWS KMS.
// See Crypt.Encrypt
func (k *KMS) Encrypt(plaintext []byte) ([]byte, error) {
	if len(k.key) == 0 {
		logrus.Debugf("Error reading kms: %v", k.key)
		return nil, ErrKmsMissing
	}

	if len(k.region) == 0 {
		logrus.Debugf("Error reading region: %v", k.region)
		return nil, ErrRegionMissing
	}

	// use AWS_DEFAULT_PROFILE environment variable to set profile
	//
	// If not set and environment variables are not set the "default"
	// will be used as the profile to load the session config from.
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: session.DefaultSharedConfigProfile,
		Config:  aws.Config{Region: aws.String(k.region)},
	}))
	svc := kms.New(awsSession, aws.NewConfig().WithRegion(k.region))
	input := &kms.EncryptInput{
		Plaintext: plaintext,
		KeyId:     aws.String(k.key),
	}
	output, err := svc.Encrypt(input)
	if err != nil {
		return nil, err
	}
	return []byte(output.CiphertextBlob), nil
}

// Decrypt is responsible for decrypting ciphertext and returning plaintext in bytes using AWS KMS.
// See Crypt.Decrypt
func (k *KMS) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(k.region) == 0 {
		logrus.Debugf("Error reading region: %v", k.region)
		return nil, ErrRegionMissing
	}

	// use AWS_DEFAULT_PROFILE environment variable to set profile
	//
	// If not set and environment variables are not set the "default"
	// will be used as the profile to load the session config from.
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: session.DefaultSharedConfigProfile,
		Config:  aws.Config{Region: aws.String(k.region)},
	}))
	svc := kms.New(awsSession)
	input := &kms.DecryptInput{
		CiphertextBlob: ciphertext,
	}
	output, err := svc.Decrypt(input)
	if err != nil {
		return nil, err
	}
	return output.Plaintext, nil
}
