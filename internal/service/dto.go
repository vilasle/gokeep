package service

import "time"

type RegisterLoginUser struct {
	Username  string
	Password  string
	PublicKey []byte
}

type AddLoginPassword struct {
	UserID   int
	Username string
	Password string
}

type UpdateLoginPassword struct {
	UserID int
	ID     int
	AddLoginPassword
}

type AddBankCard struct {
	UserID     int
	Number     string
	CVV        int
	Expiration time.Time
}

type UpdateBankCard struct {
	ID int
	AddBankCard
}

type AddTextData struct {
	UserID int
	Name   string
	Text   []byte
}

type UpdateTextData struct {
	ID int
	AddTextData
}

type AddBinaryData struct {
	UserID int
	Name   string
	Data   []byte
}

type UpdateBinaryData struct {
	ID int
	AddBinaryData
}

type AddingUpdatePrivateDataResponse struct {
	ID int
	//opened presentation of data
	View string
	//encrypted data
	Data string
	//encrypted DEK key
	Key   string
}

type ListPrivateDataResponse struct {
	Data  []map[string]any
}

type GetPrivateData struct {
	UserID int
	ID     int
}

type PrivateDataResponse struct {
	Fields map[string]any
}

type DeletePrivateData struct {
	UserID int
	ID     int
}

type SessionInfo struct {
	UserID int
	SessionID int
	PublicKey []byte
}
