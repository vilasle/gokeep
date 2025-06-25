package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"strings"
)

type AESKey struct {
	Key   string `json:"key"`
	Nonce string `json:"nonce"`
	key   []byte
	nonce []byte
	gcm   cipher.AEAD
}

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

func (k *AESKey) Encrypt(data []byte) ([]byte, error) {
	return k.gcm.Seal(nil, k.nonce, data, nil), nil
}

func (k *AESKey) Decrypt(data []byte) ([]byte, error) {
	return k.gcm.Open(nil, k.nonce, data, nil)
}

func (k AESKey) JSON() string {
	buf := strings.Builder{}
	buf.WriteString(`{"key":"`)
	buf.WriteString(k.Key)
	buf.WriteString(`","nonce":"`)
	buf.WriteString(k.Nonce)
	buf.WriteString(`"}`)
	return buf.String()
}
