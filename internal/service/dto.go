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

type PrivateDataResponse struct {
	ID int
	//opened presentation of data
	View string
	//encrypted data
	Data string
	//encrypted DEK key
	Key   string
}

type ListPrivateDataResponse struct {
	Data  []PrivateDataResponse
}

type GetPrivateData struct {
	UserID int
	ID     int
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
