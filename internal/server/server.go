package server

import (
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/deewye/users/internal/service"
)

type Server interface {
	Init(logger *logrus.Logger) error
	Start() error
	Stop()
	RegisterService(service service.Service)
}

type grpcServer struct {
	cfg *Config
	srv *grpc.Server
	log *logrus.Logger
}

func NewGRPCServer(cfg *Config, opts ...grpc.ServerOption) Server {
	return &grpcServer{
		cfg: cfg,
		srv: grpc.NewServer(opts...),
	}
}

func (gs *grpcServer) Init(logger *logrus.Logger) error {
	gs.log = logger

	return nil
}

func (gs *grpcServer) Start() error {
	lis, err := net.Listen("tcp", gs.cfg.Address)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	go func() {
		gs.log.Infof("GRPC server listening at %s", gs.cfg.Address)
		if err := gs.srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func (gs *grpcServer) Stop() {
	gs.srv.GracefulStop()
}

func (gs *grpcServer) RegisterService(s service.Service) {
	gs.srv.RegisterService(s.GRPCServiceDesc(), s.GetService())
}
