package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"errors"
)

type RSACipher struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRSACipherFromRawPublicKey creates a new RSACipher from a raw public key without private key
func NewRSACipherFromRawPublicKey(publicKey []byte) (*RSACipher, error) {
	spkiBlock, _ := pem.Decode(publicKey)
	if spkiBlock == nil {
		return nil, errors.New("public key is not valid")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(spkiBlock.Bytes)
	if err != nil {
		return nil, err
	}
	key := pubInterface.(*rsa.PublicKey)

	return &RSACipher{
		publicKey:  key,
		privateKey: nil,
	}, nil
}

// NewRSACipher creates a new RSACipher from a public key and a private key
func NewRSACipher(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *RSACipher {
	return &RSACipher{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

// Encrypt encrypts the data using the public key
func (c *RSACipher) Encrypt(data []byte) ([]byte, error) {
	if c.publicKey == nil {
		return nil, errors.New("public key is not set")
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, c.publicKey, data, nil)
}

// Decrypt decrypts the data using the private key
func (c *RSACipher) Decrypt(data []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("private key is not set")
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privateKey, data, nil)
}
