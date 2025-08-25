package client

import "time"

type MetadataValue struct {
	Key   string
	Value string
}

type EncryptedData struct {
	View string
	DEK  []byte
	Data []byte
}

type GetRequest struct {
	ID  int
	JWT string
}

type DeleteRequest struct {
	ID  int
	JWT string
}

// login-password
type LoginPasswordView struct {
	ID       int
	Login    string
	Metadata []MetadataValue
}

// requests
type LoginPasswordSaveRequest struct {
	ID              int
	Login, Password string
	JWT             string
}

type SaveResponse struct {
	ID   int
	Data EncryptedData
}

type LoginPasswordListResponse struct {
	Result []LoginPasswordView
}

//band cards

// requests
type BankCardView struct {
	ID      int
	Number  string
	Expires time.Time
}

type BankCardSaveRequest struct {
	ID      int
	Number  string
	Expires string
	CVV     int
	JWT     string
}

// response
type BankCardSaveResponse struct {
	ID   int
	Data EncryptedData
}

type BankCardListResponse struct {
	Result []BankCardView
}

// text data
type TextDataView struct {
	ID   int
	Name string
}

// requests
type TextDataSaveRequest struct {
	ID   int
	Name string
	Text []byte
	JWT  string
}

// response
type TextDataSaveResponse struct {
	ID   int
	Data EncryptedData
}

type TextDataListResponse struct {
	Result []TextDataView
}

// binary data
type BinaryDataView struct {
	ID   int
	Name string
}

//requests

type BinaryDataSaveRequest struct {
	ID   int
	Name string
	Data []byte
	JWT  string
}

// response
type BinaryDataSaveResponse struct {
	ID   int
	Data EncryptedData
}

type BinaryDataListResponse struct {
	Result []BinaryDataView
}
