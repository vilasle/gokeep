package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestBinaryDataService_New(t *testing.T) {
	conn, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer conn.Close()

	srv := grpc.NewServer()
	defer srv.Stop()

	go srv.Serve(conn)

	clientConn, err := grpc.NewClient(conn.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer clientConn.Close()

	_ = NewBinaryDataService(clientConn)
}

func TestBinaryDataService_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.BinaryDataSaveRequest{
			Name: "name",
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
			Data: []byte("data"),
		}
		expected := client.SaveResponse{
			ID: 1,
			Data: client.EncryptedData{
				View: "view",
				DEK:  []byte("dek"),
				Data: []byte("data"),
			},
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}

		mock.EXPECT().SaveBinaryData(ctx, &proto.SaveBinaryDataRequest{
			Data: []byte("data"),
			Name: "name",
			Credential: &proto.ConfirmAssess{Token: req.JWT},
			Metadata: []*proto.Metadata{
				{
					Key:   req.Metadata[0].Key,
					Value: req.Metadata[0].Value,
				},
			},
		}).Return(&proto.EncryptedDataResponse{
			Entity: &proto.EncryptedEntity{
				Id: 1,
				Data: &proto.EncryptedData{
					View: "view",
					Dek:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []*proto.Metadata{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		resp, err := svc.Save(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expected, resp)
	})

	t.Run("client error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.BinaryDataSaveRequest{
			Name: "name",
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
			Data: []byte("data"),
		}
	
		mock.EXPECT().SaveBinaryData(ctx, &proto.SaveBinaryDataRequest{
			Data: []byte("data"),
			Name: "name",
			Credential: &proto.ConfirmAssess{Token: req.JWT},
			Metadata: []*proto.Metadata{
				{
					Key:   req.Metadata[0].Key,
					Value: req.Metadata[0].Value,
				},
			},
		}).Return(nil, errors.New("client error"))

		resp, err := svc.Save(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp)
	})

	t.Run("response is nil", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.BinaryDataSaveRequest{
			Name: "name",
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
			Data: []byte("data"),
		}
	
		mock.EXPECT().SaveBinaryData(ctx, &proto.SaveBinaryDataRequest{
			Data: []byte("data"),
			Name: "name",
			Credential: &proto.ConfirmAssess{Token: req.JWT},
			Metadata: []*proto.Metadata{
				{
					Key:   req.Metadata[0].Key,
					Value: req.Metadata[0].Value,
				},
			},
		}).Return(nil, nil)

		resp, err := svc.Save(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp)
	})

	t.Run("response has error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.BinaryDataSaveRequest{
			Name: "name",
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
			Data: []byte("data"),
		}
	
		mock.EXPECT().SaveBinaryData(ctx, &proto.SaveBinaryDataRequest{
			Data: []byte("data"),
			Name: "name",
			Credential: &proto.ConfirmAssess{Token: req.JWT},
			Metadata: []*proto.Metadata{
				{
					Key:   req.Metadata[0].Key,
					Value: req.Metadata[0].Value,
				},
			},
		}).Return(&proto.EncryptedDataResponse{Error: "error"}, nil)

		resp, err := svc.Save(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp)
	})
}

func TestBinaryDataService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.GetRequest{
			JWT:  "jwt",
			ID:   1,
			Type: 1,
		}

		expected := []client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}

		mock.EXPECT().Get(ctx, &proto.GetDataRequest{
			Id:         1,
			Type:       1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(&proto.GetDataResponse{
			Data: []*proto.EncryptedEntity{
				{
					Id: 1,
					Data: &proto.EncryptedData{
						View: "view",
						Dek:  []byte("dek"),
						Data: []byte("data"),
					},
					Metadata: []*proto.Metadata{
						{
							Key:   "key",
							Value: "value",
						},
					},
				},
			},
			Error: ""}, nil)

		resp, err := svc.Get(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, resp, expected)
	})
	t.Run("client error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.GetRequest{
			JWT:  "jwt",
			ID:   1,
			Type: 1,
		}

		mock.EXPECT().Get(ctx, &proto.GetDataRequest{
			Id:         1,
			Type:       1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(nil, errors.New("error"))

		resp, err := svc.Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
	t.Run("response is nil", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.GetRequest{
			JWT:  "jwt",
			ID:   1,
			Type: 1,
		}

		mock.EXPECT().Get(ctx, &proto.GetDataRequest{
			Id:         1,
			Type:       1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(nil, nil)

		resp, err := svc.Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
	t.Run("response has error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.GetRequest{
			JWT:  "jwt",
			ID:   1,
			Type: 1,
		}

		mock.EXPECT().Get(ctx, &proto.GetDataRequest{
			Id:         1,
			Type:       1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(&proto.GetDataResponse{Error: "error"}, nil)

		resp, err := svc.Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestBinaryDataService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.DeleteRequest{
			JWT: "jwt",
			ID:  1,
		}

		mock.EXPECT().Delete(ctx, &proto.DeleteDataRequest{
			Id:         1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(&proto.DeleteDataResponse{Error: ""}, nil)

		err := svc.Delete(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("client error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.DeleteRequest{
			JWT: "jwt",
			ID:  1,
		}

		mock.EXPECT().Delete(ctx, &proto.DeleteDataRequest{
			Id:         1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(nil, errors.New("client error"))

		err := svc.Delete(ctx, req)
		assert.Error(t, err)
	})

	t.Run("response is nil", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.DeleteRequest{
			JWT: "jwt",
			ID:  1,
		}

		mock.EXPECT().Delete(ctx, &proto.DeleteDataRequest{
			Id:         1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(nil, nil)

		err := svc.Delete(ctx, req)
		assert.Error(t, err)
	})

	t.Run("response has error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockPrivateDataServiceClient(ctrl)

		svc := BinaryDataService{
			client: mock,
		}

		req := client.DeleteRequest{
			JWT: "jwt",
			ID:  1,
		}

		mock.EXPECT().Delete(ctx, &proto.DeleteDataRequest{
			Id:         1,
			Credential: &proto.ConfirmAssess{Token: req.JWT},
		}).Return(&proto.DeleteDataResponse{Error: "error"}, nil)

		err := svc.Delete(ctx, req)
		assert.Error(t, err)
	})
}
