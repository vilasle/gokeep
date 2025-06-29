package model

import "context"

type RepositoryCollector interface {
	User() UserRepository
	Private() (PrivateDataRepository)
	Encryption() (EncryptedDataRepository)
	 
}

type UserRepository interface {
	Add(context.Context, *User) error
	Update(context.Context, *User) error
	Delete(context.Context, *User) error
	Find(context.Context, string) (*User, error)
	Get(context.Context, int64) (*User, error)
	
}

type PrivateData interface {
	ID() int64
	Owner() *User
	String() string
	EncryptedData() EncryptedData
}

type PrivateDataRepository interface {
	Add(context.Context, PrivateData) error
	Update(context.Context, PrivateData) error
	Delete(context.Context, PrivateData) error
	Get(context.Context, int64) (PrivateData, error)
}

type EncryptedDataRepository interface {
	Add(context.Context, *EncryptedData) error
	Update(context.Context, *EncryptedData) error
	Delete(context.Context, *EncryptedData) error
	Get(context.Context, int64) (*EncryptedData, error)
}
