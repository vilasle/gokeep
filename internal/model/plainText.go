package model

import (
	"bytes"
	"compress/gzip"
	"context"
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
			owner:     owner,
			modelType: TypePlainText,
		},
		text: text,
	}
}

func (tp *PlainText) Save(ctx context.Context, encoder Encoder) (err error) {
	if err := tp.prepareEncryptedData(encoder); err != nil {
		//TODO wrap error with package error
		return err
	}
	return tp.model.Save(ctx)
}

func (tp *PlainText) prepareEncryptedData(encoder Encoder) error {
	data, err := tp.dataForEncryption()
	if err != nil {
		//TODO wrap error with package error
		return err
	}

	tp.encryptedData, err = encoder.Encrypt(data)
	return err
}

func (tp *PlainText) String() string {
	if len(tp.text) < 20 {
		return string(tp.text)
	}
	view := tp.text[:20]
	return string(view)
}

func (tp *PlainText) SetText(text []byte) {
	tp.text = text
}

func (tp *PlainText) decryptData(encoder Encoder) (err error) {
	data, err := encoder.Decrypt(*tp.encryptedData)
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

// FIXME add getting plain text by id and check that owner was right id
func findPlainTextByID(ctx context.Context, id int, owner *User, r PrivateDataRepository) (*PlainText, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
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
	return &PlainText{model: model}, nil
}
