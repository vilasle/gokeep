package private

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/vilasle/gokeep/internal/model"
)

func Test_BankCardService_List(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rep := NewMockRepositoryCollector(ctrl)
	userRep := NewMockUserRepository(ctrl)
	rep.EXPECT().User().Return(userRep).Times(1)
	rep.EXPECT().Private().Return(nil).Times(4)
	rep.EXPECT().Encryption().Return(nil).Times(4)

	manager := model.NewModelManager(rep)

	service := NewBankCardService(manager)

	ctx := context.Background()

	userID := int64(1)
	service.List(ctx, userID)

}
