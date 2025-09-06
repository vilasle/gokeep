package model

// Encoder is an interface for encrypting and decrypting data.
type Encoder interface {
	Encrypt([]byte) (*EncryptedData, error)
	Decrypt(EncryptedData) ([]byte, error)
}

type EncryptedData struct {
	Data  []byte
	Key   []byte
}