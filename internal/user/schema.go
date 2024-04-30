package user

import "github.com/matthewjamesboyle/golang-interview-prep/internal/model"

type AuthenticateReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthenticateResp struct {
	User        *model.User `json:"user,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}
