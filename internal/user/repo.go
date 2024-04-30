package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/model"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
)

const (
	_queryCreateUser     = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;`
	_queryUsernameExists = `SELECT exists(SELECT 1 FROM users WHERE username = $1)`
	_queryBaseFind       = `SELECT id, username, password FROM users`
	_queryFindByUsername = _queryBaseFind + " WHERE username = $1"
)

var (
	ErrUsernameExists = errors.New("db: username exists")
	ErrRecordNotFound = errors.New("db: record not found")
)

type Repo interface {
	Authenticate(ctx context.Context, username string, password string) (*model.User, error)
	Create(ctx context.Context, u *model.User) error
}

type repo struct {
	db *sql.DB
}

func (r *repo) Authenticate(ctx context.Context, username string, password string) (*model.User, error) {
	var (
		user           model.User
		hashedPassword string
	)

	err := r.db.QueryRowContext(ctx, _queryFindByUsername, username).Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("find by username: %v", err)
	}

	if err = util.CompareHashAndPassword(hashedPassword, password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) Create(ctx context.Context, u *model.User) error {
	var exists bool

	if err := r.db.QueryRowContext(ctx, _queryUsernameExists, u.Username).Scan(&exists); err != nil {
		return fmt.Errorf("username exists: %v", err)
	}

	if exists {
		return ErrUsernameExists
	}

	if err := r.db.QueryRowContext(ctx, _queryCreateUser, u.Username, u.Password).Scan(&u.ID); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}
