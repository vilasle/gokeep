package model

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"
)

var _ PrivateData = (*Usepass)(nil)

//BandCard is model for presentation bank card
type BankCard struct {
	model
	number     string
	cvv        int
	expiration time.Time
}

func newBankCard(owner UserAccess, cardNumber string, cvv int, expiration time.Time) *BankCard {
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

//SetNumber set new number for bank card and view
func (u *BankCard) SetNumber(number string) {
	u.number = number
	if len(u.number) > 4 {
		u.model.view = strings.Repeat("*", len(u.number)-4) + u.number[len(u.number)-4:]
	}
}

//SetCVV set new cvv for bank card
func (u *BankCard) SetCVV(cvv int) {
	u.cvv = cvv
}

//SetExpiration set new expiration for bank card
func (u *BankCard) SetExpiration(expiration time.Time) {
	u.expiration = expiration
}

//Save prepare view of model, encrypt it and save on repository
func (u *BankCard) Save(ctx context.Context, encoder Encoder) (err error) {
	if err := u.prepareEncryptedData(encoder); err != nil {
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

func findBankCardByID(ctx context.Context, id int, owner UserAccess, r PrivateDataRepository) (*BankCard, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.ID() {
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
