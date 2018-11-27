package crypto

import (
	"github.com/VirtusLab/crypt/files"
	"github.com/sirupsen/logrus"
)

// KMS (Key Management Service) is a common abstraction for encryption and decryption.
// A KMS must be able to decrypt the data it encrypts.
type KMS interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// Crypt type represents the crypt abstraction for simple encryption and decryption.
// A provider (e.g. AWS KMS) determines the detail of the cryptographic operations.
type Crypt struct {
	kms KMS
}

// New creates a new Crypt with the given provider
func New(kms KMS) *Crypt {
	return &Crypt{kms: kms}
}

// EncryptFile encrypts bytes from a file or stdin using a KMS provider
// and the ciphertext is saved into a file.
// If inputPath is empty, stdin is used as input
// If outputPath is empty, stdout is used as output
func (c *Crypt) EncryptFile(inputPath, outputPath string) error {
	input, err := files.ReadInput(inputPath)
	if err != nil {
		logrus.Debugf("Can't open plaintext file: %v", err)
		return err
	}
	result, err := c.Encrypt(input)
	if err != nil {
		logrus.Debugf("Encrypting failed: %s", err)
		return err
	}
	err = files.WriteOutput(outputPath, result, 0644) // 0644 - user: read&write, group: read, other: read
	if err != nil {
		logrus.Debugf("Can't save the encrypted file: %v", err)
		return err
	}
	return nil
}

// DecryptFile reads from the inputPath file or stdin if empty.
// Then decrypts content with corresponding Key Management Service.
// Plaintext is saved into outputPath file or print on stdout if empty.
func (c *Crypt) DecryptFile(inputPath, outputPath string) error {
	input, err := files.ReadInput(inputPath)
	if err != nil {
		logrus.Debugf("Can't open encrypted file: %v", err)
		return err
	}
	result, err := c.Decrypt(input)
	if err != nil {
		logrus.Debugf("Decrypting failed: %s", err)
		return err
	}
	err = files.WriteOutput(outputPath, result, 0644) // 0644 - user: read&write, group: read, other: read
	if err != nil {
		logrus.Debugf("Can't save the decrypted file: %v", err)
		return err
	}
	return nil
}

// Decrypt decrypts given bytes using the current provider
func (c *Crypt) Decrypt(input []byte) ([]byte, error) {
	return c.kms.Decrypt(input)
}

// Encrypt encrypts given bytes using the current provider
func (c *Crypt) Encrypt(input []byte) ([]byte, error) {
	return c.kms.Encrypt(input)
}
