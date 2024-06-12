package handlers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	pb "github.com/impr0ver/gophKeeper/internal/rpc"

	"github.com/impr0ver/gophKeeper/internal/userdata"
	"github.com/impr0ver/gophKeeper/internal/storage"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ClientConnGPRC keeps connection with server. Uses gRPC.
type ClientConnGPRC struct {
	pb.GokeeperClient
}

func clientLoadTLSCredentials(clientCert string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile(clientCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

// NewClientConnection connects to server and returning connection.
func newClientConn(serverAddress, clientCert string) *ClientConnGPRC {
	tlsCredentials, err := clientLoadTLSCredentials(clientCert)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.NewClient("passthrough:///"+serverAddress, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatal(err)
	}

	return &ClientConnGPRC{
		GokeeperClient: pb.NewGokeeperClient(conn),
	}
}

// Login logins user by login and password.
func (c *ClientConnGPRC) Login(credentials userdata.UserCredentials) (string, error) {
	session, err := c.GokeeperClient.Login(context.Background(), &pb.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	switch status.Code(err) {
	case codes.Unauthenticated:
		return "", storage.ErrWrongCredentials
	case codes.Internal:
		return "", storage.ErrUnknown
	case codes.InvalidArgument:
		return "", ErrEmptyField
	}

	if err != nil {
		log.Warnf("%s :: %v", "login fault", err)

		return "", err
	}

	return session.SessionToken, nil
}

// Register creates new user by login and password.
func (c *ClientConnGPRC) Register(credentials userdata.UserCredentials) (string, error) {
	session, err := c.GokeeperClient.Register(context.Background(), &pb.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	code := status.Code(err)

	switch code {
	case codes.AlreadyExists:
		return "", storage.ErrLoginExists
	case codes.Internal:
		return "", storage.ErrUnknown
	case codes.InvalidArgument:
		return "", ErrEmptyField
	}

	return session.SessionToken, nil
}

// GetRecordsInfo gets all record.
func (c *ClientConnGPRC) GetRecordsInfo(token userdata.AuthToken) ([]userdata.Record, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	gotRecords, err := c.GokeeperClient.GetRecordsInfo(ctx, &emptypb.Empty{})
	code := status.Code(err)

	switch code {
	case codes.Internal:
		return nil, storage.ErrUnknown
	case codes.Unauthenticated:
		return nil, storage.ErrUnauthenticated
	}

	records := make([]userdata.Record, 0, len(gotRecords.Records))

	for _, record := range gotRecords.Records {
		records = append(records, userdata.Record{
			ID:       record.Id,
			Metadata: record.Metadata,
			KeyHint:  record.Keyhint,
			Type:     userdata.RecordType(record.Type),
		})
	}

	return records, nil
}

// GetRecord gets record from server by ID.
func (c *ClientConnGPRC) GetRecord(token userdata.AuthToken, recordID string) (userdata.Record, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	gotRecord, err := c.GokeeperClient.GetRecord(ctx, &pb.RecordID{
		Id: recordID,
	})
	record, code := userdata.Record{}, status.Code(err)

	switch code {
	case codes.Internal:
		return record, storage.ErrUnknown
	case codes.Unauthenticated:
		return record, storage.ErrUnauthenticated
	case codes.NotFound:
		return record, storage.ErrNotFound
	}

	record = userdata.Record{
		ID:       gotRecord.Id,
		Metadata: gotRecord.Metadata,
		KeyHint:  gotRecord.Keyhint,
		Type:     userdata.RecordType(gotRecord.Type),
		Data:     gotRecord.StoredData,
	}
	return record, nil
}

// DeleteRecord deletes record from server by ID.
func (c *ClientConnGPRC) DeleteRecord(token userdata.AuthToken, recordID string) error {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	_, err := c.GokeeperClient.DeleteRecord(ctx, &pb.RecordID{
		Id: recordID,
	})
	code := status.Code(err)

	switch code {
	case codes.Internal:
		return storage.ErrUnknown
	case codes.Unauthenticated:
		return storage.ErrUnauthenticated
	case codes.NotFound:
		return storage.ErrNotFound
	}

	return nil
}

// CreateRecord creates record and saves to server.
func (c *ClientConnGPRC) CreateRecord(token userdata.AuthToken, record userdata.Record) error {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	_, err := c.GokeeperClient.CreateRecord(ctx, &pb.Record{
		Type:       pb.MessageType(record.Type),
		Keyhint:    record.KeyHint,
		Metadata:   record.Metadata,
		StoredData: record.Data,
	})

	switch status.Code(err) {
	case codes.Internal:
		return storage.ErrUnknown
	case codes.Unauthenticated:
		return storage.ErrUnauthenticated
	}

	return nil
}
