package application

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/deewye/users/gen/db"
	"github.com/deewye/users/internal/config"
	"github.com/deewye/users/internal/server"
	"github.com/deewye/users/internal/service"
	"github.com/deewye/users/internal/storage"
)

type App interface {
	Init() error
	Run() error
	Name() string
	OnShutdown()
	Config() *config.Config
	Storage() storage.Storage
}

type app struct {
	name string
	cfg  *config.Config
	log  *logrus.Logger

	masterDB *sqlx.DB
	slaveDB  *sqlx.DB

	storage storage.Storage
}

func New(name string) App {
	return &app{name: name}
}

func (a *app) Name() string {
	return a.name
}

func (a *app) Config() *config.Config {
	return a.cfg
}

func (a *app) Storage() storage.Storage {
	return a.storage
}

func (a *app) Init() error {
	conf, err := config.InitConfig(a.name)
	if err != nil {
		return fmt.Errorf("init app config: %w", err)
	}

	a.cfg = conf

	a.log = logrus.New()

	a.masterDB, a.slaveDB, err = initMasterSlaveDB(conf.Postgres.Master, conf.Postgres.Slave)
	if err != nil {
		fmt.Printf("init pg: %s", err)
	}

	a.storage = storage.New(db.New(a.masterDB), db.New(a.slaveDB))

	return nil
}

func (a *app) Run() error {
	usersService := service.New(a.log, &a.storage)

	grpcServer := server.NewGRPCServer(a.cfg.GrpcServer)
	grpcServer.Init(a.log)

	grpcServer.RegisterService(usersService)
	if err := grpcServer.Start(); err != nil {
		return err
	}

	a.gracefullyStop()
	a.OnShutdown()

	return nil
}

func (a *app) OnShutdown() {
	if err := a.masterDB.Close(); err != nil {
		fmt.Printf("[OnShutdown] close master conn: %s", err)
	}

	if err := a.slaveDB.Close(); err != nil {
		fmt.Printf("[OnShutdown] close slave conn: %s", err)
	}
}

func initMasterSlaveDB(masterConf, slaveConf *config.DatabaseConfig) (*sqlx.DB, *sqlx.DB, error) {
	masterDB, err := initDB(masterConf)
	if err != nil {
		return nil, nil, fmt.Errorf("init master conn: %w", err)
	}

	slaveDB, err := initDB(slaveConf)
	if err != nil {
		return nil, nil, fmt.Errorf("init slave conn: %w", err)
	}

	return masterDB, slaveDB, nil
}

func initDB(conf *config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", conf.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "error while opening connection")
	}

	db.SetMaxOpenConns(conf.MaxOpenConn)
	db.SetMaxIdleConns(conf.MaxIdleConn)

	return db, nil
}

func (a *app) gracefullyStop() {
	signals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, signals...)
	stop := <-sig
	fmt.Printf("stopping via signal %s", stop)
}
