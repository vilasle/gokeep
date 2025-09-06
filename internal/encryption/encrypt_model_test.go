package encryption

import (
	context "context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/model"
)

func TestEncryptionModel_Encrypt(t *testing.T) {
	jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
	DEK, err := NewAESKeyFromJSON(jsonKey)
	require.NoError(t, err)

	// //create server key
	privateKey1, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	masterKey := NewRSACipher(&privateKey1.PublicKey, privateKey1)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rep := NewMockRepositoryCollector(ctrl)

	pvRep := NewMockPrivateDataRepository(ctrl)
	rep.EXPECT().User().Return(nil)
	rep.EXPECT().Private().Return(pvRep)

	pvRep.EXPECT().Add(gomock.Any(), gomock.Any()).Return(1, nil).Times(1)

	manager := model.NewModelManager(rep)

	user := &model.User{}

	card := manager.BankCards.New(user, "123456788", 123, time.Now().Add(time.Hour*24*365))

	em := NewModelEncoding(jsonKey, masterKey, DEK)

	ctx := context.Background()
	err = card.Save(ctx, em)
	require.NoError(t, err)
}

func TestEncryptionModel_Decrypt(t *testing.T) {
	jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
	DEK, err := NewAESKeyFromJSON(jsonKey)
	require.NoError(t, err)

	ed := model.EncryptedData{
		Data: []byte("708590aed525c0c764fe3fed42cb9a3c7b16344f199dc542abdce8ddb036e956cadc57b9d3b8eee2"),
		Key:  []byte("0df991ccae8d812547cc364f4d75f35f0cf2134d21bb242a150e045e56ed9c9cd4362d93d9b56c77f8ff5e9ae9de1ebd4990e705afaf5a021ee8d4266dd129bc0788a983e4ccaf4d82487210f218991eacd5168ded72b3353b918e9c66ef99875b42b4b169cf8de673fff36d3728a59ea331a6caa323bc726ac336a8247d5be8f0c5b08dbe6eb1cab18aecd8e7f4dfb89ce4e9281951f743ffcf77cb5c465a5f6504a8a2953ab9be3778fedc5c32c2a49d484369adae5bc74637a0389acbcda08db9fde70cffe1554fc395fb11f77e8cca2e37cc168cb8e4506e16fbde10c73c0fbaf811922a3088a52788475911d3d3ae53fb8d0a7746d370e682b0511a1ac7"),
	}

	em := NewModelEncoding(jsonKey, nil, DEK)

	data, err := em.Decrypt(ed)

	require.NoError(t, err)
	require.Equal(t, "123456788\n123\n2026-06-28", string(data))

}
