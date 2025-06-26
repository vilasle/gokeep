package encryption

import "encoding/hex"

type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type EncryptedData struct {
	enc  Encryptor
	Data string `json:"data"`
	Key  string `json:"key"`
}

func (ed *EncryptedData) Encrypt(data []byte) error {
	encData, err := ed.enc.Encrypt(data)
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
	return ed.enc.Decrypt(data)
}

func (ed *EncryptedData) EncryptKey(enc Encryptor, data []byte) error {
	encData, err := enc.Encrypt(data)
	if err != nil {
		return err
	}
	ed.Key = hex.EncodeToString(encData)
	return nil
}

func (ed *EncryptedData) DecryptKey(enc Encryptor) ([]byte, error) {
	data, err := hex.DecodeString(ed.Key)
	if err != nil {
		return nil, err
	}
	return enc.Decrypt(data)
}

func (ed *EncryptedData) ReplaceKey(current, new Encryptor) error {
	key, err := ed.DecryptKey(current)
	if err != nil {
		return err
	}
	return ed.EncryptKey(new, key)
}
