package auth

import (
	"context"
	"regexp"

	authv1 "github.com/Kpatoc452/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

const (
	emptyValue = 0
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}


func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	match, err := regexp.MatchString(".+@.+[.].+", req.GetEmail())
	if err != nil || !match{
		return nil, status.Error(codes.InvalidArgument, "invalid email argument")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid password argument")
	}

	if req.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "invalid AppId argument")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO ...
		return nil, status.Error(codes.InvalidArgument, "error auth login")
	}


	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if req.GetUserId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "error auth IsAdmin")
	}

	return &authv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	match, err := regexp.MatchString(".+@.+[.].+", req.GetEmail())
	if err != nil || !match{
		return nil, status.Error(codes.InvalidArgument, "invalid email argument")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid password argument")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil { 
		// TODO: ...
		return nil, status.Error(codes.Internal, "error auth register")
	}

	return &authv1.RegisterResponse{
		UserId: userID,
	}, nil
}