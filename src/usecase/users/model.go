package users

import "time"

type GetAuthTokenReq struct {
	Provider     string `json:"provider"`
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
	SwiftToken   string `json:"swift_token"`
	DevTokenKey  bool   `json:"dev_token_key,omitempty"`
}

type GetAuthTokenResp struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

type GetUserProfileResp struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}
