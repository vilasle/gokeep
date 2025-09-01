package model

import "context"

//RepositoryCollector interface for manager
type RepositoryCollector interface {
	User() UserRepository
	Private() PrivateDataRepository
}

//UserRepository interface for user repository
type UserRepository interface {
	Add(context.Context, UserAdd) (id int, err error)
	Update(context.Context, UserUpdate) error
	Delete(context.Context, int) error
	Find(context.Context, string) (UserInfo, error)
	Get(context.Context, int) (UserInfo, error)
}

//PrivateDataRepository interface for model 
type PrivateData interface {
	ID() int
	Type() int
	Owner() UserAccess
	String() string
	EncryptedData() EncryptedData
	Metadata() map[string]string
}

//PrivateDataRepository interface for private data repository
type PrivateDataRepository interface {
	Add(ctx context.Context, data PrivateDataSave) (int, error)
	//return id and error, but id return only for having same signature with Add method
	Update(ctx context.Context, data PrivateDataSave) (int, error)
	Delete(ctx context.Context, id int) error
	Get(context.Context, int) (PrivateDataInfo, error)
	List(ctx context.Context, modelType Type, userID int) ([]PrivateDataInfo, error)
}
