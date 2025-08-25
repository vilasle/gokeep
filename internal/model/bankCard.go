package model

import (
	"bytes"
	"context"
	"errors"
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
	view := strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
	return &BankCard{
		model: model{
			view:      view,
			owner:     owner,
			modelType: TypeBankCard,
		},
		number:     cardNumber,
		expiration: expiration,
		cvv:        cvv,
	}
}

func (bc *BankCard) decryptData(encoder Encoder) (err error) {
	data, err := encoder.Decrypt(*bc.encryptedData)
	if err != nil {
		return err
	}

	lp := strings.Split(string(data), "\n")
	if len(lp) != 3 {
		//TODO use package error
		return errors.New("invalid data")
	}

	number, cvv, expiration := lp[0], lp[1], lp[2]

	bc.number = number
	bc.cvv, err = strconv.Atoi(cvv)
	if err != nil {
		//TODO use package error
		return err
	}

	bc.expiration, err = time.Parse(time.DateOnly, expiration)
	if err != nil {
		//TODO use package error
		return err
	}
	return nil
}

func (u *BankCard) SetNumber(number string) {
	u.number = number
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
		owner:          owner,
		modelType:      TypeUsepass,
		dataRepository: r,
		encryptedData: &EncryptedData{
			Data: data.Data,
			Key:  data.DEK,
		},
	}

	return &BankCard{model: model}, nil
}
