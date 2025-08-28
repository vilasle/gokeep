package model

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"
)

var _ PrivateData = (*Usepass)(nil)

type BankCard struct {
	model
	number     string
	cvv        int
	expiration time.Time
}

func newBankCard(owner *User, cardNumber string, cvv int, expiration time.Time) *BankCard {
	ent := &BankCard{
		model: model{
			owner:     owner,
			modelType: TypeBankCard,
		},
		number:     cardNumber,
		expiration: expiration,
		cvv:        cvv,
	}

	ent.SetNumber(ent.number)
	return ent
}

func (u *BankCard) SetNumber(number string) {
	u.number = number
	if len(u.number) > 4 {
		u.model.view = strings.Repeat("*", len(u.number)-4) + u.number[len(u.number)-4:]
	}
}

func (u *BankCard) SetCVV(cvv int) {
	u.cvv = cvv
}

func (u *BankCard) SetExpiration(expiration time.Time) {
	u.expiration = expiration
}

func (u *BankCard) Save(ctx context.Context, encoder Encoder) (err error) {
	if err := u.prepareEncryptedData(encoder); err != nil {
		//TODO wrap error with package error
		return err
	}
	return u.model.Save(ctx)
}

func (u *BankCard) prepareEncryptedData(encoder Encoder) (err error) {
	data := u.dataForEncryption()

	u.encryptedData, err = encoder.Encrypt(data)
	return err
}

func (bc BankCard) dataForEncryption() []byte {
	buf := bytes.Buffer{}
	buf.WriteString(bc.number)
	buf.WriteString("\n")
	buf.WriteString(strconv.Itoa(bc.cvv))
	buf.WriteString("\n")
	buf.WriteString(bc.expiration.Format("01/06"))

	return buf.Bytes()
}

// FIXME add getting band card by id and check that owner was right id
func findBankCardByID(ctx context.Context, id int, owner *User, r PrivateDataRepository) (*BankCard, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.id {
		return nil, ErrNotFound
	}

	model := model{
		id:             id,
		view:           data.View,
		owner:          owner,
		modelType:      TypeBankCard,
		dataRepository: r,
		encryptedData: &EncryptedData{
			Data: data.Data,
			Key:  data.DEK,
		},
		metadata: data.Metadata,
	}

	return &BankCard{model: model}, nil
}
