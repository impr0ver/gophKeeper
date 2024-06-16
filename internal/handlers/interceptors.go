package handlers

import (
	"context"

	"github.com/impr0ver/gophKeeper/internal/logger"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logging some data on interceptor.
func (s *ServerConn) LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var sLogger = logger.NewSugarLogger()

	if s.ServerConsoleLog {
		sLogger.Infof("FullMethod: %s, Received request: %v", info.FullMethod, req)
	}

	resp, err := handler(ctx, req)
	return resp, err
}

// VerifyAuth check authentication token on interceptor.
func (s *ServerConn) VerifyAuth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var sLogger = logger.NewSugarLogger()
		var token userdata.AuthToken

		md, ok := metadata.FromIncomingContext(ctx)

		if s.ServerConsoleLog {
			sLogger.Info("Hooked MD on interceptor: ", md.String())
		}

		if ok && len(md.Get("authToken")) > 0 {
			token = userdata.AuthToken(md.Get("authToken")[0])

			if s.ServerConsoleLog {
				sLogger.Info("Hooked token: ", token)
			}

			userIDValid, err := s.Authenticator.ValidateToken(userdata.AuthToken(token))
			if err != nil {
				log.Warnf("%s :: %v", "interceptor validate token error", err)
				if s.ServerConsoleLog {
					sLogger.Infof("%s :: %v", "interceptor validate token error", err)
				}

				return nil, status.Errorf(codes.Unauthenticated, "validate token error :: %v", err)
			}

			// Add validated userID in context
			md.Append("userID", string(userIDValid))
			ctx = metadata.NewIncomingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}
