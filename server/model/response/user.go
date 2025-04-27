package response

import (
	"github.com/gofrs/uuid"
	"server/model/database"
)

type Login struct {
	User                 database.User `json:"user"`
	AccessToken          string        `json:"access_token"`
	AccessTokenExpiresAt int64         `json:"access_token_expires_at"`
}

type UserCard struct {
	UUID      uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	Address   string    `json:"address"`
	Signature string    `json:"signature"`
}
