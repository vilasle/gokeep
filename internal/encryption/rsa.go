package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"errors"
)

/*
	var aesCipher, rsaCipher Cipher
	var data []byte

	dek := aes.NewCipher(key)
	encryptedData := dek.Encrypt(data)

	serverDEKKey := aesCipher.Encrypt(dek)
	encData := createEncData(encryptedData, serverDEKKey)

	err := encData.ReplaceMasterKey(aesCipher, rsaCipher)

*/

type RSACipher struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRSACipher(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *RSACipher {
	return &RSACipher{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (c *RSACipher) Encrypt(data []byte) ([]byte, error) {
	if c.publicKey == nil {
		//TODO use a package error
		return nil, errors.New("public key is not set")
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, c.publicKey, data, nil)
}

func (c *RSACipher) Decrypt(data []byte) ([]byte, error) {
	if c.privateKey == nil {
		//TODO use a package error
		return nil, errors.New("private key is not set")
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privateKey, data, nil)
}
