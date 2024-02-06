package service

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	users "github.com/deewye/users/gen/proto"
	"github.com/deewye/users/internal/storage"
)

type Service interface {
	GRPCServiceDesc() *grpc.ServiceDesc
	GetService() any
}

type userService struct {
	users.UnimplementedUsersServer

	storage *storage.Storage
	log     *logrus.Logger
}

func New(log *logrus.Logger, storage *storage.Storage) Service {
	return &userService{
		storage: storage,
		log:     log,
	}
}

func (s userService) GRPCServiceDesc() *grpc.ServiceDesc {
	return &users.Users_ServiceDesc
}

func (s *userService) GetService() any {
	return s
}
