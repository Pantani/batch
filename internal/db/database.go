package db

import (
	"github.com/Pantani/batch/internal/config"
	"github.com/Pantani/batch/internal/db/database"
	"github.com/Pantani/batch/internal/db/memory"
	"github.com/Pantani/batch/internal/db/redis"

	"github.com/Pantani/errors"
)

type (
	// Type represents the database type.
	Type string
	// Db represents the database object.
	Db struct {
		client database.IDatabase
	}
)

const (
	// InMemory represents the in memory database type.
	InMemory Type = "memory"
	// Redis represents the redis database type.
	Redis Type = "redis"
)

// Init create a new database object based in type.
// It returns a database object and an error if occurs.
func Init(dbType Type) (database.IDatabase, error) {
	switch dbType {
	case InMemory:
		return memory.New(), nil
	case Redis:
		return redis.New(
			config.Configuration.Database.Redis.Host,
			config.Configuration.Database.Redis.Password,
			config.Configuration.Database.Redis.Index,
		)
	default:
		return nil, errors.E("invalid database type", errors.Params{"type": dbType})
	}
}
