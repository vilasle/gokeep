package encryption

import "encoding/hex"

type Encoder interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type EncryptedData struct {
	dek  Encoder
	Data string `json:"data"`
	Key  string `json:"key"`
}

func NewEncryptedDataFromReadyData(dek Encoder, data, key []byte) *EncryptedData {
	return &EncryptedData{
		dek:  dek,
		Data: string(data),
		Key:  string(key),
	}
}

func NewEncryptedData(dek Encoder) *EncryptedData {
	return &EncryptedData{dek: dek}
}

func (ed *EncryptedData) Encrypt(data []byte) error {
	encData, err := ed.dek.Encrypt(data)
	if err != nil {
		return err
	}
	ed.Data = hex.EncodeToString(encData)
	return nil
}

func (ed *EncryptedData) Decrypt() ([]byte, error) {
	data, err := hex.DecodeString(ed.Data)
	if err != nil {
		return nil, err
	}
	return ed.dek.Decrypt(data)
}

func (ed *EncryptedData) EncryptKey(enc Encoder, data []byte) error {
	encData, err := enc.Encrypt(data)
	if err != nil {
		return err
	}
	ed.Key = hex.EncodeToString(encData)
	return nil
}

func (ed *EncryptedData) DecryptKey(enc Encoder) ([]byte, error) {
	data, err := hex.DecodeString(ed.Key)
	if err != nil {
		return nil, err
	}
	return enc.Decrypt(data)
}

func (ed *EncryptedData) ReplaceKey(current, new Encoder) error {
	key, err := ed.DecryptKey(current)
	if err != nil {
		return err
	}
	return ed.EncryptKey(new, key)
}
