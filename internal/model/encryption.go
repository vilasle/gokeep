package model

type Encrypter interface {
	Encrypt([]byte) (*EncryptedData, error)
	Decrypt(EncryptedData) ([]byte, error)
}

type EncryptedData struct {
	ID    int64
	Owner PrivateData
	Data  []byte
	Key   []byte
}