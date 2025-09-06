package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

type AESKey struct {
	Key   string `json:"key"`
	Nonce string `json:"nonce"`
	key   []byte
	nonce []byte
	gcm   cipher.AEAD
}

//GenerateNewAESKey generates a new AES key
func GenerateNewAESKey() (*AESKey, error) {
	key := make([]byte, 32)
	if _, err := rand.Reader.Read(key); err != nil {
		return nil, errors.Join(errors.New("error generating random encryption key"), err)
	}

	aesKey := &AESKey{
		Key:   hex.EncodeToString(key),
		Nonce: "",
	}

	if err := aesKey.initGCM(); err != nil {
		return nil, err
	}

	nonce := make([]byte, aesKey.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, errors.Join(errors.New("error generating random nonce"), err)
	}
	aesKey.Nonce = hex.EncodeToString(nonce)
	aesKey.nonce = nonce
	return aesKey, nil
}

//NewAESKeyFromJSON creates a new AES key from a JSON string
func NewAESKeyFromJSON(content []byte) (*AESKey, error) {
	key := &AESKey{}
	err := json.Unmarshal(content, key)
	if err != nil {
		return nil, err
	}

	if err := key.initGCM(); err != nil {
		return nil, err
	}

	return key, nil
}

func (k *AESKey) initGCM() error {
	key, err := hex.DecodeString(k.Key)
	if err != nil {
		return err
	}

	nonce, err := hex.DecodeString(k.Nonce)
	if err != nil {
		return err
	}

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return err
	}

	k.gcm = aesGcm
	k.key = key
	k.nonce = nonce
	return nil
}

//Encrypt encrypts the data using the AES key
func (k *AESKey) Encrypt(data []byte) ([]byte, error) {
	return k.gcm.Seal(nil, k.nonce, data, nil), nil
}

//Decrypt decrypts the data using the AES key
func (k *AESKey) Decrypt(data []byte) ([]byte, error) {
	return k.gcm.Open(nil, k.nonce, data, nil)
}

//JSON returns the JSON representation of the AES key
func (k AESKey) JSON() string {
	buf := strings.Builder{}
	buf.WriteString(`{"key":"`)
	buf.WriteString(k.Key)
	buf.WriteString(`","nonce":"`)
	buf.WriteString(k.Nonce)
	buf.WriteString(`"}`)
	return buf.String()
}
