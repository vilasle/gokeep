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
	Users      *userManager
	BankCards  *bankCardManager
	PlainTexts *plainTextManager
	BinaryData *binaryDataManager
	Usepass    *usepassManager
}

// NewModelManager creates a new model manager
func NewModelManager(repository RepositoryCollector) ModelManager {
	p := repository.Private()
	return ModelManager{
		Users: &userManager{
			repository: repository.User(),
		},
		BankCards: &bankCardManager{
			pvRepository: p,
		},
		PlainTexts: &plainTextManager{
			pvRepository: p,
		},
		BinaryData: &binaryDataManager{
			pvRepository: p,
		},
		Usepass: &usepassManager{
			pvRepository: p,
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
func (c *userManager) FindByLogin(ctx context.Context, login string) (UserAccess, error) {
	logger.Debug("search user by login", "login", login)
	return findUserByLogin(ctx, login, c.repository)
}

// Get - return UserAccess by id
func (c *userManager) Get(ctx context.Context, id int) (UserAccess, error) {
	logger.Debug("get user by id", "id", id)
	return getUserByID(ctx, id, c.repository)
}

type usepassManager struct {
	pvRepository PrivateDataRepository
}

// New creates a new usepass with login and password
func (c *usepassManager) New(owner UserAccess, login, password string) *Usepass {
	usepass := newUsepass(owner, login, password)
	usepass.dataRepository = c.pvRepository
	return usepass
}

// List returns list of usepass for user
func (c *usepassManager) List(ctx context.Context, owner UserAccess) ([]*Usepass, error) {
	logger.Debug("getting list of entities", "owner", owner.Login(), "entity", "usepass")
	return fillListOfPrivateData[*Usepass](ctx, owner, TypeUsepass, c.pvRepository)
}

// Get returns usepass by id
func (c *usepassManager) Get(ctx context.Context, owner UserAccess, id int) (*Usepass, error) {
	logger.Debug("getting entity by id", "owner", owner.Login(), "id", id, "entity", "usepass")

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

// New creates a new bank card with number, cvv and expiration date
func (c *bankCardManager) New(owner UserAccess, number string, cvv int, expirationDate time.Time) *BankCard {
	bankCard := newBankCard(owner, number, cvv, expirationDate)
	bankCard.dataRepository = c.pvRepository
	return bankCard
}

// Get returns bank card by id
func (c *bankCardManager) Get(ctx context.Context, owner UserAccess, id int) (*BankCard, error) {
	logger.Debug("getting entity by id", "owner", owner.Login(), "id", id, "entity", "bank card")
	usepass, err := findBankCardByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	usepass.dataRepository = c.pvRepository
	return usepass, nil
}

// List returns list of bank cards for user
func (c *bankCardManager) List(ctx context.Context, owner UserAccess) ([]*BankCard, error) {
	logger.Debug("getting list of entities", "owner", owner.Login(), "entity", "bank card")
	return fillListOfPrivateData[*BankCard](ctx, owner, TypeBankCard, c.pvRepository)
}

type plainTextManager struct {
	pvRepository PrivateDataRepository
}

// New creates a new plain text with text and view
func (c *plainTextManager) New(owner UserAccess, text []byte, view string) *PlainText {
	plainText := newPlainText(owner, text, view)
	plainText.dataRepository = c.pvRepository
	return plainText
}

// Get returns plain text by id
func (c *plainTextManager) Get(ctx context.Context, owner UserAccess, id int) (*PlainText, error) {
	logger.Info("getting entity by id", "owner", owner.Login(), "id", id, "entity", "plain text")
	plainText, err := findPlainTextByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	plainText.dataRepository = c.pvRepository
	return plainText, nil
}

// List returns list of plain texts for user
func (c *plainTextManager) List(ctx context.Context, owner UserAccess) ([]*PlainText, error) {
	logger.Debug("getting list of entities", "owner", owner.Login(), "entity", "plain text")
	return fillListOfPrivateData[*PlainText](ctx, owner, TypePlainText, c.pvRepository)
}

type binaryDataManager struct {
	pvRepository PrivateDataRepository
}

// New creates a new binary data with data and name
func (c *binaryDataManager) New(owner UserAccess, data []byte, name string) *BinaryData {
	binaryData := newBinaryData(owner, data, name)
	binaryData.dataRepository = c.pvRepository
	return binaryData
}

// Get returns binary data by id
func (c *binaryDataManager) Get(ctx context.Context, owner UserAccess, id int) (*BinaryData, error) {
	logger.Info("getting entity by id", "owner", owner.Login(), "id", id, "entity", "binary")
	binaryData, err := findBinaryDataByID(ctx, id, owner, c.pvRepository)
	if err != nil {
		return nil, err
	}
	binaryData.dataRepository = c.pvRepository
	return binaryData, nil
}

// List returns list of binary data for user
func (c *binaryDataManager) List(ctx context.Context, owner UserAccess) ([]*BinaryData, error) {
	logger.Debug("getting list of entities", "owner", owner.Login(), "entity", "binary")
	return fillListOfPrivateData[*BinaryData](ctx, owner, TypeBinaryData, c.pvRepository)
}

func fillListOfPrivateData[T *Usepass | *BankCard | *PlainText | *BinaryData](
	ctx context.Context,
	owner UserAccess,
	modelType Type,
	repository PrivateDataRepository) ([]T, error) {

	ls, err := repository.List(ctx, modelType, owner.ID())
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

func newPrivateDate(pd PrivateDataInfo, modelType Type, owner UserAccess, r PrivateDataRepository) (PrivateData, error) {
	var pv PrivateData

	meta := make(map[string]string)
	for k, v := range pd.Metadata {
		meta[k] = v
	}

	model := model{
		id:             pd.ID,
		view:           pd.View,
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
