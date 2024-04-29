package user

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
)

var ErrUsernameExists = errors.New("username exists")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) (*Service, error) {
	return &Service{
		db: db,
	}, nil
}

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s *Service) AddUser(u *User) error {
	usernameExists := `SELECT exists(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	if err := s.db.QueryRow(usernameExists, u.Username).Scan(&exists); err != nil {
		return fmt.Errorf("username exists: %v", err)
	}

	if exists {
		return ErrUsernameExists
	}

	hashedPassword, err := util.HashString(u.Password)
	if err != nil {
		return fmt.Errorf("hash password: %v", err)
	}

	u.Password = hashedPassword

	q := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id;`

	if err = s.db.QueryRow(q, u.Username, u.Password).Scan(&u.ID); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}
