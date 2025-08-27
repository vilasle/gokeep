package service

import "time"

type RegisterLoginUser struct {
	Username  string
	Password  string
	PublicKey []byte
}

type MetadataValue struct {
	Key   string
	Value string
}

type AddLoginPassword struct {
	UserID   int
	Username string
	Password string
	Metadata []MetadataValue
}

type UpdateLoginPassword struct {
	UserID int
	ID     int
	AddLoginPassword
	Metadata []MetadataValue
}

type AddBankCard struct {
	UserID     int
	Number     string
	CVV        int
	Expiration time.Time
	Metadata   []MetadataValue
}

type UpdateBankCard struct {
	ID int
	AddBankCard
	Metadata []MetadataValue
}

type AddTextData struct {
	UserID   int
	Name     string
	Text     []byte
	Metadata []MetadataValue
}

type UpdateTextData struct {
	ID int
	AddTextData
	Metadata []MetadataValue
}

type AddBinaryData struct {
	UserID   int
	Name     string
	Data     []byte
	Metadata []MetadataValue
}

type UpdateBinaryData struct {
	ID int
	AddBinaryData
	Metadata []MetadataValue
}

type PrivateDataResponse struct {
	ID int
	//opened presentation of data
	View string
	//encrypted data
	Data string
	//encrypted DEK key
	Key      string
	Metadata []MetadataValue
}

type ListPrivateDataResponse struct {
	Data []PrivateDataResponse
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
	UserID    int
	SessionID int
	PublicKey []byte
}
