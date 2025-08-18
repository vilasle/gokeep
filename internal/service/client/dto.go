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
}

type LoginPasswordGetRequest struct {
	ID int
}

type LoginPasswordDeleteRequest struct {
	ID int
}

// response
type LoginPasswordSaveResponse struct {
	ID    int
	Data  EncryptedData
	Error string
}

type LoginPasswordListResponse struct {
	Result []LoginPasswordView
	Error  string
}

type LoginPasswordGetResponse struct {
	Data  EncryptedData
	Error string
}

type LoginPasswordDeleteResponse struct {
	Error string
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
	Expires time.Time
	CVV     int
}

type BankCardGetRequest struct {
	ID int
}

type BankCardDeleteRequest struct {
	ID int
}

// response
type BankCardSaveResponse struct {
	ID    int
	Data  EncryptedData
	Error string
}

type BankCardListResponse struct {
	Result []BankCardView
	Error  string
}

type BankCardGetResponse struct {
	Data  EncryptedData
	Error string
}

type BankCardDeleteResponse struct {
	Error string
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
}

type TextDataGetRequest struct {
	ID int
}

type TextDataDeleteRequest struct {
	ID int
}

// response
type TextDataSaveResponse struct {
	ID    int
	Data  EncryptedData
	Error string
}

type TextDataListResponse struct {
	Result []TextDataView
	Error  string
}

type TextDataGetResponse struct {
	Data  EncryptedData
	Error string
}

type TextDataDeleteResponse struct {
	Error string
}

// binary data
type BinaryDataView struct {
	ID   int
	Name string
}

//requests

type BinaryDataSaveRequest struct {
	Name string
	Data []byte
}

type BinaryDataGetRequest struct {
	ID int
}

type BinaryDataDeleteRequest struct {
	ID int
}

// response
type BinaryDataSaveResponse struct {
	ID    int
	Data  EncryptedData
	Error string
}

type BinaryDataListResponse struct {
	Result []BinaryDataView
	Error  string
}

type BinaryDataGetResponse struct {
	Data  EncryptedData
	Error string
}

type BinaryDataDeleteResponse struct {
	Error string
}
