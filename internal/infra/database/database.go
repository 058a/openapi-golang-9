package database

import (
	"database/sql"
	"openapi/internal/infra/env"

	_ "github.com/lib/pq"
)

func Open() (*sql.DB, error) {
	driverName := env.GetDbDriver()
	dataSourceName := env.GetDbDataSourceName()

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
