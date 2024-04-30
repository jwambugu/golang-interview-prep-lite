package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/model"
	"time"
)

var (
	ErrTokenExpired          = errors.New("token is expired")
	ErrTokenMalformed        = errors.New("token is malformed")
	ErrTokenSignatureInvalid = errors.New("token signature is invalid")
)

type JwtManager interface {
	Generate(expiresAt time.Duration, u *model.User) (string, error)
	Validate(token string) (*Claims, error)
}

type jwtToken struct {
	secret []byte
}

type Claims struct {
	jwt.RegisteredClaims
}

func (j *jwtToken) Generate(expiresAt time.Duration, u *model.User) (string, error) {
	var (
		now    = time.Now()
		claims = Claims{
			jwt.RegisteredClaims{
				Issuer:    "app",
				Subject:   u.Username,
				ExpiresAt: jwt.NewNumericDate(now.Add(expiresAt)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        u.ID,
			},
		}
	)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %v", err)
	}

	return token, nil
}

func (j *jwtToken) Validate(accessToken string) (*Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrTokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenExpired
		default:
			return nil, fmt.Errorf("parse token: %v", err)
		}
	}

	if _, ok := token.Claims.(*Claims); !ok {
		return nil, fmt.Errorf("claims assertion: %v", err)
	}

	return &claims, nil
}

func NewJwtTokenManager(secret string) JwtManager {
	return &jwtToken{
		secret: []byte(secret),
	}
}
