package azure

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/VirtusLab/crypt/crypto"
	"github.com/sirupsen/logrus"
)

var (
	// ErrVaultURLMissing - this is the custom error, returned when vault url is missing
	ErrVaultURLMissing = errors.New("vault url is empty or missing")
	// ErrKeyMissing = this is the custom error, returned when the key is missing
	ErrKeyMissing = errors.New("key is empty or missing")
	// ErrKeyVersionMissing = this is the custom error, returned when the key version is missing
	ErrKeyVersionMissing = errors.New("key version is empty or missing")
)

// KMS struct represents Azure Key Vault
type KMS struct {
	vaultURL   string
	key        string
	keyVersion string
}

// New creates Azure Key Vault KMS
func New(vaultURL, key, keyVersion string) crypto.KMS {
	return &KMS{
		vaultURL:   vaultURL,
		key:        key,
		keyVersion: keyVersion,
	}
}

func newKeyVaultClient() (keyvault.BaseClient, error) {
	var err error
	vaultClient := keyvault.New()
	vaultClient.Authorizer, err = auth.NewAuthorizerFromEnvironment()
	if err != nil {
		logrus.WithError(err).Error("Failed to create Azure Authorizer")
		return vaultClient, err
	}
	return vaultClient, nil
}

// Encrypt is encrypts plaintext using Azure Key Vault and returns ciphertext
// See Crypt.Encrypt
func (k *KMS) Encrypt(plaintext []byte) ([]byte, error) {
	err := k.validate()
	if err != nil {
		return nil, err
	}

	client, err := newKeyVaultClient()
	if err != nil {
		return nil, err
	}
	data := base64.RawURLEncoding.EncodeToString(plaintext)
	p := keyvault.KeyOperationsParameters{Value: &data, Algorithm: keyvault.RSAOAEP256}

	ctx := context.Background()
	res, err := client.Encrypt(ctx, k.vaultURL, k.key, k.keyVersion, p)
	if err != nil {
		return nil, err
	}

	result, err := base64.RawURLEncoding.DecodeString(*res.Result)
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"key":        k.key,
		"keyVersion": k.keyVersion,
	}).Info("Encryption succeeded")
	return result, nil
}

// Decrypt is responsible for decrypting ciphertext by Azure Key Vault encryption key and returning plaintext in bytes.
// See Crypt.EncryptFile
func (k *KMS) Decrypt(ciphertext []byte) ([]byte, error) {
	// FIXME k.validateParams()
	client, err := newKeyVaultClient()
	if err != nil {
		return nil, err
	}
	data := base64.RawURLEncoding.EncodeToString(ciphertext)
	p := keyvault.KeyOperationsParameters{Value: &data, Algorithm: keyvault.RSAOAEP256}

	ctx := context.Background()
	res, err := client.Decrypt(ctx, k.vaultURL, k.key, k.keyVersion, p)
	if err != nil {
		return nil, err
	}

	plaintext, err := base64.RawURLEncoding.DecodeString(*res.Result)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"key":        k.key,
		"keyVersion": k.keyVersion,
	}).Info("Decryption succeeded")

	return plaintext, nil
}

func (k *KMS) validate() error {
	if len(k.vaultURL) == 0 {
		logrus.Debugf("Error reading vaultURL: %v", k.vaultURL)
		return ErrVaultURLMissing
	}
	if len(k.key) == 0 {
		logrus.Debugf("Error reading key: %v", k.key)
		return ErrKeyMissing
	}
	if len(k.keyVersion) == 0 {
		logrus.Debugf("Error reading keyVersion: %v", k.keyVersion)
		return ErrKeyVersionMissing
	}
	return nil
}
