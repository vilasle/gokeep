package service

import "time"

type RegisterLoginUser struct {
	Username string
	Password string
}

type RegisterLoginResponse struct {
	Error      string
	Credential string
}

type AddLoginPassword struct {
	UserID   int64
	Username string
	Password string
}

type UpdateLoginPassword struct {
	UserID int64
	ID     int64
	AddLoginPassword
}

type AddBankCard struct {
	UserID     int64
	Number     string
	CVV        int
	Expiration time.Time
}

type UpdateBankCard struct {
	ID int64
	AddBankCard
}

type AddTextData struct {
	UserID int64
	Name   string
	Text   []byte
}

type UpdateTextData struct {
	ID int64
	AddTextData
}

type AddBinaryData struct {
	UserID int64
	Name   string
	Data   []byte
}

type UpdateBinaryData struct {
	ID int64
	AddBinaryData
}

type AddingUpdatePrivateDataResponse struct {
	ID int64
	//encrypted data
	Data string
	//encrypted DEK key
	Key   string
	Error string
}

type ListPrivateDataResponse struct {
	Data  []map[string]any
	Error string
}

type GetPrivateData struct {
	UserID int64
	ID     int64
}

type PrivateDataResponse struct {
	Fields map[string]any
	Error  string
}

type DeletePrivateData struct {
	UserID int64
	ID     int64
}
