package model

import (
	"context"
	"errors"
	"time"

	"github.com/vilasle/gokeep/internal/logger"
)

type Type = int

const (
	TypeUser Type = iota + 1
	TypeUsepass
	TypeBankCard
	TypePlainText
	TypeBinaryData
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
			pvRepository: repository.Private(),
		},
		PlainTexts: &plainTextManager{
			pvRepository: repository.Private(),
		},
		BinaryData: &binaryDataManager{
			pvRepository: repository.Private(),
		},
		Usepass: &usepassManager{
			pvRepository: repository.Private(),
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

func (c *userManager) Get(ctx context.Context, id int) (*User, error) {
	//TODO add logger
	return getUserByID(ctx, id, c.repository)
}

type usepassManager struct {
	pvRepository PrivateDataRepository
}

func (c *usepassManager) New(owner *User, login, password string) *Usepass {
	usepass := newUsepass(owner, login, password)
	usepass.dataRepository = c.pvRepository
	return usepass
}

func (c *usepassManager) List(ctx context.Context, owner *User) ([]*Usepass, error) {
	return fillListOfPrivateData[*Usepass](ctx, owner, TypeUsepass, c.pvRepository)
}

func (c *usepassManager) Get(ctx context.Context, owner *User, id int) (*Usepass, error) {
	//TODO add logger
	usepass, err := findUsepassByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	usepass.dataRepository = c.pvRepository
	return usepass, nil
}

type bankCardManager struct {
	pvRepository PrivateDataRepository
}

func (c *bankCardManager) New(owner *User, number string, cvv int, expirationDate time.Time) *BankCard {
	bankCard := newBankCard(owner, number, cvv, expirationDate)
	bankCard.dataRepository = c.pvRepository
	return bankCard
}

func (c *bankCardManager) Get(ctx context.Context, owner *User, id int) (*BankCard, error) {
	//TODO add logger
	usepass, err := findBankCardByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	usepass.dataRepository = c.pvRepository
	return usepass, nil
}

func (c *bankCardManager) List(ctx context.Context, owner *User) ([]*BankCard, error) {
	return fillListOfPrivateData[*BankCard](ctx, owner, TypeBankCard, c.pvRepository)
}

type plainTextManager struct {
	pvRepository PrivateDataRepository
}

func (c *plainTextManager) New(owner *User, text []byte, view string) *PlainText {
	plainText := newPlainText(owner, text, view)
	plainText.dataRepository = c.pvRepository
	return plainText
}

func (c *plainTextManager) Get(ctx context.Context, owner *User, id int) (*PlainText, error) {
	//TODO add logger
	plainText, err := findPlainTextByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	plainText.dataRepository = c.pvRepository
	return plainText, nil
}

func (c *plainTextManager) List(ctx context.Context, owner *User) ([]*PlainText, error) {
	return fillListOfPrivateData[*PlainText](ctx, owner, TypePlainText, c.pvRepository)
}

type binaryDataManager struct {
	pvRepository PrivateDataRepository
}

func (c *binaryDataManager) New(owner *User, data []byte, name string) *BinaryData {
	binaryData := newBinaryData(owner, data, name)
	binaryData.dataRepository = c.pvRepository
	return binaryData
}

func (c *binaryDataManager) Get(ctx context.Context, owner *User, id int) (*BinaryData, error) {
	//TODO add logger
	binaryData, err := findBinaryDataByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	binaryData.dataRepository = c.pvRepository
	return binaryData, nil
}

func (c *binaryDataManager) List(ctx context.Context, owner *User) ([]*BinaryData, error) {
	return fillListOfPrivateData[*BinaryData](ctx, owner, TypeBinaryData, c.pvRepository)
}

func fillListOfPrivateData[T *Usepass | *BankCard | *PlainText | *BinaryData](
	ctx context.Context,
	owner *User,
	modelType Type,
	repository PrivateDataRepository) ([]T, error) {

	ls, err := repository.List(ctx, modelType, owner)
	if err != nil {
		return nil, err
	}

	result := make([]T, len(ls))
	for i, data := range ls {
		pv, err := newPrivateDate(data, modelType, owner, repository)
		if err != nil {
			logger.Error("can not create private data from PrivateDataInfo", "modelType", modelType, "data", data)
			continue
		}

		if entity, ok := pv.(T); ok {
			result[i] = entity
		} else {
			logger.Error("wrong type of data", "modelType", modelType, "data", data)

		}
	}
	return result, nil
}

func newPrivateDate(pd PrivateDataInfo, modelType Type, owner *User, r PrivateDataRepository) (PrivateData, error) {

	var pv PrivateData

	meta := make(map[string]string)
	for k, v := range pd.Metadata {
		meta[k] = v
	}

	model := model{
		id:             pd.ID,
		owner:          owner,
		dataRepository: r,
		modelType:      modelType,
		encryptedData: &EncryptedData{
			Data: pd.Data,
			Key:  pd.DEK,
		},
		metadata: meta,
	}

	switch modelType {
	case TypeUsepass:
		pv = &Usepass{model: model}
	case TypeBankCard:
		pv = &BankCard{model: model}
	case TypePlainText:
		pv = &PlainText{model: model}
	case TypeBinaryData:
		pv = &BinaryData{model: model}
	default:
		return nil, errors.New("unknown model type")
	}

	return pv, nil
}
