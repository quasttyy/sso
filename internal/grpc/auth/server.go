package auth

import (
	"context"

	ssov1 "github.com/quasttyy/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		app_id int32,
	) (string, error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (int, error)
	IsAdmin(ctx context.Context, user_id int) (bool error)
}

// serverAPI будет реализовывать интерфейс AuthService
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

// Register связывает RPC методы с методами структуры serverAPI
func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
} 

// Login - хэндлер для авторизации
func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	// Валидация:
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	// todo: implement login via auth service
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		// todo: error handling
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me!")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me!")
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "mail is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}