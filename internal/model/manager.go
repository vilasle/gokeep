package model

import (
	"context"
	"time"
)

type ModelType = int8

const (
	ModelTypeUser ModelType = iota + 1
	ModelTypeUsepass
	ModelTypeBankCard
	ModelTypePlainText
	ModelTypeBinaryData
)

// ModelManager is a manager for all models.
type ModelManager struct {
	repository RepositoryCollector
	Users      *userManager
	BankCards  *bankCardManager
	PlainTexts *plainTextManager
	BinaryData *binaryDataManager
}

// NewModelManager creates a new model manager
func NewModelManager(repository RepositoryCollector) ModelManager {
	return ModelManager{
		repository: repository,
		Users: &userManager{
			repository: repository.User(),
		},
		BankCards: &bankCardManager{
			pvRepository:  repository.Private(),
			encRepository: repository.Encryption(),
		},
		PlainTexts: &plainTextManager{
			pvRepository:  repository.Private(),
			encRepository: repository.Encryption(),
		},
		BinaryData: &binaryDataManager{
			pvRepository:  repository.Private(),
			encRepository: repository.Encryption(),
		},
	}
}

// userManager is a manager for users
type userManager struct {
	repository UserRepository
}

// New creates a new user with login and hashed password
func (c *userManager) New(login, password string) *User {
	return newUser(login, password, c.repository)
}

// FindByLogin  finds a user by login on repository, return ErrUserNotFound if not found
func (c *userManager) FindByLogin(ctx context.Context, login string) (*User, error) {
	//TODO add logger
	return findUserByLogin(ctx, login, c.repository)
}

func (c *userManager) GetByID(ctx context.Context, id int64) (*User, error) {
	//TODO add logger
	return getUserByID(ctx, id, c.repository)
}

type usepassManager struct {
	pvRepository  PrivateDataRepository
	encRepository EncryptedDataRepository
}

func (c *usepassManager) New(owner *User, login, password string) *Usepass {
	usepass := newUsepass(owner, login, password)
	usepass.dataRepository = c.pvRepository
	usepass.encryptionRepository = c.encRepository
	return usepass
}

func (c *usepassManager) FindByID(ctx context.Context, id int64) (*Usepass, error) {
	//TODO add logger
	usepass, err := findUsepassByID(ctx, id, c.pvRepository)
	if err != nil {
		return nil, err
	}
	usepass.dataRepository = c.pvRepository
	usepass.encryptionRepository = c.encRepository
	return usepass, nil
}

type bankCardManager struct {
	pvRepository  PrivateDataRepository
	encRepository EncryptedDataRepository
}

func (c *bankCardManager) New(owner *User, number string, cvv int, expirationDate time.Time) *BankCard {
	bankCard := newBankCard(owner, number, cvv, expirationDate)
	bankCard.dataRepository = c.pvRepository
	bankCard.encryptionRepository = c.encRepository
	return bankCard
}

func (c *bankCardManager) FindByID(ctx context.Context, id int64) (*BankCard, error) {
	//TODO add logger
	usepass, err := findBankCardByID(ctx, id, c.pvRepository)
	if err != nil {
		return nil, err
	}
	usepass.dataRepository = c.pvRepository
	usepass.encryptionRepository = c.encRepository
	return usepass, nil
}

type plainTextManager struct {
	pvRepository  PrivateDataRepository
	encRepository EncryptedDataRepository
}

type binaryDataManager struct {
	pvRepository  PrivateDataRepository
	encRepository EncryptedDataRepository
}
