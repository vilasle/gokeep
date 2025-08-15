package cli

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	repository "github.com/vilasle/gokeep/internal/repository/client"
	"github.com/vilasle/gokeep/internal/repository/client/sqlite"
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/internal/service/client/grpc"
	"gopkg.in/yaml.v3"
)

const (
	pubKeyName  = "public.key"
	privKeyName = "private.key"
)

type externalServices struct {
	credentials client.LoginPasswordDataService
	bankCard    client.BankCardDataService
	text        client.TextDataDataService
	binary      client.BinaryDataDataService
}

func newDataServices(grpcSocket string) *externalServices {
	return &externalServices{
		credentials: grpc.NewLoginPasswordDataService(grpcSocket),
		bankCard:    grpc.NewBankCardService(grpcSocket),
		text:        grpc.NewTextDataService(grpcSocket),
		binary:      grpc.NewBinaryDataService(grpcSocket),
	}
}

type Client struct {
	workspace WorkplaceConfig
	config    Config
	//private key use for decryption messages from server
	privateKey *rsa.PrivateKey
	//publicKeyContent - need only bytes for pass it to server
	publicKeyContent []byte
	auth             client.AuthService
	externalServices *externalServices
	localStorage     repository.ClientRepository
	credential       []byte
}

func NewClient(workspace WorkplaceConfig) (*Client, error) {
	client := &Client{
		workspace: workspace,
	}

	if err := client.loadConfiguration(); err != nil {
		return nil, err
	}

	if err := client.loadCredentialIfExists(); err != nil {
		return nil, err
	}

	if err := client.loadRSAKeys(); err != nil {
		return nil, err
	}

	client.auth = grpc.NewGRPCAuthService(client.config.ServerSocket)
	client.externalServices = newDataServices(client.config.ServerSocket)
	client.localStorage = sqlite.NewSQLiteClient(client.config.DBPath)

	return client, nil
}

func (c *Client) loadConfiguration() error {
	//try load information about grpc server and local database
	configFd, err := os.Open(c.workspace.Config.Path)
	if err != nil {
		//TODO add error context
		return err
	}
	defer configFd.Close()

	dec := yaml.NewDecoder(configFd)

	if err := dec.Decode(&c.config); err != nil {
		//TODO add error context
		return err
	}

	return nil
}

func (c *Client) loadRSAKeys() error {
	publicPath, privatePath, err := findKeysPath(c.workspace.Certificate.Path)
	if publicPath == "" || privatePath == "" {
		//TODO add error context
		return fmt.Errorf("public and private keys not found")
	}
	//save raw content of public key
	if c.publicKeyContent, err = os.ReadFile(publicPath); err != nil {
		//TODO add error context
		return err
	}

	//get raw content of private key
	privateKeyContent, err := os.ReadFile(privatePath)
	if err != nil {
		//TODO add error context
		return err
	}
	//parse and store private key
	block, _ := pem.Decode(privateKeyContent)
	if block == nil {
		//TODO add error context
		return fmt.Errorf("failed to decode private key")
	}

	c.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		//TODO add error context
		return err
	}
	return nil
}

// loadCredentialIfExists - load credential from file if it exists. if file does not exists method will not return error
func (c *Client) loadCredentialIfExists() error {
	if c.workspace.Credentials.Error != nil {
		return nil
	}

	credential, err := os.ReadFile(c.workspace.Credentials.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	c.credential = credential
	return nil
}

func (c *Client) saveCredential(account string) error {
	if len(c.credential) == 0 {
		return errors.New("credential is empty")
	}

	path := filepath.Join(c.workspace.Credentials.Path, account)
	path += creadExt

	return os.WriteFile(path, c.credential, 0600)
}



func findKeysPath(path string) (string, string, error) {
	var publicPath, privatePath string

	ls, err := os.ReadDir(path)
	if err != nil {
		//TODO add error context
		return "", "", err
	}

	for _, file := range ls {
		if file.Name() == pubKeyName {
			publicPath = filepath.Join(path, file.Name())
		} else if file.Name() == privKeyName {
			privatePath = filepath.Join(path, file.Name())
		}
	}

	return publicPath, privatePath, nil
}
