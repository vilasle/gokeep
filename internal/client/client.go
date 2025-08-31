package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Client interface {
	Close() error

	CreateAccount(ctx context.Context, accountName, password string) error
	Login(ctx context.Context, accountName, password string) (err error)

	DeleteLoginPassword(ctx context.Context, id int) error
	DeleteBankCard(ctx context.Context, id int) error
	DeleteTextData(ctx context.Context, id int) error
	DeleteBinaryData(ctx context.Context, id int) error

	SaveLoginPassword(ctx context.Context, login, password string, id int, meta map[string]string) (err error)
	SaveBankCard(ctx context.Context, number, expires string, cvv, id int, meta map[string]string) (err error)
	SaveTextData(ctx context.Context, text string, name string, id int, meta map[string]string) (err error)
	SaveBinaryData(ctx context.Context, content []byte, name string, id int, meta map[string]string) error

	GetLoginPassword(ctx context.Context, id int) error
	GetBankCard(ctx context.Context, id int) error
	GetTextData(ctx context.Context, id int) error
	GetBinaryData(ctx context.Context, id int) error

	Sync(ctx context.Context) error
}

const (
	PubKeyName  = "public.key"
	PrivKeyName = "private.key"

	MainDir                  = ".gokeep"
	ConfigName               = "config.yaml"
	UploadName               = "upload"
	CertificateNameDirectory = "cert"
	CreadExt                 = ".cread"
)

type Config struct {
	ServerSocket string `yaml:"grpc-server" env:"GOKEEP_GRPC_SERVER" env-required:"true"`
	DBPath       string `yaml:"database-path" env:"GOKEEP_DATABASE_PATH" env-required:"true"`
}

type WorkplaceConfig struct {
	ConfigDirectory PathInfo
	UploadDirectory PathInfo
	Config          PathInfo
	Certificate     PathInfo
	Credentials     PathInfo
}

func GetCurrentConfiguration(customConfigPath string) (WorkplaceConfig, error) {
	var workplace WorkplaceConfig
	if len(customConfigPath) > 0 {
		workplace = customConfiguration(customConfigPath)
	} else {
		workplace = defaultConfiguration()
	}

	return workplace, workplace.PathError()
}

func CreateNewConfiguration(configPath, serverSocket, dbPath string) (err error) {
	workplaceConfig := WorkplaceConfig{}
	if len(configPath) > 0 {
		workplaceConfig = customConfiguration(configPath)
	} else {
		workplaceConfig = defaultConfiguration()
	}

	if err := createDirectories(workplaceConfig.ConfigDirectory, workplaceConfig.Certificate, workplaceConfig.UploadDirectory); err != nil {
		return err
	}

	if err := generateRSAKeys(workplaceConfig.Certificate.Path); err != nil {
		return err
	}

	return generateConfig(workplaceConfig.Config.Path, serverSocket, dbPath)
}

func (w WorkplaceConfig) Report() {
	fmt.Println("workplace directory:", w.ConfigDirectory.Path)
	fmt.Println("upload directory:", w.UploadDirectory.Path)
	fmt.Println("config:", w.Config.Path)
	fmt.Println("certificate:", w.Certificate.Path)
}

func (w WorkplaceConfig) PathError() error {
	return errors.Join(w.ConfigDirectory.Error, w.Config.Error,
		w.Certificate.Error)
}

type PathInfo struct {
	Path  string
	Error error
}

func defaultConfiguration() WorkplaceConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return WorkplaceConfig{}
	}
	return customConfiguration(homeDir)

}

func customConfiguration(path string) WorkplaceConfig {
	return WorkplaceConfig{
		ConfigDirectory: getPathInfoCheckOnlyExisting(path),
		UploadDirectory: getPathInfoCheckOnlyExisting(filepath.Join(path, UploadName)),
		Config:          getPathInfoCheckOnlyExisting(filepath.Join(path, ConfigName)),
		Certificate:     getCertificatesPathInfo(filepath.Join(path, CertificateNameDirectory)),
		Credentials:     getCredentialsPathInfo(path),
	}
}

func getPathInfoCheckOnlyExisting(path string) PathInfo {
	_, err := os.Stat(path)
	return PathInfo{Path: path, Error: err}
}

func getCertificatesPathInfo(path string) PathInfo {
	pathInfo := getPathInfoCheckOnlyExisting(path)
	if pathInfo.Error != nil {
		return pathInfo
	}

	ls, err := os.ReadDir(path)
	if err != nil {
		pathInfo.Error = errors.New("reading directory")
	} else if len(ls) == 0 {
		pathInfo.Error = errors.New("no certificates found")
	}
	return pathInfo
}

func getCredentialsPathInfo(path string) PathInfo {
	ls, err := os.ReadDir(path)
	if err != nil {
		return PathInfo{Path: path, Error: err}
	}
	for _, file := range ls {
		ext := filepath.Ext(file.Name())
		if ext == CreadExt {
			return PathInfo{Path: filepath.Join(path, file.Name()), Error: nil}
		}
	}
	return PathInfo{Path: path, Error: errors.New("no credentials found")}
}

func createDirectories(path ...PathInfo) error {
	for _, path := range path {
		if path.Error != nil {
			if err := os.MkdirAll(path.Path, os.ModePerm); err != nil {
				return errors.Join(err, errors.New("create directory"))
			}
		}
	}
	return nil
}

func generateRSAKeys(savePath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Join(err, errors.New("generate private key"))
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM := pem.EncodeToMemory(privateKeyBlock)

	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return errors.Join(err, errors.New("marshal public key"))
	}

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)

	privKeyPath := filepath.Join(savePath, PrivKeyName)
	pubKeyPath := filepath.Join(savePath, PubKeyName)

	if err := os.WriteFile(privKeyPath, privateKeyPEM, 0644); err != nil {
		return errors.Join(err, errors.New("write private key"))
	}

	if err := os.WriteFile(pubKeyPath, publicKeyPEM, 0644); err != nil {
		return errors.Join(err, errors.New("write public key"))
	}
	return nil
}

func generateConfig(savePath, serverSocket, dbPath string) error {
	db, err := filepath.Abs(dbPath)
	if err != nil {
		return err
	}
	config := Config{
		ServerSocket: serverSocket,
		DBPath:       db,
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return errors.Join(err, errors.New("marshal config"))
	}
	if err := os.WriteFile(savePath, configBytes, 0644); err != nil {
		return errors.Join(err, errors.New("write config"))
	}
	return nil
}
