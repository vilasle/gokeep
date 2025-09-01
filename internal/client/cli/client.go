package cli

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/vilasle/gokeep/internal/client"
	"github.com/vilasle/gokeep/internal/encryption"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	"github.com/vilasle/gokeep/internal/repository/client/sqlite"
	"github.com/vilasle/gokeep/internal/service/client"
	grpcSvc "github.com/vilasle/gokeep/internal/service/client/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type externalServices struct {
	credentials client.LoginPasswordDataService
	bankCard    client.BankCardDataService
	text        client.TextDataDataService
	binary      client.BinaryDataDataService
}

func newDataServices(socket *grpc.ClientConn) *externalServices {
	return &externalServices{
		credentials: grpcSvc.NewLoginPasswordDataService(socket),
		bankCard:    grpcSvc.NewBankCardService(socket),
		text:        grpcSvc.NewTextDataService(socket),
		binary:      grpcSvc.NewBinaryDataService(socket),
	}
}

//CommandLineClient - client for command line interface
type CommandLineClient struct {
	workspace config.WorkplaceConfig
	config    config.Config
	//private key use for decryption messages from server
	encoder encryption.Encoder
	//publicKeyContent - need only bytes for pass it to server
	publicKeyContent []byte
	auth             client.AuthService
	externalServices *externalServices
	localStorage     repository.ClientRepository
	credential       []byte
	conn             *grpc.ClientConn
}

//NewClient - create new client, loads config files and path from workspace
func NewClient(workspace config.WorkplaceConfig, caFile string) (client *CommandLineClient, err error) {
	client = &CommandLineClient{
		workspace: workspace,
	}

	if err = client.loadConfiguration(); err != nil {
		return nil, err
	}

	if err = client.loadCredentialIfExists(); err != nil {
		return nil, err
	}

	if err = client.loadRSAKeys(); err != nil {
		return nil, err
	}

	var transportOpt credentials.TransportCredentials
	if caFile != "" {
		transportOpt, err = credentials.NewClientTLSFromFile(caFile, "")
		if err != nil {
			return nil, err
		}
	} else {
		transportOpt = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(client.config.ServerSocket, grpc.WithTransportCredentials(transportOpt))
	if err != nil {
		return nil, err
	}

	client.conn = conn

	client.auth = grpcSvc.NewGRPCAuthService(conn)
	client.externalServices = newDataServices(conn)
	client.localStorage, err = sqlite.NewSQLiteClient(client.config.DBPath)

	return client, err
}

//Close - close client
func (c *CommandLineClient) Close() error {
	return errors.Join(c.conn.Close(), c.localStorage.Close())
}

func (c *CommandLineClient) loadConfiguration() error {
	//try load information about grpc server and local database
	configFd, err := os.Open(c.workspace.Config.Path)
	if err != nil {
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

func (c *CommandLineClient) loadRSAKeys() error {
	publicPath, privatePath, err := findKeysPath(c.workspace.Certificate.Path)
	if err != nil {
		return err
	}

	if publicPath == "" || privatePath == "" {
		return fmt.Errorf("public and private keys not found")
	}
	//save raw content of public key
	if c.publicKeyContent, err = os.ReadFile(publicPath); err != nil {
		return err
	}

	//get raw content of private key
	privateKeyContent, err := os.ReadFile(privatePath)
	if err != nil {
		return err
	}
	//parse and store private key
	block, _ := pem.Decode(privateKeyContent)
	if block == nil {
		return fmt.Errorf("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	c.encoder = encryption.NewRSACipher(&privateKey.PublicKey, privateKey)
	return nil
}

// loadCredentialIfExists - load credential from file if it exists. if file does not exists method will not return error
func (c *CommandLineClient) loadCredentialIfExists() error {
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

func (c *CommandLineClient) saveCredential(account string) error {
	if len(c.credential) == 0 {
		return errors.New("credential is empty")
	}

	stat, err := os.Stat(c.workspace.Credentials.Path)
	if err != nil {
		return err
	}

	path := c.workspace.Credentials.Path
	if stat.IsDir() {
		path = filepath.Join(c.workspace.Credentials.Path, account)
		path += config.CreadExt
	}

	os.RemoveAll(path)

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
		if file.Name() == config.PubKeyName {
			publicPath = filepath.Join(path, file.Name())
		} else if file.Name() == config.PrivKeyName {
			privatePath = filepath.Join(path, file.Name())
		}
	}

	return publicPath, privatePath, nil
}

func (c *CommandLineClient) getExternalID(ctx context.Context, req repository.GetRequest) int {
	if req.ID == 0 {
		return 0
	}

	response, err := c.localStorage.Get(ctx, req)
	if err != nil || len(response) == 0 {
		return 0
	}
	return response[0].ExternalID
}
