package auth

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/auth_service/internal/services/auth"
	authv1 "github.com/Sleeps17/events-planning-service-backend/auth_service/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const emptyValue = 0

type Auth interface {
	Login(ctx context.Context, login string, password string, appId uint32) (token string, err error)
	RegisterNewUser(ctx context.Context, login string, password string) (userID uint64, err error)
	IsAdmin(ctx context.Context, userID uint64) (bool, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	login := req.GetLogin()
	if len(login) < 8 {
		return nil, status.Error(codes.InvalidArgument, "login should be at least 8 characters long")
	}

	password := req.GetPassword()
	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "len password required will be grate then 8")
	}

	appId := req.GetAppId()
	if appId == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "invalid app id")
	}

	token, err := s.auth.Login(ctx, login, password, appId)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		if errors.Is(err, auth.ErrInvalidAppID) {
			return nil, status.Error(codes.InvalidArgument, "invalid app id")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	login := req.GetLogin()
	if len(login) < 8 {
		return nil, status.Error(codes.InvalidArgument, "login should be at least 8 characters long")
	}

	password := req.GetPassword()
	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "len password required will be grate then 8")
	}

	userID, err := s.auth.RegisterNewUser(ctx, login, password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	userID := req.GetUserId()
	if userID == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
