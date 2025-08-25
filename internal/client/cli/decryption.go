package cli

import "github.com/vilasle/gokeep/internal/encryption"

// TODO implement decryption
func (c *Client) decrypt(dek []byte, content []byte) ([]byte, error) {
	//decryption DEK before decryption content
	dekED := encryption.NewEncryptedDataFromReadyData(c.encoder, dek, nil)
	dekOpened, err := dekED.Decrypt()
	if err != nil {
		return nil, err
	}
	//create encoder for data
	dekEnc, err := encryption.NewAESKeyFromJSON(dekOpened)
	if err != nil {
		return nil, err
	}

	//decryption content with help DEK encoder
	contentED := encryption.NewEncryptedDataFromReadyData(dekEnc, content, nil)
	contentOpened, err := contentED.Decrypt()
	if err != nil {
		return nil, err
	}
	return contentOpened, nil
}
