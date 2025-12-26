package auth

import (
	"context"
	"log/slog"
	"sso/internal/domain/models"
	"time"
)

type AuthService struct {
	log *slog.Logger
	userSaver UserSaver
	userProvider UserProvider
	appProvider AppProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New возвращает новый инстанс AuthService
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		log: log,
		userSaver: userSaver,
		userProvider: userProvider,
		appProvider: appProvider,
		tokenTTL: tokenTTL,
	}
} 

func (a *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	panic("not implemented")
}

func (a *AuthService) RegisterUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	panic("not implemented")
}

func (a *AuthService) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	panic("not implemented")
} 