package model

type Encrypter interface {
	Encrypt([]byte) (*EncryptedData, error)
	Decrypt(EncryptedData) ([]byte, error)

	EncryptKey(EncryptedData) ([]byte, error)
	DecryptKey([]byte) (*EncryptedData, error)
}

type EncryptedData struct {
	id    int64
	owner PrivateData
	data  []byte
	key   []byte
}

func (e EncryptedData) ID() int64 {
	return e.id
}

func (e EncryptedData) Owner() PrivateData {
	return e.owner
}

func (e EncryptedData) Data() []byte {
	return e.data
}

func (e EncryptedData) Key() []byte {
	return e.key
}
