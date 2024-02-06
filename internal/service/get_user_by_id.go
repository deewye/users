package service

import (
	"context"

	users "github.com/deewye/users/gen/proto"
)

func (s *userService) GetUserByID(context.Context, *users.GetUserByIDRequest) (*users.User, error) {
	return nil, nil
}
