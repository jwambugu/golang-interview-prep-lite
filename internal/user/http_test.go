package user_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/config"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/db"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/testutil"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/user"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testSetup struct {
	db *sql.DB
}

func setupTest(t *testing.T) *testSetup {
	err := config.Load(".env")
	require.NoError(t, err)

	conn, err := db.NewConnection(config.Config.DbDSN)
	require.NoError(t, err)

	return &testSetup{db: conn}
}

func TestHandler_AddUser_Success(t *testing.T) {

	payload := testutil.NewUser()

	reqBody, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	setup := setupTest(t)

	svc, err := user.NewService(setup.db)
	require.NoError(t, err)

	mux := user.Routes(svc)
	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var resp user.User
	err = json.NewDecoder(rr.Body).Decode(&resp)

	require.NoError(t, err)
	require.NotEmpty(t, resp.ID)
	require.Equal(t, payload.Username, resp.Username)

	reqBody, err = json.Marshal(payload)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
	rr = httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestHandler_AddUser_BadRequest_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/user", nil)
	rr := httptest.NewRecorder()

	setup := setupTest(t)

	svc, err := user.NewService(setup.db)
	require.NoError(t, err)

	mux := user.Routes(svc)
	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	rr = httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
