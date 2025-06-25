package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type encryptedData struct {
	Data string `json:"data"`
	Key  string `json:"key"`
}

type aesKey struct {
	Key   string `json:"key"`
	Nonce string `json:"nonce"`
	key   []byte
	nonce []byte
	gcm   cipher.AEAD
}

func (k *aesKey) initGCM() error {
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

func (k *aesKey) encrypt(data []byte) ([]byte, error) {
	return k.gcm.Seal(nil, k.nonce, data, nil), nil
}

func (k *aesKey) decrypt(data []byte) ([]byte, error) {
	return k.gcm.Open(nil, k.nonce, data, nil)
}

func (k aesKey) String() string {
	return fmt.Sprintf(`{"key":"%s","nonce":"%s"}`, k.Key, k.Nonce)
}

func TestCreateMasterAESKey(t *testing.T) {
	key := make([]byte, 2*aes.BlockSize)
	_, err := rand.Read(key)
	require.NoError(t, err)

	aesBlock, err := aes.NewCipher(key)
	require.NoError(t, err)

	aesGcm, err := cipher.NewGCM(aesBlock)

	require.NoError(t, err)

	nonce := make([]byte, aesGcm.NonceSize())
	_, err = rand.Read(nonce)
	require.NoError(t, err)

	objKey := aesKey{
		Key:   hex.EncodeToString(key),
		Nonce: hex.EncodeToString(nonce),
	}

	content, err := json.Marshal(objKey)
	require.NoError(t, err)

	err = os.WriteFile("key.json", content, 0644)
	require.NoError(t, err)
}

func TestCreateDEKAESKey(t *testing.T) {
	key := make([]byte, 2*aes.BlockSize)
	_, err := rand.Read(key)
	require.NoError(t, err)

	aesBlock, err := aes.NewCipher(key)
	require.NoError(t, err)

	aesGcm, err := cipher.NewGCM(aesBlock)

	require.NoError(t, err)

	nonce := make([]byte, aesGcm.NonceSize())
	_, err = rand.Read(nonce)
	require.NoError(t, err)

	objKey := aesKey{
		Key:   hex.EncodeToString(key),
		Nonce: hex.EncodeToString(nonce),
	}

	content, err := json.Marshal(objKey)
	require.NoError(t, err)

	err = os.WriteFile("dek.json", content, 0644)
	require.NoError(t, err)
}

func TestCreateClientRSAKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	fdPriv, err := os.OpenFile("private.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	require.NoError(t, err)

	fdPub, err := os.OpenFile("public.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	require.NoError(t, err)

	pem.Encode(fdPriv, privateKeyBlock)
	pem.Encode(fdPub, publicKeyBlock)
}

func TestXxx(t *testing.T) {
	source := []byte("some text which need to be encrypted")
	fmt.Println("source", source)

	//read master key
	content, err := os.ReadFile("key.json")
	require.NoError(t, err)

	objKey := aesKey{}
	err = json.Unmarshal(content, &objKey)
	require.NoError(t, err)

	err = objKey.initGCM()
	require.NoError(t, err)

	//create DEK certificate
	dekKey := aesKey{}
	err = json.Unmarshal(content, &dekKey)
	require.NoError(t, err)

	err = dekKey.initGCM()
	require.NoError(t, err)

	publicKeyContent, err := os.ReadFile("public.key")
	require.NoError(t, err)

	publicBlock, _ := pem.Decode(publicKeyContent)
	publicKey, err := x509.ParsePKCS1PublicKey(publicBlock.Bytes)
	require.NoError(t, err)

	privateKeyContent, err := os.ReadFile("private.key")
	require.NoError(t, err)

	privateBlock, _ := pem.Decode(privateKeyContent)
	require.NoError(t, err)

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	require.NoError(t, err)

	_, _ = publicKey, privateKey
	//encrypt data with DEK certificate
	encData, err := dekKey.encrypt(source)
	require.NoError(t, err)

	//encrypt DEK certificate with master key
	encryptedDEK, err := objKey.encrypt([]byte(dekKey.String()))
	require.NoError(t, err)

	serverData := encryptedData{
		Data: hex.EncodeToString(encData),
		Key:  hex.EncodeToString(encryptedDEK),
	}

	fmt.Println(serverData)

	clientEncDEKKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(dekKey.String()), nil)
	require.NoError(t, err)

	//decrypt data via master key and private key
	clientData := encryptedData{
		Data: hex.EncodeToString(encData),
		Key:  hex.EncodeToString(clientEncDEKKey),
	}

	fmt.Println(clientData)

	//decrypt data via master key and private key

	//server part
	sdKey, err := hex.DecodeString(serverData.Key)
	require.NoError(t, err)

	decDEK, err := objKey.decrypt([]byte(sdKey))
	require.NoError(t, err)

	decDEKObj := aesKey{}
	err = json.Unmarshal(decDEK, &decDEKObj)
	require.NoError(t, err)

	err = decDEKObj.initGCM()
	require.NoError(t, err)

	sdd, err := hex.DecodeString(serverData.Data)
	require.NoError(t, err)
	decData, err := decDEKObj.decrypt([]byte(sdd))
	require.NoError(t, err)

	fmt.Println("server part", string(decData))

	//client part
	cdKey, err := hex.DecodeString(clientData.Key)
	require.NoError(t, err)
	decDEK1, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, []byte(cdKey), nil)
	require.NoError(t, err)

	decDEKObj1 := aesKey{}
	err = json.Unmarshal(decDEK1, &decDEKObj1)
	require.NoError(t, err)

	err = decDEKObj1.initGCM()
	require.NoError(t, err)

	cdd, err := hex.DecodeString(serverData.Data)
	require.NoError(t, err)
	decData1, err := decDEKObj1.decrypt([]byte(cdd))
	require.NoError(t, err)

	fmt.Println("client part", string(decData1))
}
