package auth_test

import (
	"github.com/matthewjamesboyle/golang-interview-prep/internal/auth"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/testutil"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwtToken_Generate_Validate(t *testing.T) {
	var (
		jwtToken  = auth.NewJwtTokenManager("secret")
		testUser  = testutil.NewUser()
		expiresAt = 10 * time.Second
	)

	token, err := jwtToken.Generate(expiresAt, testUser)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwtToken.Validate(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims)
	require.Equal(t, testUser.ID, claims.ID)
	require.WithinDuration(t, time.Now(), claims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, time.Now().Add(expiresAt), claims.ExpiresAt.Time, time.Second)
}

func TestJwtToken_Generate_ExpiredToken(t *testing.T) {
	var (
		jwtToken = auth.NewJwtTokenManager("secret")
		testUser = testutil.NewUser()
	)

	token, err := jwtToken.Generate(-time.Hour, testUser)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwtToken.Validate(token)
	require.EqualError(t, err, auth.ErrTokenExpired.Error())
	require.Nil(t, claims)

	claims, err = jwtToken.Validate("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	require.EqualError(t, err, auth.ErrTokenMalformed.Error())
	require.Nil(t, claims)
}
