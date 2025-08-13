package client

type AuthService interface {
	CreateAccount(accountName, password string) error
	Login(accountName, password string) (cread string, err error)
}

type PrivateDataService interface {}

