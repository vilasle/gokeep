package client

type AuthService interface {
	CreateAccount(accountName, password string, publicKey []byte) error
	Login(accountName, password string) (credential []byte, err error)
}

type PrivateDataService interface{}
