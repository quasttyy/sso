package storage

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		db: pool,
	}
}

const saveUserQuery = `
INSERT INTO users(email, pass_hash)
VALUES ($1, $2)
RETURNING id;
`

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (int64, error) {
	const op = "Storage.SaveUser"
	var userID int64

	err := s.db.QueryRow(
		ctx,
		saveUserQuery,
		email,
		passHash).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

const getUserQuery = `
SELECT id, email, pass_hash
FROM users
WHERE email = $1;
`

func (s *Storage) User(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "Storage.User"
	var user models.User

	err := s.db.QueryRow(
		ctx,
		getUserQuery,
		email).Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

const isAdminQuery = `
SELECT is_admin
FROM users
WHERE id = $1;
`

func (s *Storage) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "Storage.IsAdmin"
	var isAdmin bool

	err := s.db.QueryRow(ctx, isAdminQuery, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

const getAppQuery = `
SELECT id, app_name, app_secret
FROM apps
WHERE id = $1;
`

func (s *Storage) App(
	ctx context.Context,
	appID int64,
) (models.App, error) {
	const op = "Storage.App"
	var app models.App

	err := s.db.QueryRow(ctx,
		getAppQuery,
		appID).Scan(&app.App_id, &app.Name, &app.Secret)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
