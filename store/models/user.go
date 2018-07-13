package models

import (
	"time"

	"github.com/smartcontractkit/chainlink/utils"
)

type User struct {
	Email          string `json:"email" storm:"id,unique"`
	HashedPassword string `json:"hashedPassword"`
	SessionID      string `json:"sessionId", storm:"index,unique"`
	CreatedAt      Time   `json:"createdAt"`
}

func NewUser(email, plainPwd string) (User, error) {
	pwd, err := utils.HashPassword(plainPwd)
	if err != nil {
		return User{}, err
	}

	return User{
		Email:          email,
		HashedPassword: pwd,
		CreatedAt:      Time{Time: time.Now()},
	}, nil
}
