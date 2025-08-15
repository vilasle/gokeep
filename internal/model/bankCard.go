package model

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	return &BankCard{
		model: model{
			owner: owner,
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

func (u *BankCard) String() string {
	view := strings.Repeat("*", len(u.number)-4) + u.number[len(u.number)-4:]
	return string(view)
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

func (u *BankCard) prepareEncryptedData(encoder Encoder) error {
	data := u.dataForEncryption()

	encryptedData, err := encoder.Encrypt(data)
	if err != nil {
		return err
	}

	u.encryptedData = encryptedData

	u.encryptedData.Owner = u

	return nil
}

func (bc BankCard) dataForEncryption() []byte {
	buf := bytes.Buffer{}
	buf.WriteString(bc.number)
	buf.WriteString("\n")
	buf.WriteString(strconv.Itoa(bc.cvv))
	buf.WriteString("\n")
	buf.WriteString(bc.expiration.Format(time.DateOnly))

	return buf.Bytes()
}

// FIXME add getting band card by id and check that owner was right id
func findBankCardByID(ctx context.Context, id int64, owner *User, r PrivateDataRepository) (*BankCard, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	u, ok := data.(*BankCard)
	if !ok {
		return nil, fmt.Errorf("invalid data. private data(id=%d) does not have type '*BankCard'", id)
	}

	return u, nil
}
