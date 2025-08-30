package model

import "context"

type RepositoryCollector interface {
	User() UserRepository
	Private() PrivateDataRepository
}

type UserRepository interface {
	Add(context.Context, UserAdd) (id int, err error)
	Update(context.Context, UserUpdate) error
	Delete(context.Context, int) error
	Find(context.Context, string) (UserInfo, error)
	Get(context.Context, int) (UserInfo, error)
}

type PrivateData interface {
	ID() int
	Type() int
	Owner() UserAccess
	String() string
	EncryptedData() EncryptedData
	Metadata() map[string]string
}

type PrivateDataRepository interface {
	Add(ctx context.Context, data PrivateDataSave) (int, error)
	//return id and error, but id return only for having same signature with Add method
	Update(ctx context.Context, data PrivateDataSave) (int, error)
	Delete(ctx context.Context, id int) error
	Get(context.Context, int) (PrivateDataInfo, error)
	List(ctx context.Context, modelType Type, owner UserAccess) ([]PrivateDataInfo, error)
}
