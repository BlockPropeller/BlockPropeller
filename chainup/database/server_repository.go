package database

import (
	"context"

	"chainup.dev/chainup/infrastructure"
	"github.com/pkg/errors"
)

// ServerRepository is a databased backed implementation of a infrastructure.ServerRepository.
type ServerRepository struct {
	db *DB
}

// NewServerRepository returns a new ServerRepository instance.
func NewServerRepository(db *DB) *ServerRepository {
	return &ServerRepository{db: db}
}

// Find a Server given a ServerID.
func (repo *ServerRepository) Find(ctx context.Context, id infrastructure.ServerID) (*infrastructure.Server, error) {
	var server infrastructure.Server
	err := repo.db.Model(ctx, &server).Where("id = ?", id).First(&server).Error
	if err != nil {
		return nil, errors.Wrap(err, "find server by ID")
	}

	return &server, nil
}

// Create a new Server.
func (repo *ServerRepository) Create(ctx context.Context, server *infrastructure.Server) error {
	err := repo.db.Model(ctx, server).Create(server).Error
	if err != nil {
		return errors.Wrap(err, "create server")
	}

	return nil
}

// Update an existing Server.
func (repo *ServerRepository) Update(ctx context.Context, server *infrastructure.Server) error {
	err := repo.db.Model(ctx, server).Save(server).Error
	if err != nil {
		return errors.Wrap(err, "update server")
	}

	return nil
}
