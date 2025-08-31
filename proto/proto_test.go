package proto

import (
	context "context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type mockSrv struct {
	UnimplementedAccountServiceServer
	UnimplementedPrivateDataServiceServer
}

func (mockSrv) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return &CreateAccountResponse{Error: ""}, nil
}
func (mockSrv) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return &LoginResponse{Token: "token", Error: ""}, nil
}

func (mockSrv) SaveLoginPassword(context.Context, *SaveLoginPasswordRequest) (*EncryptedDataResponse, error) {
	return &EncryptedDataResponse{}, nil
}
func (mockSrv) SaveBankCard(context.Context, *SaveBankCardRequest) (*EncryptedDataResponse, error) {
	return &EncryptedDataResponse{}, nil
}
func (mockSrv) SaveTextData(context.Context, *SaveTextDataRequest) (*EncryptedDataResponse, error) {
	return &EncryptedDataResponse{}, nil
}
func (mockSrv) SaveBinaryData(context.Context, *SaveBinaryDataRequest) (*EncryptedDataResponse, error) {
	return &EncryptedDataResponse{}, nil
}
func (mockSrv) Get(context.Context, *GetDataRequest) (*GetDataResponse, error) {
	return &GetDataResponse{}, nil
}
func (mockSrv) Delete(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error) {
	return &DeleteDataResponse{}, nil
}

func Test_AuthService(t *testing.T) {
	//listen free port
	conn, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer conn.Close()

	srv := grpc.NewServer()
	defer srv.Stop()

	mock := mockSrv{}

	RegisterAccountServiceServer(srv, mock)

	go srv.Serve(conn)

	clientConn, err := grpc.NewClient(conn.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer clientConn.Close()

	client := NewAccountServiceClient(clientConn)

	ctx := context.Background()
	{
		resp, err := client.CreateAccount(ctx, &CreateAccountRequest{Login: "test", Password: "test"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)

	}

	{
		resp, err := client.Login(ctx, &LoginRequest{Login: "test", Password: "test"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

}

func Test_PrivateService(t *testing.T) {
	//listen free port
	conn, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer conn.Close()

	srv := grpc.NewServer()
	defer srv.Stop()

	mock := mockSrv{}

	RegisterPrivateDataServiceServer(srv, mock)

	go srv.Serve(conn)

	clientConn, err := grpc.NewClient(conn.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer clientConn.Close()

	client := NewPrivateDataServiceClient(clientConn)

	ctx := context.Background()
	{
		resp, err := client.SaveBankCard(ctx, &SaveBankCardRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

	{
		resp, err := client.SaveBinaryData(ctx, &SaveBinaryDataRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

	{
		resp, err := client.SaveLoginPassword(ctx, &SaveLoginPasswordRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

	{
		resp, err := client.SaveTextData(ctx, &SaveTextDataRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

	{
		resp, err := client.Get(ctx, &GetDataRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

	{
		resp, err := client.Delete(ctx, &DeleteDataRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	}

}
