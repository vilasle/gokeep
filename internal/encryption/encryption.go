package encryption

import "encoding/hex"

// Encoder is an interface for encrypting and decrypting data.
type Encoder interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

// EncryptedData is a struct for encrypting and decrypting data.
type EncryptedData struct {
	dek  Encoder
	Data string `json:"data"`
	Key  string `json:"key"`
}

// NewEncryptedDataFromReadyData creates a new EncryptedData struct from ready data.
func NewEncryptedDataFromReadyData(dek Encoder, data, key []byte) *EncryptedData {
	return &EncryptedData{
		dek:  dek,
		Data: string(data),
		Key:  string(key),
	}
}

// NewEncryptedData creates a new EncryptedData struct.
func NewEncryptedData(dek Encoder) *EncryptedData {
	return &EncryptedData{dek: dek}
}

// Encrypt encrypts the data using the DEK and encode it to hex string
func (ed *EncryptedData) Encrypt(data []byte) error {
	encData, err := ed.dek.Encrypt(data)
	if err != nil {
		return err
	}
	ed.Data = hex.EncodeToString(encData)
	return nil
}

// Decrypt decrypts the data using the DEK
func (ed *EncryptedData) Decrypt() ([]byte, error) {
	data, err := hex.DecodeString(ed.Data)
	if err != nil {
		return nil, err
	}
	return ed.dek.Decrypt(data)
}

//EncryptKey encrypts the key using the KEK and encode it to hex string
func (ed *EncryptedData) EncryptKey(enc Encoder, data []byte) error {
	encData, err := enc.Encrypt(data)
	if err != nil {
		return err
	}
	ed.Key = hex.EncodeToString(encData)
	return nil
}

//DecryptKey decrypts the key using the KEK
func (ed *EncryptedData) DecryptKey(enc Encoder) ([]byte, error) {
	data, err := hex.DecodeString(ed.Key)
	if err != nil {
		return nil, err
	}
	return enc.Decrypt(data)
}

//ReplaceKey replaces the key using the current and new encoders. New encoder will be valid key for decode DEK
func (ed *EncryptedData) ReplaceKey(current, new Encoder) error {
	key, err := ed.DecryptKey(current)
	if err != nil {
		return err
	}
	return ed.EncryptKey(new, key)
}
