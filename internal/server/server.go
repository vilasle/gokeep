package server

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/logger"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

type Option func(*Server) error

func WithLogger(s *Server) error {
	s.opts = append(s.opts, grpc.UnaryInterceptor(
		logging.UnaryServerInterceptor(logger.InterceptorLogger()),
	))
	return nil
}

type Config struct {
	Addr string
	service.AuthService
	service.LoginPasswordService
	service.BankCardService
	service.TextDataService
	service.BinaryDataService
}

type Server struct {
	auth service.AuthService
	//data services
	cread  service.LoginPasswordService
	bank   service.BankCardService
	text   service.TextDataService
	binary service.BinaryDataService
	//grpc fields
	conn net.Listener
	srv  *grpc.Server
	opts []grpc.ServerOption
	pb.UnimplementedPrivateDataServiceServer
	pb.UnimplementedAccountServiceServer
}

func NewServer(config Config, opts ...Option) (*Server, error) {
	conn, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		conn:   conn,
		opts:   make([]grpc.ServerOption, 0, 2),
		auth:   config.AuthService,
		cread:  config.LoginPasswordService,
		bank:   config.BankCardService,
		text:   config.TextDataService,
		binary: config.BinaryDataService,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	s.srv = grpc.NewServer(s.opts...)
	return s, nil
}

func (s *Server) Listen() error {
	pb.RegisterPrivateDataServiceServer(s.srv, s)
	pb.RegisterAccountServiceServer(s.srv, s)
	logger.Info("starting server", "addr", s.conn.Addr().String())
	return s.srv.Serve(s.conn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
	s.conn.Close()
}

// CreateAccount - add new account via service.AuthService
func (s *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	resp := &pb.CreateAccountResponse{}
	dto := service.RegisterLoginUser{
		Username: req.Login,
		Password: req.Password,
	}

	if err := s.auth.Register(ctx, dto); err != nil {
		resp.Error = err.Error()
	}

	return resp, nil
}

// Login - login via service.AuthService, create new session and save current public key for session which will be used for encrypting data
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp := &pb.LoginResponse{}

	dto := service.RegisterLoginUser{
		Username:  req.Login,
		Password:  req.Password,
		PublicKey: req.PublicKey,
	}

	if result, err := s.auth.Login(ctx, dto); err != nil {
		resp.Error = err.Error()
	} else {
		resp.Token = result
	}

	return resp, nil
}

func (s *Server) SaveLoginPassword(ctx context.Context, req *pb.SaveLoginPasswordRequest) (*pb.EncryptedDataResponse, error) {
	resp := &pb.EncryptedDataResponse{}
	//get session and client key by token
	credential := req.Credential

	ses, err := s.getSessionByToken(ctx, credential.Token)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	dto := service.AddLoginPassword{
		UserID:   ses.userID,
		Username: req.Login,
		Password: req.Password,
		Metadata: castMetadata(req.Metadata),
	}

	if result, err := s.cread.Add(ctx, dto, ses.encoder); err != nil {
		resp.Error = err.Error()
	} else {
		resp.Entity = &pb.EncryptedEntity{
			Id: int64(result.ID),
			Data: &pb.EncryptedData{
				Data: []byte(result.Data),
				Dek:  []byte(result.Key),
				View: result.View,
			},
			Metadata: castMetadataPB(result.Metadata),
		}

	}

	return resp, nil
}

func (s *Server) SaveBankCard(ctx context.Context, req *pb.SaveBankCardRequest) (*pb.EncryptedDataResponse, error) {
	resp := &pb.EncryptedDataResponse{}
	//get session and client key by token
	credential := req.Credential

	ses, err := s.getSessionByToken(ctx, credential.Token)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	expiration, err := time.Parse("01/06", req.Expires)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	dto := service.AddBankCard{
		UserID:     ses.userID,
		Number:     req.Number,
		CVV:        int(req.Cvv),
		Expiration: expiration,
		Metadata:   castMetadata(req.Metadata),
	}

	if result, err := s.bank.Add(ctx, dto, ses.encoder); err != nil {
		resp.Error = err.Error()
	} else {
		resp.Entity = &pb.EncryptedEntity{
			Id: int64(result.ID),
			Data: &pb.EncryptedData{
				Data: []byte(result.Data),
				Dek:  []byte(result.Key),
				View: result.View,
			},
			Metadata: castMetadataPB(result.Metadata),
		}

	}

	return resp, nil
}

func (s *Server) SaveTextData(ctx context.Context, req *pb.SaveTextDataRequest) (*pb.EncryptedDataResponse, error) {
	resp := &pb.EncryptedDataResponse{}
	//get session and client key by token
	credential := req.Credential

	ses, err := s.getSessionByToken(ctx, credential.Token)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	dto := service.AddTextData{
		UserID:   ses.userID,
		Name:     req.Name,
		Text:     req.Data,
		Metadata: castMetadata(req.Metadata),
	}

	if result, err := s.text.Add(ctx, dto, ses.encoder); err != nil {
		resp.Error = err.Error()
	} else {
		resp.Entity = &pb.EncryptedEntity{
			Id: int64(result.ID),
			Data: &pb.EncryptedData{
				Data: []byte(result.Data),
				Dek:  []byte(result.Key),
				View: result.View,
			},
			Metadata: castMetadataPB(result.Metadata),
		}

	}

	return resp, nil
}

func (s *Server) SaveBinaryData(ctx context.Context, req *pb.SaveBinaryDataRequest) (*pb.EncryptedDataResponse, error) {
	resp := &pb.EncryptedDataResponse{}
	//get session and client key by token
	credential := req.Credential

	ses, err := s.getSessionByToken(ctx, credential.Token)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	dto := service.AddBinaryData{
		UserID:   ses.userID,
		Name:     req.Name,
		Data:     req.Data,
		Metadata: castMetadata(req.Metadata),
	}

	if result, err := s.binary.Add(ctx, dto, ses.encoder); err != nil {
		resp.Error = err.Error()
	} else {
		resp.Entity = &pb.EncryptedEntity{
			Id: int64(result.ID),
			Data: &pb.EncryptedData{
				Data: []byte(result.Data),
				Dek:  []byte(result.Key),
				View: result.View,
			},
			Metadata: castMetadataPB(result.Metadata),
		}

	}

	return resp, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteDataRequest) (resp *pb.DeleteDataResponse, err error) {
	resp = &pb.DeleteDataResponse{}
	//get session and client key by token
	credential := req.Credential

	ses, err := s.getSessionByToken(ctx, credential.Token)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	dto := service.DeletePrivateData{
		UserID: ses.userID,
		ID:     int(req.Id),
	}
	switch req.Type {
	case int32(model.TypeUsepass):
		err = s.cread.Delete(ctx, dto)
	case int32(model.TypeBankCard):
		err = s.bank.Delete(ctx, dto)
	case int32(model.TypePlainText):
		err = s.text.Delete(ctx, dto)
	case int32(model.TypeBinaryData):
		err = s.binary.Delete(ctx, dto)
	default:
		err = errors.New("unknown type")
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *pb.GetDataRequest) (resp *pb.GetDataResponse, err error) {
	resp = &pb.GetDataResponse{}
	//get session and client key by token
	credential := req.Credential

	ses, err := s.getSessionByToken(ctx, credential.Token)
	if err != nil {
		resp.Error = err.Error()
		return resp, nil
	}

	var result []*pb.EncryptedEntity
	if req.Id == 0 {
		result, err = s.list(ctx, req, ses)
	} else {
		result, err = s.get(ctx, req, ses)
	}

	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Data = result
	}

	return resp, nil
}

func (s *Server) get(ctx context.Context, req *pb.GetDataRequest, ses session) ([]*pb.EncryptedEntity, error) {
	var data service.PrivateDataResponse
	var err error

	dto := service.GetPrivateData{
		UserID: ses.userID,
		ID:     int(req.Id),
	}

	switch req.Type {
	case int32(model.TypeUsepass):
		data, err = s.cread.Get(ctx, dto, ses.encoder)
	case int32(model.TypeBankCard):
		data, err = s.bank.Get(ctx, dto, ses.encoder)
	case int32(model.TypePlainText):
		data, err = s.text.Get(ctx, dto, ses.encoder)
	case int32(model.TypeBinaryData):
		data, err = s.binary.Get(ctx, dto, ses.encoder)
	default:
		err = errors.New("unknown type")
	}

	if err != nil {
		return nil, err
	}

	entity := &pb.EncryptedEntity{
		Id: int64(data.ID),
		Data: &pb.EncryptedData{
			Data: []byte(data.Data),
			Dek:  []byte(data.Key),
			View: data.View,
		},
	}

	return []*pb.EncryptedEntity{entity}, err
}

func (s *Server) list(ctx context.Context, req *pb.GetDataRequest, ses session) ([]*pb.EncryptedEntity, error) {
	var data service.ListPrivateDataResponse
	var err error

	switch req.Type {
	case int32(model.TypeUsepass):
		data, err = s.cread.List(ctx, ses.userID, ses.encoder)
	case int32(model.TypeBankCard):
		data, err = s.bank.List(ctx, ses.userID, ses.encoder)
	case int32(model.TypePlainText):
		data, err = s.text.List(ctx, ses.userID, ses.encoder)
	case int32(model.TypeBinaryData):
		data, err = s.binary.List(ctx, ses.userID, ses.encoder)
	default:
		err = errors.New("unknown type")
	}

	if err != nil {
		return nil, err
	}

	response := make([]*pb.EncryptedEntity, len(data.Data))
	for i, d := range data.Data {
		response[i] = &pb.EncryptedEntity{
			Id: int64(d.ID),
			Data: &pb.EncryptedData{
				Data: []byte(d.Data),
				Dek:  []byte(d.Key),
				View: d.View,
			},
		}
	}
	return response, err
}

type session struct {
	userID  int
	encoder encryption.Encoder
}

func (s *Server) getSessionByToken(ctx context.Context, token string) (session, error) {
	ses, err := s.auth.GetSessionByCredentialToken(ctx, token)
	if err != nil {
		return session{}, err
	}
	//create encoder from session public key
	encoder, err := encryption.NewRSACipherFroRawPublicKey(ses.PublicKey)
	return session{ses.UserID, encoder}, err
}

func castMetadata(src []*pb.Metadata) []service.MetadataValue {
	if src == nil {
		return []service.MetadataValue{}
	}

	dst := make([]service.MetadataValue, len(src))
	for i, v := range src {
		dst[i] = service.MetadataValue{Key: v.Key, Value: v.Value}
	}
	return dst
}

func castMetadataPB(src []service.MetadataValue) []*pb.Metadata {
	if src == nil {
		return []*pb.Metadata{}
	}

	dst := make([]*pb.Metadata, len(src))
	for i, v := range src {
		dst[i] = &pb.Metadata{Key: v.Key, Value: v.Value}
	}
	return dst
}
