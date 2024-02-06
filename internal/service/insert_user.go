package service

import (
	"context"

	users "github.com/deewye/users/gen/proto"
)

func (s *userService) InsertUser(context.Context, users.InsertUserRequest) (*users.InsertUserResponse, error) {
	return nil, nil
}
