package user_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/auth"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/config"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/db"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/model"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/testutil"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/user"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testSetup struct {
	db         *sql.DB
	jwtManager auth.JwtManager
}

func setupTest(t *testing.T) *testSetup {
	err := config.Load(".env")
	require.NoError(t, err)

	conn, err := db.NewConnection(config.Config.DbDSN)
	require.NoError(t, err)

	return &testSetup{
		db:         conn,
		jwtManager: auth.NewJwtTokenManager("secret"),
	}
}

func TestHandler_AddUser_Success(t *testing.T) {
	var (
		setup    = setupTest(t)
		userRepo = user.NewRepo(setup.db)
	)

	svc, err := user.NewService(setup.jwtManager, userRepo)
	require.NoError(t, err)

	userHandler := user.NewHandler(setup.jwtManager, svc, userRepo)

	mux := http.NewServeMux()
	userHandler.Routes(mux)

	var (
		testUser = testutil.NewUser()
		admin    = testutil.NewUser()
		ctx      = context.Background()
	)

	err = userRepo.Create(ctx, admin)
	require.NoError(t, err)

	reqBody, err := json.Marshal(testUser)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))

	accessToken, err := setup.jwtManager.Generate(time.Minute, admin)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var resp model.User
	err = json.NewDecoder(rr.Body).Decode(&resp)

	require.NoError(t, err)
	require.NotEmpty(t, resp.ID)
	require.Equal(t, testUser.Username, resp.Username)

	reqBody, err = json.Marshal(testUser)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr = httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestHandler_AddUser_BadRequest_MethodNotAllowed(t *testing.T) {
	setup := setupTest(t)
	userRepo := user.NewRepo(setup.db)

	svc, err := user.NewService(setup.jwtManager, userRepo)
	require.NoError(t, err)

	userHandler := user.NewHandler(setup.jwtManager, svc, userRepo)

	mux := http.NewServeMux()
	userHandler.Routes(mux)

	req := httptest.NewRequest(http.MethodPost, "/user", nil)

	var (
		admin = testutil.NewUser()
		ctx   = context.Background()
	)

	err = userRepo.Create(ctx, admin)
	require.NoError(t, err)

	accessToken, err := setup.jwtManager.Generate(time.Minute, admin)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr = httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestHandler_Authenticate(t *testing.T) {
	var (
		setup    = setupTest(t)
		userRepo = user.NewRepo(setup.db)
		testUser = testutil.NewUser()
	)

	password := testUser.Password

	var err error

	testUser.Password, err = util.HashString(testUser.Password)
	require.NoError(t, err)

	ctx := context.Background()

	err = userRepo.Create(ctx, testUser)
	require.NoError(t, err)

	svc, err := user.NewService(setup.jwtManager, userRepo)
	require.NoError(t, err)

	userHandler := user.NewHandler(setup.jwtManager, svc, userRepo)

	mux := http.NewServeMux()
	userHandler.Routes(mux)

	authReq := user.AuthenticateReq{
		Username: testUser.Username,
		Password: password,
	}

	reqBody, err := json.Marshal(authReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp user.AuthenticateResp

	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, resp.User.Username)
	require.Empty(t, resp.User.Password)
}
