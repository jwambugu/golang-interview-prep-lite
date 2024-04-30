package user

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/auth"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/model"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
	"time"
)

type Service interface {
	Authenticate(ctx context.Context, username string, password string) (*AuthenticateResp, error)
	Create(ctx context.Context, u *model.User) error
}

type service struct {
	jwtManager auth.JwtManager
	repo       Repo
}

func (s *service) Authenticate(ctx context.Context, username string, password string) (*AuthenticateResp, error) {
	user, err := s.repo.Authenticate(ctx, username, password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtManager.Generate(time.Hour, user)
	if err != nil {
		return nil, err
	}

	return &AuthenticateResp{
		AccessToken: accessToken,
		User:        user,
	}, nil
}

func (s *service) Create(ctx context.Context, u *model.User) error {
	hashedPassword, err := util.HashString(u.Password)
	if err != nil {
		return fmt.Errorf("hash password: %v", err)
	}

	u.Password = hashedPassword

	if err = s.repo.Create(ctx, u); err != nil {
		return err
	}

	return nil
}

func NewService(jwtManager auth.JwtManager, repo Repo) (Service, error) {
	return &service{
		jwtManager: jwtManager,
		repo:       repo,
	}, nil
}
