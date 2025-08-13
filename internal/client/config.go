package client

import (
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

const (
	mainDir                  = ".gokeep"
	configName               = "config.yaml"
	certificateNameDirectory = "cert"
	creadExt                 = ".cread"
)

type Config struct {
	ServerSocket string `yaml:"grpc-server" env:"GOKEEP_GRPC_SERVER" env-required:"true"`
	DBPath       string `yaml:"database-path" env:"GOKEEP_DATABASE_PATH" env-required:"true"`
}

type WorkplaceConfig struct {
	ConfigDirectoryPath PathInfo
	ConfigPath          PathInfo
	CertificatePath     PathInfo
	CredentialsPath     PathInfo
}

func (w WorkplaceConfig) Report() {
	fmt.Println("workplace directory:", w.ConfigDirectoryPath.Path)
	fmt.Println("config:", w.ConfigPath.Path)
	fmt.Println("certificate:", w.CertificatePath.Path)
}

type PathInfo struct {
	Path  string
	Error error
}

func GetCurrentConfiguration(customConfigPath string) (WorkplaceConfig, error) {
	if len(customConfigPath) > 0 {
		return customConfiguration(customConfigPath)
	}
	return defaultConfiguration(), nil
}

func CreateNewConfiguration(configPath, serverSocket, dbPath string) (err error) {
	workplaceConfig := WorkplaceConfig{}
	if len(configPath) > 0 {
		workplaceConfig, err = customConfiguration(configPath)
		if err != nil {
			return err
		}
	} else {
		workplaceConfig = defaultConfiguration()
	}

	if err := createDirectories(workplaceConfig.ConfigDirectoryPath, workplaceConfig.CertificatePath); err != nil {
		return err
	}

	if err := generateRSAKeys(workplaceConfig.CertificatePath.Path); err != nil {
		return err
	}

	return generateConfig(workplaceConfig.ConfigPath.Path, serverSocket, dbPath)
}

func defaultConfiguration() WorkplaceConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return WorkplaceConfig{}
	}

	configDir := filepath.Join(homeDir, mainDir)
	configPath := filepath.Join(configDir, configName)
	certPath := filepath.Join(configDir, certificateNameDirectory)

	return WorkplaceConfig{
		ConfigDirectoryPath: getPathInfoCheckOnlyExisting(configDir),
		ConfigPath:          getPathInfoCheckOnlyExisting(configPath),
		CertificatePath:     getCertificatesPathInfo(certPath),
		CredentialsPath:     getCredentialsPathInfo(configDir),
	}
}

func customConfiguration(path string) (WorkplaceConfig, error) {
	return WorkplaceConfig{
		ConfigDirectoryPath: getPathInfoCheckOnlyExisting(path),
		ConfigPath:          getPathInfoCheckOnlyExisting(filepath.Join(path, configName)),
		CertificatePath:     getCertificatesPathInfo(filepath.Join(path, certificateNameDirectory)),
		CredentialsPath:     getCredentialsPathInfo(path),
	}, nil

}

func isExists(path string) error {
	_, err := os.Stat(path)
	return err
}

func getPathInfoCheckOnlyExisting(path string) PathInfo {
	if err := isExists(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return PathInfo{Path: path, Error: errors.New("does not exists")}
		} else {
			return PathInfo{Path: path, Error: err}
		}
	}
	return PathInfo{Path: path, Error: nil}
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
		if ext == creadExt {
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

	privKeyPath := filepath.Join(savePath, privKeyName)
	pubKeyPath := filepath.Join(savePath, pubKeyName)

	if err := os.WriteFile(privKeyPath, privateKeyPEM, 0644); err != nil {
		return errors.Join(err, errors.New("write private key"))
	}

	if err := os.WriteFile(pubKeyPath, publicKeyPEM, 0644); err != nil {
		return errors.Join(err, errors.New("write public key"))
	}
	return nil
}

func generateConfig(savePath, serverSocket, dbPath string) error {
	config := Config{
		ServerSocket: serverSocket,
		DBPath:       dbPath,
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
