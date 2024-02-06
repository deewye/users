package handlers

import users "github.com/deewye/users/gen/proto"

type Handler interface {
	users.UnimplementedUsersServer
	users.UsersServer
}
