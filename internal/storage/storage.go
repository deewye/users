package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/deewye/users/gen/db"
)

type Storage interface {
	InsertUser(ctx context.Context, params db.InsertUserParams) error
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type storage struct {
	masterQueries *db.Queries
	slaveQueries  *db.Queries
}

func New(masterDB, slaveDB *db.Queries) Storage {
	return &storage{
		masterQueries: masterDB,
		slaveQueries:  slaveDB,
	}
}

func (s *storage) InsertUser(ctx context.Context, params db.InsertUserParams) error {
	return s.masterQueries.InsertUser(ctx, params)
}

func (s *storage) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return s.slaveQueries.GetUserByID(ctx, id)
}
