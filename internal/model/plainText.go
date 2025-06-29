package model

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
)

var _ PrivateData = (*PlainText)(nil)

type PlainText struct {
	model
	text []byte
}

func newPlainText(owner *User, text []byte) *PlainText {
	return &PlainText{
		model: model{
			owner: owner,
		},
		text: text,
	}
}

func (tp *PlainText) Save(ctx context.Context, encrypter Encrypter) (err error) {
	if err := tp.prepareEncryptedData(encrypter); err != nil {
		//TODO wrap error with package error
		return err
	}
	return tp.model.Save(ctx)
}

func (tp *PlainText) prepareEncryptedData(encrypter Encrypter) error {
	data, err := tp.dataForEncryption()
	if err != nil {
		//TODO wrap error with package error
		return err
	}

	encryptedData, err := encrypter.Encrypt(data)
	if err != nil {
		return err
	}
	encryptedData.Owner = tp

	tp.encryptedData = encryptedData

	return nil
}

func (tp *PlainText) String() (string) {
	view := tp.text[:20]
	return string(view)
}

func (tp *PlainText) SetText(text []byte) {
	tp.text = text
}

func (tp *PlainText) decryptData(encrypter Encrypter) (err error) {
	data, err := encrypter.Decrypt(*tp.encryptedData)
	if err != nil {
		return err
	}

	//if NewReader or ReadAll got error we can same error on Close
	rd, _ := gzip.NewReader(bytes.NewReader(data))
	tp.text, _ = io.ReadAll(rd)

	return rd.Close()
}

func (tp PlainText) dataForEncryption() ([]byte, error) {
	buf := bytes.Buffer{}

	//if NewWriterLevel or Write got error we can same error on Close
	wrt, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	wrt.Write(tp.text)

	//if Write got error we can same on Close
	err := wrt.Close()
	return buf.Bytes(), err
}

//FIXME add getting plain text by id and check that owner was right id
func findPlainTextByID(ctx context.Context, id int64, owner *User, r PrivateDataRepository) (*PlainText, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	u, ok := data.(*PlainText)
	if !ok {
		return nil, fmt.Errorf("invalid data. private data(id=%d) does not have type '*PlainText'", id)
	}

	return u, nil
}
