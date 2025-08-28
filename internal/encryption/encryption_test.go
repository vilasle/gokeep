package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RSACipher_EncryptDecrypt(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		key     string
		public  *rsa.PublicKey
		private *rsa.PrivateKey
		data    []byte
		encErr  error
		decErr  error
	}{
		{
			name:    "success",
			public:  &privateKey.PublicKey,
			private: privateKey,
			data:    []byte("SomeContentForEncryptionAndDecryption"),
			encErr:  nil,
			decErr:  nil,
		},
		{
			name:    "without public key",
			public:  nil,
			private: privateKey,
			data:    []byte("SomeContentForEncryptionAndDecryption"),
			encErr:  errors.New("error"),
			decErr:  nil,
		},
		{
			name:    "without private key",
			public:  &privateKey.PublicKey,
			private: nil,
			data:    []byte("SomeContentForEncryptionAndDecryption"),
			encErr:  nil,
			decErr:  errors.New("error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := NewRSACipher(tt.public, tt.private)

			encData, err := c.Encrypt(tt.data)
			if tt.encErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			decData, err := c.Decrypt(encData)
			if tt.decErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.data, decData)

		})
	}
}

func Test_AESKey_NewAESKeyFromJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonContent []byte
		err         error
	}{
		{
			name: "success",
			jsonContent: []byte(
				`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`),
			err: nil,
		},
		{
			name: "wrong json",
			jsonContent: []byte(
				`"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"`),
			err: errors.New("error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAESKeyFromJSON(tt.jsonContent)
			if tt.err != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}

}

func TestAESKey_JSON(t *testing.T) {
	expected := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)

	key, err := NewAESKeyFromJSON(expected)
	require.NoError(t, err)

	json := key.JSON()

	assert.Equal(t, string(expected), json)
}

func TestAESKey_EncryptDecrypt(t *testing.T) {
	jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
	content := []byte("SomeContentForEncryptionAndDecryption")
	key, err := NewAESKeyFromJSON(jsonKey)
	require.NoError(t, err)

	encData, err := key.Encrypt(content)
	require.NoError(t, err)

	decData, err := key.Decrypt(encData)
	require.NoError(t, err)

	require.Equal(t, content, decData)
}

func Test_EncryptedData_EncryptDecrypt(t *testing.T) {
	//valid case
	jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
	content := []byte("SomeContentForEncryptionAndDecryption")
	DEK, err := NewAESKeyFromJSON(jsonKey)
	require.NoError(t, err)

	data := NewEncryptedData(DEK)

	err = data.Encrypt(content)
	require.NoError(t, err)

	newData := NewEncryptedDataFromReadyData(DEK, []byte(data.Data), DEK.key)

	decData, err := newData.Decrypt()
	require.NoError(t, err)
	require.Equal(t, content, decData)

	edWrongHexData := EncryptedData{dek: DEK, Data: "not hex string"}

	_, err = edWrongHexData.Decrypt()
	require.Error(t, err)

	edWrongHexData.Key = "not hex string"

	_, err = edWrongHexData.DecryptKey(DEK)
	require.Error(t, err)


}

func Test_EncryptedData_ReplaceKey(t *testing.T) {
	//create DEK key
	jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
	content := []byte("SomeContentForEncryptionAndDecryption")
	DEK, err := NewAESKeyFromJSON(jsonKey)
	require.NoError(t, err)

	//create server key
	privateKey1, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	masterKey := NewRSACipher(&privateKey1.PublicKey, privateKey1)

	//create client key
	privateKey2, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	clientKey := NewRSACipher(&privateKey2.PublicKey, privateKey2)

	data := EncryptedData{dek: DEK}
	//encrypt key with created DEK key
	err = data.Encrypt(content)
	require.NoError(t, err)

	//encrypt key with server key
	err = data.EncryptKey(masterKey, jsonKey)
	require.NoError(t, err)

	//replace key, decrypt key with master key and encrypt key with client key
	err = data.ReplaceKey(masterKey, clientKey)
	require.NoError(t, err)

	//decrypt key with client key
	decDataKey, err := data.DecryptKey(clientKey)
	require.NoError(t, err)

	//create DEK key from decrypted key
	decDEK, err := NewAESKeyFromJSON(decDataKey)
	require.NoError(t, err)

	data.dek = decDEK
	//decrypt data with decrypted key
	decData, err := data.Decrypt()
	require.NoError(t, err)

	require.Equal(t, content, decData)
}

func Test_GenerateNewAESKey(t *testing.T) {
	key, err := GenerateNewAESKey()
	require.NoError(t, err)
	require.NotNil(t, key)
}

func Test_NewRSACipherFromRawPublicKey(t *testing.T) {
	//valid public key
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKey := pem.EncodeToMemory(publicKeyBlock)

	cipher, err := NewRSACipherFromRawPublicKey(publicKey)
	require.NoError(t, err)
	require.NotNil(t, cipher)

	//not valid public key

	cipher, err = NewRSACipherFromRawPublicKey([]byte("not valid public key"))
	require.Error(t, err)
	require.Nil(t, cipher)

}
