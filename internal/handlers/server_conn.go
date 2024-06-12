package handlers

import (
	"context"
	"crypto/tls"
	"errors"
	"net"

	"github.com/impr0ver/gophKeeper/internal/logger"
	pb "github.com/impr0ver/gophKeeper/internal/rpc"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ServerConn keeps server endpoints alive.
type ServerConn struct {
	pb.UnimplementedGokeeperServer
	Handlers         ServerHandlers
	server           *grpc.Server
	Authenticator    Authenticator
	ServerCert       string
	ServerKey        string
	ServerConsoleLog bool
}

// NewServerConn returns new server connection.
func NewServerConn(h ServerHandlers, a Authenticator, serverCert string, serverKey string, serverConsoleLog bool) *ServerConn {
	return &ServerConn{
		Handlers:         h,
		Authenticator:    a,
		ServerCert:       serverCert,
		ServerKey:        serverKey,
		ServerConsoleLog: serverConsoleLog,
	}
}

// loadTLSCredentials load certificates and provate key.
func (s *ServerConn) loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(s.ServerCert, s.ServerKey)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

// Run runs server listener.
func (s *ServerConn) Run(_ context.Context, runAddress string) {
	sLogger := logger.NewSugarLogger()
	listen, err := net.Listen("tcp", runAddress)
	if err != nil {
		log.Fatal(err)
	}

	tlsCredentials, err := s.loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	grpcServ := grpc.NewServer(grpc.Creds(tlsCredentials), grpc.ChainUnaryInterceptor(grpc.UnaryServerInterceptor(s.LoggingInterceptor),
		grpc.UnaryServerInterceptor(s.VerifyAuth())))

	pb.RegisterGokeeperServer(grpcServ, s)

	go func() {
		log.Println("gRPC server is start...")
		sLogger.Info("gRPC server is start...")
		
		if err := grpcServ.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	s.server = grpcServ
}

func (s *ServerConn) Stop() {
	s.server.GracefulStop()
	log.Println("Shutdown server gracefully.")
}

// Register process register endpoint.
func (s *ServerConn) Register(_ context.Context, credentials *pb.UserCredentials) (*pb.Session, error) {
	token, err := s.Handlers.CreateUser(userdata.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	if errors.Is(err, ErrEmptyField) {
		log.Infoln(err)

		return nil, status.Errorf(codes.InvalidArgument, "login or password is empty.")
	}

	if errors.Is(err, storage.ErrLoginExists) {
		log.Infoln(err)

		return nil, status.Errorf(codes.AlreadyExists, "login already exists.")
	}

	if err != nil {
		log.Warnf("%s %s :: %v", "register new user fault", credentials.Login, err)

		return nil, status.Errorf(codes.Internal, "internal server error.")
	}

	return &pb.Session{SessionToken: string(token)}, nil
}

// Login process login endpoint.
func (s *ServerConn) Login(_ context.Context, credentials *pb.UserCredentials) (*pb.Session, error) {
	token, err := s.Handlers.LoginUser(userdata.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	if errors.Is(err, ErrEmptyField) {
		log.Infoln(err)

		return nil, status.Errorf(codes.InvalidArgument, "login or password is empty.")
	}

	if errors.Is(err, storage.ErrWrongCredentials) {
		log.Infoln(err)

		return nil, status.Errorf(codes.Unauthenticated, "wrong login or password.")
	}

	if err != nil {
		log.Warnf("%s %s :: %v", "login fault", credentials.Login, err)

		return nil, status.Errorf(codes.Internal, "Internal server error.")
	}

	return &pb.Session{SessionToken: string(token)}, nil
}

// GetRecordsInfo process get all records endpoint.
func (s *ServerConn) GetRecordsInfo(ctx context.Context, _ *emptypb.Empty) (*pb.RecordsList, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "metadata for authentication not found.")
	}

	records, err := s.Handlers.GetRecordsInfo(ctx)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return nil, status.Errorf(codes.Unauthenticated, "bad token.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "get record info fault", err)

		return nil, status.Errorf(codes.Internal, "internal server error.")
	}

	recordsList := make([]*pb.Record, 0, len(records))

	for _, record := range records {
		recordsList = append(recordsList, &pb.Record{
			Id:       record.ID,
			Metadata: record.Metadata,
			Keyhint:  record.KeyHint,
			Type:     pb.MessageType(record.Type),
		})
	}

	return &pb.RecordsList{Records: recordsList}, nil
}

// GetRecord process get record endpoint.
func (s *ServerConn) GetRecord(ctx context.Context, recordID *pb.RecordID) (*pb.Record, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "metadata for authentication not found.")
	}

	record, err := s.Handlers.GetRecord(ctx, recordID.Id)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return nil, status.Errorf(codes.Unauthenticated, "bad token.")
	}

	if errors.Is(err, storage.ErrNotFound) {
		log.Infoln(err)

		return nil, status.Errorf(codes.NotFound, "not found record by id.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "get record fault", err)

		return nil, status.Errorf(codes.Internal, "internal server error.")
	}

	return &pb.Record{
		Id:         record.ID,
		Type:       pb.MessageType(record.Type),
		Keyhint:    record.KeyHint,
		Metadata:   record.Metadata,
		StoredData: record.Data,
	}, nil
}

// CreateRecord process create record endpoint.
func (s *ServerConn) CreateRecord(ctx context.Context, record *pb.Record) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "metadata for authentication not found.")
	}

	err := s.Handlers.CreateRecord(ctx, userdata.Record{
		Metadata: record.Metadata,
		KeyHint:  record.Keyhint,
		Type:     userdata.RecordType(record.Type),
		Data:     record.StoredData,
	})

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "bad token.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "create record fault", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "internal server error.")
	}

	return &emptypb.Empty{}, nil
}

// DeleteRecord process delete record endpoint.
func (s *ServerConn) DeleteRecord(ctx context.Context, recordID *pb.RecordID) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "metadata for authentication not found.")
	}

	err := s.Handlers.DeleteRecord(ctx, recordID.Id)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "bad token.")
	}

	if errors.Is(err, storage.ErrNotFound) {
		log.Infoln(err)

		return nil, status.Errorf(codes.NotFound, "not found record by id.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "delete record fault", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "internal server error.")
	}

	return &emptypb.Empty{}, nil
}
