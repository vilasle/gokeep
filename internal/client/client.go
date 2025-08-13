package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	pubKeyName  = "public.key"
	privKeyName = "private.key"
)

type Client struct {
	workspace WorkplaceConfig
	config    Config
	//private key use for decryption messages from server
	privateKey *rsa.PrivateKey
	//publicKeyContent - need only bytes for pass it to server
	publicKeyContent []byte
}

func NewClient(workspace WorkplaceConfig) (*Client, error) {
	client := &Client{
		workspace: workspace,
	}

	if err := client.loadConfiguration(); err != nil {
		return nil, err
	}

	if err := client.loadRSAKeys(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) loadConfiguration() error {
	//try load information about grpc server and local database
	configFd, err := os.Open(c.workspace.ConfigPath.Path)
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
	publicPath, privatePath, err := findKeysPath(c.workspace.CredentialsPath.Path)
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

func findKeysPath(path string) (string, string, error) {
	var publicPath, privatePath string

	ls, err := os.ReadDir(path)
	if err != nil {
		//TODO add error context
		return "", "", err
	}

	for _, file := range ls {
		if file.Name() == pubKeyName {
			publicPath = file.Name()
		} else if file.Name() == privKeyName {
			privatePath = file.Name()
		}
	}

	return publicPath, privatePath, nil
}
