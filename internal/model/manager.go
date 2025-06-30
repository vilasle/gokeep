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
	Usepass    *usepassManager
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
		Usepass: &usepassManager{
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

func (c *userManager) Get(ctx context.Context, id int64) (*User, error) {
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

func (c *usepassManager) List(ctx context.Context, owner *User) ([]*Usepass, error) {
	return fillListOfPrivateData[*Usepass](ctx, owner, ModelTypeUsepass, c.pvRepository)
}

func (c *usepassManager) Get(ctx context.Context, owner *User, id int64) (*Usepass, error) {
	//TODO add logger
	usepass, err := findUsepassByID(ctx, id, owner, c.pvRepository)
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

func (c *bankCardManager) Get(ctx context.Context, owner *User, id int64) (*BankCard, error) {
	//TODO add logger
	usepass, err := findBankCardByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	usepass.dataRepository = c.pvRepository
	usepass.encryptionRepository = c.encRepository
	return usepass, nil
}

func (c *bankCardManager) List(ctx context.Context, owner *User) ([]*BankCard, error) {
	return fillListOfPrivateData[*BankCard](ctx, owner, ModelTypeBankCard, c.pvRepository)
}

type plainTextManager struct {
	pvRepository  PrivateDataRepository
	encRepository EncryptedDataRepository
}

func (c *plainTextManager) New(owner *User, text []byte) *PlainText {
	plainText := newPlainText(owner, text)
	plainText.dataRepository = c.pvRepository
	plainText.encryptionRepository = c.encRepository
	return plainText
}

func (c *plainTextManager) Get(ctx context.Context, owner *User, id int64) (*PlainText, error) {
	//TODO add logger
	plainText, err := findPlainTextByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	plainText.dataRepository = c.pvRepository
	plainText.encryptionRepository = c.encRepository
	return plainText, nil
}

func (c *plainTextManager) List(ctx context.Context, owner *User) ([]*PlainText, error) {
	return fillListOfPrivateData[*PlainText](ctx, owner, ModelTypePlainText, c.pvRepository)
}

type binaryDataManager struct {
	pvRepository  PrivateDataRepository
	encRepository EncryptedDataRepository
}

func (c *binaryDataManager) New(owner *User, name string, data []byte) *BinaryData {
	binaryData := newBinaryData(owner, data, name)
	binaryData.dataRepository = c.pvRepository
	binaryData.encryptionRepository = c.encRepository
	return binaryData
}

func (c *binaryDataManager) Get(ctx context.Context, owner *User, id int64) (*BinaryData, error) {
	//TODO add logger
	binaryData, err := findBinaryDataByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	binaryData.dataRepository = c.pvRepository
	binaryData.encryptionRepository = c.encRepository
	return binaryData, nil
}

func (c *binaryDataManager) List(ctx context.Context, owner *User) ([]*BinaryData, error) {
	return fillListOfPrivateData[*BinaryData](ctx, owner, ModelTypeBinaryData, c.pvRepository)
}

func fillListOfPrivateData[T *Usepass | *BankCard | *PlainText | *BinaryData](ctx context.Context, owner *User, modelType ModelType, repository PrivateDataRepository) ([]T, error) {
	privateData, err := repository.List(ctx, modelType, owner)
	if err != nil {
		return nil, err
	}

	result := make([]T, len(privateData))
	for i, data := range privateData {
		if entity, ok := data.(T); ok {
			result[i] = entity
		} else {
			//TODO add logger about wrong type of data
			//TODO implement it
			panic("not implemented")
		}
	}
	return result, nil
}
